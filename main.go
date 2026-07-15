package main

import (
	"context"
	"flag"
	"fmt"
	"os"
)

func main() {

	// 1. 优先注册 Action
	RegisterActions(
		os.Getenv("BARK_SERVER"),
		os.Getenv("BARK_KEY"),
	)

	// 如果什么参数都没传，默认调用 bark
	if len(os.Args) < 2 {
		input := ReadInput()
		runSingleAction("bark", []string{}, input)
		return
	}

	firstArg := os.Args[1]

	// 【新增分流】：API 服务模式 -> .\nb-action.exe server 8080
	if firstArg == "server" {
		port := ""
		if len(os.Args) > 2 {
			port = os.Args[2]
		}
		StartServer(port)
		return
	}

	// 2. 管道模式：nb-action.exe pipe random 64 , bark
	if firstArg == "pipe" {
		input := ReadInput()
		// 切割参数，剥离 "pipe" 关键字
		steps := parseByComma(os.Args[2:])
		if len(steps) < 2 {
			WriteError(fmt.Errorf("pipeline requires at least 2 actions separated by ','"))
			os.Exit(1)
		}

		// 启动内存接力管道
		runPipeline(context.Background(), steps, input)
		return
	}

	// 3. 智能直达模式：如果第一个参数就是已注册的 Action 名字（如 random、bark 等）
	if _, ok := GetAction(firstArg); ok {
		input := ReadInput()
		runSingleAction(firstArg, os.Args[2:], input)
		return
	}

	// 4. 兼容老旧的 -action 旗标模式
	actionName := flag.String(
		"action",
		"bark",
		"action name",
	)

	flag.Parse()

	args := flag.Args()
	input := ReadInput()

	runSingleAction(*actionName, args, input)
}

// runSingleAction 抽取出来的单步执行辅助函数
func runSingleAction(name string, args []string, input map[string]interface{}) {
	action, ok := GetAction(name)

	if !ok {
		WriteError(
			fmt.Errorf(
				"action not found: %s",
				name,
			),
		)
		os.Exit(1)
	}

	result, err := action.Execute(
		context.Background(),
		args,
		input,
	)

	if err != nil {
		WriteError(err)
		os.Exit(1)
	}

	WriteOutput(result)
}
