package services

import (
	"fmt"
	"strings"
	"gongChang/models"
	"gorm.io/gorm"
)

type DesignerSearchService struct {
	db *gorm.DB
}

func NewDesignerSearchService(db *gorm.DB) *DesignerSearchService {
	return &DesignerSearchService{db: db}
}

// SearchDesigners 搜索设计师
func (s *DesignerSearchService) SearchDesigners(req *models.DesignerSearchRequest) (*models.DesignerSearchResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}
	if req.SortBy == "" {
		req.SortBy = "rating"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	// 构建基础查询
	query := s.db.Model(&models.DesignerProfile{}).
		Joins("JOIN users ON designer_profiles.user_id = users.id").
		Where("users.role = ? AND users.deleted_at IS NULL", "designer")

	// 添加搜索条件
	if req.Query != "" {
		searchQuery := "%" + req.Query + "%"
		query = query.Where(
			"designer_profiles.company_name LIKE ? OR designer_profiles.address LIKE ?",
			searchQuery, searchQuery,
		)
	}

	// 地区筛选
	if req.Region != "" {
		regionQuery := "%" + req.Region + "%"
		query = query.Where("designer_profiles.address LIKE ?", regionQuery)
	}

	// 专业领域筛选
	if len(req.Specialties) > 0 {
		query = query.Joins("JOIN designer_specialties ON designer_profiles.id = designer_specialties.designer_id").
			Where("designer_specialties.specialty IN ?", req.Specialties)
	}

	// 评分筛选 - 使用子查询获取平均评分
	if req.MinRating > 0 {
		query = query.Where("(SELECT COALESCE(AVG(rating), 0) FROM designer_ratings WHERE designer_ratings.designer_id = designer_profiles.id) >= ?", req.MinRating)
	}
	if req.MaxRating > 0 && req.MaxRating <= 5.0 {
		query = query.Where("(SELECT COALESCE(AVG(rating), 0) FROM designer_ratings WHERE designer_ratings.designer_id = designer_profiles.id) <= ?", req.MaxRating)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("获取总数失败: %v", err)
	}

	// 排序
	sortField := s.getSortField(req.SortBy)
	sortOrder := "DESC"
	if req.SortOrder == "asc" {
		sortOrder = "ASC"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortField, sortOrder))

	// 分页
	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	// 执行查询
	var designerProfiles []models.DesignerProfile
	if err := query.Preload("User").Find(&designerProfiles).Error; err != nil {
		return nil, fmt.Errorf("查询设计师失败: %v", err)
	}

	// 转换为搜索结果
	designers := make([]models.DesignerSearchResult, 0, len(designerProfiles))
	for _, profile := range designerProfiles {
		designer := models.DesignerSearchResult{
			ID:          profile.ID,
			Name:        profile.CompanyName,
			Address:     profile.Address,
			Specialties: s.getDesignerSpecialties(profile.ID),
			Rating:      s.getDesignerRating(profile.ID),
			Description: profile.Bio,
			ContactInfo: models.ContactInfo{
				Phone: profile.User.Email,
				Email: profile.User.Email,
			},
			CreatedAt: profile.CreatedAt,
			UpdatedAt: profile.UpdatedAt,
		}
		designers = append(designers, designer)
	}

	response := &models.DesignerSearchResponse{
		Success: true,
		Data: models.DesignerSearchResultData{
			Designers: designers,
			Total:     total,
			Page:      req.Page,
			PageSize:  req.PageSize,
		},
	}

	return response, nil
}

// GetSearchSuggestions 获取搜索建议
func (s *DesignerSearchService) GetSearchSuggestions(req *models.DesignerSearchSuggestionRequest) (*models.DesignerSearchSuggestionResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 20 {
		req.Limit = 20
	}

	var suggestions []models.DesignerSearchSuggestion

	// 搜索设计师名称建议
	nameSuggestions, err := s.getNameSuggestions(req.Query, req.Limit/3)
	if err == nil {
		suggestions = append(suggestions, nameSuggestions...)
	}

	// 搜索地址建议
	addressSuggestions, err := s.getAddressSuggestions(req.Query, req.Limit/3)
	if err == nil {
		suggestions = append(suggestions, addressSuggestions...)
	}

	// 搜索专业领域建议
	specialtySuggestions, err := s.getSpecialtySuggestions(req.Query, req.Limit/3)
	if err == nil {
		suggestions = append(suggestions, specialtySuggestions...)
	}

	// 限制总数
	if len(suggestions) > req.Limit {
		suggestions = suggestions[:req.Limit]
	}

	response := &models.DesignerSearchSuggestionResponse{
		Success: true,
		Data: models.DesignerSearchSuggestionData{
			Suggestions: suggestions,
		},
	}

	return response, nil
}

