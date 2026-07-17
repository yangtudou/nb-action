package bark

import (
	"fmt"
	"os"
)

type Config struct {
	Key []byte

	IV []byte
}

func LoadConfig() (
	Config,
	error,
) {

	// 从环境变量读取 key

	keyString := os.Getenv(
		"BARK_AES_KEY",
	)

	// 从环境变量读取 iv

	ivString := os.Getenv(
		"BARK_AES_IV",
	)

	// 检查 AES-256 key 长度

	if len(keyString) != 32 {

		return Config{},
			fmt.Errorf(
				"BARK_AES_KEY 长度错误，需要32字节",
			)

	}

	// 检查 GCM nonce 长度

	if len(ivString) != 12 {

		return Config{},
			fmt.Errorf(
				"BARK_AES_IV 长度错误，需要12字节",
			)

	}

	// 转换成加密需要的 []byte

	config := Config{

		Key: []byte(keyString),

		IV: []byte(ivString),
	}

	return config, nil
}
