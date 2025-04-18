package services

import (
	"aneworder.com/backend/models"
	"errors"
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
	err := s.db.Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

func (s *OrderService) GetOrderByID(orderID uint) (*models.Order, error) {
	var order models.Order
	err := s.db.First(&order, orderID).Error
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

// UpdateOrder updates an existing order
func (s *OrderService) UpdateOrder(orderID uint, order *models.Order, userID uint) error {
	// Check if order exists and user has permission
	var existingOrder models.Order
	if err := s.db.First(&existingOrder, orderID).Error; err != nil {
		return err
	}

	// Check if user has permission to update the order
	if existingOrder.DesignerID != userID && existingOrder.CustomerID != userID {
		return errors.New("unauthorized to update this order")
	}

	// Update order
	return s.db.Model(&models.Order{}).Where("id = ?", orderID).Updates(order).Error
}

// DeleteOrder deletes an order
func (s *OrderService) DeleteOrder(orderID uint, userID uint) error {
	// Check if order exists and user has permission
	var order models.Order
	if err := s.db.First(&order, orderID).Error; err != nil {
		return err
	}

	// Check if user has permission to delete the order
	if order.DesignerID != userID && order.CustomerID != userID {
		return errors.New("unauthorized to delete this order")
	}

	// Delete order
	return s.db.Delete(&models.Order{}, orderID).Error
}

// GetOrderWithFiles gets an order with its associated files
func (s *OrderService) GetOrderWithFiles(orderID uint) (*models.Order, error) {
	var order models.Order
	err := s.db.Preload("ModelFiles").Preload("DetailImages").
		First(&order, orderID).Error
	return &order, err
}

// GetOrdersWithFilters gets orders with pagination and filters
func (s *OrderService) GetOrdersWithFilters(page, pageSize int, status string, startDate, endDate string) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := s.db.Model(&models.Order{})

	// Apply status filter
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Apply date range filter
	if startDate != "" && endDate != "" {
		query = query.Where("order_date BETWEEN ? AND ?", startDate, endDate)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).
		Preload("ModelFiles").Preload("DetailImages").
		Find(&orders).Error

	return orders, total, err
}

// GetAllOrders 获取所有订单
func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	err := s.db.Preload("Designer").Preload("Customer").Preload("Product").
		Preload("ModelFiles").Preload("DetailImages").
		Find(&orders).Error
	return orders, err
} 