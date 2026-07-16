package actions

import (
	"context"
	"flag"
	"fmt"
	"io"
	"sync" // 确保并发 append 绝对安全
	"time"

	"github.com/yangtudou/nb-action/internal/logger"
	"github.com/yangtudou/nb-action/internal/parser"
	"github.com/yangtudou/nb-action/internal/resolver"
	"github.com/yangtudou/nb-action/internal/result"
	"github.com/yangtudou/nb-action/internal/retry"
	"github.com/yangtudou/nb-action/internal/syncer"
	"github.com/yangtudou/nb-action/internal/worker"
)

type RegistrySync struct{}

// Name 实现 Action 接口规范，声明当前动作的命令名称
func (a *RegistrySync) Name() string {
	return "registry-sync"
}

func (r *RegistrySync) Description() string {
	return "Docker 镜像仓库同步"
}

func (r *RegistrySync) Help() string {
	return `
registry-sync

Docker 镜像仓库同步工具

Usage:
  nb-action registry-sync [options]

Required:
  --dst-prefix <registry>
        目标镜像仓库

Options:
  --src-prefix <registry>
        源镜像仓库

  --file <path>
        镜像列表文件

  --src-flatten
        源仓库路径扁平化

  --dst-flatten
        目标仓库路径扁平化

  --concurrency <number>
        并发数量

  --retries <number>
        重试次数

  --dry-run
        仅显示同步计划，不执行

Examples:
  nb-action registry-sync \
    --src-prefix docker.io \
    --dst-prefix registry.example.com \
    --file images.txt

`
}

func (a *RegistrySync) Execute(ctx context.Context, args []string, input map[string]interface{}) (map[string]interface{}, error) {
	// 1. 使用局部的 FlagSet，命令名对应 "registry-sync"
	fs := flag.NewFlagSet("registry-sync", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	src := fs.String("src", "", "Source image")
	base := fs.String("base", "images.txt", "Base image manifest file")
	srcPrefix := fs.String("src-prefix", "", "Source registry prefix")
	srcFlatten := fs.Bool("src-flatten", false, "Flatten source image path")
	dstPrefix := fs.String("dst-prefix", "", "Destination registry prefix")
	dstFlatten := fs.Bool("dst-flatten", false, "Flatten destination image path")
	platform := fs.String("platform", "", "Target platform (e.g., linux/amd64)")
	concurrency := fs.Int("concurrency", 4, "Concurrent workers")
	retries := fs.Int("retries", 3, "Retry attempts")
	dryRun := fs.Bool("dry-run", false, "Only show actions")

	if err := fs.Parse(args); err != nil {
		return nil, fmt.Errorf("invalid registry-sync arguments: %w", err)
	}

	// 2. 管道智能适配：如果命令行没给 --src，自动从上游拿到 value 或 image 字段
	if *src == "" {
		if val, ok := input["value"].(string); ok && val != "" {
			*src = val
		} else if val, ok := input["image"].(string); ok && val != "" {
			*src = val
		}
	}

	// 3. 防御性校验
	if *dstPrefix == "" {
		return nil, fmt.Errorf("missing required flag: --dst-prefix")
	}
	if *concurrency < 1 {
		*concurrency = 1
	}

	// 4. 加载镜像列表
	var images []string
	if *src != "" {
		images = []string{*src}
	} else {
		var err error
		images, err = parser.ReadImageList(*base)
		if err != nil {
			return nil, fmt.Errorf("load image list from '%s' failed: %w", *base, err)
		}
	}

	if len(images) == 0 {
		return map[string]interface{}{
			"status":  "skipped",
			"message": "no images found to sync",
		}, nil
	}

	total := len(images)
	stats := result.New(total)
	var tasks []worker.Task

	// 初始化互斥锁和收集成功目标的切片
	var mu sync.Mutex
	syncedDestinations := make([]string, 0, total)

	for i, img := range images {
		index := i + 1
		image := img // 闭包安全

		tasks = append(tasks, func() error {
			source := resolver.Resolve(image, resolver.Rule{Prefix: *srcPrefix, Flatten: *srcFlatten})
			target := resolver.Resolve(image, resolver.Rule{Prefix: *dstPrefix, Flatten: *dstFlatten})

			if *dryRun {
				logger.Printf("[%d/%d] Dry Run: %s -> %s\n", index, total, source, target)
				stats.AddSuccess()

				// 加锁写入
				mu.Lock()
				syncedDestinations = append(syncedDestinations, target)
				mu.Unlock()
				return nil
			}

			err := retry.Do(func() error {
				return syncer.CopyWithPlatform(source, target, *platform)
			}, *retries, 2*time.Second)

			if err != nil {
				stats.AddFailed()
				logger.Printf("[%d/%d] Failed: %s (%v)\n", index, total, image, err)
				return err
			}

			stats.AddSuccess()
			logger.Printf("[%d/%d] Success: %s\n", index, total, image)

			// 加锁写入
			mu.Lock()
			syncedDestinations = append(syncedDestinations, target)
			mu.Unlock()
			return nil
		})
	}

	// 5. 并发执行
	_ = worker.Run(tasks, *concurrency)

	// 6. 整理需要放入 value 的嵌套 Map 数据
	outputData := map[string]interface{}{
		"status":      "ok",
		"total":       stats.Total,
		"success":     stats.Success,
		"failed":      stats.Failed,
		"duration_ms": stats.Duration().Milliseconds(),
		"images":      syncedDestinations,
	}

	// 7. 返回结果
	res := map[string]interface{}{
		"status":      "ok",
		"total":       stats.Total,
		"success":     stats.Success,
		"failed":      stats.Failed,
		"duration_ms": stats.Duration().Milliseconds(),

		// 💡 重点：不转 String，直接把 Map 塞给 value，消灭反斜杠转义
		"value": outputData,
	}

	if stats.Failed > 0 {
		return res, fmt.Errorf("sync completed with %d failures", stats.Failed)
	}

	return res, nil
}
