package bark

import (
	"context"
	"fmt"
)

type Bark struct{}

func New() *Bark {
	return &Bark{}
}

func (b *Bark) Name() string {
	return "bark"
}

func (b *Bark) Description() string {
	return "发送 Bark 推送通知"
}

func (b *Bark) Help() string {
	return helpText
}

func (b *Bark) Execute(
	ctx context.Context,
	args []string,
	input map[string]interface{},
) (map[string]interface{}, error) {

	title, body := parseInput(args, input)

	if body == "" {
		return nil, fmt.Errorf("缺少推送内容！")
	}

	ciphertext, err := encrypt(
		BarkPayload{
			Title: title,
			Body:  body,
		},
	)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ciphertext": ciphertext,
	}, nil
}

func parseInput(
	args []string,
	input map[string]interface{},
) (
	string,
	string,
) {

	var title string
	var body string

	if len(args) >= 2 {
		title = args[0]
		body = args[1]

	} else if len(args) == 1 {
		body = args[0]
	}

	if title == "" {
		title, _ = input["title"].(string)
	}

	if body == "" {

		if value, ok := input["body"].(string); ok {
			body = value
		} else {
			body, _ = input["value"].(string)
		}
	}

	return title, body
}
