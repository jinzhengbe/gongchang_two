package models

import (
	"gorm.io/gorm"
	"time"
	"gorm.io/datatypes"
)

type OrderStatus string

const (
	OrderStatusDraft     OrderStatus = "draft"
	OrderStatusPublished OrderStatus = "published"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt *time.Time     `json:"created_at" gorm:"autoCreateTime:false"`
	UpdatedAt *time.Time     `json:"updated_at" gorm:"autoUpdateTime:false"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	Title             string      `json:"title" gorm:"not null"`
	Description       string      `json:"description"`
	Fabric            string      `json:"fabric"`
	Quantity          int         `json:"quantity"`
	FactoryID         *string     `json:"factory_id"`
	Status            OrderStatus `json:"status" gorm:"type:varchar(191);default:'draft'"`
	Factory           FactoryProfile `json:"factory" gorm:"foreignKey:FactoryID;references:UserID"`
	DesignerID        string      `json:"designer_id"`
	CustomerID        string      `json:"customer_id"`
	UnitPrice         float64     `json:"unit_price"`
	TotalPrice        float64     `json:"total_price"`
	PaymentStatus     string      `json:"payment_status"`
	ShippingAddress   string      `json:"shipping_address"`
	OrderType         string      `json:"order_type"`
	Fabrics           string      `json:"fabrics"`
	DeliveryDate      *time.Time  `json:"delivery_date"`
	OrderDate         *time.Time  `json:"order_date"`
	SpecialRequirements string    `json:"special_requirements"`

	Attachments       *datatypes.JSON `json:"attachments" gorm:"type:jsonb;default:null"`
	Models            *datatypes.JSON `json:"models" gorm:"type:jsonb;default:null"`
	Images            *datatypes.JSON `json:"images" gorm:"type:jsonb;default:null"`
	Videos            *datatypes.JSON `json:"videos" gorm:"type:jsonb;default:null"`
	
	// 添加文件关联
	Files             []File         `json:"files" gorm:"foreignKey:OrderID"`
}

type OrderRequest struct {
	Title             string    `json:"title" binding:"required"`
	Description       string    `json:"description"`
	Fabric            string    `json:"fabric"`
	Quantity          int       `json:"quantity" binding:"required,min=1"`
	DesignerID        string    `json:"designer_id" binding:"required"`
	CustomerID        string    `json:"customer_id" binding:"required"`
	UnitPrice         float64   `json:"unit_price"`
	TotalPrice        float64   `json:"total_price"`
	Status            string    `json:"status"`
	PaymentStatus     string    `json:"payment_status"`
	ShippingAddress   string    `json:"shipping_address"`
	OrderType         string    `json:"orderType"`
	Fabrics           string    `json:"fabrics"`
	DeliveryDate      *time.Time `json:"deliveryDate"`
	OrderDate         *time.Time `json:"order_date"`
	SpecialRequirements string  `json:"specialRequirements"`
	Attachments       []string  `json:"attachments"`
	Models            []string  `json:"models"`
	Images            []string  `json:"images"`
	Videos            []string  `json:"videos"`
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

type OrderUpdateRequest struct {
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	Fabric            string    `json:"fabric"`
	Quantity          int       `json:"quantity"`
	Status            string    `json:"status"`
	PaymentStatus     string    `json:"payment_status"`
	ShippingAddress   string    `json:"shipping_address"`
	OrderType         string    `json:"orderType"`
	Fabrics           string    `json:"fabrics"`
	DeliveryDate      *time.Time `json:"deliveryDate"`
	OrderDate         *time.Time `json:"order_date"`
	SpecialRequirements string  `json:"specialRequirements"`
	Attachments       []string  `json:"attachments"`
	Models            []string  `json:"models"`
	Images            []string  `json:"images"`
	Videos            []string  `json:"videos"`
} 

// OrderDetailResponse 订单详情响应（包含布料详细信息）
type OrderDetailResponse struct {
	ID                 uint                    `json:"id"`
	Title              string                  `json:"title"`
	Description        string                  `json:"description"`
	Fabric             string                  `json:"fabric"`
	Quantity           int                     `json:"quantity"`
	FactoryID          *string                 `json:"factory_id"`
	Status             OrderStatus             `json:"status"`
	DesignerID         string                  `json:"designer_id"`
	CustomerID         string                  `json:"customer_id"`
	UnitPrice          float64                 `json:"unit_price"`
	TotalPrice         float64                 `json:"total_price"`
	PaymentStatus      string                  `json:"payment_status"`
	ShippingAddress    string                  `json:"shipping_address"`
	OrderType          string                  `json:"order_type"`
	Fabrics            []Fabric                `json:"fabrics"`           // 布料详细信息数组
	FabricsIDs         string                  `json:"fabrics_ids"`       // 原始布料ID字符串
	DeliveryDate       *time.Time              `json:"delivery_date"`
	OrderDate          *time.Time              `json:"order_date"`
	SpecialRequirements string                 `json:"special_requirements"`
	Attachments        []string                `json:"attachments"`
	Models             []string                `json:"models"`
	Images             []string                `json:"images"`
	Videos             []string                `json:"videos"`
	Files              []File                  `json:"files"`
	CreatedAt          *time.Time              `json:"created_at"`
	UpdatedAt          *time.Time              `json:"updated_at"`
} 

// RemoveFileFromOrderRequest 从订单移除文件的请求
type RemoveFileFromOrderRequest struct {
	FileID   string `json:"fileId" binding:"required"`   // 文件ID（UUID格式）
	FileType string `json:"fileType" binding:"required"` // 文件类型：image, attachment, model, video
}

// RemoveFileFromOrderResponse 从订单移除文件的响应
type RemoveFileFromOrderResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Order   *Order `json:"order,omitempty"`
} 