// getNameSuggestions 获取设计师名称建议
func (s *DesignerSearchService) getNameSuggestions(query string, limit int) ([]models.DesignerSearchSuggestion, error) {
	var suggestions []models.DesignerSearchSuggestion
	
	err := s.db.Model(&models.DesignerProfile{}).
		Joins("JOIN users ON designer_profiles.user_id = users.id").
		Where("users.role = ? AND users.deleted_at IS NULL", "designer").
		Where("designer_profiles.company_name LIKE ?", "%"+query+"%").
		Select("DISTINCT designer_profiles.company_name").
		Limit(limit).
		Scan(&suggestions).Error

	if err != nil {
		return nil, err
	}

	// 添加类型和高亮
	for i := range suggestions {
		suggestions[i].Type = "designer_name"
		suggestions[i].Text = suggestions[i].Text
		suggestions[i].Highlight = strings.Replace(suggestions[i].Text, query, "<em>"+query+"</em>", -1)
	}

	return suggestions, nil
}

// getAddressSuggestions 获取地址建议
func (s *DesignerSearchService) getAddressSuggestions(query string, limit int) ([]models.DesignerSearchSuggestion, error) {
	var suggestions []models.DesignerSearchSuggestion
	
	err := s.db.Model(&models.DesignerProfile{}).
		Joins("JOIN users ON designer_profiles.user_id = users.id").
		Where("users.role = ? AND users.deleted_at IS NULL", "designer").
		Where("designer_profiles.address LIKE ?", "%"+query+"%").
		Select("DISTINCT designer_profiles.address").
		Limit(limit).
		Scan(&suggestions).Error

	if err != nil {
		return nil, err
	}

	// 添加类型和高亮
	for i := range suggestions {
		suggestions[i].Type = "designer_address"
		suggestions[i].Text = suggestions[i].Text
		suggestions[i].Highlight = strings.Replace(suggestions[i].Text, query, "<em>"+query+"</em>", -1)
	}

	return suggestions, nil
}

// getSpecialtySuggestions 获取专业领域建议
func (s *DesignerSearchService) getSpecialtySuggestions(query string, limit int) ([]models.DesignerSearchSuggestion, error) {
	var suggestions []models.DesignerSearchSuggestion
	
	err := s.db.Model(&models.DesignerSpecialty{}).
		Where("specialty LIKE ?", "%"+query+"%").
		Select("DISTINCT specialty").
		Limit(limit).
		Scan(&suggestions).Error

	if err != nil {
		return nil, err
	}

	// 添加类型和高亮
	for i := range suggestions {
		suggestions[i].Type = "specialty"
		suggestions[i].Text = suggestions[i].Text
		suggestions[i].Highlight = strings.Replace(suggestions[i].Text, query, "<em>"+query+"</em>", -1)
	}

	return suggestions, nil
}

// getDesignerSpecialties 获取设计师专业领域
func (s *DesignerSearchService) getDesignerSpecialties(designerID uint) []string {
	var specialties []models.DesignerSpecialty
	err := s.db.Where("designer_id = ?", designerID).Find(&specialties).Error
	if err != nil {
		return []string{}
	}

	result := make([]string, 0, len(specialties))
	for _, specialty := range specialties {
		result = append(result, specialty.Specialty)
	}
	return result
}

// getDesignerRating 获取设计师评分
func (s *DesignerSearchService) getDesignerRating(designerID uint) float64 {
	var avgRating float64
	err := s.db.Model(&models.DesignerRating{}).
		Where("designer_id = ?", designerID).
		Select("COALESCE(AVG(rating), 0)").
		Scan(&avgRating).Error
	
	if err != nil {
		return 0.0
	}
	return avgRating
}

// getSortField 获取排序字段
func (s *DesignerSearchService) getSortField(sortBy string) string {
	switch sortBy {
	case "name":
		return "designer_profiles.company_name"
	case "created_at":
		return "designer_profiles.created_at"
	case "rating":
		return "(SELECT COALESCE(AVG(rating), 0) FROM designer_ratings WHERE designer_ratings.designer_id = designer_profiles.id)"
	default:
		return "(SELECT COALESCE(AVG(rating), 0) FROM designer_ratings WHERE designer_ratings.designer_id = designer_profiles.id)"
	}
}

// CreateDesignerSpecialty 创建设计师专业领域
func (s *DesignerSearchService) CreateDesignerSpecialty(designerID uint, specialty string) error {
	specialtyRecord := models.DesignerSpecialty{
		DesignerID: designerID,
		Specialty:  specialty,
	}
	return s.db.Create(&specialtyRecord).Error
}

// CreateDesignerRating 创建设计师评分
func (s *DesignerSearchService) CreateDesignerRating(designerID uint, rating float64, comment string, raterID string) error {
	ratingRecord := models.DesignerRating{
		DesignerID: designerID,
		Rating:     rating,
		Comment:    comment,
		RaterID:    raterID,
	}
	return s.db.Create(&ratingRecord).Error
} 