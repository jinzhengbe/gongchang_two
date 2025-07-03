package services

import (
	"fmt"
	"gongChang/models"
	"gongChang/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

type OrderSearchService struct {
	db *gorm.DB
}

func NewOrderSearchService(db *gorm.DB) *OrderSearchService {
	return &OrderSearchService{db: db}
}

// SearchOrders 高级订单搜索
func (s *OrderSearchService) SearchOrders(req *models.OrderSearchRequest) (*models.OrderSearchResponse, error) {
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
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	// 构建基础查询
	query := s.db.Model(&models.Order{}).Preload("Factory")

	// 根据用户角色和ID添加权限过滤
	query = s.addPermissionFilter(query, req.UserID, req.UserRole)

	// 添加搜索条件
	query = s.addSearchConditions(query, req)

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("获取总数失败: %w", err)
	}

	// 添加排序
	query = s.addSorting(query, req.SortBy, req.SortOrder)

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	var orders []models.Order
	if err := query.Offset(offset).Limit(req.PageSize).Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("查询订单失败: %w", err)
	}

	// 转换为响应格式
	orderItems := s.convertToSearchItems(orders)

	response := &models.OrderSearchResponse{
		Success: true,
	}
	response.Data.Orders = orderItems
	response.Data.Total = total
	response.Data.Page = req.Page
	response.Data.PageSize = req.PageSize

	return response, nil
}

// GetSearchSuggestions 获取搜索建议
func (s *OrderSearchService) GetSearchSuggestions(req *models.SearchSuggestionRequest) (*models.SearchSuggestionResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 20 {
		req.Limit = 20
	}

	var suggestions []models.SearchSuggestion

	// 搜索订单标题建议
	titleSuggestions, err := s.getTitleSuggestions(req.Query, req.Limit/3)
	if err == nil {
		suggestions = append(suggestions, titleSuggestions...)
	}

	// 搜索面料名称建议
	fabricSuggestions, err := s.getFabricSuggestions(req.Query, req.Limit/3)
	if err == nil {
		suggestions = append(suggestions, fabricSuggestions...)
	}

	// 搜索工厂名称建议
	factorySuggestions, err := s.getFactorySuggestions(req.Query, req.Limit/3)
	if err == nil {
		suggestions = append(suggestions, factorySuggestions...)
	}

	// 限制总数
	if len(suggestions) > req.Limit {
		suggestions = suggestions[:req.Limit]
	}

	response := &models.SearchSuggestionResponse{
		Success: true,
	}
	response.Data.Suggestions = suggestions

	return response, nil
}

// addPermissionFilter 添加权限过滤
func (s *OrderSearchService) addPermissionFilter(query *gorm.DB, userID, userRole string) *gorm.DB {
	switch userRole {
	case "factory":
		// 工厂用户可以看到分配给自己的订单和未分配的订单
		return query.Where("factory_id = ? OR factory_id IS NULL", userID)
	case "designer":
		// 设计师只能看到自己创建的订单
		return query.Where("designer_id = ?", userID)
	case "admin":
		// 管理员可以看到所有订单
		return query
	default:
		// 默认只能看到公开的订单
		return query.Where("status = ?", "published")
	}
}

