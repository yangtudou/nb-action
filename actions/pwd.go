package actions

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"

	// 引入 Curve25519 库，用来真正计算出公钥
	"golang.org/x/crypto/curve25519"
)

type Password struct{}

func NewPassword() *Password {
	return &Password{}
}

func (p *Password) Name() string {
	return "pwd"
}

func (p *Password) Execute(
	ctx context.Context,
	args []string,
	input map[string]interface{},
) (map[string]interface{}, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("pwd requires a subcommand: wg-keypair or rand")
	}

	subcommand := args[0]

	switch subcommand {
	case "wg-keypair":
		return p.generateWGKeypair()
	case "rand":
		if len(args) < 2 {
			return nil, fmt.Errorf("rand requires size argument")
		}
		return p.generateRand(args[1:])
	default:
		return nil, fmt.Errorf("unknown subcommand: %s (supported: wg-keypair, rand)", subcommand)
	}
}

// generateWGKeypair 生成真实的 WireGuard 密钥对
func (p *Password) generateWGKeypair() (map[string]interface{}, error) {
	var privateKey [32]byte
	_, err := rand.Read(privateKey[:])
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// WireGuard 标准：针对 Curve25519 算法进行特殊位处理 (Clamping)
	privateKey[0] &= 248
	privateKey[31] = (privateKey[31] & 127) | 64

	// 计算公钥 (通过私钥和 Basepoint 计算)
	publicKey, err := curve25519.X25519(privateKey[:], curve25519.Basepoint)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate public key: %w", err)
	}

	return map[string]interface{}{
		// 标准的 WireGuard 密钥都是 Base64 编码的
		"private_key": base64.StdEncoding.EncodeToString(privateKey[:]),
		"public_key":  base64.StdEncoding.EncodeToString(publicKey),
		// 顺便保留一份十六进制格式备用
		"private_key_hex": hex.EncodeToString(privateKey[:]),
		"public_key_hex":  hex.EncodeToString(publicKey),
	}, nil
}

// generateRand 生成随机密钥
func (p *Password) generateRand(args []string) (map[string]interface{}, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("rand requires size argument")
	}

	// 解析大小参数
	size, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, fmt.Errorf("invalid size: %s", args[0])
	}

	if size <= 0 {
		return nil, fmt.Errorf("size must be positive")
	}

	// 检查是否有 --base64 flag
	useBase64 := false
	for i := 1; i < len(args); i++ {
		if args[i] == "--base64" {
			useBase64 = true
			break
		}
	}

	// 生成随机数据
	buf := make([]byte, size)
	_, err = rand.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random data: %w", err)
	}

	var value string
	if useBase64 {
		// Base64 编码
		value = base64.StdEncoding.EncodeToString(buf)
	} else {
		// 十六进制编码
		value = hex.EncodeToString(buf)
	}

	return map[string]interface{}{
		"value": value,
	}, nil
}
