package bark

import (
	"crypto/aes"
	"crypto/cipher"
)

func aesGCMEncrypt(
	key []byte,
	iv []byte,
	data []byte,
) (
	[]byte,
	error,
) {

	// 第一步：
	// 使用 key 创建 AES 加密器
	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	// 第二步：
	// 使用 AES 加密器创建 GCM 模式
	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, err
	}

	// 第三步：
	// 使用 GCM 对数据进行加密
	ciphertext := gcm.Seal(
		nil,
		iv,
		data,
		nil,
	)

	// 第四步：
	// 返回加密后的数据
	return ciphertext, nil
}
