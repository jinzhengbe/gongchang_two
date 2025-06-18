package services

import (
	"backend/models"
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
	err := s.db.First(&product, id).Error
	return &product, err
}

func (s *ProductService) UpdateProduct(id uint, req *models.ProductUpdateRequest) error {
	return s.db.Model(&models.Product{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"category":    req.Category,
		"price":       req.Price,
		"stock":       req.Stock,
	}).Error
}

func (s *ProductService) DeleteProduct(id uint) error {
	return s.db.Delete(&models.Product{}, id).Error
}

func (s *ProductService) GetProducts(page, pageSize int, category string) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := s.db.Model(&models.Product{})
	if category != "" {
		query = query.Where("category = ?", category)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&products).Error
	return products, total, err
}

func (s *ProductService) SearchProducts(query string, page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	err := s.db.Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.db.Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&products).Error
	return products, total, err
}

func (s *ProductService) GetProductsByCategory(category string, page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	err := s.db.Where("category = ?", category).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.db.Where("category = ?", category).
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&products).Error
	return products, total, err
}

func (s *ProductService) GetLatestProducts(limit int) ([]models.Product, error) {
	var products []models.Product
	err := s.db.Order("id desc").Limit(limit).Find(&products).Error
	return products, err
}

func (s *ProductService) GetHotProducts(limit int) ([]models.Product, error) {
	var products []models.Product
	err := s.db.Order("views desc").Limit(limit).Find(&products).Error
	return products, err
} 