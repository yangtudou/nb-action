package bootstrap

import (
	"os"

	"github.com/yangtudou/nb-action/actions"
	"github.com/yangtudou/nb-action/core"
)

func Load(runtime *core.Runtime) {

	actions.RegisterAll(
		runtime.Registry,
		os.Getenv("BARK_SERVER"),
		os.Getenv("BARK_KEY"),
	)

}
