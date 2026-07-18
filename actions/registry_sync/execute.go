package registry_sync

import (
	"context"
	"fmt"
)

func (r *RegistrySync) Execute(
	ctx context.Context,
	args []string,
	input map[string]interface{},
) (map[string]interface{}, error) {

	opt, err := ParseOptions(args)

	if err != nil {
		return nil, fmt.Errorf(
			"invalid registry-sync arguments: %w",
			err,
		)
	}

	// 支持 pipeline 输入
	if opt.Src == "" {

		if v, ok := input["value"].(string); ok {
			opt.Src = v
		}

		if v, ok := input["image"].(string); ok {
			opt.Src = v
		}
	}

	if opt.DstPrefix == "" {
		return nil, fmt.Errorf(
			"missing required flag: --dst-prefix",
		)
	}

	return Run(
		ctx,
		opt,
	)
}
