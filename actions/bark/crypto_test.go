package bark

import (
	"testing"
)

func TestAESGCMEncrypt(t *testing.T) {

	// AES-256 要求 key 长度必须是 32 bytes
	key := []byte(
		"12345678901234567890123456789012",
	)

	// GCM 推荐 nonce 长度为 12 bytes
	iv := []byte(
		"123456789012",
	)

	// 模拟需要加密的数据
	data := []byte(
		"hello bark",
	)

	// 调用我们的 AES-GCM 加密函数
	ciphertext, err := aesGCMEncrypt(
		key,
		iv,
		data,
	)

	// 检查是否报错
	if err != nil {

		t.Fatal(err)

	}

	// 检查是否生成密文
	if len(ciphertext) == 0 {

		t.Fatal("ciphertext is empty")

	}

	// 打印密文，方便观察
	t.Log("plaintext:")
	t.Log(string(data))

	t.Log("ciphertext:")
	t.Log(ciphertext)

}
