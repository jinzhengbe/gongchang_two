package services

import (
	"gongChang/models"
	"gorm.io/gorm"
	"gorm.io/datatypes"
	"encoding/json"
	"strings"
	"path/filepath"
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
	// 使用事务确保数据一致性
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 创建订单
		if err := tx.Omit("delivery_date", "order_date").Create(order).Error; err != nil {
			return err
		}

		// 如果有文件ID，创建文件关联
		if order.Attachments != nil || order.Models != nil || order.Images != nil || order.Videos != nil {
			// 获取所有文件ID
			var fileIDs []string
			if order.Attachments != nil {
				var attachments []string
				if err := json.Unmarshal(*order.Attachments, &attachments); err == nil {
					fileIDs = append(fileIDs, attachments...)
				}
			}
			if order.Models != nil {
				var models []string
				if err := json.Unmarshal(*order.Models, &models); err == nil {
					fileIDs = append(fileIDs, models...)
				}
			}
			if order.Images != nil {
				var images []string
				if err := json.Unmarshal(*order.Images, &images); err == nil {
					fileIDs = append(fileIDs, images...)
				}
			}
			if order.Videos != nil {
				var videos []string
				if err := json.Unmarshal(*order.Videos, &videos); err == nil {
					fileIDs = append(fileIDs, videos...)
				}
			}

			// 更新文件记录，设置OrderID
			if len(fileIDs) > 0 {
				if err := tx.Model(&models.File{}).Where("id IN ?", fileIDs).Update("order_id", order.ID).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (s *OrderService) GetOrderByID(orderID uint) (*models.Order, error) {
	var order models.Order
	err := s.db.Preload("Files").First(&order, orderID).Error
	if err != nil {
		return nil, err
	}

	// 处理文件关联
	if len(order.Files) > 0 {
		// 将文件按类型分类
		var attachments, models, images, videos []string
		for _, file := range order.Files {
			ext := strings.ToLower(filepath.Ext(file.Name))
			switch {
			case strings.HasPrefix(ext, ".doc") || strings.HasPrefix(ext, ".pdf") || strings.HasPrefix(ext, ".txt"):
				attachments = append(attachments, file.ID)
			case strings.HasPrefix(ext, ".stl") || strings.HasPrefix(ext, ".obj") || strings.HasPrefix(ext, ".3ds"):
				models = append(models, file.ID)
			case strings.HasPrefix(ext, ".jpg") || strings.HasPrefix(ext, ".png") || strings.HasPrefix(ext, ".gif"):
				images = append(images, file.ID)
			case strings.HasPrefix(ext, ".mp4") || strings.HasPrefix(ext, ".avi") || strings.HasPrefix(ext, ".mov"):
				videos = append(videos, file.ID)
			}
		}

		// 更新 JSON 字段
		if len(attachments) > 0 {
			jsonData, _ := json.Marshal(attachments)
			order.Attachments = (*datatypes.JSON)(&jsonData)
		} else {
			emptyArray := []string{}
			jsonData, _ := json.Marshal(emptyArray)
			order.Attachments = (*datatypes.JSON)(&jsonData)
		}

		if len(models) > 0 {
			jsonData, _ := json.Marshal(models)
			order.Models = (*datatypes.JSON)(&jsonData)
		} else {
			emptyArray := []string{}
			jsonData, _ := json.Marshal(emptyArray)
			order.Models = (*datatypes.JSON)(&jsonData)
		}

		if len(images) > 0 {
			jsonData, _ := json.Marshal(images)
			order.Images = (*datatypes.JSON)(&jsonData)
		} else {
			emptyArray := []string{}
			jsonData, _ := json.Marshal(emptyArray)
			order.Images = (*datatypes.JSON)(&jsonData)
		}

		if len(videos) > 0 {
			jsonData, _ := json.Marshal(videos)
			order.Videos = (*datatypes.JSON)(&jsonData)
		} else {
			emptyArray := []string{}
			jsonData, _ := json.Marshal(emptyArray)
			order.Videos = (*datatypes.JSON)(&jsonData)
		}
	} else {
		// 如果没有文件，设置空数组
		emptyArray := []string{}
		jsonData, _ := json.Marshal(emptyArray)
		order.Attachments = (*datatypes.JSON)(&jsonData)
		order.Models = (*datatypes.JSON)(&jsonData)
		order.Images = (*datatypes.JSON)(&jsonData)
		order.Videos = (*datatypes.JSON)(&jsonData)
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
	err := s.db.Preload("Factory").Order("id desc").Limit(limit).Find(&orders).Error
	return orders, err
}

func (s *OrderService) GetOrdersByUserID(userID string, status string, page int, pageSize int) ([]models.Order, error) {
	var orders []models.Order
	query := s.db.Model(&models.Order{}).Where("designer_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	// 增加SQL调试日志
	query = query.Debug()
	err := query.Preload("Factory").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("id desc").
		Find(&orders).Error
	return orders, err
}

func (s *OrderService) GetOrdersCount(userID string, status string) (int64, error) {
	query := s.db.Model(&models.Order{}).Where("customer_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	var count int64
	err := query.Count(&count).Error
	return count, err
}

func (s *OrderService) UpdateOrder(orderID uint, req *models.OrderUpdateRequest) error {
	order := &models.Order{
		Title:             req.Title,
		Description:       req.Description,
		Fabric:            req.Fabric,
		Quantity:          req.Quantity,
		PaymentStatus:     req.PaymentStatus,
		ShippingAddress:   req.ShippingAddress,
		OrderType:         req.OrderType,
		Fabrics:           req.Fabrics,
		DeliveryDate:      req.DeliveryDate,
		OrderDate:         req.OrderDate,
		SpecialRequirements: req.SpecialRequirements,
	}

	if req.Status != "" {
		order.Status = models.OrderStatus(req.Status)
	}

	// 处理文件关联
	if req.Attachments != nil && len(req.Attachments) > 0 {
		attachmentsJSON, _ := json.Marshal(req.Attachments)
		jsonData := datatypes.JSON(attachmentsJSON)
		order.Attachments = &jsonData
	}

	if req.Models != nil && len(req.Models) > 0 {
		modelsJSON, _ := json.Marshal(req.Models)
		jsonData := datatypes.JSON(modelsJSON)
		order.Models = &jsonData
	}

	if req.Images != nil && len(req.Images) > 0 {
		imagesJSON, _ := json.Marshal(req.Images)
		jsonData := datatypes.JSON(imagesJSON)
		order.Images = &jsonData
	}

	if req.Videos != nil && len(req.Videos) > 0 {
		videosJSON, _ := json.Marshal(req.Videos)
		jsonData := datatypes.JSON(videosJSON)
		order.Videos = &jsonData
	}

	return s.db.Model(&models.Order{}).Where("id = ?", orderID).Updates(order).Error
}

func (s *OrderService) DeleteOrder(orderID uint) error {
	return s.db.Delete(&models.Order{}, orderID).Error
}

func (s *OrderService) GetPublicOrders(page, pageSize int) ([]models.PublicOrder, error) {
	var orders []models.PublicOrder
	offset := (page - 1) * pageSize
	
	err := s.db.Model(&models.Order{}).
		Select("orders.id, orders.title, orders.description, orders.fabric, orders.quantity, factory_profiles.company_name as factory, orders.status, orders.created_at as create_time").
		Joins("LEFT JOIN factory_profiles ON orders.factory_id = factory_profiles.user_id").
		Where("orders.status = ?", models.OrderStatusPublished).
		Order("orders.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&orders).Error
	
	return orders, err
}

func (s *OrderService) GetPublicOrdersCount() (int64, error) {
	var count int64
	err := s.db.Model(&models.Order{}).
		Where("status = ?", models.OrderStatusPublished).
		Count(&count).Error
	return count, err
}

// AddFabricToOrder 添加布料到订单
func (s *OrderService) AddFabricToOrder(orderID uint, req *models.AddFabricToOrderRequest, fabricService *FabricService) (*models.AddFabricToOrderResponse, error) {
	// 使用事务确保数据一致性
	var response *models.AddFabricToOrderResponse
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 验证订单是否存在
		var order models.Order
		if err := tx.First(&order, orderID).Error; err != nil {
			return err
		}

		// 2. 创建布料请求
		fabricReq := &models.FabricRequest{
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
		}

		// 3. 创建布料
		fabric, err := fabricService.CreateFabric(fabricReq)
		if err != nil {
			return err
		}

		// 4. 更新订单的Fabrics字段
		var fabricIDList models.FabricIDList
		if err := fabricIDList.FromCommaString(order.Fabrics); err != nil {
			// 如果解析失败，创建空列表
			fabricIDList = make(models.FabricIDList, 0)
		}

		// 添加新布料ID
		fabricIDList.AddFabricID(fabric.ID)

		// 更新订单的Fabrics字段
		newFabricsStr := fabricIDList.ToCommaString()
		if err := tx.Model(&order).Update("fabrics", newFabricsStr).Error; err != nil {
			return err
		}

		// 5. 构建响应
		response = &models.AddFabricToOrderResponse{
			Message:           "布料添加成功",
			Fabric:            fabric,
			OrderID:           orderID,
			AssociationCreated: true,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetOrderFabrics 获取订单关联的布料列表
func (s *OrderService) GetOrderFabrics(orderID uint) ([]models.Fabric, error) {
	// 1. 获取订单
	var order models.Order
	if err := s.db.First(&order, orderID).Error; err != nil {
		return nil, err
	}

	// 2. 解析Fabrics字段
	var fabricIDList models.FabricIDList
	if err := fabricIDList.FromCommaString(order.Fabrics); err != nil {
		// 如果解析失败，返回空列表
		return make([]models.Fabric, 0), nil
	}

	// 3. 如果没有布料ID，返回空列表
	if len(fabricIDList) == 0 {
		return make([]models.Fabric, 0), nil
	}

	// 4. 查询布料信息
	var fabrics []models.Fabric
	err := s.db.Where("id IN ?", fabricIDList).Find(&fabrics).Error
	return fabrics, err
} 