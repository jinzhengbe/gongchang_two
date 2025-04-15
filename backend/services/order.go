package services

import (
	"backend/models"
	"time"

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
	err := s.db.Where("designer_id = ? OR customer_id = ?", userID, userID).
		Preload("Designer").
		Preload("Customer").
		Preload("Product").
		Find(&orders).Error
	return orders, err
}

func (s *OrderService) GetOrderByID(orderID uint) (*models.Order, error) {
	var order models.Order
	err := s.db.Preload("Designer").
		Preload("Customer").
		Preload("Product").
		First(&order, orderID).Error
	return &order, err
}

// 更新订单状态
func (s *OrderService) UpdateOrderStatus(orderID uint, status string) error {
	return s.db.Model(&models.Order{}).Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

// 搜索订单
func (s *OrderService) SearchOrders(query string, userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := s.db.Where("(designer_id = ? OR customer_id = ?) AND (title LIKE ? OR description LIKE ?)",
		userID, userID, "%"+query+"%", "%"+query+"%").
		Preload("Designer").
		Preload("Customer").
		Preload("Product").
		Find(&orders).Error
	return orders, err
}

// GetRecentOrders 获取最近的订单
func (s *OrderService) GetRecentOrders(limit int) ([]models.Order, error) {
	var orders []models.Order
	err := s.db.Order("created_at desc").
		Limit(limit).
		Preload("Designer").
		Preload("Customer").
		Preload("Product").
		Find(&orders).Error
	return orders, err
}

// GetOrderStatistics 获取订单统计信息
func (s *OrderService) GetOrderStatistics(userID uint) (*models.OrderStatistics, error) {
	var stats models.OrderStatistics
	
	// 获取总订单数
	if err := s.db.Model(&models.Order{}).
		Where("designer_id = ? OR customer_id = ?", userID, userID).
		Count(&stats.TotalOrders).Error; err != nil {
		return nil, err
	}

	// 获取各状态订单数
	statuses := []string{"pending", "processing", "completed", "cancelled"}
	stats.StatusCounts = make(map[string]int64)
	for _, status := range statuses {
		var count int64
		if err := s.db.Model(&models.Order{}).
			Where("(designer_id = ? OR customer_id = ?) AND status = ?", userID, userID, status).
			Count(&count).Error; err != nil {
			return nil, err
		}
		stats.StatusCounts[status] = count
	}

	// 获取最近30天的订单趋势
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	var trendData []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	if err := s.db.Model(&models.Order{}).
		Where("(designer_id = ? OR customer_id = ?) AND created_at >= ?", userID, userID, thirtyDaysAgo).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Group("DATE(created_at)").
		Order("date").
		Scan(&trendData).Error; err != nil {
		return nil, err
	}
	stats.TrendData = trendData

	return &stats, nil
}

// GetLatestOrders 获取最新订单列表
func (s *OrderService) GetLatestOrders(limit int) ([]models.Order, error) {
	var orders []models.Order
	err := s.db.Order("created_at desc").
		Limit(limit).
		Preload("Designer").
		Preload("Customer").
		Preload("Product").
		Find(&orders).Error
	return orders, err
}

// GetHotOrders 获取热门订单列表
func (s *OrderService) GetHotOrders(limit int) ([]models.Order, error) {
	var orders []models.Order
	// 这里可以根据实际需求定义"热门"的标准
	// 例如：根据订单金额、浏览次数等
	err := s.db.Order("total_price desc").
		Where("status != ?", "cancelled").
		Limit(limit).
		Preload("Designer").
		Preload("Customer").
		Preload("Product").
		Find(&orders).Error
	return orders, err
} 