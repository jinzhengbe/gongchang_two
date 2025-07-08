package services

import (
	"gongChang/models"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"github.com/google/uuid"
	"log"
)

// Custom error types
var (
	ErrUsernameExists = errors.New("username already exists")
	ErrEmailExists    = errors.New("email already exists")
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) Register(req models.RegisterRequest) error {
	// 使用事务确保数据一致性
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 检查用户名和邮箱是否已存在
	var existingUser models.User
		if err := tx.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return ErrUsernameExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
		if err := tx.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			return ErrEmailExists
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// 2. 创建 User 对象
		hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return err
	}
		user := &models.User{
			ID:       uuid.New().String(),
			Username: req.Username,
			Password: hashedPassword,
			Email:    req.Email,
			Role:     models.UserRole(req.Role),
		}

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// 3. 根据角色创建对应的 Profile
		switch user.Role {
		case models.RoleDesigner:
			profile := models.DesignerProfile{UserID: user.ID, CompanyName: req.CompanyName, Bio: req.Bio}
			if err := tx.Create(&profile).Error; err != nil {
				return err
			}
		case models.RoleFactory:
			profile := models.FactoryProfile{UserID: user.ID, CompanyName: req.CompanyName, Address: req.Address}
			if err := tx.Create(&profile).Error; err != nil {
				return err
			}
		case models.RoleSupplier:
			profile := models.SupplierProfile{UserID: user.ID, CompanyName: req.CompanyName, MainProducts: req.MainProducts}
			if err := tx.Create(&profile).Error; err != nil {
				return err
			}
		default:
			return errors.New("invalid user role")
		}

		return nil
	})
}

func (s *UserService) Login(username, password string) (*models.User, interface{}, error) {
	log.Printf("Attempting login for user: %s", username)
	
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("User not found: %s", username)
			return nil, nil, errors.New("user not found")
		}
		log.Printf("Database error: %v", err)
		return nil, nil, err
	}

	log.Printf("User found, stored password hash: %s", user.Password)
	log.Printf("Attempting to verify password...")

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("Password verification failed: %v", err)
		return nil, nil, errors.New("invalid password")
	}

	log.Printf("Password verified successfully for user %s with role %s", username, user.Role)

	// 根据用户角色获取相应的档案信息
	var profile interface{}
	var err error
	switch user.Role {
	case models.RoleDesigner:
		var designerProfile models.DesignerProfile
		err = s.db.Where("user_id = ?", user.ID).First(&designerProfile).Error
			profile = designerProfile
	case models.RoleFactory:
		var factoryProfile models.FactoryProfile
		err = s.db.Where("user_id = ?", user.ID).First(&factoryProfile).Error
			profile = factoryProfile
	case models.RoleSupplier:
		var supplierProfile models.SupplierProfile
		err = s.db.Where("user_id = ?", user.ID).First(&supplierProfile).Error
			profile = supplierProfile
	default:
		log.Printf("User %s has an unknown role: %s", username, user.Role)
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Failed to retrieve profile for user %s: %v", username, err)
		return nil, nil, err
	}

	return &user, profile, nil
}

func (s *UserService) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(user *models.User) error {
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}
	return s.db.Model(user).Updates(user).Error
}

func (s *UserService) DeleteUser(userID string) error {
	return s.db.Delete(&models.User{}, userID).Error
}

// HashPassword 使用 bcrypt 对密码进行哈希处理
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
} 

// ChangePassword 修改用户密码
func (s *UserService) ChangePassword(userID, oldPassword, newPassword string) error {
	// 获取用户信息
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("旧密码错误")
	}

	// 哈希新密码
	hashedNewPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	// 更新密码
	if err := s.db.Model(&user).Update("password", hashedNewPassword).Error; err != nil {
		return err
	}

	return nil
} 