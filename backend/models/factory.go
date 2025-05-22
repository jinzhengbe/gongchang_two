package models

import (
	"time"
	"gorm.io/gorm"
)

type Factory struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"type:varchar(191);not null"`
	Username    string         `json:"username" gorm:"type:varchar(191);uniqueIndex"`
	Password    string         `json:"-" gorm:"type:varchar(191);not null"`
	Address     string         `json:"address" gorm:"type:varchar(191)"`
	Contact     string         `json:"contact" gorm:"type:varchar(191)"`
	Phone       string         `json:"phone" gorm:"type:varchar(191)"`
	Email       string         `json:"email" gorm:"type:varchar(191)"`
	License     string         `json:"license" gorm:"type:varchar(191)"`
	Description string         `json:"description" gorm:"type:text"`
	Status      int           `json:"status" gorm:"default:1"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type FactoryRegisterRequest struct {
	Name        string `json:"name" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Address     string `json:"address"`
	Contact     string `json:"contact"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	License     string `json:"license"`
	Description string `json:"description"`
}

type FactoryLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type FactoryResponse struct {
	Factory Factory `json:"factory"`
	Token   string  `json:"token"`
} 