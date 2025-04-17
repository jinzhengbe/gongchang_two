package database

import (
	"backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
)

// 初始化测试数据
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
			ID:       "1",
			Username: "designer1",
			Password: string(hashedPassword),
			Email:    "designer1@test.com",
			Role:     string(models.RoleDesigner),
		},
		{
			ID:       "2",
			Username: "factory1",
			Password: string(hashedPassword),
			Email:    "factory1@test.com",
			Role:     string(models.RoleFactory),
		},
		{
			ID:       "3",
			Username: "supplier1",
			Password: string(hashedPassword),
			Email:    "supplier1@test.com",
			Role:     string(models.RoleSupplier),
		},
	}

	// 创建用户
	for _, user := range testUsers {
		if err := db.Create(&user).Error; err != nil {
			log.Printf("Error creating test user %s: %v", user.Username, err)
			continue
		}
	}

	// 创建用户档案
	designerProfile := models.DesignerProfile{
		UserID:      1,
		CompanyName: "设计工作室1",
		Address:     "北京市朝阳区",
		Website:     "http://designer1.com",
		Bio:         "专业服装设计工作室，专注于高端定制",
	}

	factoryProfile := models.FactoryProfile{
		UserID:       2,
		CompanyName:  "服装厂1",
		Address:      "广东省深圳市",
		Capacity:     1000,
		Equipment:    "全自动裁剪机,工业缝纫机",
		Certificates: "ISO9001,质量管理体系认证",
	}

	supplierProfile := models.SupplierProfile{
		UserID:       3,
		CompanyName:  "面料供应商1",
		Address:      "浙江省绍兴市",
		MainProducts: "棉料,丝绸,化纤",
		Certificates: "环保认证,质量认证",
	}

	// 创建档案
	if err := db.Create(&designerProfile).Error; err != nil {
		log.Printf("Error creating designer profile: %v", err)
	}
	if err := db.Create(&factoryProfile).Error; err != nil {
		log.Printf("Error creating factory profile: %v", err)
	}
	if err := db.Create(&supplierProfile).Error; err != nil {
		log.Printf("Error creating supplier profile: %v", err)
	}

	log.Println("Test data initialized successfully")
	return nil
} 