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
	Role      string         `json:"role" gorm:"default:'user'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type DesignerProfile struct {
	gorm.Model
	UserID      uint   `gorm:"uniqueIndex"`
	CompanyName string
	Address     string
	Website     string
	Bio         string
}

type FactoryProfile struct {
	gorm.Model
	UserID      uint   `gorm:"uniqueIndex"`
	CompanyName string
	Address     string
	Capacity    int
	Equipment   string
	Certificates string
}

type SupplierProfile struct {
	gorm.Model
	UserID      uint   `gorm:"uniqueIndex"`
	CompanyName string
	Address     string
	MainProducts string
	Certificates string
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginData struct {
	User    User        `json:"user"`
	Profile interface{} `json:"profile"`
}

type LoginResponse struct {
	Data LoginData `json:"data"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof=designer factory supplier"`
} 