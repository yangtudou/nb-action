package registry_sync

import (
	"context"
	"sync"
	"time"

	"github.com/yangtudou/nb-action/internal/logger"
	"github.com/yangtudou/nb-action/internal/parser"
	"github.com/yangtudou/nb-action/internal/resolver"
	"github.com/yangtudou/nb-action/internal/result"
	"github.com/yangtudou/nb-action/internal/retry"
	"github.com/yangtudou/nb-action/internal/syncer"
	"github.com/yangtudou/nb-action/internal/worker"
)

func Run(
	ctx context.Context,
	opt *Options,
) (map[string]interface{}, error) {

	err := PrepareAuth()

	if err != nil {
		return nil, err
	}

	if opt.Concurrency < 1 {
		opt.Concurrency = 1
	}

	var images []string

	if opt.Src != "" {

		images = []string{
			opt.Src,
		}

	} else {

		images, err = parser.ReadImageList(
			opt.Base,
		)

		if err != nil {
			return nil, err
		}
	}

	stats := result.New(
		len(images),
	)

	var mu sync.Mutex

	targets := make(
		[]string,
		0,
		len(images),
	)

	tasks := make(
		[]worker.Task,
		0,
		len(images),
	)

	for i, img := range images {

		index := i + 1
		image := img

		tasks = append(
			tasks,
			func() error {

				source := resolver.Resolve(
					image,
					resolver.Rule{
						Prefix:  opt.SrcPrefix,
						Flatten: opt.SrcFlatten,
					},
				)

				target := resolver.Resolve(
					image,
					resolver.Rule{
						Prefix:  opt.DstPrefix,
						Flatten: opt.DstFlatten,
					},
				)

				if opt.DryRun {

					logger.Printf(
						"[%d/%d] %s -> %s",
						index,
						len(images),
						source,
						target,
					)

				} else {

					err := retry.Do(
						func() error {

							return syncer.CopyWithPlatform(
								source,
								target,
								opt.Platform,
							)

						},
						opt.Retries,
						2*time.Second,
					)

					if err != nil {

						stats.AddFailed()

						return err
					}
				}

				stats.AddSuccess()

				mu.Lock()

				targets = append(
					targets,
					target,
				)

				mu.Unlock()

				return nil
			},
		)
	}

	err = worker.Run(
		tasks,
		opt.Concurrency,
	)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"status": "ok",

		"total": stats.Total,

		"success": stats.Success,

		"failed": stats.Failed,

		"duration_ms": stats.Duration().Milliseconds(),

		"value": map[string]interface{}{
			"images": targets,
		},
	}, nil
}
