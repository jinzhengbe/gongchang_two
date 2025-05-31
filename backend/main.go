package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"gongChang/config"
	"gongChang/database"
	"gongChang/routes"
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

	// 自动迁移数据库表结构（已稳定，暂时注释）
	/*
	if err := database.MigrateData(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	*/

	// 启动数据库监控
	go database.MonitorDatabase(db, 5*time.Minute)
	database.MonitorSlowQueries(db, 1*time.Second)

	// 设置定时备份
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			if err := database.BackupDatabase(cfg); err != nil {
				log.Printf("Failed to backup database: %v", err)
			} else {
				log.Println("Database backup completed successfully")
			}

			// 清理30天前的备份
			if err := database.CleanOldBackups(30 * 24 * time.Hour); err != nil {
				log.Printf("Failed to clean old backups: %v", err)
			}
		}
	}()

	// 设置 Gin 模式
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 设置路由
	router := routes.SetupRouter(db, cfg)

	// 打印所有已注册的路由
	for _, route := range router.Routes() {
		log.Printf("[ROUTE] %s %s", route.Method, route.Path)
	}

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