package main

import (
	"fmt"
	"log"
	"aneworder.com/backend/config"
	"aneworder.com/backend/database"
	"aneworder.com/backend/routes"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 初始化数据库
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 设置路由
	router := routes.SetupRouter(db, cfg)

	// 启动服务器
	addr := fmt.Sprintf("0.0.0.0:%s", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 