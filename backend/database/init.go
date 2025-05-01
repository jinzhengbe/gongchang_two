package database

import (
	"aneworder.com/backend/config"
	"aneworder.com/backend/models"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
	"github.com/google/uuid"
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
		})
		if err == nil {
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

	// 自动迁移数据库表结构
	err = db.AutoMigrate(
		&models.User{},
		&models.DesignerProfile{},
		&models.FactoryProfile{},
		&models.SupplierProfile{},
		&models.Product{},
		&models.Order{},
		&models.OrderProgress{},
		&models.OrderAttachment{},
		&models.File{},
	)
	if err != nil {
		return nil, err
	}

	// 初始化测试数据
	if err := InitTestData(db); err != nil {
		log.Printf("Warning: Failed to initialize test data: %v", err)
	}

	return db, nil
}

// InitTestData 初始化测试数据
func InitTestData(db *gorm.DB) error {
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
		var existingUser models.User
		if err := db.First(&existingUser, "username = ?", user.Username).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&user).Error; err != nil {
					log.Printf("Error creating test user %s: %v", user.Username, err)
					continue
				}
				log.Printf("Created test user: %s", user.Username)
			} else {
				log.Printf("Error checking existing user %s: %v", user.Username, err)
				continue
			}
		} else {
			log.Printf("Test user already exists: %s", user.Username)
		}
	}

	// 创建用户档案
	designerProfile := models.DesignerProfile{
		UserID:      "1",
		CompanyName: "设计工作室1",
		Address:     "北京市朝阳区",
		Website:     "http://designer1.com",
		Bio:         "专业服装设计工作室，专注于高端定制",
	}

	factoryProfile := models.FactoryProfile{
		UserID:       "2",
		CompanyName:  "服装厂1",
		Address:      "广东省深圳市",
		Capacity:     1000,
		Equipment:    "全自动裁剪机,工业缝纫机",
		Certificates: "ISO9001,质量管理体系认证",
	}

	supplierProfile := models.SupplierProfile{
		UserID:       "3",
		CompanyName:  "面料供应商1",
		Address:      "浙江省绍兴市",
		MainProducts: "棉料,丝绸,化纤",
		Certificates: "环保认证,质量认证",
	}

	// 创建档案
	var existingDesignerProfile models.DesignerProfile
	if err := db.First(&existingDesignerProfile, "user_id = ?", designerProfile.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&designerProfile).Error; err != nil {
				log.Printf("Error creating designer profile: %v", err)
			} else {
				log.Printf("Created designer profile for user ID: %s", designerProfile.UserID)
			}
		}
	} else {
		log.Printf("Designer profile already exists for user ID: %s", designerProfile.UserID)
	}

	var existingFactoryProfile models.FactoryProfile
	if err := db.First(&existingFactoryProfile, "user_id = ?", factoryProfile.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&factoryProfile).Error; err != nil {
				log.Printf("Error creating factory profile: %v", err)
			} else {
				log.Printf("Created factory profile for user ID: %s", factoryProfile.UserID)
			}
		}
	} else {
		log.Printf("Factory profile already exists for user ID: %s", factoryProfile.UserID)
	}

	var existingSupplierProfile models.SupplierProfile
	if err := db.First(&existingSupplierProfile, "user_id = ?", supplierProfile.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&supplierProfile).Error; err != nil {
				log.Printf("Error creating supplier profile: %v", err)
			} else {
				log.Printf("Created supplier profile for user ID: %s", supplierProfile.UserID)
			}
		}
	} else {
		log.Printf("Supplier profile already exists for user ID: %s", supplierProfile.UserID)
	}

	log.Println("Test data initialization completed")
	return nil
} 