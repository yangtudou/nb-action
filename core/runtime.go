package core

import (
	"context"
	"fmt"
)

type Runtime struct {
	Registry *Registry
}

func NewRuntime() *Runtime {

	return &Runtime{
		Registry: NewRegistry(),
	}
}

func (r *Runtime) Execute(
	name string,
	args []string,
	input map[string]interface{},
) {

	action, ok := r.Registry.GetAction(name)

	if !ok {

		WriteError(
			fmt.Errorf("action not found: %s", name),
		)

		return
	}

	result, err := action.Execute(
		context.Background(),
		args,
		input,
	)

	if err != nil {

		WriteError(err)

		return
	}

	WriteOutput(result)
}
