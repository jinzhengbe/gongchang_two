package services

import (
	"aneworder.com/backend/models"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

// Custom error types
var (
	ErrUsernameExists = errors.New("username already exists")
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) Register(user *models.User) error {
	// 检查用户名是否已存在
	var existingUser models.User
	if err := s.db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return ErrUsernameExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 生成唯一 ID
	user.ID = uuid.New().String()

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// 创建用户
	return s.db.Create(user).Error
}

func (s *UserService) Login(username, password string) (*models.LoginResponse, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	// 根据用户角色获取相应的档案信息
	var profile interface{}
	switch user.Role {
	case models.RoleDesigner:
		var designerProfile models.DesignerProfile
		if err := s.db.Where("user_id = ?", user.ID).First(&designerProfile).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
		} else {
			profile = designerProfile
		}
	case models.RoleFactory:
		var factoryProfile models.FactoryProfile
		if err := s.db.Where("user_id = ?", user.ID).First(&factoryProfile).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
		} else {
			profile = factoryProfile
		}
	case models.RoleSupplier:
		var supplierProfile models.SupplierProfile
		if err := s.db.Where("user_id = ?", user.ID).First(&supplierProfile).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
		} else {
			profile = supplierProfile
		}
	}

	return &models.LoginResponse{
		Data: models.LoginData{
			User:    user,
			Profile: profile,
		},
	}, nil
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