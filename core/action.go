package core

import "context"

type Action interface {

	// Action 名称
	Name() string

	// 简短描述
	Description() string

	// 详细帮助
	Help() string

	// 执行
	Execute(
		ctx context.Context,
		args []string,
		input map[string]interface{},
	) (
		map[string]interface{},
		error,
	)
}
