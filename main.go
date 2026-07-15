package main

import (
	"context"
	"flag"
	"fmt"
	"os"
)

func main() {
	// 1. 注册 Action (只做映射，不做环境变量校验)
	RegisterActions(
		os.Getenv("BARK_SERVER"),
		os.Getenv("BARK_KEY"),
	)

	// 2. 无参数时打印帮助
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	firstArg := os.Args[1]

	// 3. API 服务模式
	if firstArg == "server" {
		port := ""
		if len(os.Args) > 2 {
			port = os.Args[2]
		}
		StartServer(port)
		return
	}

	// 4. 管道模式
	if firstArg == "pipe" {
		input := ReadInput()
		steps := parseByComma(os.Args[2:])
		if len(steps) < 2 {
			WriteError(fmt.Errorf("pipeline requires at least 2 actions separated by ','"))
			os.Exit(1)
		}
		runPipeline(context.Background(), steps, input)
		return
	}

	// 5. 智能直达模式
	if _, ok := GetAction(firstArg); ok {
		input := ReadInput()
		runSingleAction(firstArg, os.Args[2:], input)
		return
	}

	// 6. 旗标模式 (彻底移除默认 bark，改为按需触发)
	actionName := flag.String("action", "", "action name")
	flag.Parse()

	if *actionName != "" {
		input := ReadInput()
		runSingleAction(*actionName, flag.Args(), input)
		return
	}

	// 如果所有匹配都失败，提示错误并给出帮助
	fmt.Fprintf(os.Stderr, "Unknown command or action: %s\n", firstArg)
	printUsage()
}

// 使用方法
func printUsage() {
	fmt.Println(`
使用: nb-action <command> [args]

命令:
  server [port]              启动 HTTP API 服务 (默认 8080)
  pipe <action1>,<action2>   管道模式：执行链式操作
  <action> [args...]         执行单个动作 (如: random 64, bark title content)

举例:
  nb-action server 8080
  nb-action random 64
  nb-action pipe random 64, bark
	`)
}

// 执行辅助函数
func runSingleAction(name string, args []string, input map[string]interface{}) {
	action, ok := GetAction(name)
	if !ok {
		WriteError(fmt.Errorf("action not found: %s", name))
		os.Exit(1)
	}

	result, err := action.Execute(context.Background(), args, input)
	if err != nil {
		WriteError(err)
		os.Exit(1)
	}

	WriteOutput(result)
}
