package models

import (
	"time"
	"gorm.io/gorm"
)

// ProgressType 进度类型
type ProgressType string

const (
	ProgressTypeDesign      ProgressType = "design"      // 设计阶段
	ProgressTypeMaterial    ProgressType = "material"    // 材料准备
	ProgressTypeProduction  ProgressType = "production"  // 生产阶段
	ProgressTypeQuality     ProgressType = "quality"     // 质检阶段
	ProgressTypePackaging   ProgressType = "packaging"   // 包装阶段
	ProgressTypeShipping    ProgressType = "shipping"    // 发货阶段
	ProgressTypeCustom      ProgressType = "custom"      // 自定义阶段
)

// ProgressStatus 进度状态
type ProgressStatus string

const (
	ProgressStatusNotStarted ProgressStatus = "not_started" // 未开始
	ProgressStatusInProgress ProgressStatus = "in_progress" // 进行中
	ProgressStatusCompleted  ProgressStatus = "completed"   // 已完成
	ProgressStatusDelayed    ProgressStatus = "delayed"     // 延期
	ProgressStatusOnHold     ProgressStatus = "on_hold"     // 暂停
)

// OrderProgress 订单进度模型（符合要求文档）
type OrderProgress struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	OrderID       uint           `json:"order_id" gorm:"not null;index"`
	FactoryID     string         `json:"factory_id" gorm:"type:varchar(191);not null;index"`
	Type          ProgressType   `json:"type" gorm:"type:varchar(50);not null;comment:进度类型"`
	Status        ProgressStatus `json:"status" gorm:"type:varchar(50);not null;default:'not_started';comment:进度状态"`
	Description   string         `json:"description" gorm:"type:text;comment:进度描述"`
	StartTime     *time.Time     `json:"start_time" gorm:"comment:开始时间"`
	CompletedTime *time.Time     `json:"completed_time" gorm:"comment:完成时间"`
	Images        string         `json:"images" gorm:"type:text;comment:图片URL数组(JSON格式)"`
	CreatedAt     *time.Time     `json:"created_at" gorm:"autoCreateTime:false"`
	UpdatedAt     *time.Time     `json:"updated_at" gorm:"autoUpdateTime:false"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	
	// 关联关系
	Order   Order          `json:"order" gorm:"foreignKey:OrderID"`
	Factory FactoryProfile `json:"factory" gorm:"foreignKey:FactoryID;references:UserID"`
}

// TableName 指定表名
func (OrderProgress) TableName() string {
	return "order_progress"
}

// CreateProgressRequest 创建进度请求（符合要求文档）
type CreateProgressRequest struct {
	OrderID       uint         `json:"order_id" binding:"required"`
	FactoryID     string       `json:"factory_id" binding:"required"`
	Type          ProgressType `json:"type" binding:"required"`
	Status        ProgressStatus `json:"status" binding:"required"`
	Description   string       `json:"description"`
	StartTime     *time.Time   `json:"start_time"`
	CompletedTime *time.Time   `json:"completed_time"`
	Images        []string     `json:"images"`
}

// UpdateProgressRequest 更新进度请求（符合要求文档）
type UpdateProgressRequest struct {
	Type          ProgressType `json:"type"`
	Status        ProgressStatus `json:"status"`
	Description   string       `json:"description"`
	StartTime     *time.Time   `json:"start_time"`
	CompletedTime *time.Time   `json:"completed_time"`
	Images        []string     `json:"images"`
}

// ProgressResponse 进度响应（符合要求文档）
type ProgressResponse struct {
	ID            uint         `json:"id"`
	OrderID       uint         `json:"order_id"`
	FactoryID     string       `json:"factory_id"`
	Type          ProgressType `json:"type"`
	Status        ProgressStatus `json:"status"`
	Description   string       `json:"description"`
	StartTime     *time.Time   `json:"start_time"`
	CompletedTime *time.Time   `json:"completed_time"`
	Images        []string     `json:"images"`
	CreatedAt     *time.Time   `json:"created_at"`
	UpdatedAt     *time.Time   `json:"updated_at"`
	
	// 关联数据
	Order   *Order          `json:"order,omitempty"`
	Factory *FactoryProfile `json:"factory,omitempty"`
}

// ProgressListResponse 进度列表响应
type ProgressListResponse struct {
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
	Progress []ProgressResponse `json:"progress"`
}

// ProgressStatistics 进度统计
type ProgressStatistics struct {
	NotStarted  int64 `json:"not_started"`
	InProgress  int64 `json:"in_progress"`
	Completed   int64 `json:"completed"`
	Delayed     int64 `json:"delayed"`
	OnHold      int64 `json:"on_hold"`
	Total       int64 `json:"total"`
} 