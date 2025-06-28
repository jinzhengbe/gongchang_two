package models

import (
	"time"
)

type File struct {
	ID        string     `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	OrderID   *uint      `json:"order_id,omitempty" gorm:"index"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// AddFileToOrderRequest 添加文件到订单的请求模型
type AddFileToOrderRequest struct {
	Type        string `form:"type" binding:"required,oneof=image attachment model video" json:"type"`
	Description string `form:"description" json:"description"`
}

// AddFileToOrderResponse 添加文件到订单的响应模型
type AddFileToOrderResponse struct {
	Success bool        `json:"success"`
	Order   *Order     `json:"order"`
	File    *FileInfo  `json:"file"`
}

// FileInfo 文件信息
type FileInfo struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}