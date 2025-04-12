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
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Role      string    `json:"role" gorm:"not null"` // designer, factory, supplier
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof=designer factory supplier"`
} 