package models

import "time"

// DesignerSearchRequest 设计师搜索请求
type DesignerSearchRequest struct {
	Query             string   `json:"query" form:"query"`                           // 搜索关键词
	Region            string   `json:"region" form:"region"`                         // 地区筛选
	Specialties       []string `json:"specialties" form:"specialties"`               // 专业领域数组
	MinRating         float64  `json:"min_rating" form:"min_rating"`                 // 最低评分
	MaxRating         float64  `json:"max_rating" form:"max_rating"`                 // 最高评分
	Page              int      `json:"page" form:"page"`                             // 页码
	PageSize          int      `json:"page_size" form:"page_size"`                   // 每页数量
	SortBy            string   `json:"sort_by" form:"sort_by"`                       // 排序字段
	SortOrder         string   `json:"sort_order" form:"sort_order"`                 // 排序方向
}

// DesignerSearchSuggestionRequest 设计师搜索建议请求
type DesignerSearchSuggestionRequest struct {
	Query string `json:"query" form:"query" binding:"required"` // 搜索关键词
	Limit int    `json:"limit" form:"limit"`                    // 建议数量
}

// DesignerSearchResponse 设计师搜索响应
type DesignerSearchResponse struct {
	Success bool                      `json:"success"`
	Data    DesignerSearchResultData  `json:"data"`
}

// DesignerSearchResultData 设计师搜索结果数据
type DesignerSearchResultData struct {
	Designers []DesignerSearchResult `json:"designers"`
	Total     int64                  `json:"total"`
	Page      int                    `json:"page"`
	PageSize  int                    `json:"page_size"`
}

// DesignerSearchResult 设计师搜索结果
type DesignerSearchResult struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	Specialties []string  `json:"specialties"`
	Rating      float64   `json:"rating"`
	Description string    `json:"description"`
	ContactInfo ContactInfo `json:"contact_info"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// DesignerSearchSuggestionResponse 设计师搜索建议响应
type DesignerSearchSuggestionResponse struct {
	Success bool                            `json:"success"`
	Data    DesignerSearchSuggestionData    `json:"data"`
}

// DesignerSearchSuggestionData 设计师搜索建议数据
type DesignerSearchSuggestionData struct {
	Suggestions []DesignerSearchSuggestion `json:"suggestions"`
}

// DesignerSearchSuggestion 设计师搜索建议
type DesignerSearchSuggestion struct {
	Type      string `json:"type"`       // 建议类型：designer_name, designer_address, specialty
	Text      string `json:"text"`       // 建议文本
	Highlight string `json:"highlight"`  // 高亮显示的文本
}

// DesignerSpecialty 设计师专业领域
type DesignerSpecialty struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	DesignerID uint   `json:"designer_id"`
	Specialty  string `json:"specialty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// DesignerRating 设计师评分
type DesignerRating struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	DesignerID uint      `json:"designer_id"`
	Rating     float64   `json:"rating"`
	Comment    string    `json:"comment"`
	RaterID    string    `json:"rater_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
} 