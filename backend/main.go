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
	go func() {
		log.Printf("HTTP Server starting on port 8080")
		if err := router.Run(":8080"); err != nil {
			log.Printf("Error starting HTTP server: %v", err)
		}
	}()

	// 启动 HTTPS 服务器
	certFile := "/app/ssl/cert.pem"
	keyFile := "/app/ssl/key.pem"

	log.Printf("HTTPS Server starting on port 443")
	if err := router.RunTLS(":443", certFile, keyFile); err != nil {
		log.Printf("Error starting HTTPS server: %v", err)
	}
} 