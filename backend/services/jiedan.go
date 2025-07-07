package services

import (
	"fmt"
	"time"
	"gongChang/models"
	"gorm.io/gorm"
)

type JiedanService struct {
	db *gorm.DB
}

func NewJiedanService(db *gorm.DB) *JiedanService {
	return &JiedanService{
		db: db,
	}
}

// CreateJiedan 创建接单记录
func (s *JiedanService) CreateJiedan(req *models.CreateJiedanRequest) (*models.Jiedan, error) {
	// 检查订单是否存在
	var order models.Order
	if err := s.db.First(&order, req.OrderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("订单不存在")
		}
		return nil, err
	}

	// 检查是否已经存在该工厂对该订单的接单记录
	var existingJiedan models.Jiedan
	if err := s.db.Where("order_id = ? AND factory_id = ?", req.OrderID, req.FactoryID).First(&existingJiedan).Error; err == nil {
		return nil, fmt.Errorf("该工厂已对该订单进行过接单操作")
	}

	// 创建接单记录
	now := time.Now()
	jiedan := &models.Jiedan{
		OrderID:    req.OrderID,
		FactoryID:  req.FactoryID,
		Status:     models.JiedanStatusPending,
		Price:      req.Price,
		JiedanTime: &now,
	}

	if err := s.db.Create(jiedan).Error; err != nil {
		return nil, err
	}

	return jiedan, nil
}

// GetJiedanByID 根据ID获取接单记录
func (s *JiedanService) GetJiedanByID(id uint) (*models.Jiedan, error) {
	var jiedan models.Jiedan
	if err := s.db.Preload("Order").Preload("Factory").First(&jiedan, id).Error; err != nil {
		return nil, err
	}
	return &jiedan, nil
}

// GetJiedansByOrderID 根据订单ID获取接单记录列表
func (s *JiedanService) GetJiedansByOrderID(orderID uint) ([]models.Jiedan, error) {
	var jiedans []models.Jiedan
	if err := s.db.Where("order_id = ?", orderID).Preload("Order").Preload("Factory").Find(&jiedans).Error; err != nil {
		return nil, err
	}
	return jiedans, nil
}

// GetJiedansByFactoryID 根据工厂ID获取接单记录列表
func (s *JiedanService) GetJiedansByFactoryID(factoryID string, page, pageSize int) ([]models.Jiedan, int64, error) {
	var jiedans []models.Jiedan
	var total int64

	// 获取总数
	if err := s.db.Model(&models.Jiedan{}).Where("factory_id = ?", factoryID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := s.db.Where("factory_id = ?", factoryID).
		Preload("Order").Preload("Factory").
		Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&jiedans).Error; err != nil {
		return nil, 0, err
	}

	return jiedans, total, nil
}

// AcceptJiedan 同意接单
func (s *JiedanService) AcceptJiedan(id uint, req *models.AcceptJiedanRequest) (*models.Jiedan, error) {
	var jiedan models.Jiedan
	if err := s.db.First(&jiedan, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("接单记录不存在")
		}
		return nil, err
	}

	// 检查状态
	if jiedan.Status != models.JiedanStatusPending {
		return nil, fmt.Errorf("只能对待处理的接单进行同意操作")
	}

	// 更新接单状态
	now := time.Now()
	updates := map[string]interface{}{
		"status":        models.JiedanStatusAccepted,
		"agree_time":    &now,
		"agree_user_id": req.AgreeUserID,
	}

	if err := s.db.Model(&jiedan).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 重新获取更新后的记录
	if err := s.db.Preload("Order").Preload("Factory").First(&jiedan, id).Error; err != nil {
		return nil, err
	}

	return &jiedan, nil
}

// RejectJiedan 拒绝接单
func (s *JiedanService) RejectJiedan(id uint, req *models.RejectJiedanRequest) (*models.Jiedan, error) {
	var jiedan models.Jiedan
	if err := s.db.First(&jiedan, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("接单记录不存在")
		}
		return nil, err
	}

	// 检查状态
	if jiedan.Status != models.JiedanStatusPending {
		return nil, fmt.Errorf("只能对待处理的接单进行拒绝操作")
	}

	// 更新接单状态
	updates := map[string]interface{}{
		"status": models.JiedanStatusRejected,
	}

	if err := s.db.Model(&jiedan).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 重新获取更新后的记录
	if err := s.db.Preload("Order").Preload("Factory").First(&jiedan, id).Error; err != nil {
		return nil, err
	}

	return &jiedan, nil
}

// UpdateJiedan 更新接单记录
func (s *JiedanService) UpdateJiedan(id uint, req *models.UpdateJiedanRequest) (*models.Jiedan, error) {
	var jiedan models.Jiedan
	if err := s.db.First(&jiedan, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("接单记录不存在")
		}
		return nil, err
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Price != nil {
		updates["price"] = req.Price
	}
	if req.AgreeUserID != "" {
		updates["agree_user_id"] = req.AgreeUserID
	}

	if len(updates) > 0 {
		if err := s.db.Model(&jiedan).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	// 重新获取更新后的记录
	if err := s.db.Preload("Order").Preload("Factory").First(&jiedan, id).Error; err != nil {
		return nil, err
	}

	return &jiedan, nil
}

// DeleteJiedan 删除接单记录
func (s *JiedanService) DeleteJiedan(id uint) error {
	return s.db.Delete(&models.Jiedan{}, id).Error
}

// GetJiedanStatistics 获取接单统计信息
func (s *JiedanService) GetJiedanStatistics(factoryID string) (map[string]int64, error) {
	stats := make(map[string]int64)
	
	// 统计各状态的接单数量
	statuses := []models.JiedanStatus{
		models.JiedanStatusPending,
		models.JiedanStatusAccepted,
		models.JiedanStatusRejected,
	}

	for _, status := range statuses {
		var count int64
		if err := s.db.Model(&models.Jiedan{}).Where("factory_id = ? AND status = ?", factoryID, status).Count(&count).Error; err != nil {
			return nil, err
		}
		stats[string(status)] = count
	}

	return stats, nil
}

// GetJiedanByOrderIDAndFactoryID 根据订单ID和工厂ID获取接单记录
func (s *JiedanService) GetJiedanByOrderIDAndFactoryID(orderID uint, factoryID string) (*models.Jiedan, error) {
	var jiedan models.Jiedan
	if err := s.db.Where("order_id = ? AND factory_id = ?", orderID, factoryID).
		Preload("Order").Preload("Factory").
		First(&jiedan).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 返回nil表示没有找到记录
		}
		return nil, err
	}
	return &jiedan, nil
} 