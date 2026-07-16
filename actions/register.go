package actions

import (
	"github.com/yangtudou/nb-action/core"
)

func RegisterAll(
	registry *core.Registry,
	server string,
	deviceKey string,
) {

	registry.Register(
		NewBark(
			server,
			deviceKey,
		),
		NewTest(),
		NewPassword(),
		&RegistrySync{},
	)
}
