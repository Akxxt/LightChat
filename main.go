// LightChat 主程序入口文件
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"LightChat/routes"
	"LightChat/utils"
	"LightChat/websocket"
)

func main() {
	// 创建日志目录
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}

	// 初始化Redis
	fmt.Println("初始化Redis连接...")
	utils.InitRedis()

	// 初始化MinIO
	fmt.Println("初始化MinIO连接...")
	utils.InitMinIO()

	// 启动WebSocket Hub
	fmt.Println("启动WebSocket集线器...")
	go websocket.HubInstance.Run()

	// 配置路由
	fmt.Println("配置路由规则...")
	routes.SetupRoutes()

	// 获取端口配置
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// 启动HTTP服务
	fmt.Println("=====================================")
	fmt.Println("LightChat 聊天服务启动成功！")
	fmt.Println("服务地址：http://0.0.0.0:" + port)
	fmt.Println("=====================================")

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("服务启动失败：%v", err)
	}
}