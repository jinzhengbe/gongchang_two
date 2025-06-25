package services

import (
	"gongChang/models"
	"gorm.io/gorm"
	"errors"
	"log"
)

type FabricService struct {
	db *gorm.DB
}

func NewFabricService(db *gorm.DB) *FabricService {
	return &FabricService{db: db}
}

// CreateFabric 创建布料
func (s *FabricService) CreateFabric(req *models.FabricRequest) (*models.Fabric, error) {
	log.Printf("CreateFabric service called with req.DesignerID=%s, req.SupplierID=%s, req.FactoryID=%s", req.DesignerID, req.SupplierID, req.FactoryID)
	
	fabric := &models.Fabric{
		Name:         req.Name,
		Category:     req.Category,
		Material:     req.Material,
		Color:        req.Color,
		Pattern:      req.Pattern,
		Weight:       req.Weight,
		Width:        req.Width,
		Price:        req.Price,
		Unit:         req.Unit,
		Stock:        req.Stock,
		MinOrder:     req.MinOrder,
		Description:  req.Description,
		ImageURL:     req.ImageURL,
		ThumbnailURL: req.ThumbnailURL,
		Tags:         req.Tags,
		Status:       1, // 默认启用
	}

	// 设置设计师ID
	if req.DesignerID != "" {
		fabric.DesignerID = &req.DesignerID
		log.Printf("CreateFabric service: setting fabric.DesignerID=%s", *fabric.DesignerID)
	}

	// 设置供应商ID
	if req.SupplierID != "" {
		fabric.SupplierID = &req.SupplierID
		log.Printf("CreateFabric service: setting fabric.SupplierID=%s", *fabric.SupplierID)
	}

	// 设置工厂ID
	if req.FactoryID != "" {
		fabric.FactoryID = &req.FactoryID
		log.Printf("CreateFabric service: setting fabric.FactoryID=%s", *fabric.FactoryID)
	}

	log.Printf("CreateFabric service: final fabric.DesignerID=%v, fabric.SupplierID=%v, fabric.FactoryID=%v", fabric.DesignerID, fabric.SupplierID, fabric.FactoryID)

	if err := s.db.Create(fabric).Error; err != nil {
		log.Printf("CreateFabric service: database error: %v", err)
		return nil, err
	}

	log.Printf("CreateFabric service: fabric created successfully with ID=%d", fabric.ID)
	return fabric, nil
}

// GetFabricByID 根据ID获取布料
func (s *FabricService) GetFabricByID(id uint) (*models.Fabric, error) {
	var fabric models.Fabric
	if err := s.db.First(&fabric, id).Error; err != nil {
		return nil, err
	}
	return &fabric, nil
}

// UpdateFabric 更新布料
func (s *FabricService) UpdateFabric(id uint, req *models.FabricUpdateRequest) (*models.Fabric, error) {
	fabric, err := s.GetFabricByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Name != "" {
		fabric.Name = req.Name
	}
	if req.Category != "" {
		fabric.Category = req.Category
	}
	if req.Material != "" {
		fabric.Material = req.Material
	}
	if req.Color != "" {
		fabric.Color = req.Color
	}
	if req.Pattern != "" {
		fabric.Pattern = req.Pattern
	}
	if req.Weight > 0 {
		fabric.Weight = req.Weight
	}
	if req.Width > 0 {
		fabric.Width = req.Width
	}
	if req.Price >= 0 {
		fabric.Price = req.Price
	}
	if req.Unit != "" {
		fabric.Unit = req.Unit
	}
	if req.Stock >= 0 {
		fabric.Stock = req.Stock
	}
	if req.MinOrder > 0 {
		fabric.MinOrder = req.MinOrder
	}
	if req.Description != "" {
		fabric.Description = req.Description
	}
	if req.ImageURL != "" {
		fabric.ImageURL = req.ImageURL
	}
	if req.ThumbnailURL != "" {
		fabric.ThumbnailURL = req.ThumbnailURL
	}
	if req.Tags != "" {
		fabric.Tags = req.Tags
	}
	if req.Status >= 0 {
		fabric.Status = req.Status
	}
	if req.DesignerID != "" {
		fabric.DesignerID = &req.DesignerID
	}
	if req.SupplierID != "" {
		fabric.SupplierID = &req.SupplierID
	}
	if req.FactoryID != "" {
		fabric.FactoryID = &req.FactoryID
	}

	if err := s.db.Save(fabric).Error; err != nil {
		return nil, err
	}

	return fabric, nil
}

// DeleteFabric 删除布料
func (s *FabricService) DeleteFabric(id uint) error {
	return s.db.Delete(&models.Fabric{}, id).Error
}

