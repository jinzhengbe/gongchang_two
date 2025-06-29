package models

import (
	"time"
	"gorm.io/gorm"
)

// JiedanStatus 接单状态
type JiedanStatus string

const (
	JiedanStatusPending  JiedanStatus = "pending"  // 待处理
	JiedanStatusAccepted JiedanStatus = "accepted" // 已同意
	JiedanStatusRejected JiedanStatus = "rejected" // 已拒绝
)

// Jiedan 接单模型
type Jiedan struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	OrderID      uint           `json:"order_id" gorm:"not null;index"`
	FactoryID    string         `json:"factory_id" gorm:"type:varchar(191);not null;index"`
	Status       JiedanStatus   `json:"status" gorm:"type:varchar(50);not null;default:'pending';index"`
	Price        *float64       `json:"price" gorm:"type:decimal(10,2);comment:接单价格"`
	JiedanTime   *time.Time     `json:"jiedan_time" gorm:"comment:接单时间"`
	AgreeTime    *time.Time     `json:"agree_time" gorm:"comment:同意时间"`
	AgreeUserID  *string        `json:"agree_user_id" gorm:"type:varchar(191);comment:同意的用户ID"`
	CreatedAt    *time.Time     `json:"created_at" gorm:"autoCreateTime:false"`
	UpdatedAt    *time.Time     `json:"updated_at" gorm:"autoUpdateTime:false"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	
	// 关联关系
	Order   Order  `json:"order" gorm:"foreignKey:OrderID"`
	Factory FactoryProfile `json:"factory" gorm:"foreignKey:FactoryID;references:UserID"`
}

// TableName 指定表名
func (Jiedan) TableName() string {
	return "jiedan"
}

// CreateJiedanRequest 创建接单请求
type CreateJiedanRequest struct {
	OrderID   uint     `json:"order_id" binding:"required"`
	FactoryID string   `json:"factory_id" binding:"required"`
	Price     *float64 `json:"price"`
}

// UpdateJiedanRequest 更新接单请求
type UpdateJiedanRequest struct {
	Status      JiedanStatus `json:"status"`
	Price       *float64     `json:"price"`
	AgreeUserID string       `json:"agree_user_id"`
}

// JiedanResponse 接单响应
type JiedanResponse struct {
	ID           uint         `json:"id"`
	OrderID      uint         `json:"order_id"`
	FactoryID    string       `json:"factory_id"`
	Status       JiedanStatus `json:"status"`
	Price        *float64     `json:"price"`
	JiedanTime   *time.Time   `json:"jiedan_time"`
	AgreeTime    *time.Time   `json:"agree_time"`
	AgreeUserID  *string      `json:"agree_user_id"`
	CreatedAt    *time.Time   `json:"created_at"`
	UpdatedAt    *time.Time   `json:"updated_at"`
	
	// 关联数据
	Order   *Order  `json:"order,omitempty"`
	Factory *FactoryProfile `json:"factory,omitempty"`
}

// JiedanListResponse 接单列表响应
type JiedanListResponse struct {
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
	Jiedans  []JiedanResponse `json:"jiedans"`
}

// AcceptJiedanRequest 同意接单请求
type AcceptJiedanRequest struct {
	AgreeUserID string `json:"agree_user_id" binding:"required"`
}

// RejectJiedanRequest 拒绝接单请求
type RejectJiedanRequest struct {
	Reason string `json:"reason"`
} 