// addSearchConditions 添加搜索条件
func (s *OrderSearchService) addSearchConditions(query *gorm.DB, req *models.OrderSearchRequest) *gorm.DB {
	// 关键词搜索
	if req.Query != "" {
		searchTerm := "%" + req.Query + "%"
		query = query.Where(
			"title LIKE ? OR description LIKE ? OR fabric LIKE ? OR fabrics LIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	// 状态筛选
	if req.Status != "" && req.Status != "all" {
		query = query.Where("status = ?", req.Status)
	}

	// 时间范围筛选
	if req.StartDate != "" {
		startTime, err := time.Parse("2006-01-02", req.StartDate)
		if err == nil {
			query = query.Where("created_at >= ?", startTime)
		}
	}

	if req.EndDate != "" {
		endTime, err := time.Parse("2006-01-02", req.EndDate)
		if err == nil {
			// 结束日期包含当天
			endTime = endTime.Add(24 * time.Hour)
			query = query.Where("created_at < ?", endTime)
		}
	}

	return query
}

// addSorting 添加排序
func (s *OrderSearchService) addSorting(query *gorm.DB, sortBy, sortOrder string) *gorm.DB {
	// 验证排序字段
	validSortFields := map[string]string{
		"created_at": "created_at",
		"updated_at": "updated_at",
		"id":         "id",
		"title":      "title",
		"status":     "status",
	}

	if field, exists := validSortFields[sortBy]; exists {
		if sortOrder == "asc" {
			query = query.Order(field + " ASC")
		} else {
			query = query.Order(field + " DESC")
		}
	} else {
		// 默认按创建时间倒序
		query = query.Order("created_at DESC")
	}

	return query
}

// convertToSearchItems 转换为搜索项
func (s *OrderSearchService) convertToSearchItems(orders []models.Order) []models.OrderSearchItem {
	items := make([]models.OrderSearchItem, len(orders))
	for i, order := range orders {
		// 解析面料信息
		var fabrics []string
		if order.Fabrics != "" {
			fabrics = strings.Split(order.Fabrics, ",")
			// 清理空白字符
			for j, fabric := range fabrics {
				fabrics[j] = strings.TrimSpace(fabric)
			}
		}

		// 构建工厂信息
		factoryInfo := models.FactoryInfo{}
		if order.Factory.ID != 0 {
			if order.FactoryID != nil {
				factoryInfo.ID = *order.FactoryID
			}
			factoryInfo.Name = order.Factory.CompanyName
		}

		// 处理时间字段
		var createdAt, updatedAt time.Time
		if order.CreatedAt != nil {
			createdAt = *order.CreatedAt
		}
		if order.UpdatedAt != nil {
			updatedAt = *order.UpdatedAt
		}

		items[i] = models.OrderSearchItem{
			ID:        order.ID,
			Title:     order.Title,
			OrderNo:   fmt.Sprintf("ORD%06d", order.ID), // 生成订单号
			Status:    string(order.Status),
			Fabrics:   fabrics,
			Factory:   factoryInfo,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
	}
	return items
}

// getTitleSuggestions 获取标题建议
func (s *OrderSearchService) getTitleSuggestions(query string, limit int) ([]models.SearchSuggestion, error) {
	var titles []string
	err := s.db.Model(&models.Order{}).
		Where("title LIKE ?", "%"+query+"%").
		Distinct().
		Pluck("title", &titles).
		Limit(limit).
		Error

	if err != nil {
		return nil, err
	}

	suggestions := make([]models.SearchSuggestion, len(titles))
	for i, title := range titles {
		highlight := utils.HighlightText(title, query)
		suggestions[i] = models.SearchSuggestion{
			Type:      "order_title",
			Text:      title,
			Highlight: highlight,
		}
	}

	return suggestions, nil
}

// getFabricSuggestions 获取面料建议
func (s *OrderSearchService) getFabricSuggestions(query string, limit int) ([]models.SearchSuggestion, error) {
	var fabrics []string
	err := s.db.Model(&models.Order{}).
		Where("fabric LIKE ? OR fabrics LIKE ?", "%"+query+"%", "%"+query+"%").
		Distinct().
		Pluck("fabric", &fabrics).
		Limit(limit).
		Error

	if err != nil {
		return nil, err
	}

	suggestions := make([]models.SearchSuggestion, len(fabrics))
	for i, fabric := range fabrics {
		if fabric != "" {
			highlight := utils.HighlightText(fabric, query)
			suggestions[i] = models.SearchSuggestion{
				Type:      "fabric_name",
				Text:      fabric,
				Highlight: highlight,
			}
		}
	}

	return suggestions, nil
}

// getFactorySuggestions 获取工厂建议
func (s *OrderSearchService) getFactorySuggestions(query string, limit int) ([]models.SearchSuggestion, error) {
	var factories []models.FactoryProfile
	err := s.db.Model(&models.FactoryProfile{}).
		Where("company_name LIKE ?", "%"+query+"%").
		Limit(limit).
		Find(&factories).Error

	if err != nil {
		return nil, err
	}

	suggestions := make([]models.SearchSuggestion, len(factories))
	for i, factory := range factories {
		highlight := utils.HighlightText(factory.CompanyName, query)
		suggestions[i] = models.SearchSuggestion{
			Type:      "factory_name",
			Text:      factory.CompanyName,
			Highlight: highlight,
		}
	}

	return suggestions, nil
} 