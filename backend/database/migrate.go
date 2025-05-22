package database

import (
	"backend/models"
	"backend/internal/factory"
	"gorm.io/gorm"
)

func MigrateData(db *gorm.DB) error {
	// 自动迁移数据库表结构
	err := db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.File{},
		&models.DesignerProfile{},
		&models.FactoryProfile{},
		&models.SupplierProfile{},
		&models.OrderProgress{},
		&models.OrderAttachment{},
		&factory.Factory{},
	)
	if err != nil {
		return err
	}

	return nil
} 