package models

import (
	"time"
)

// OrderSearchRequest 订单搜索请求
type OrderSearchRequest struct {
	Query      string `form:"query" json:"query"`           // 搜索关键词
	Status     string `form:"status" json:"status"`         // 订单状态筛选
	StartDate  string `form:"start_date" json:"start_date"` // 开始日期
	EndDate    string `form:"end_date" json:"end_date"`     // 结束日期
	Page       int    `form:"page" json:"page"`             // 页码
	PageSize   int    `form:"page_size" json:"page_size"`   // 每页数量
	SortBy     string `form:"sort_by" json:"sort_by"`       // 排序字段
	SortOrder  string `form:"sort_order" json:"sort_order"` // 排序方向
	UserID     string `form:"user_id" json:"user_id"`       // 用户ID（用于权限控制）
	UserRole   string `form:"user_role" json:"user_role"`   // 用户角色
}

// OrderSearchResponse 订单搜索响应
type OrderSearchResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Orders   []OrderSearchItem `json:"orders"`
		Total    int64             `json:"total"`
		Page     int               `json:"page"`
		PageSize int               `json:"page_size"`
	} `json:"data"`
}

// OrderSearchItem 订单搜索项
type OrderSearchItem struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	OrderNo   string    `json:"order_no"`
	Status    string    `json:"status"`
	Fabrics   []string  `json:"fabrics"`
	Factory   FactoryInfo `json:"factory"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FactoryInfo 工厂信息
type FactoryInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SearchSuggestionRequest 搜索建议请求
type SearchSuggestionRequest struct {
	Query string `form:"query" json:"query" binding:"required"` // 搜索关键词
	Limit int    `form:"limit" json:"limit"`                    // 建议数量
}

// SearchSuggestionResponse 搜索建议响应
type SearchSuggestionResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Suggestions []SearchSuggestion `json:"suggestions"`
	} `json:"data"`
}

// SearchSuggestion 搜索建议
type SearchSuggestion struct {
	Type      string `json:"type"`      // 建议类型：order_title, fabric_name, factory_name
	Text      string `json:"text"`      // 建议文本
	Highlight string `json:"highlight"` // 高亮显示的文本
}

// OrderSearchResult 订单搜索结果（内部使用）
type OrderSearchResult struct {
	Order       Order
	Relevance   float64 // 相关性分数
	MatchFields []string // 匹配的字段
} 