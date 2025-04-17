package main

import (
	"log"
	"os"
	"time"

	"backend/config"
	"backend/database"
	"backend/models"
	"backend/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// 初始化数据库连接（带重试）
	var db *gorm.DB
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		dsn := cfg.DBUser + ":" + cfg.DBPassword + "@tcp(" + cfg.DBHost + ":" + cfg.DBPort + ")/" + cfg.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Error connecting to database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to database after %d attempts: %v", maxRetries, err)
	}

	// 自动迁移数据库表
	err = db.AutoMigrate(
		&models.User{},
		&models.DesignerProfile{},
		&models.FactoryProfile{},
		&models.SupplierProfile{},
	)
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	// 初始化测试数据
	if err := database.InitTestData(db); err != nil {
		log.Printf("Error initializing test data: %v", err)
	}

	// 设置路由
	router := routes.SetupRouter(db)

	// 配置 CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 启动服务器
	log.Printf("Server starting on port %s", cfg.Port)
	
	// 检查是否提供了 SSL 证书
	certFile := os.Getenv("SSL_CERT_FILE")
	keyFile := os.Getenv("SSL_KEY_FILE")
	
	// 使用 HTTP 模式
	go func() {
		if err := router.Run(":" + cfg.Port); err != nil {
			log.Printf("Error starting HTTP server: %v", err)
		}
	}()

	// 如果有证书，同时启动 HTTPS
	if certFile != "" && keyFile != "" {
		if err := router.RunTLS(":443", certFile, keyFile); err != nil {
			log.Printf("Error starting HTTPS server: %v", err)
		}
	}

	// 保持主程序运行
	select {}
} 