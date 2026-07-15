package actions

import (
	"context"
)

type Test struct{}

func NewTest() *Test {
	return &Test{}
}

func (t *Test) Name() string {
	return "test"
}

func (t *Test) Execute(
	ctx context.Context,
	args []string,
	input map[string]interface{},
) (
	map[string]interface{},
	error,
) {
	// 默认拿传入的第一个参数作为 value
	val := ""
	if len(args) > 0 {
		val = args[0]
	}

	return map[string]interface{}{
		"value":          val,
		"received_count": len(input),
		"args":           args,
		"status":         "ok",
	}, nil
}
