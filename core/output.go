package core

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// maskString 负责对敏感字符串进行局部打码，保留首尾各 4 位用于比对
func maskString(s string) string {
	runes := []rune(s)
	n := len(runes)
	if n == 0 {
		return ""
	}

	// 如果太短（小于等于 8 个字符），为了安全，只保留第一位
	if n <= 8 {
		return string(runes[0:1]) + "..."
	}

	// 保留前 4 位和后 4 位，中间用 ... 连起来
	return string(runes[:4]) + "..." + string(runes[n-4:])
}

// maskSensitive 递归地给敏感字段打码，保护隐私安全
func maskSensitive(data interface{}) interface{} {
	// 1. 如果是 Map，遍历并给敏感 key 进行局部脱敏
	if m, ok := data.(map[string]interface{}); ok {
		masked := make(map[string]interface{})
		// 敏感词黑名单
		sensitiveKeys := []string{"value", "key", "secret", "password", "token", "device_key", "device_token", "auth"}

		for k, v := range m {
			isSensitive := false
			lowerKey := strings.ToLower(k)
			for _, s := range sensitiveKeys {
				if strings.Contains(lowerKey, s) {
					isSensitive = true
					break
				}
			}

			if isSensitive {
				// 转为字符串进行首尾局部保留打码
				strVal := fmt.Sprintf("%v", v)
				masked[k] = maskString(strVal)
			} else {
				// 递归调用，防止嵌套的 map 漏网
				masked[k] = maskSensitive(v)
			}
		}
		return masked
	}

	// 2. 如果是 Slice/数组，递归处理里面的每一个元素
	if s, ok := data.([]interface{}); ok {
		maskedSlice := make([]interface{}, len(s))
		for i, val := range s {
			maskedSlice[i] = maskSensitive(val)
		}
		return maskedSlice
	}

	// 3. 基本类型直接返回
	return data
}

func WriteOutput(data interface{}) {
	// 在打印到屏幕前，先进行脱敏处理
	safeData := maskSensitive(data)

	output, err := json.Marshal(safeData)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}

func WriteError(err error) {
	output := map[string]interface{}{
		"error": err.Error(),
	}

	data, marshalErr := json.Marshal(output)

	if marshalErr != nil {
		fmt.Fprintln(os.Stderr, marshalErr)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, string(data))
}
