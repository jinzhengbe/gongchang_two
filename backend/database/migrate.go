package database

import (
	"aneworder.com/backend/models"
	"gorm.io/gorm"
)

func MigrateData(db *gorm.DB) error {
	// 自动迁移数据库表结构
	err := db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.File{},
	)
	if err != nil {
		return err
	}

	return nil
} 