package bark

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {

	// 设置正确的 AES key

	t.Setenv(
		"BARK_AES_KEY",
		"12345678901234567890123456789012",
	)

	// 设置正确的 IV

	t.Setenv(
		"BARK_AES_IV",
		"123456789012",
	)

	// 加载配置

	config, err := LoadConfig()

	// 检查错误

	if err != nil {

		t.Fatal(err)

	}

	// 检查 key 长度

	if len(config.Key) != 32 {

		t.Fatalf(
			"key length error: %d",
			len(config.Key),
		)

	}

	// 检查 IV 长度

	if len(config.IV) != 12 {

		t.Fatalf(
			"iv length error: %d",
			len(config.IV),
		)

	}

	t.Log("config load success")

}
