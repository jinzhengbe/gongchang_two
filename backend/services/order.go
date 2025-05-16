package services

import (
	"aneworder.com/backend/models"
	"gorm.io/gorm"
	"gorm.io/datatypes"
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
	// 显式排除 delivery_date 和 order_date 字段，避免插入零值时间字段
	return s.db.Omit("delivery_date", "order_date").Create(order).Error
}

func (s *OrderService) GetOrderByID(orderID uint) (*models.Order, error) {
	var order models.Order
	err := s.db.Preload("Factory").
		Preload("Files").
		First(&order, orderID).Error
	if err != nil {
		return nil, err
	}

	// 确保 JSON 字段不为 nil
	if order.Attachments == nil {
		emptyJSON := datatypes.JSON("[]")
		order.Attachments = &emptyJSON
	}
	if order.Models == nil {
		emptyJSON := datatypes.JSON("[]")
		order.Models = &emptyJSON
	}
	if order.Images == nil {
		emptyJSON := datatypes.JSON("[]")
		order.Images = &emptyJSON
	}
	if order.Videos == nil {
		emptyJSON := datatypes.JSON("[]")
		order.Videos = &emptyJSON
	}

	return &order, nil
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