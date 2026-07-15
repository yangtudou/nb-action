package actions

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strconv"
)

type Random struct{}

func NewRandom() *Random {
	return &Random{}
}

func (r *Random) Name() string {
	return "random"
}

func (r *Random) Execute(
	ctx context.Context,
	args []string,
	input map[string]interface{},
) (
	map[string]interface{},
	error,
) {

	size := 32

	// 1. 优先尝试从命令行参数获取长度（如：nb-action -action random 64）
	if len(args) > 0 {
		if v, err := strconv.Atoi(args[0]); err == nil {
			size = v
		}
	} else {
		// 2. 如果没有命令行参数，则回退到原有的 JSON 读取逻辑
		if value, ok := input["bytes"]; ok {
			if v, ok := value.(float64); ok {
				size = int(v)
			}
		}
	}

	buf := make([]byte, size)

	_, err := rand.Read(buf)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"value": base64.StdEncoding.EncodeToString(buf),
	}, nil
}
