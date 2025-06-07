package factory

import (
	"time"
	"gorm.io/gorm"
)

// Factory 工厂信息模型
type Factory struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"type:varchar(191);not null"`         // 工厂名称
	Username    string         `json:"username" gorm:"type:varchar(191);uniqueIndex"`  // 登录用户名
	Password    string         `json:"-" gorm:"type:longtext;not null"`           // 登录密码（不返回给前端）
	Address     string         `json:"address" gorm:"type:varchar(191)"`              // 工厂地址
	Contact     string         `json:"contact" gorm:"type:varchar(191)"`              // 联系人
	Phone       string         `json:"phone" gorm:"type:varchar(191)"`                // 联系电话
	Email       string         `json:"email" gorm:"type:varchar(191)"`                // 电子邮箱
	License     string         `json:"license" gorm:"type:varchar(191)"`              // 营业执照号
	Description string         `json:"description" gorm:"type:text"`                  // 工厂描述
	Status      int           `json:"status" gorm:"default:1"`                        // 状态：1-正常 2-禁用
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// RegisterRequest 工厂注册请求
type RegisterRequest struct {
	Name        string `json:"name"`
	Username    string `json:"username" binding:"required,min=4,max=20"`
	Password    string `json:"password" binding:"required,min=6,max=20"`
	Address     string `json:"address"`
	Contact     string `json:"contact"`
	Phone       string `json:"phone"`
	Email       string `json:"email" binding:"omitempty,email"`
	License     string `json:"license"`
	Description string `json:"description"`
}

// LoginRequest 工厂登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

// LoginResponse 工厂登录响应
type LoginResponse struct {
	Token   string  `json:"token"`
	Factory Factory `json:"factory"`
}

// Order 订单模型
type Order struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	Title            string         `json:"title" gorm:"type:longtext;not null"`
	Description      string         `json:"description" gorm:"type:longtext"`
	Fabric           string         `json:"fabric" gorm:"type:longtext"`
	Quantity         int64          `json:"quantity"`
	FactoryID        string         `json:"factory_id" gorm:"type:varchar(191)"`
	Status           string         `json:"status" gorm:"type:varchar(191);default:draft"`
	DesignerID       string         `json:"designer_id" gorm:"type:longtext"`
	CustomerID       string         `json:"customer_id" gorm:"type:longtext"`
	UnitPrice        float64        `json:"unit_price"`
	TotalPrice       float64        `json:"total_price"`
	PaymentStatus    string         `json:"payment_status" gorm:"type:longtext"`
	ShippingAddress  string         `json:"shipping_address" gorm:"type:longtext"`
	OrderType        string         `json:"order_type" gorm:"type:longtext"`
	Fabrics          string         `json:"fabrics" gorm:"type:longtext"`
	DeliveryDate     *time.Time     `json:"delivery_date"`
	OrderDate        *time.Time     `json:"order_date"`
	SpecialRequirements string      `json:"special_requirements" gorm:"type:longtext"`
	Attachments      string         `json:"attachments" gorm:"type:json"`
	Models           string         `json:"models" gorm:"type:json"`
	Images           string         `json:"images" gorm:"type:json"`
	Videos           string         `json:"videos" gorm:"type:json"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// OrderItem 订单项模型
type OrderItem struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	OrderID     uint    `json:"order_id" gorm:"index"`                // 订单ID
	ProductName string  `json:"product_name" gorm:"type:varchar(191)"` // 产品名称
	Quantity    int     `json:"quantity"`                             // 数量
	Price       float64 `json:"price" gorm:"type:decimal(10,2)"`      // 单价
}

// OrderListRequest 订单列表请求
type OrderListRequest struct {
	Page      int    `form:"page" binding:"required,min=1"`
	PageSize  int    `form:"page_size" binding:"required,min=1,max=100"`
	Status    string `form:"status"`                                  // 订单状态（可选）
	StartDate string `form:"start_date"`                              // 开始日期（可选）
	EndDate   string `form:"end_date"`                                // 结束日期（可选）
	Title     string `form:"title"`                                   // 订单标题（可选）
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=created_at updated_at total_price"` // 排序字段
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`                        // 排序方式
}

// OrderListResponse 订单列表响应
type OrderListResponse struct {
	Total       int64   `json:"total"`        // 总数
	CurrentPage int     `json:"current_page"` // 当前页
	PageSize    int     `json:"page_size"`    // 每页数量
	Orders      []Order `json:"orders"`       // 订单列表
} 