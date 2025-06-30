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

// OrderProgress 订单进度模型（扩展版）
type OrderProgress struct {
	ID                      uint           `json:"id" gorm:"primaryKey"`
	OrderID                 uint           `json:"order_id" gorm:"not null;index"`
	FactoryID               string         `json:"factory_id" gorm:"type:varchar(191);not null;index"`
	ProgressType            ProgressType   `json:"progress_type" gorm:"type:varchar(50);not null;comment:进度类型"`
	Percentage              *int           `json:"percentage" gorm:"comment:完成百分比(0-100)"`
	Status                  ProgressStatus `json:"status" gorm:"type:varchar(50);not null;default:'not_started';comment:进度状态"`
	Description             string         `json:"description" gorm:"type:text;comment:进度描述"`
	EstimatedCompletionTime *time.Time     `json:"estimated_completion_time" gorm:"comment:预计完成时间"`
	ActualCompletionTime    *time.Time     `json:"actual_completion_time" gorm:"comment:实际完成时间"`
	CreatorID               string         `json:"creator_id" gorm:"type:varchar(191);not null;comment:创建者ID"`
	CreatedAt               *time.Time     `json:"created_at" gorm:"autoCreateTime:false"`
	UpdatedAt               *time.Time     `json:"updated_at" gorm:"autoUpdateTime:false"`
	DeletedAt               gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	
	// 关联关系
	Order   Order          `json:"order" gorm:"foreignKey:OrderID"`
	Factory FactoryProfile `json:"factory" gorm:"foreignKey:FactoryID;references:UserID"`
}

// TableName 指定表名
func (OrderProgress) TableName() string {
	return "order_progress"
}

// CreateProgressRequest 创建进度请求
type CreateProgressRequest struct {
	OrderID                 uint         `json:"order_id" binding:"required"`
	FactoryID               string       `json:"factory_id" binding:"required"`
	ProgressType            ProgressType `json:"progress_type" binding:"required"`
	Percentage              *int         `json:"percentage"`
	Status                  ProgressStatus `json:"status" binding:"required"`
	Description             string       `json:"description"`
	EstimatedCompletionTime *time.Time   `json:"estimated_completion_time"`
	ActualCompletionTime    *time.Time   `json:"actual_completion_time"`
	CreatorID               string       `json:"creator_id" binding:"required"`
}

// UpdateProgressRequest 更新进度请求
type UpdateProgressRequest struct {
	ProgressType            ProgressType `json:"progress_type"`
	Percentage              *int         `json:"percentage"`
	Status                  ProgressStatus `json:"status"`
	Description             string       `json:"description"`
	EstimatedCompletionTime *time.Time   `json:"estimated_completion_time"`
	ActualCompletionTime    *time.Time   `json:"actual_completion_time"`
}

// ProgressResponse 进度响应
type ProgressResponse struct {
	ID                      uint         `json:"id"`
	OrderID                 uint         `json:"order_id"`
	FactoryID               string       `json:"factory_id"`
	ProgressType            ProgressType `json:"progress_type"`
	Percentage              *int         `json:"percentage"`
	Status                  ProgressStatus `json:"status"`
	Description             string       `json:"description"`
	EstimatedCompletionTime *time.Time   `json:"estimated_completion_time"`
	ActualCompletionTime    *time.Time   `json:"actual_completion_time"`
	CreatorID               string       `json:"creator_id"`
	CreatedAt               *time.Time   `json:"created_at"`
	UpdatedAt               *time.Time   `json:"updated_at"`
	
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