// SearchFabrics 搜索布料
func (s *FabricService) SearchFabrics(req *models.FabricSearchRequest) (*models.FabricListResponse, error) {
	query := s.db.Model(&models.Fabric{})

	// 搜索关键词
	if req.Query != "" {
		searchTerm := "%" + req.Query + "%"
		query = query.Where(
			"name LIKE ? OR description LIKE ? OR material LIKE ? OR color LIKE ? OR pattern LIKE ? OR tags LIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	// 分类筛选
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}

	// 材质筛选
	if req.Material != "" {
		query = query.Where("material = ?", req.Material)
	}

	// 颜色筛选
	if req.Color != "" {
		query = query.Where("color = ?", req.Color)
	}

	// 价格范围筛选
	if req.MinPrice > 0 {
		query = query.Where("price >= ?", req.MinPrice)
	}
	if req.MaxPrice > 0 {
		query = query.Where("price <= ?", req.MaxPrice)
	}

	// 库存筛选
	if req.MinStock > 0 {
		query = query.Where("stock >= ?", req.MinStock)
	}

	// 状态筛选
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	// 按创建时间倒序排列
	query = query.Order("created_at DESC")

	var fabrics []models.Fabric
	if err := query.Find(&fabrics).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	fabricResponses := make([]models.FabricResponse, len(fabrics))
	for i, fabric := range fabrics {
		fabricResponses[i] = models.FabricResponse{
			ID:           fabric.ID,
			Name:         fabric.Name,
			Category:     fabric.Category,
			Material:     fabric.Material,
			Color:        fabric.Color,
			Pattern:      fabric.Pattern,
			Weight:       fabric.Weight,
			Width:        fabric.Width,
			Price:        fabric.Price,
			Unit:         fabric.Unit,
			Stock:        fabric.Stock,
			MinOrder:     fabric.MinOrder,
			Description:  fabric.Description,
			ImageURL:     fabric.ImageURL,
			ThumbnailURL: fabric.ThumbnailURL,
			Tags:         fabric.Tags,
			Status:       fabric.Status,
			DesignerID:   fabric.DesignerID,
			SupplierID:   fabric.SupplierID,
			FactoryID:    fabric.FactoryID,
			CreatedAt:    fabric.CreatedAt,
			UpdatedAt:    fabric.UpdatedAt,
		}
	}

	return &models.FabricListResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Fabrics:  fabricResponses,
	}, nil
}

// GetAllFabrics 获取所有布料（用于前端下拉选择）
func (s *FabricService) GetAllFabrics() ([]models.FabricResponse, error) {
	var fabrics []models.Fabric
	if err := s.db.Where("status = ?", 1).Order("name ASC").Find(&fabrics).Error; err != nil {
		return nil, err
	}

	fabricResponses := make([]models.FabricResponse, len(fabrics))
	for i, fabric := range fabrics {
		fabricResponses[i] = models.FabricResponse{
			ID:           fabric.ID,
			Name:         fabric.Name,
			Category:     fabric.Category,
			Material:     fabric.Material,
			Color:        fabric.Color,
			Pattern:      fabric.Pattern,
			Weight:       fabric.Weight,
			Width:        fabric.Width,
			Price:        fabric.Price,
			Unit:         fabric.Unit,
			Stock:        fabric.Stock,
			MinOrder:     fabric.MinOrder,
			Description:  fabric.Description,
			ImageURL:     fabric.ImageURL,
			ThumbnailURL: fabric.ThumbnailURL,
			Tags:         fabric.Tags,
			Status:       fabric.Status,
			DesignerID:   fabric.DesignerID,
			SupplierID:   fabric.SupplierID,
			FactoryID:    fabric.FactoryID,
			CreatedAt:    fabric.CreatedAt,
			UpdatedAt:    fabric.UpdatedAt,
		}
	}

	return fabricResponses, nil
}

// GetFabricCategories 获取布料分类
func (s *FabricService) GetFabricCategories() ([]models.FabricCategory, error) {
	var categories []models.FabricCategory
	if err := s.db.Where("status = ?", 1).Order("sort ASC, name ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetFabricsByCategory 根据分类获取布料
func (s *FabricService) GetFabricsByCategory(category string, page, pageSize int) (*models.FabricListResponse, error) {
	status := 1
	req := &models.FabricSearchRequest{
		Category: category,
		Page:     page,
		PageSize: pageSize,
		Status:   &status,
	}
	return s.SearchFabrics(req)
}

// GetFabricsByMaterial 根据材质获取布料
func (s *FabricService) GetFabricsByMaterial(material string, page, pageSize int) (*models.FabricListResponse, error) {
	status := 1
	req := &models.FabricSearchRequest{
		Material: material,
		Page:     page,
		PageSize: pageSize,
		Status:   &status,
	}
	return s.SearchFabrics(req)
}

// UpdateFabricStock 更新布料库存
func (s *FabricService) UpdateFabricStock(id uint, quantity int) error {
	fabric, err := s.GetFabricByID(id)
	if err != nil {
		return err
	}

	newStock := fabric.Stock + quantity
	if newStock < 0 {
		return errors.New("库存不足")
	}

	return s.db.Model(fabric).Update("stock", newStock).Error
}

// GetFabricStatistics 获取布料统计信息
func (s *FabricService) GetFabricStatistics() (map[string]interface{}, error) {
	var totalFabrics, availableFabrics, lowStockFabrics int64

	// 总布料数量
	if err := s.db.Model(&models.Fabric{}).Count(&totalFabrics).Error; err != nil {
		return nil, err
	}

	// 可用布料数量
	if err := s.db.Model(&models.Fabric{}).Where("status = ?", 1).Count(&availableFabrics).Error; err != nil {
		return nil, err
	}

	// 库存不足的布料数量（库存小于10）
	if err := s.db.Model(&models.Fabric{}).Where("stock < ?", 10).Count(&lowStockFabrics).Error; err != nil {
		return nil, err
	}

	// 按分类统计
	var categoryStats []struct {
		Category string `json:"category"`
		Count    int64  `json:"count"`
	}
	if err := s.db.Model(&models.Fabric{}).
		Select("category, count(*) as count").
		Group("category").
		Scan(&categoryStats).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_fabrics":     totalFabrics,
		"available_fabrics": availableFabrics,
		"low_stock_fabrics": lowStockFabrics,
		"category_stats":    categoryStats,
	}, nil
} 