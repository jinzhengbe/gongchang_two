package models

import (
	"time"
)

type File struct {
	ID        string     `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	Type      string     `json:"type"`                    // 文件类型：image, attachment, model, video
	OrderID   *uint      `json:"order_id,omitempty" gorm:"index"`
	FactoryID string     `json:"factory_id,omitempty" gorm:"index"` // 新增：关联工厂
	Category  string     `json:"category,omitempty"`                 // 新增：图片分类
	Size      int64      `json:"size,omitempty"`                     // 新增：文件大小
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

// 工厂图片相关模型

// BatchUploadFactoryPhotosRequest 批量上传工厂图片请求
type BatchUploadFactoryPhotosRequest struct {
	Category string `form:"category" json:"category"` // 可选：图片分类
}

// BatchUploadFactoryPhotosResponse 批量上传工厂图片响应
type BatchUploadFactoryPhotosResponse struct {
	Success      bool                `json:"success"`
	Message      string              `json:"message"`
	UploadedCount int                `json:"uploaded_count"`
	FailedCount  int                 `json:"failed_count"`
	Photos       []*FactoryPhotoInfo `json:"photos"`
	FailedFiles  []*FailedFileInfo   `json:"failed_files,omitempty"`
}

// FactoryPhotoInfo 工厂图片信息
type FactoryPhotoInfo struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	URL           string `json:"url"`
	ThumbnailURL  string `json:"thumbnail_url,omitempty"`
	Category      string `json:"category,omitempty"`
	Size          int64  `json:"size"`
	FactoryID     string `json:"factory_id"`
	Status        string `json:"status"` // success, failed
	CreatedAt     string `json:"created_at"`
}

// FailedFileInfo 失败文件信息
type FailedFileInfo struct {
	Name  string `json:"name"`
	Error string `json:"error"`
}

// GetFactoryPhotosRequest 获取工厂图片列表请求
type GetFactoryPhotosRequest struct {
	Category string `form:"category" json:"category"` // 可选：按分类筛选
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
}

// GetFactoryPhotosResponse 获取工厂图片列表响应
type GetFactoryPhotosResponse struct {
	Success     bool                `json:"success"`
	Total       int64               `json:"total"`
	Photos      []*FactoryPhotoInfo `json:"photos"`
	Categories  []*PhotoCategory    `json:"categories,omitempty"`
}

// PhotoCategory 图片分类
type PhotoCategory struct {
	ID        uint   `json:"id"`
	FactoryID string `json:"factory_id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	Count     int    `json:"count"`
}

// BatchDeleteFactoryPhotosRequest 批量删除工厂图片请求
type BatchDeleteFactoryPhotosRequest struct {
	PhotoIDs []string `json:"photo_ids" binding:"required"`
}

// BatchDeleteFactoryPhotosResponse 批量删除工厂图片响应
type BatchDeleteFactoryPhotosResponse struct {
	Success       bool     `json:"success"`
	Message       string   `json:"message"`
	DeletedCount  int      `json:"deleted_count"`
	FailedCount   int      `json:"failed_count"`
	FailedPhotoIDs []string `json:"failed_photo_ids,omitempty"`
}