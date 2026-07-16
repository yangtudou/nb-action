package core

import (
	"encoding/json"
	"io"
	"os"
)

type Input struct {
	Data map[string]interface{} `json:"data"`
	Meta map[string]interface{} `json:"meta"`
}

func ReadInput() map[string]interface{} {

	input := make(map[string]interface{})

	info, err := os.Stdin.Stat()

	if err != nil {
		return input
	}

	if info.Mode()&os.ModeNamedPipe == 0 {
		return input
	}

	data, err := io.ReadAll(os.Stdin)

	if err != nil || len(data) == 0 {
		return input
	}

	var raw map[string]interface{}

	err = json.Unmarshal(
		data,
		&raw,
	)

	if err != nil {
		return input
	}

	// 新格式
	if value, ok := raw["data"]; ok {

		if dataMap, ok := value.(map[string]interface{}); ok {

			if meta, ok := raw["meta"].(map[string]interface{}); ok {
				_ = meta
			}

			return dataMap
		}
	}

	// 兼容旧格式
	return raw
}
