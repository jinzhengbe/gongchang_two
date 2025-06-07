package database

import (
	"gongChang/config"
	"gongChang/models"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm/logger"
)

// InitDB 初始化数据库连接
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FShanghai",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	var db *gorm.DB
	var err error
	maxRetries := 5
	retryInterval := time.Second * 5

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			// 配置连接池
			sqlDB, err := db.DB()
			if err == nil {
				// 设置最大空闲连接数
				sqlDB.SetMaxIdleConns(10)
				// 设置最大打开连接数
				sqlDB.SetMaxOpenConns(100)
				// 设置连接最大生命周期
				sqlDB.SetConnMaxLifetime(time.Hour)
				// 设置空闲连接最大生命周期
				sqlDB.SetConnMaxIdleTime(time.Minute * 10)
			}
			break
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			log.Printf("Retrying in %v...", retryInterval)
			time.Sleep(retryInterval)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
	}

	// 删除旧表（仅开发调试时使用，生产环境请勿启用）
	/*
	err = db.Migrator().DropTable(
		&models.File{},
		&models.Order{},
		&models.Product{},
		&models.User{},
		&models.DesignerProfile{},
		&models.FactoryProfile{},
		&models.SupplierProfile{},
		&models.OrderProgress{},
		&models.OrderAttachment{},
	)
	if err != nil {
		log.Printf("Warning: Failed to drop tables: %v", err)
	}
	*/

	// 执行自动迁移
	if err := MigrateData(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	// 初始化测试数据
	if err := InitTestData(db); err != nil {
		log.Printf("Warning: Failed to initialize test data: %v", err)
	}

	return db, nil
}

// InitTestData 初始化测试数据
func InitTestData(db *gorm.DB) error {
	// 检查数据库是否为空
	var count int64
	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check if database is empty: %v", err)
	}

	// 如果数据库不为空，跳过初始化
	if count > 0 {
		log.Printf("Database is not empty, skipping test data initialization")
		return nil
	}

	log.Printf("Database is empty, initializing test data...")

	// 创建测试用户密码
	password := "test123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 测试用户数据
	testUsers := []models.User{
		{
			ID:       uuid.New().String(),
			Username: "designer1",
			Password: string(hashedPassword),
			Email:    "designer1@test.com",
			Role:     models.RoleDesigner,
		},
		{
			ID:       uuid.New().String(),
			Username: "factory1",
			Password: string(hashedPassword),
			Email:    "factory1@test.com",
			Role:     models.RoleFactory,
		},
		{
			ID:       uuid.New().String(),
			Username: "supplier1",
			Password: string(hashedPassword),
			Email:    "supplier1@test.com",
			Role:     models.RoleSupplier,
		},
	}

	// 创建用户
	for _, user := range testUsers {
		if err := db.Create(&user).Error; err != nil {
			log.Printf("Error creating test user %s: %v", user.Username, err)
			continue
		}
		log.Printf("Created test user: %s", user.Username)

		// 根据用户角色创建对应的档案
		switch user.Role {
		case models.RoleDesigner:
			designerProfile := models.DesignerProfile{
				UserID:      user.ID,
				CompanyName: "设计工作室1",
				Address:     "北京市朝阳区",
				Website:     "http://designer1.com",
				Bio:         "专业服装设计工作室，专注于高端定制",
			}
			if err := db.Create(&designerProfile).Error; err != nil {
				log.Printf("Error creating designer profile for %s: %v", user.Username, err)
			}
		case models.RoleFactory:
			factoryProfile := models.FactoryProfile{
				UserID:      user.ID,
				CompanyName: "服装厂1",
				Address:     "广东省深圳市",
				Capacity:    1000,
				Equipment:   "全自动裁剪机,工业缝纫机",
				Certificates: "ISO9001,质量管理体系认证",
			}
			if err := db.Create(&factoryProfile).Error; err != nil {
				log.Printf("Error creating factory profile for %s: %v", user.Username, err)
			}
		case models.RoleSupplier:
			supplierProfile := models.SupplierProfile{
				UserID:      user.ID,
				CompanyName: "面料供应商1",
				Address:     "浙江省绍兴市",
				MainProducts: "棉料,丝绸,化纤",
				Certificates: "环保认证,质量认证",
			}
			if err := db.Create(&supplierProfile).Error; err != nil {
				log.Printf("Error creating supplier profile for %s: %v", user.Username, err)
			}
		}
	}

	return nil
} 