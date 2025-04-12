package models

import (
	"time"
)

type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Category    string    `json:"category" gorm:"not null"`
	Price       float64   `json:"price" gorm:"not null"`
	Stock       int       `json:"stock" gorm:"not null"`
	Status      string    `json:"status" gorm:"default:'active'"`
	CreatedBy   uint      `json:"created_by" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// 关联
	Creator     User      `json:"creator" gorm:"foreignKey:CreatedBy"`
}

type ProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Category    string  `json:"category" binding:"required"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Stock       int     `json:"stock" binding:"required,min=0"`
}

type ProductUpdateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Price       float64 `json:"price" binding:"min=0"`
	Stock       int     `json:"stock" binding:"min=0"`
	Status      string  `json:"status"`
} 