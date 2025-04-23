package services

import (
	"aneworder.com/backend/models"
	"gorm.io/gorm"
)

type OrderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{
		db: db,
	}
}

func (s *OrderService) CreateOrder(order *models.Order) error {
	return s.db.Create(order).Error
}

func (s *OrderService) GetOrdersByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := s.db.Preload("Designer").Preload("Customer").Preload("Product").
		Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

func (s *OrderService) GetOrderByID(orderID uint) (*models.Order, error) {
	var order models.Order
	err := s.db.Preload("Designer").Preload("Customer").Preload("Product").
		First(&order, orderID).Error
	return &order, err
}

func (s *OrderService) UpdateOrderStatus(orderID uint, status string) error {
	return s.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func (s *OrderService) SearchOrders(query string, userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := s.db.Where("user_id = ? AND (description LIKE ? OR status LIKE ?)", 
		userID, "%"+query+"%", "%"+query+"%").
		Find(&orders).Error
	return orders, err
}

func (s *OrderService) GetOrderStatistics(userID uint) (*models.OrderStatistics, error) {
	var stats models.OrderStatistics
	
	// 获取总订单数
	err := s.db.Model(&models.Order{}).Where("user_id = ?", userID).Count(&stats.TotalOrders).Error
	if err != nil {
		return nil, err
	}
	
	// 获取待处理订单数
	err = s.db.Model(&models.Order{}).Where("user_id = ? AND status = ?", userID, "pending").Count(&stats.PendingOrders).Error
	if err != nil {
		return nil, err
	}
	
	// 获取已完成订单数
	err = s.db.Model(&models.Order{}).Where("user_id = ? AND status = ?", userID, "completed").Count(&stats.CompletedOrders).Error
	if err != nil {
		return nil, err
	}
	
	return &stats, nil
}

func (s *OrderService) GetRecentOrders(limit int) ([]models.Order, error) {
	var orders []models.Order
	err := s.db.Order("created_at desc").Limit(limit).Find(&orders).Error
	return orders, err
}

func (s *OrderService) GetLatestOrders(limit int) ([]models.Order, error) {
	var orders []models.Order
	err := s.db.Order("created_at desc").Limit(limit).Find(&orders).Error
	return orders, err
}

func (s *OrderService) GetHotOrders(limit int) ([]models.Order, error) {
	var orders []models.Order
	err := s.db.Order("views desc").Limit(limit).Find(&orders).Error
	return orders, err
} 