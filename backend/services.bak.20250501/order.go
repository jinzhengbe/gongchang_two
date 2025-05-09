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

func (s *OrderService) GetOrderByID(orderID uint) (*models.Order, error) {
	var order models.Order
	err := s.db.Preload("Factory").First(&order, orderID).Error
	return &order, err
}

func (s *OrderService) UpdateOrderStatus(orderID uint, status models.OrderStatus) error {
	return s.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func (s *OrderService) SearchOrders(query string, factoryID string) ([]models.Order, error) {
	var orders []models.Order
	err := s.db.Where("factory_id = ? AND (description LIKE ? OR title LIKE ?)", 
		factoryID, "%"+query+"%", "%"+query+"%").
		Preload("Factory").
		Find(&orders).Error
	return orders, err
}

func (s *OrderService) GetOrderStatistics(factoryID string) (*models.OrderStatistics, error) {
	var stats models.OrderStatistics
	
	// 获取总订单数
	err := s.db.Model(&models.Order{}).Where("factory_id = ?", factoryID).Count(&stats.TotalOrders).Error
	if err != nil {
		return nil, err
	}
	
	// 获取活跃订单数（已发布状态）
	err = s.db.Model(&models.Order{}).Where("factory_id = ? AND status = ?", factoryID, models.OrderStatusPublished).Count(&stats.ActiveOrders).Error
	if err != nil {
		return nil, err
	}
	
	// 获取已完成订单数
	err = s.db.Model(&models.Order{}).Where("factory_id = ? AND status = ?", factoryID, models.OrderStatusCompleted).Count(&stats.CompletedOrders).Error
	if err != nil {
		return nil, err
	}

	// 获取待处理订单数
	err = s.db.Model(&models.Order{}).Where("factory_id = ? AND status = ?", factoryID, models.OrderStatusDraft).Count(&stats.PendingOrders).Error
	if err != nil {
		return nil, err
	}

	// 获取各状态订单数量
	stats.StatusCounts = make(map[string]int64)
	for _, status := range []models.OrderStatus{models.OrderStatusDraft, models.OrderStatusPublished, models.OrderStatusCompleted, models.OrderStatusCancelled} {
		var count int64
		err = s.db.Model(&models.Order{}).Where("factory_id = ? AND status = ?", factoryID, status).Count(&count).Error
		if err != nil {
			return nil, err
		}
		stats.StatusCounts[string(status)] = count
	}
	
	return &stats, nil
}

func (s *OrderService) GetRecentOrders(limit int) ([]models.Order, error) {
	var orders []models.Order
	err := s.db.Preload("Factory").Order("created_at desc").Limit(limit).Find(&orders).Error
	return orders, err
}

func (s *OrderService) GetOrdersByUserID(factoryID string, status string, page int, pageSize int) ([]models.Order, error) {
	var orders []models.Order
	query := s.db.Model(&models.Order{}).Where("factory_id = ?", factoryID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Preload("Factory").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at desc").
		Find(&orders).Error
	return orders, err
}

func (s *OrderService) GetOrdersCount(factoryID string, status string) (int64, error) {
	var total int64
	query := s.db.Model(&models.Order{}).Where("factory_id = ?", factoryID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	return total, err
} 