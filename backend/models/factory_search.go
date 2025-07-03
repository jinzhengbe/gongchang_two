package models

import "time"

// FactorySearchRequest 工厂搜索请求
type FactorySearchRequest struct {
	Query             string   `json:"query" form:"query"`                           // 搜索关键词
	Region            string   `json:"region" form:"region"`                         // 地区筛选
	Specialties       []string `json:"specialties" form:"specialties"`               // 专业领域数组
	CooperationStatus string   `json:"cooperation_status" form:"cooperation_status"` // 合作状态
	MinRating         float64  `json:"min_rating" form:"min_rating"`                 // 最低评分
	MaxRating         float64  `json:"max_rating" form:"max_rating"`                 // 最高评分
	Page              int      `json:"page" form:"page"`                             // 页码
	PageSize          int      `json:"page_size" form:"page_size"`                   // 每页数量
	SortBy            string   `json:"sort_by" form:"sort_by"`                       // 排序字段
	SortOrder         string   `json:"sort_order" form:"sort_order"`                 // 排序方向
}

// FactorySearchSuggestionRequest 工厂搜索建议请求
type FactorySearchSuggestionRequest struct {
	Query string `json:"query" form:"query" binding:"required"` // 搜索关键词
	Limit int    `json:"limit" form:"limit"`                    // 建议数量
}

// FactorySearchResponse 工厂搜索响应
type FactorySearchResponse struct {
	Success bool                    `json:"success"`
	Data    FactorySearchResultData `json:"data"`
}

// FactorySearchResultData 工厂搜索结果数据
type FactorySearchResultData struct {
	Factories []FactorySearchResult `json:"factories"`
	Total     int64                 `json:"total"`
	Page      int                   `json:"page"`
	PageSize  int                   `json:"page_size"`
}

// FactorySearchResult 工厂搜索结果
type FactorySearchResult struct {
	ID                uint      `json:"id"`
	Name              string    `json:"name"`
	Address           string    `json:"address"`
	Specialties       []string  `json:"specialties"`
	Rating            float64   `json:"rating"`
	CooperationStatus string    `json:"cooperation_status"`
	Description       string    `json:"description"`
	ContactInfo       ContactInfo `json:"contact_info"`
	Capacity          Capacity    `json:"capacity"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

// ContactInfo 联系信息
type ContactInfo struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
}

// Capacity 产能信息
type Capacity struct {
	MonthlyOrders int `json:"monthly_orders"`
	MaxOrderSize  int `json:"max_order_size"`
}

// FactorySearchSuggestionResponse 工厂搜索建议响应
type FactorySearchSuggestionResponse struct {
	Success bool                           `json:"success"`
	Data    FactorySearchSuggestionData    `json:"data"`
}

// FactorySearchSuggestionData 工厂搜索建议数据
type FactorySearchSuggestionData struct {
	Suggestions []FactorySearchSuggestion `json:"suggestions"`
}

// FactorySearchSuggestion 工厂搜索建议
type FactorySearchSuggestion struct {
	Type      string `json:"type"`       // 建议类型：factory_name, factory_address, specialty
	Text      string `json:"text"`       // 建议文本
	Highlight string `json:"highlight"`  // 高亮显示的文本
}

// FactorySpecialty 工厂专业领域
type FactorySpecialty struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	FactoryID  uint   `json:"factory_id"`
	Specialty  string `json:"specialty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// FactoryRating 工厂评分
type FactoryRating struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	FactoryID  uint      `json:"factory_id"`
	Rating     float64   `json:"rating"`
	Comment    string    `json:"comment"`
	RaterID    string    `json:"rater_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
} 