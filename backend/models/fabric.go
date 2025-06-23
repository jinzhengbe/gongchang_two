package models

import (
	"time"
	"gorm.io/gorm"
)

// Fabric 布料模型
type Fabric struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"type:varchar(191);not null"`         // 布料名称
	Category    string         `json:"category" gorm:"type:varchar(191)"`              // 布料类别
	Material    string         `json:"material" gorm:"type:varchar(191)"`              // 材质
	Color       string         `json:"color" gorm:"type:varchar(191)"`                 // 颜色
	Pattern     string         `json:"pattern" gorm:"type:varchar(191)"`               // 花纹/图案
	Weight      float64        `json:"weight" gorm:"type:decimal(8,2)"`               // 克重 (g/m²)
	Width       float64        `json:"width" gorm:"type:decimal(8,2)"`                // 幅宽 (cm)
	Price       float64        `json:"price" gorm:"type:decimal(10,2)"`               // 单价 (元/米)
	Unit        string         `json:"unit" gorm:"type:varchar(50);default:'米'"`      // 单位
	Stock       int            `json:"stock" gorm:"default:0"`                        // 库存数量
	MinOrder    int            `json:"min_order" gorm:"default:1"`                    // 最小订购量
	Description string         `json:"description" gorm:"type:text"`                  // 描述
	ImageURL    string         `json:"image_url" gorm:"type:varchar(500)"`            // 图片URL
	ThumbnailURL string        `json:"thumbnail_url" gorm:"type:varchar(500)"`        // 缩略图URL
	Tags        string         `json:"tags" gorm:"type:varchar(500)"`                 // 标签，逗号分隔
	Status      int            `json:"status" gorm:"default:1"`                       // 状态：1-可用 0-停用
	SupplierID  *string        `json:"supplier_id" gorm:"type:varchar(191)"`          // 供应商ID
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// FabricCategory 布料分类
type FabricCategory struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"type:varchar(191);not null"`        // 分类名称
	Description string         `json:"description" gorm:"type:text"`                  // 分类描述
	Icon        string         `json:"icon" gorm:"type:varchar(191)"`                 // 分类图标
	Sort        int            `json:"sort" gorm:"default:0"`                         // 排序
	Status      int            `json:"status" gorm:"default:1"`                       // 状态：1-启用 0-禁用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// FabricRequest 创建布料请求
type FabricRequest struct {
	Name        string  `json:"name" binding:"required"`
	Category    string  `json:"category"`
	Material    string  `json:"material"`
	Color       string  `json:"color"`
	Pattern     string  `json:"pattern"`
	Weight      float64 `json:"weight"`
	Width       float64 `json:"width"`
	Price       float64 `json:"price"`
	Unit        string  `json:"unit"`
	Stock       int     `json:"stock"`
	MinOrder    int     `json:"min_order"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	ThumbnailURL string `json:"thumbnail_url"`
	Tags        string  `json:"tags"`
	SupplierID  string  `json:"supplier_id"`
}

// FabricUpdateRequest 更新布料请求
type FabricUpdateRequest struct {
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Material    string  `json:"material"`
	Color       string  `json:"color"`
	Pattern     string  `json:"pattern"`
	Weight      float64 `json:"weight"`
	Width       float64 `json:"width"`
	Price       float64 `json:"price"`
	Unit        string  `json:"unit"`
	Stock       int     `json:"stock"`
	MinOrder    int     `json:"min_order"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	ThumbnailURL string `json:"thumbnail_url"`
	Tags        string  `json:"tags"`
	Status      int     `json:"status"`
	SupplierID  string  `json:"supplier_id"`
}

// FabricSearchRequest 布料搜索请求
type FabricSearchRequest struct {
	Query    string `json:"query" form:"q"`                    // 搜索关键词
	Category string `json:"category" form:"category"`          // 分类筛选
	Material string `json:"material" form:"material"`          // 材质筛选
	Color    string `json:"color" form:"color"`               // 颜色筛选
	MinPrice float64 `json:"min_price" form:"min_price"`      // 最低价格
	MaxPrice float64 `json:"max_price" form:"max_price"`      // 最高价格
	MinStock int     `json:"min_stock" form:"min_stock"`      // 最低库存
	Status   int     `json:"status" form:"status"`            // 状态筛选
	Page     int     `json:"page" form:"page"`                // 页码
	PageSize int     `json:"page_size" form:"page_size"`      // 每页数量
}

// FabricResponse 布料响应
type FabricResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Material    string    `json:"material"`
	Color       string    `json:"color"`
	Pattern     string    `json:"pattern"`
	Weight      float64   `json:"weight"`
	Width       float64   `json:"width"`
	Price       float64   `json:"price"`
	Unit        string    `json:"unit"`
	Stock       int       `json:"stock"`
	MinOrder    int       `json:"min_order"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	ThumbnailURL string   `json:"thumbnail_url"`
	Tags        string    `json:"tags"`
	Status      int       `json:"status"`
	SupplierID  *string   `json:"supplier_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// FabricListResponse 布料列表响应
type FabricListResponse struct {
	Total   int64            `json:"total"`
	Page    int              `json:"page"`
	PageSize int             `json:"page_size"`
	Fabrics []FabricResponse `json:"fabrics"`
}

// FabricCategoryResponse 布料分类响应
type FabricCategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status"`
} 