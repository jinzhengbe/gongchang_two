package services

import (
	"backend/models"
	"errors"
	"gorm.io/gorm"
)

type ProductService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{
		db: db,
	}
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	return s.db.Create(product).Error
}

func (s *ProductService) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	err := s.db.Preload("Creator").First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

func (s *ProductService) UpdateProduct(id uint, updates *models.ProductUpdateRequest) error {
	product := &models.Product{}
	if err := s.db.First(product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	// 更新字段
	if updates.Name != "" {
		product.Name = updates.Name
	}
	if updates.Description != "" {
		product.Description = updates.Description
	}
	if updates.Category != "" {
		product.Category = updates.Category
	}
	if updates.Price >= 0 {
		product.Price = updates.Price
	}
	if updates.Stock >= 0 {
		product.Stock = updates.Stock
	}
	if updates.Status != "" {
		product.Status = updates.Status
	}

	return s.db.Save(product).Error
}

func (s *ProductService) DeleteProduct(id uint) error {
	return s.db.Delete(&models.Product{}, id).Error
}

func (s *ProductService) GetProducts(page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	// 获取总数
	if err := s.db.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err := s.db.Preload("Creator").
		Offset(offset).
		Limit(pageSize).
		Find(&products).Error

	return products, total, err
}

func (s *ProductService) SearchProducts(query string, page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	// 构建搜索条件
	searchQuery := "%" + query + "%"
	
	// 获取总数
	if err := s.db.Model(&models.Product{}).
		Where("name LIKE ? OR description LIKE ? OR category LIKE ?", 
			searchQuery, searchQuery, searchQuery).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err := s.db.Preload("Creator").
		Where("name LIKE ? OR description LIKE ? OR category LIKE ?", 
			searchQuery, searchQuery, searchQuery).
		Offset(offset).
		Limit(pageSize).
		Find(&products).Error

	return products, total, err
}

func (s *ProductService) GetProductsByCategory(category string, page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	// 获取总数
	if err := s.db.Model(&models.Product{}).
		Where("category = ?", category).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err := s.db.Preload("Creator").
		Where("category = ?", category).
		Offset(offset).
		Limit(pageSize).
		Find(&products).Error

	return products, total, err
} 