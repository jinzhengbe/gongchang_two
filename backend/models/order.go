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
	DesignerID         string    `json:"designer_id" gorm:"type:varchar(64);not null;index"`
	CustomerID         string    `json:"customer_id" gorm:"type:varchar(64);not null;index"`
	ProductID          *uint     `json:"product_id" gorm:"index;constraint:OnDelete:SET NULL"`
	Quantity           int       `json:"quantity" gorm:"not null"`
	UnitPrice          float64   `json:"unit_price" gorm:"not null;default:0"`
	TotalPrice         float64   `json:"total_price" gorm:"not null;default:0"`
	Status             string    `json:"status" gorm:"not null;default:'pending'"`
	PaymentStatus      string    `json:"payment_status" gorm:"not null;default:'unpaid'"`
	ShippingAddress    string    `json:"shipping_address" gorm:"type:longtext;not null"`
	OrderDate          time.Time `json:"order_date" gorm:"not null"`
	Title              string    `json:"title" gorm:"type:varchar(255);not null"`
	Description        string    `json:"description" gorm:"type:text"`
	OrderType          string    `json:"orderType" gorm:"type:varchar(50);not null"`
	Fabrics            string    `json:"fabrics" gorm:"type:text"`
	DeliveryDate       time.Time `json:"deliveryDate" gorm:"not null"`
	SpecialRequirements string    `json:"specialRequirements" gorm:"type:text"`
	Designer           User      `json:"designer" gorm:"foreignKey:DesignerID;references:ID"`
	Customer           User      `json:"customer" gorm:"foreignKey:CustomerID;references:ID"`
	Product            *Product  `json:"product" gorm:"foreignKey:ProductID"`
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