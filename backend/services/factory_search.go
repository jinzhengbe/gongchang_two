package services

import (
	"fmt"
	"strings"
	"gongChang/models"
	"gorm.io/gorm"
)

type FactorySearchService struct {
	db *gorm.DB
}

func NewFactorySearchService(db *gorm.DB) *FactorySearchService {
	return &FactorySearchService{db: db}
}

// SearchFactories 搜索工厂
func (s *FactorySearchService) SearchFactories(req *models.FactorySearchRequest) (*models.FactorySearchResponse, error) {
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
	query := s.db.Model(&models.FactoryProfile{}).
		Joins("JOIN users ON factory_profiles.user_id = users.id").
		Where("users.role = ? AND users.deleted_at IS NULL", "factory")

	// 添加搜索条件
	if req.Query != "" {
		searchQuery := "%" + req.Query + "%"
		query = query.Where(
			"factory_profiles.company_name LIKE ? OR factory_profiles.address LIKE ?",
			searchQuery, searchQuery,
		)
	}

	// 地区筛选
	if req.Region != "" {
		regionQuery := "%" + req.Region + "%"
		query = query.Where("factory_profiles.address LIKE ?", regionQuery)
	}

	// 专业领域筛选
	if len(req.Specialties) > 0 {
		query = query.Joins("JOIN factory_specialties ON factory_profiles.id = factory_specialties.factory_id").
			Where("factory_specialties.specialty IN ?", req.Specialties)
	}

	// 合作状态筛选
	if req.CooperationStatus != "" && req.CooperationStatus != "all" {
		// 这里可以根据实际业务逻辑实现合作状态筛选
		// 暂时使用工厂状态作为合作状态
		if req.CooperationStatus == "cooperating" {
			query = query.Where("factory_profiles.status = ?", 1)
		} else if req.CooperationStatus == "not_cooperating" {
			query = query.Where("factory_profiles.status = ?", 0)
		}
	}

	// 评分筛选 - 使用子查询获取平均评分
	if req.MinRating > 0 {
		query = query.Where("(SELECT COALESCE(AVG(rating), 0) FROM factory_ratings WHERE factory_ratings.factory_id = factory_profiles.id) >= ?", req.MinRating)
	}
	if req.MaxRating > 0 && req.MaxRating <= 5.0 {
		query = query.Where("(SELECT COALESCE(AVG(rating), 0) FROM factory_ratings WHERE factory_ratings.factory_id = factory_profiles.id) <= ?", req.MaxRating)
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
	var factoryProfiles []models.FactoryProfile
	if err := query.Preload("User").Find(&factoryProfiles).Error; err != nil {
		return nil, fmt.Errorf("查询工厂失败: %v", err)
	}

	// 转换为搜索结果
	factories := make([]models.FactorySearchResult, 0, len(factoryProfiles))
	for _, profile := range factoryProfiles {
		// 获取专业领域
		specialties, _ := s.getFactorySpecialties(profile.ID)
		
		// 获取评分
		rating := s.getFactoryRating(profile.ID)
		
		// 获取联系信息
		contactInfo := models.ContactInfo{
			Phone: profile.User.Email, // 暂时使用邮箱作为联系方式
			Email: profile.User.Email,
		}
		
		// 获取产能信息
		capacity := models.Capacity{
			MonthlyOrders: profile.Capacity,
			MaxOrderSize:  profile.Capacity * 10, // 简单估算
		}

		factory := models.FactorySearchResult{
			ID:                profile.ID,
			Name:              profile.CompanyName,
			Address:           profile.Address,
			Specialties:       specialties,
			Rating:            rating,
			CooperationStatus: s.getCooperationStatus(profile.Status),
			Description:       profile.Equipment, // 使用设备信息作为描述
			ContactInfo:       contactInfo,
			Capacity:          capacity,
			CreatedAt:         profile.CreatedAt,
			UpdatedAt:         profile.UpdatedAt,
		}
		factories = append(factories, factory)
	}

	return &models.FactorySearchResponse{
		Success: true,
		Data: models.FactorySearchResultData{
			Factories: factories,
			Total:     total,
			Page:      req.Page,
			PageSize:  req.PageSize,
		},
	}, nil
}

// GetSearchSuggestions 获取搜索建议
func (s *FactorySearchService) GetSearchSuggestions(req *models.FactorySearchSuggestionRequest) (*models.FactorySearchSuggestionResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 20 {
		req.Limit = 20
	}

	suggestions := make([]models.FactorySearchSuggestion, 0)

	// 工厂名称建议
	factoryNameSuggestions, err := s.getFactoryNameSuggestions(req.Query, req.Limit/3)
	if err == nil {
		suggestions = append(suggestions, factoryNameSuggestions...)
	}

	// 地址建议
	addressSuggestions, err := s.getAddressSuggestions(req.Query, req.Limit/3)
	if err == nil {
		suggestions = append(suggestions, addressSuggestions...)
	}

	// 专业领域建议
	specialtySuggestions, err := s.getSpecialtySuggestions(req.Query, req.Limit/3)
	if err == nil {
		suggestions = append(suggestions, specialtySuggestions...)
	}

	// 限制返回数量
	if len(suggestions) > req.Limit {
		suggestions = suggestions[:req.Limit]
	}

	return &models.FactorySearchSuggestionResponse{
		Success: true,
		Data: models.FactorySearchSuggestionData{
			Suggestions: suggestions,
		},
	}, nil
}

// getFactoryNameSuggestions 获取工厂名称建议
func (s *FactorySearchService) getFactoryNameSuggestions(query string, limit int) ([]models.FactorySearchSuggestion, error) {
	var names []string
	err := s.db.Model(&models.FactoryProfile{}).
		Where("company_name LIKE ?", "%"+query+"%").
		Limit(limit).
		Pluck("DISTINCT company_name", &names).Error
	
	if err != nil {
		return nil, err
	}

	suggestions := make([]models.FactorySearchSuggestion, 0, len(names))
	for _, name := range names {
		highlight := strings.ReplaceAll(name, query, "<em>"+query+"</em>")
		suggestions = append(suggestions, models.FactorySearchSuggestion{
			Type:      "factory_name",
			Text:      name,
			Highlight: highlight,
		})
	}
	return suggestions, nil
}

// getAddressSuggestions 获取地址建议
func (s *FactorySearchService) getAddressSuggestions(query string, limit int) ([]models.FactorySearchSuggestion, error) {
	var addresses []string
	err := s.db.Model(&models.FactoryProfile{}).
		Where("address LIKE ?", "%"+query+"%").
		Limit(limit).
		Pluck("DISTINCT address", &addresses).Error
	
	if err != nil {
		return nil, err
	}

	suggestions := make([]models.FactorySearchSuggestion, 0, len(addresses))
	for _, address := range addresses {
		highlight := strings.ReplaceAll(address, query, "<em>"+query+"</em>")
		suggestions = append(suggestions, models.FactorySearchSuggestion{
			Type:      "factory_address",
			Text:      address,
			Highlight: highlight,
		})
	}
	return suggestions, nil
}

// getSpecialtySuggestions 获取专业领域建议
func (s *FactorySearchService) getSpecialtySuggestions(query string, limit int) ([]models.FactorySearchSuggestion, error) {
	var specialties []string
	err := s.db.Model(&models.FactorySpecialty{}).
		Where("specialty LIKE ?", "%"+query+"%").
		Limit(limit).
		Pluck("DISTINCT specialty", &specialties).Error
	
	if err != nil {
		return nil, err
	}

	suggestions := make([]models.FactorySearchSuggestion, 0, len(specialties))
	for _, specialty := range specialties {
		highlight := strings.ReplaceAll(specialty, query, "<em>"+query+"</em>")
		suggestions = append(suggestions, models.FactorySearchSuggestion{
			Type:      "specialty",
			Text:      specialty,
			Highlight: highlight,
		})
	}
	return suggestions, nil
}

// getFactorySpecialties 获取工厂专业领域
func (s *FactorySearchService) getFactorySpecialties(factoryID uint) ([]string, error) {
	var specialties []string
	err := s.db.Model(&models.FactorySpecialty{}).
		Where("factory_id = ?", factoryID).
		Pluck("specialty", &specialties).Error
	return specialties, err
}

// getFactoryRating 获取工厂评分
func (s *FactorySearchService) getFactoryRating(factoryID uint) float64 {
	var avgRating float64
	err := s.db.Model(&models.FactoryRating{}).
		Where("factory_id = ?", factoryID).
		Select("COALESCE(AVG(rating), 0)").
		Scan(&avgRating).Error
	
	if err != nil {
		return 0.0
	}
	return avgRating
}

// getCooperationStatus 获取合作状态
func (s *FactorySearchService) getCooperationStatus(status int) string {
	if status == 1 {
		return "合作中"
	}
	return "未合作"
}

// getSortField 获取排序字段
func (s *FactorySearchService) getSortField(sortBy string) string {
	switch sortBy {
	case "name":
		return "factory_profiles.company_name"
	case "created_at":
		return "factory_profiles.created_at"
	case "rating":
		return "(SELECT COALESCE(AVG(rating), 0) FROM factory_ratings WHERE factory_ratings.factory_id = factory_profiles.id)"
	default:
		return "(SELECT COALESCE(AVG(rating), 0) FROM factory_ratings WHERE factory_ratings.factory_id = factory_profiles.id)"
	}
}

// CreateFactorySpecialty 创建工厂专业领域
func (s *FactorySearchService) CreateFactorySpecialty(factoryID uint, specialty string) error {
	specialtyRecord := models.FactorySpecialty{
		FactoryID: factoryID,
		Specialty:  specialty,
	}
	return s.db.Create(&specialtyRecord).Error
}

// CreateFactoryRating 创建工厂评分
func (s *FactorySearchService) CreateFactoryRating(factoryID uint, rating float64, comment string, raterID string) error {
	ratingRecord := models.FactoryRating{
		FactoryID: factoryID,
		Rating:    rating,
		Comment:   comment,
		RaterID:   raterID,
	}
	return s.db.Create(&ratingRecord).Error
}

// GetFactoryRatings 获取工厂评分列表
func (s *FactorySearchService) GetFactoryRatings(factoryID uint, page, pageSize int) ([]map[string]interface{}, int64, error) {
	var ratings []models.FactoryRating
	var total int64

	// 获取总数
	if err := s.db.Model(&models.FactoryRating{}).Where("factory_id = ?", factoryID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := s.db.Where("factory_id = ?", factoryID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&ratings).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	result := make([]map[string]interface{}, 0, len(ratings))
	for _, rating := range ratings {
		// 获取评分者信息
		var user models.User
		s.db.Where("id = ?", rating.RaterID).First(&user)

		result = append(result, map[string]interface{}{
			"id":         rating.ID,
			"rating":     rating.Rating,
			"comment":    rating.Comment,
			"rater": map[string]interface{}{
				"id":       user.ID,
				"username": user.Username,
				"role":     user.Role,
			},
			"created_at": rating.CreatedAt,
		})
	}

	return result, total, nil
}

// GetFactoryRatingStats 获取工厂评分统计
func (s *FactorySearchService) GetFactoryRatingStats(factoryID uint) (map[string]interface{}, error) {
	var stats struct {
		TotalRatings   int64   `json:"total_ratings"`
		AverageRating  float64 `json:"average_rating"`
		MaxRating      float64 `json:"max_rating"`
		MinRating      float64 `json:"min_rating"`
		RatingCounts   map[int]int `json:"rating_counts"`
	}

	// 基础统计
	if err := s.db.Model(&models.FactoryRating{}).
		Where("factory_id = ?", factoryID).
		Select("COUNT(*) as total_ratings, AVG(rating) as average_rating, MAX(rating) as max_rating, MIN(rating) as min_rating").
		Scan(&stats).Error; err != nil {
		return nil, err
	}

	// 评分分布统计
	stats.RatingCounts = make(map[int]int)
	for i := 1; i <= 5; i++ {
		var count int64
		s.db.Model(&models.FactoryRating{}).
			Where("factory_id = ? AND rating = ?", factoryID, float64(i)).
			Count(&count)
		stats.RatingCounts[i] = int(count)
	}

	// 计算评分等级
	var ratingLevel string
	switch {
	case stats.AverageRating >= 4.5:
		ratingLevel = "优秀"
	case stats.AverageRating >= 4.0:
		ratingLevel = "良好"
	case stats.AverageRating >= 3.0:
		ratingLevel = "一般"
	default:
		ratingLevel = "较差"
	}

	return map[string]interface{}{
		"total_ratings":  stats.TotalRatings,
		"average_rating": stats.AverageRating,
		"max_rating":     stats.MaxRating,
		"min_rating":     stats.MinRating,
		"rating_level":   ratingLevel,
		"rating_counts":  stats.RatingCounts,
	}, nil
} 