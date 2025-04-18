package models

import (
	"time"
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	FileName    string    `json:"file_name"`
	FilePath    string    `json:"file_path"`
	FileType    string    `json:"file_type"`
	UploadedBy  uint      `json:"uploaded_by"`
	OrderID     uint      `json:"order_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
} 