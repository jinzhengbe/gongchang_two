package database

import (
	"gongChang/models"
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
		&models.Fabric{},
		&models.FabricCategory{},
	)
	if err != nil {
		return err
	}

	// 执行额外的迁移
	if err := db.Exec("ALTER TABLE users MODIFY COLUMN role varchar(191) NOT NULL").Error; err != nil {
		return err
	}

	return nil
} 