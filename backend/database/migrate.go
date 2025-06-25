package database

import (
	"gongChang/models"
	"gorm.io/gorm"
	"log"
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

	// 手动执行SQL迁移文件
	log.Println("Executing SQL migrations...")
	
	// 添加designer_id字段（如果不存在）
	if err := db.Exec("ALTER TABLE fabrics ADD COLUMN IF NOT EXISTS designer_id VARCHAR(191) NULL").Error; err != nil {
		log.Printf("Warning: Failed to add designer_id column: %v", err)
	}
	
	// 添加factory_id字段（如果不存在）
	if err := db.Exec("ALTER TABLE fabrics ADD COLUMN IF NOT EXISTS factory_id VARCHAR(191) NULL").Error; err != nil {
		log.Printf("Warning: Failed to add factory_id column: %v", err)
	}
	
	// 添加索引（如果不存在）
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_fabrics_designer_id ON fabrics(designer_id)").Error; err != nil {
		log.Printf("Warning: Failed to create designer_id index: %v", err)
	}
	
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_fabrics_supplier_id ON fabrics(supplier_id)").Error; err != nil {
		log.Printf("Warning: Failed to create supplier_id index: %v", err)
	}
	
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_fabrics_factory_id ON fabrics(factory_id)").Error; err != nil {
		log.Printf("Warning: Failed to create factory_id index: %v", err)
	}

	log.Println("SQL migrations completed")

	return nil
} 