package models

import (
	"gorm.io/gorm"
	"time"
)

type OrderStatus string

const (
	OrderStatusDraft     OrderStatus = "draft"
	OrderStatusPublished OrderStatus = "published"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	gorm.Model
	Title       string      `json:"title" gorm:"not null"`
	Description string      `json:"description"`
	Fabric      string      `json:"fabric"`
	Quantity    int         `json:"quantity"`
	FactoryID   uint        `json:"factory_id" gorm:"not null"`
	Status      OrderStatus `json:"status" gorm:"type:varchar(191);default:'draft'"`
	Factory     FactoryProfile `json:"factory" gorm:"foreignKey:FactoryID"`
}

type PublicOrder struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Fabric      string    `json:"fabric"`
	Quantity    int       `json:"quantity"`
	Factory     string    `json:"factory"`
	Status      string    `json:"status"`
	CreateTime  time.Time `json:"createTime"`
}

type PublicOrderResponse struct {
	Orders     []PublicOrder `json:"orders"`
	Total      int          `json:"total"`
	Page       int          `json:"page"`
	PageSize   int          `json:"pageSize"`
	TotalPages int          `json:"totalPages"`
}

type OrderProgress struct {
	gorm.Model
	OrderID     uint      `gorm:"not null"`
	Status      string    `gorm:"not null"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"not null"`
	CreatedBy   uint      `gorm:"not null"`
	Order       Order     `gorm:"foreignKey:OrderID"`
	User        User      `gorm:"foreignKey:CreatedBy"`
}

type OrderAttachment struct {
	gorm.Model
	OrderID     uint   `gorm:"not null"`
	FileName    string `gorm:"not null"`
	FilePath    string `gorm:"not null"`
	FileType    string `gorm:"not null"`
	UploadedBy  uint   `gorm:"not null"`
	Order       Order  `gorm:"foreignKey:OrderID"`
	User        User   `gorm:"foreignKey:UploadedBy"`
}

type OrderStatistics struct {
	TotalOrders     int64            `json:"totalOrders"`
	ActiveOrders    int64            `json:"activeOrders"`
	CompletedOrders int64           `json:"completedOrders"`
	PendingOrders   int64           `json:"pendingOrders"`
	StatusCounts    map[string]int64 `json:"statusCounts"`
	TrendData       []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	} `json:"trendData"`
} 