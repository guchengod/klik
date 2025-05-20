package main

import (
	"klik/server/config"
	"klik/server/router"
	"log"
)

func main() {
	// 初始化配置
	config.Init()

	// 初始化路由
	r := router.InitRouter()

	// 启动服务器
	log.Println("服务器启动在 http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
