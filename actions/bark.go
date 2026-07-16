package actions

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

type Bark struct{}

// 保留签名，防止破坏 main.go 里的注册调用
func NewBark(server string, deviceKey string) *Bark {
	return &Bark{}
}

func (b *Bark) Name() string {
	return "bark"
}

func (b *Bark) Description() string {
	return "发送 Bark 推送通知"
}

func (b *Bark) Help() string {
	return `
bark

发送 Bark 推送通知

Usage:
  nb-action bark <title> <message>

Arguments:
  title       通知标题
  message     通知内容

Examples:
  nb-action bark "服务器报警" "Docker 服务停止"

Environment:
  BARK_SERVER
  BARK_KEY
`
}

type BarkPayload struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body"`
	Sound string `json:"sound,omitempty"`
}

func (b *Bark) Execute(
	ctx context.Context,
	args []string,
	input map[string]interface{},
) (map[string]interface{}, error) {
	key := os.Getenv("BARK_AES_KEY")
	iv := os.Getenv("BARK_AES_IV")

	if len(key) != 32 || len(iv) != 12 {
		return nil, fmt.Errorf("BARK_AES_KEY(32位) 或 BARK_AES_IV(12位) 环境变量配置错误")
	}

	// 1. 提取参数
	var title, body, sound string
	if len(args) >= 2 {
		title = args[0]
		body = args[1]
		if len(args) >= 3 {
			sound = args[2]
		}
	} else if len(args) == 1 {
		body = args[0]
	}

	// 2. 管道接力与默认值兜底
	if title == "" {
		title, _ = input["title"].(string)
	}
	if body == "" {
		if v, ok := input["body"].(string); ok {
			body = v
		} else {
			body, _ = input["value"].(string) // 完美兼容 random
		}
	}
	if sound == "" {
		if v, ok := input["sound"].(string); ok {
			sound = v
		} else {
			sound = "birdsong"
		}
	}

	if body == "" {
		return nil, fmt.Errorf("缺少消息内容 (body)")
	}

	// 3. 序列化明文
	jsonBytes, err := json.Marshal(BarkPayload{Title: title, Body: body, Sound: sound})
	if err != nil {
		return nil, err
	}

	// 4. 原生 AES-256-GCM 加密
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	encrypted := gcm.Seal(nil, []byte(iv), jsonBytes, nil)
	base64Cipher := base64.StdEncoding.EncodeToString(encrypted)

	// 5. 直接进行 URL 编码并作为 ciphertext 字段输出
	return map[string]interface{}{
		"ciphertext": url.QueryEscape(base64Cipher),
	}, nil
}
