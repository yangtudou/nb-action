package core

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// parseByComma 升级版：完美支持 "wg-keypair,bark"、"wg-keypair, bark"、"wg-keypair , bark" 等所有人类迷惑行为
func ParseByComma(args []string) [][]string {
	var steps [][]string
	var current []string

	for _, arg := range args {
		// 核心改进：直接用逗号对每个参数进行二次切分
		parts := strings.Split(arg, ",")
		for i, part := range parts {
			trimmed := strings.TrimSpace(part)

			// 如果 index > 0，说明我们刚刚跨过了一个“逗号”
			if i > 0 {
				if len(current) > 0 {
					steps = append(steps, current)
					current = []string{} // 重置，准备收集下一个 Action
				}
			}

			if trimmed != "" {
				current = append(current, trimmed)
			}
		}
	}
	if len(current) > 0 {
		steps = append(steps, current)
	}
	return steps
}

// formatValue 智能格式化传给下游占位符的 value
func formatValue(val interface{}) string {
	if val == nil {
		return ""
	}

	switch v := val.(type) {
	case string:
		return v
	case map[string]interface{}, []interface{}, []string:
		// 🟢 核心魔法：如果上游传下来的是 Map 或 数组，在 CLI 替换渲染时自动转成无换行的 JSON 字符串
		bytes, err := json.Marshal(v)
		if err == nil {
			return string(bytes)
		}
		return fmt.Sprintf("%v", v)
	default:
		// 其他基础类型（数字、布尔）直接按原样输出字符串
		return fmt.Sprintf("%v", v)
	}
}

// resolveArgs 负责在执行前，将参数中形如 {key} 的占位符，动态替换为上游输出的值
func resolveArgs(args []string, input map[string]interface{}) []string {
	resolved := make([]string, len(args))

	for i, arg := range args {
		temp := arg
		// 遍历上游输出的每一个 key，动态替换占位符
		for k, v := range input {
			placeholder := "{" + k + "}"
			if strings.Contains(temp, placeholder) {
				// 💡 采用智能格式化函数，避免把 map 渲染成极其难看的 "map[...]"
				strVal := formatValue(v)
				temp = strings.ReplaceAll(temp, placeholder, strVal)
			}
		}
		resolved[i] = temp
	}

	return resolved
}

// runPipeline 内存直接传递 Map，并且【合并】每一步的输出，防止数据在传递中丢失
func (r *Runtime) RunPipeline(
	ctx context.Context,
	steps [][]string,
	initialInput map[string]interface{},
) {
	// 初始化累积 Context，防止空指针，并将初始输入深拷贝进去
	currentInput := make(map[string]interface{})
	if initialInput != nil {
		for k, v := range initialInput {
			currentInput[k] = v
		}
	}

	for _, step := range steps {
		if len(step) == 0 {
			continue
		}
		actionName := step[0]
		actionArgs := step[1:]

		action, ok := r.Registry.GetAction(actionName)
		if !ok {
			WriteError(fmt.Errorf("action not found in pipeline: %s", actionName))
			os.Exit(1)
		}

		// 1. 用累积的 currentInput 动态渲染当前步骤 learnings (参数替换)
		resolvedArgs := resolveArgs(actionArgs, currentInput)

		// 2. 接力执行（此时的 currentInput 包含了前面所有步骤累积下来的所有 key-value 对）
		result, err := action.Execute(ctx, resolvedArgs, currentInput)
		if err != nil {
			WriteError(fmt.Errorf("pipeline step [%s] failed: %w", actionName, err))
			os.Exit(1)
		}

		// 3. 增量合并
		if result != nil {
			for k, v := range result {
				currentInput[k] = v
			}
		}
	}

	// 管道执行完毕，打印最终累积的所有结果
	WriteOutput(currentInput)
}
