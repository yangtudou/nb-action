package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/yangtudou/nb-action/bootstrap"
	"github.com/yangtudou/nb-action/core"
)

func main() {

	runtime := core.NewRuntime()

	registry := runtime.Registry

	bootstrap.Load(runtime)

	// 无参数直接显示帮助
	if len(os.Args) < 2 {
		printHelp(registry)
		return
	}

	firstArg := os.Args[1]

	if firstArg == "help" {
		if len(os.Args) > 2 {
			actionName := os.Args[2]
			action, ok := registry.GetAction(actionName)

			if !ok {
				fmt.Fprintf(
					os.Stderr,
					"Unknown action: %s\n",
					actionName,
				)
				return
			}

			fmt.Println(action.Help())
			return
		}

		printHelp(registry)
		return
	}

	// API 服务模式
	if firstArg == "server" {
		port := ""
		if len(os.Args) > 2 {
			port = os.Args[2]
		}

		core.StartServer(runtime, port)
		return
	}

	// 管道模式
	if firstArg == "pipe" {

		input := core.ReadInput()

		steps := core.ParseByComma(os.Args[2:])

		if len(steps) < 2 {
			core.WriteError(
				fmt.Errorf("pipeline requires at least 2 actions separated by ','"),
			)
			os.Exit(1)
		}

		runtime.RunPipeline(
			context.Background(),
			steps,
			input,
		)

		return
	}

	// 直接执行 Action
	if action, ok := registry.GetAction(firstArg); ok {

		args := os.Args[2:]

		for _, arg := range args {
			if arg == "--help" || arg == "-h" {
				fmt.Println(action.Help())
				return
			}
		}

		input := core.ReadInput()
		runSingleAction(registry, firstArg, args, input)
		return
	}

	// flag 模式保留，后续再处理
	actionName := flag.String(
		"action",
		"",
		"action name",
	)

	flag.Parse()

	if *actionName != "" {

		input := core.ReadInput()

		runSingleAction(
			registry,
			*actionName,
			flag.Args(),
			input,
		)

		return
	}

	fmt.Fprintf(
		os.Stderr,
		"Unknown command or action: %s\n\n",
		firstArg,
	)

	printHelp(registry)
}

// 执行辅助函数
func runSingleAction(
	registry *core.Registry,
	name string,
	args []string,
	input map[string]interface{},
) {

	action, ok := registry.GetAction(name)

	if !ok {

		core.WriteError(
			fmt.Errorf("action not found: %s", name),
		)

		os.Exit(1)
	}

	result, err := action.Execute(
		context.Background(),
		args,
		input,
	)

	if err != nil {

		core.WriteError(err)

		os.Exit(1)
	}

	core.WriteOutput(result)
}

// 动态帮助
func printHelp(
	registry *core.Registry,
) {

	fmt.Println("nb-action")

	fmt.Println()

	fmt.Println("Usage:")
	fmt.Println("  nb-action <action> [args]")
	fmt.Println()

	fmt.Println("Actions:")

	list := registry.List()

	sort.Slice(
		list,
		func(i, j int) bool {
			return list[i].Name() < list[j].Name()
		},
	)

	for _, action := range list {

		fmt.Printf(
			"  %-16s %s\n",
			action.Name(),
			action.Description(),
		)
	}

	fmt.Println()

	fmt.Println("Commands:")

	fmt.Println("  server [port]              Start HTTP API Server")

	fmt.Println("  pipe <action1>,<action2>   Execute action pipeline")

	fmt.Println("  help                       Show this help")
}
