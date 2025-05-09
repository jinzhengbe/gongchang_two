package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"aneworder.com/backend/config"
	"aneworder.com/backend/database"
	"aneworder.com/backend/routes"
	"github.com/gin-gonic/gin"
	"net/http"
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

	// 自动迁移数据库表结构
	if err := database.MigrateData(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 设置 Gin 模式
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 设置路由
	router := routes.SetupRouter(db, cfg)

	// 配置服务器
	server := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 启动服务器
	log.Printf("Server starting on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 