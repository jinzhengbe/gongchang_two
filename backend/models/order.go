package models

import (
	"gorm.io/gorm"
	"time"
)

type OrderStatus string

const (
	StatusPending    OrderStatus = "pending"
	StatusAccepted   OrderStatus = "accepted"
	StatusInProgress OrderStatus = "in_progress"
	StatusCompleted  OrderStatus = "completed"
	StatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	gorm.Model
	DesignerID    uint      `json:"designer_id" gorm:"not null"`
	CustomerID    uint      `json:"customer_id" gorm:"not null"`
	ProductID     uint      `json:"product_id" gorm:"not null"`
	Quantity      int       `json:"quantity" gorm:"not null"`
	UnitPrice     float64   `json:"unit_price" gorm:"not null"`
	TotalPrice    float64   `json:"total_price" gorm:"not null"`
	Status        string    `json:"status" gorm:"not null;default:'pending'"`
	PaymentStatus string    `json:"payment_status" gorm:"not null;default:'unpaid'"`
	ShippingAddress string  `json:"shipping_address" gorm:"not null"`
	OrderDate     time.Time `json:"order_date" gorm:"not null"`
	Designer      User      `json:"designer" gorm:"foreignKey:DesignerID"`
	Customer      User      `json:"customer" gorm:"foreignKey:CustomerID"`
	Product       Product   `json:"product" gorm:"foreignKey:ProductID"`
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
	TotalOrders    int64            `json:"totalOrders"`
	ActiveOrders   int64            `json:"activeOrders"`
	CompletedOrders int64           `json:"completedOrders"`
	StatusCounts   map[string]int64 `json:"statusCounts"`
	TrendData      []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	} `json:"trendData"`
} 