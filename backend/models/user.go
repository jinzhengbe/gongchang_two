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
	CompanyName string
	Address     string
	Website     string
	Bio         string
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