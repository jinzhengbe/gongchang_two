package main

import (
	"fmt"
	"log"

	"aneworder.com/backend/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 数据库配置
	dbConfig := struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}{
		Host:     "localhost",  // 使用 localhost 连接容器
		Port:     "3307",       // 使用映射到主机的端口
		User:     "gongchang",
		Password: "gongchang",
		DBName:   "gongchang",
	}

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 执行数据迁移
	if err := database.MigrateData(db); err != nil {
		log.Fatalf("Failed to migrate data: %v", err)
	}

	log.Println("Data migration completed successfully")
} 