package services

import (
	"encoding/json"
	"fmt"
	"time"
	"gongChang/models"
	"gorm.io/gorm"
)

type ProgressService struct {
	db *gorm.DB
}

func NewProgressService(db *gorm.DB) *ProgressService {
	return &ProgressService{
		db: db,
	}
}

// CreateProgress 创建进度记录
func (s *ProgressService) CreateProgress(req *models.CreateProgressRequest) (*models.OrderProgress, error) {
	// 检查订单是否存在
	var order models.Order
	if err := s.db.First(&order, req.OrderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("订单不存在")
		}
		return nil, err
	}

	// 处理图片数组
	imagesJSON := ""
	if len(req.Images) > 0 {
		imagesBytes, err := json.Marshal(req.Images)
		if err != nil {
			return nil, fmt.Errorf("图片数据格式错误: %v", err)
		}
		imagesJSON = string(imagesBytes)
	}

	// 创建进度记录
	now := time.Now()
	progress := &models.OrderProgress{
		OrderID:       req.OrderID,
		FactoryID:     req.FactoryID,
		Type:          req.Type,
		Status:        req.Status,
		Description:   req.Description,
		StartTime:     req.StartTime,
		CompletedTime: req.CompletedTime,
		Images:        imagesJSON,
		CreatedAt:     &now,
	}

	if err := s.db.Create(progress).Error; err != nil {
		return nil, err
	}

	return progress, nil
}

// GetProgressByID 根据ID获取进度记录
func (s *ProgressService) GetProgressByID(id uint) (*models.OrderProgress, error) {
	var progress models.OrderProgress
	if err := s.db.Preload("Order").Preload("Factory").First(&progress, id).Error; err != nil {
		return nil, err
	}
	return &progress, nil
}

// GetProgressByOrderID 根据订单ID获取进度记录列表
func (s *ProgressService) GetProgressByOrderID(orderID uint) ([]models.OrderProgress, error) {
	var progress []models.OrderProgress
	if err := s.db.Where("order_id = ?", orderID).
		Preload("Order").Preload("Factory").
		Order("created_at DESC").
		Find(&progress).Error; err != nil {
		return nil, err
	}
	return progress, nil
}

// GetProgressByFactoryID 根据工厂ID获取进度记录列表
func (s *ProgressService) GetProgressByFactoryID(factoryID string, page, pageSize int) ([]models.OrderProgress, int64, error) {
	var progress []models.OrderProgress
	var total int64

	// 获取总数
	if err := s.db.Model(&models.OrderProgress{}).Where("factory_id = ?", factoryID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := s.db.Where("factory_id = ?", factoryID).
		Preload("Order").Preload("Factory").
		Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&progress).Error; err != nil {
		return nil, 0, err
	}

	return progress, total, nil
}

// UpdateProgress 更新进度记录
func (s *ProgressService) UpdateProgress(id uint, req *models.UpdateProgressRequest) (*models.OrderProgress, error) {
	var progress models.OrderProgress
	if err := s.db.First(&progress, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("进度记录不存在")
		}
		return nil, err
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.StartTime != nil {
		updates["start_time"] = req.StartTime
	}
	if req.CompletedTime != nil {
		updates["completed_time"] = req.CompletedTime
	}
	
	// 处理图片数组
	if req.Images != nil {
		imagesJSON := ""
		if len(req.Images) > 0 {
			imagesBytes, err := json.Marshal(req.Images)
			if err != nil {
				return nil, fmt.Errorf("图片数据格式错误: %v", err)
			}
			imagesJSON = string(imagesBytes)
		}
		updates["images"] = imagesJSON
	}

	if len(updates) > 0 {
		if err := s.db.Model(&progress).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	// 重新获取更新后的记录
	if err := s.db.Preload("Order").Preload("Factory").First(&progress, id).Error; err != nil {
		return nil, err
	}

	return &progress, nil
}

// DeleteProgress 删除进度记录
func (s *ProgressService) DeleteProgress(id uint) error {
	return s.db.Delete(&models.OrderProgress{}, id).Error
}

// GetProgressStatistics 获取进度统计信息
func (s *ProgressService) GetProgressStatistics(factoryID string) (map[string]int64, error) {
	stats := make(map[string]int64)
	
	// 统计各状态的进度数量
	statuses := []models.ProgressStatus{
		models.ProgressStatusNotStarted,
		models.ProgressStatusInProgress,
		models.ProgressStatusCompleted,
		models.ProgressStatusDelayed,
		models.ProgressStatusOnHold,
	}

	for _, status := range statuses {
		var count int64
		if err := s.db.Model(&models.OrderProgress{}).Where("factory_id = ? AND status = ?", factoryID, status).Count(&count).Error; err != nil {
			return nil, err
		}
		stats[string(status)] = count
	}

	return stats, nil
} 