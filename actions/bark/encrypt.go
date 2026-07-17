package bark

import (
	"encoding/json"
)

func encrypt(
	payload BarkPayload,
) (
	string,
	error,
) {

	// 第一步：
	// 加载 Bark 加密配置
	//
	// 包括：
	// AES key
	// GCM nonce(iv)

	config, err := LoadConfig()

	if err != nil {

		return "",
			err

	}

	// 第二步：
	// 将 BarkPayload 转换成 JSON 数据

	jsonData, err := json.Marshal(
		payload,
	)

	if err != nil {

		return "",
			err

	}

	// 第三步：
	// 调用 AES-GCM 加密

	ciphertext, err := aesGCMEncrypt(
		config.Key,
		config.IV,
		jsonData,
	)

	if err != nil {

		return "",
			err

	}

	// 第四步：
	// 将二进制密文转换成 URL 可用字符串

	encodedText := encodeCiphertext(
		ciphertext,
	)

	// 第五步：
	// 返回最终结果

	return encodedText, nil
}
