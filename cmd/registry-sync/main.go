package main

import (
	"context"
	"fmt"
	"os"

	"github.com/yangtudou/nb-action/actions/registry_sync"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println(
			"Usage: registry-sync sync [options]",
		)

		os.Exit(1)
	}

	switch os.Args[1] {

	case "sync":

		opt, err := registry_sync.ParseOptions(
			os.Args[2:],
		)

		if err != nil {
			fmt.Println(
				"invalid arguments:",
				err,
			)

			os.Exit(1)
		}

		if opt.DstPrefix == "" {

			fmt.Println(
				"missing required flag: --dst-prefix",
			)

			os.Exit(1)
		}

		result, err := registry_sync.Run(
			context.Background(),
			opt,
		)

		if err != nil {

			fmt.Println(
				"registry-sync failed:",
				err,
			)

			os.Exit(1)
		}

		fmt.Println(
			result,
		)

	default:

		fmt.Printf(
			"unknown command: %s\n",
			os.Args[1],
		)

		os.Exit(1)
	}
}
