package factory

import (
	"time"
	"gorm.io/gorm"
)

// Factory 工厂信息模型
type Factory struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"type:varchar(191);not null"`         // 工厂名称
	Username    string         `json:"username" gorm:"type:varchar(191);uniqueIndex"`  // 登录用户名
	Password    string         `json:"-" gorm:"type:varchar(191);not null"`           // 登录密码（不返回给前端）
	Address     string         `json:"address" gorm:"type:varchar(191)"`              // 工厂地址
	Contact     string         `json:"contact" gorm:"type:varchar(191)"`              // 联系人
	Phone       string         `json:"phone" gorm:"type:varchar(191)"`                // 联系电话
	Email       string         `json:"email" gorm:"type:varchar(191)"`                // 电子邮箱
	License     string         `json:"license" gorm:"type:varchar(191)"`              // 营业执照号
	Description string         `json:"description" gorm:"type:text"`                  // 工厂描述
	Status      int           `json:"status" gorm:"default:1"`                        // 状态：1-正常 2-禁用
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// RegisterRequest 工厂注册请求
type RegisterRequest struct {
	Name        string `json:"name" binding:"required"`
	Username    string `json:"username" binding:"required,min=4,max=20"`
	Password    string `json:"password" binding:"required,min=6,max=20"`
	Address     string `json:"address" binding:"required"`
	Contact     string `json:"contact" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	License     string `json:"license" binding:"required"`
	Description string `json:"description"`
}

// LoginRequest 工厂登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 工厂登录响应
type LoginResponse struct {
	Token   string  `json:"token"`
	Factory Factory `json:"factory"`
} 