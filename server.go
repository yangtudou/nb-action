package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// StartServer 启动极简 HTTP API 服务
func StartServer(port string) {
	if port == "" {
		port = "8080"
	}

	// 强制只有 /api/v1/action/ 下的请求才进入处理逻辑
	http.HandleFunc("/api/v1/action/", handleAction)

	fmt.Printf("🚀 Nb-Action API Server 正在运行在 http://127.0.0.1:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("❌ 服务器启动失败: %v\n", err)
	}
}

func handleAction(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/action/")
	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) == 0 || parts[0] == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	actionName := parts[0]
	args := parts[1:]

	// --- 极简参数合法性校验 ---
	argCount := len(args)
	switch actionName {
	case "bark":
		if argCount < 1 || argCount > 2 { // 只允许 1 或 2
			http.Error(w, "Bark 仅支持 1-2 个参数", http.StatusBadRequest)
			return
		}
	case "random":
		if argCount != 1 { // 假设 random 只允许 1 个参数
			http.Error(w, "Random 仅支持 1 个参数", http.StatusBadRequest)
			return
		}
		// 其他 Action 可以在这里继续扩展
	}
	// -------------------------

	action, ok := GetAction(actionName)
	if !ok {
		http.Error(w, "Action Not Found", http.StatusNotFound)
		return
	}

	// 执行逻辑...
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	result, err := action.Execute(ctx, args, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
