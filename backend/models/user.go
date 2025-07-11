package models

import (
	"gorm.io/gorm"
	"time"
)

type UserRole string

const (
	RoleDesigner UserRole = "designer"
	RoleFactory  UserRole = "factory"
	RoleSupplier UserRole = "supplier"
)

type User struct {
	ID        string         `json:"id" gorm:"primaryKey;type:varchar(191)"`
	Username  string         `json:"username" gorm:"unique;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Email     string         `json:"email" gorm:"not null"`
	Role      UserRole       `json:"role" gorm:"type:varchar(191)"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type DesignerProfile struct {
	gorm.Model
	UserID      string `gorm:"uniqueIndex;type:varchar(191)"`
	User        User   `gorm:"foreignKey:UserID"`
	CompanyName string
	Address     string
	Website     string
	Bio         string
	Rating      float64 `gorm:"default:0"` // 设计师评分
	Status      int     `gorm:"default:1"` // 设计师状态：1-正常，0-停用
}

type FactoryProfile struct {
	gorm.Model
	UserID      string `gorm:"uniqueIndex;type:varchar(191)"`
	User        User   `gorm:"foreignKey:UserID"`
	CompanyName string
	Address     string
	Capacity    int
	Equipment   string
	Certificates string
	Photos      string `gorm:"type:json"` // 工厂照片，JSON格式存储多个照片URL
	Videos      string `gorm:"type:json"` // 工厂视频，JSON格式存储多个视频URL
	EmployeeCount int  `gorm:"default:0"` // 员工数量
	Rating      float64 `gorm:"default:0"` // 工厂评分
	Status      int     `gorm:"default:1"` // 工厂状态：1-正常，0-停用
}

type SupplierProfile struct {
	gorm.Model
	UserID      string `gorm:"uniqueIndex;type:varchar(191)"`
	CompanyName string
	Address     string
	MainProducts string
	Certificates string
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Role         string `json:"role" binding:"required,oneof=designer factory supplier"`
	CompanyName  string `json:"company_name"`
	Address      string `json:"address"`
	Bio          string `json:"bio"`
	MainProducts string `json:"main_products"`
}

type UpdateProfileRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
} 

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
} 

// UpdateFactoryProfileRequest 更新工厂信息请求
type UpdateFactoryProfileRequest struct {
	CompanyName   string   `json:"company_name"`
	Address       string   `json:"address"`
	Capacity      *int     `json:"capacity"`
	Equipment     string   `json:"equipment"`
	Certificates  string   `json:"certificates"`
	Photos        []string `json:"photos"` // 工厂照片URL数组
	Videos        []string `json:"videos"` // 工厂视频URL数组
	EmployeeCount *int     `json:"employee_count"`
} 