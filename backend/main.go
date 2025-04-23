package main

import (
	"aneworder.com/backend/config"
	"aneworder.com/backend/database"
	"aneworder.com/backend/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// 初始化数据库
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	// 初始化测试数据
	if err := database.InitTestData(db); err != nil {
		log.Printf("Error initializing test data: %v", err)
	}

	// 创建 Gin 引擎
	router := gin.Default()

	// 设置受信任的代理
	router.SetTrustedProxies([]string{"aneworders.com"})

	// 设置路由
	routes.SetupRouter(router, db)

	// 启动 HTTP 服务器
	log.Printf("HTTP Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
} 