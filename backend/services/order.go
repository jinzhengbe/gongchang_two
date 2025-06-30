package services

import (
	"encoding/json"
	"fmt"
	"gongChang/models"
	"gorm.io/datatypes"
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

	// 如果数据库中的JSON字段为空，则设置为空数组
	if order.Attachments == nil {
		emptyArray := []string{}
		jsonData, _ := json.Marshal(emptyArray)
		order.Attachments = (*datatypes.JSON)(&jsonData)
	}

	if order.Models == nil {
		emptyArray := []string{}
		jsonData, _ := json.Marshal(emptyArray)
		order.Models = (*datatypes.JSON)(&jsonData)
	}

	if order.Images == nil {
		emptyArray := []string{}
		jsonData, _ := json.Marshal(emptyArray)
		order.Images = (*datatypes.JSON)(&jsonData)
	}

	if order.Videos == nil {
		emptyArray := []string{}
		jsonData, _ := json.Marshal(emptyArray)
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

func (s *OrderService) GetRecentOrders(limit int, status string) ([]models.Order, error) {
	var orders []models.Order
	query := s.db.Preload("Factory").Order("id desc").Limit(limit)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	// 兼容老用法
	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
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
	// 首先获取现有订单数据
	var existingOrder models.Order
	if err := s.db.First(&existingOrder, orderID).Error; err != nil {
		return err
	}

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

	// 处理文件关联 - 总是更新这些字段，即使为空
	if req.Attachments != nil {
		attachmentsJSON, _ := json.Marshal(req.Attachments)
		jsonData := datatypes.JSON(attachmentsJSON)
		order.Attachments = &jsonData
	} else {
		// 如果请求中没有attachments字段，设置为空数组
		emptyArray := []string{}
		attachmentsJSON, _ := json.Marshal(emptyArray)
		jsonData := datatypes.JSON(attachmentsJSON)
		order.Attachments = &jsonData
	}

	if req.Models != nil {
		modelsJSON, _ := json.Marshal(req.Models)
		jsonData := datatypes.JSON(modelsJSON)
		order.Models = &jsonData
	} else {
		// 如果请求中没有models字段，设置为空数组
		emptyArray := []string{}
		modelsJSON, _ := json.Marshal(emptyArray)
		jsonData := datatypes.JSON(modelsJSON)
		order.Models = &jsonData
	}

	// 处理Images字段 - 实现合并逻辑而不是覆盖
	if req.Images != nil {
		// 获取现有的图片ID列表
		var existingImages []string
		if existingOrder.Images != nil {
			json.Unmarshal(*existingOrder.Images, &existingImages)
		}
		
		// 合并现有图片ID和新图片ID，去重
		imageMap := make(map[string]bool)
		for _, img := range existingImages {
			imageMap[img] = true
		}
		for _, img := range req.Images {
			imageMap[img] = true
		}
		
		// 转换回切片
		mergedImages := make([]string, 0, len(imageMap))
		for img := range imageMap {
			mergedImages = append(mergedImages, img)
		}
		
		imagesJSON, _ := json.Marshal(mergedImages)
		jsonData := datatypes.JSON(imagesJSON)
		order.Images = &jsonData
	} else {
		// 如果请求中没有images字段，保持现有图片不变
		order.Images = existingOrder.Images
	}

	if req.Videos != nil {
		videosJSON, _ := json.Marshal(req.Videos)
		jsonData := datatypes.JSON(videosJSON)
		order.Videos = &jsonData
	} else {
		// 如果请求中没有videos字段，设置为空数组
		emptyArray := []string{}
		videosJSON, _ := json.Marshal(emptyArray)
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

// RemoveFabricFromOrder 从订单移除布料
func (s *OrderService) RemoveFabricFromOrder(orderID uint, req *models.RemoveFabricFromOrderRequest, fabricService *FabricService) (*models.RemoveFabricFromOrderResponse, error) {
	var response *models.RemoveFabricFromOrderResponse

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 检查订单是否存在
		var order models.Order
		if err := tx.First(&order, orderID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("订单不存在")
			}
			return err
		}

		// 2. 检查布料是否存在
		var fabric models.Fabric
		if err := tx.First(&fabric, req.FabricID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("布料不存在")
			}
			return err
		}

		// 3. 解析订单的布料ID列表
		var fabricIDList models.FabricIDList
		if err := fabricIDList.FromCommaString(order.Fabrics); err != nil {
			// 如果解析失败，创建空列表
			fabricIDList = make(models.FabricIDList, 0)
		}

		// 4. 检查订单是否包含该布料
		if !fabricIDList.ContainsFabricID(req.FabricID) {
			return fmt.Errorf("订单中不包含该布料")
		}

		// 5. 从列表中移除布料ID
		fabricIDList.RemoveFabricID(req.FabricID)

		// 6. 更新订单的布料字段
		newFabricsStr := fabricIDList.ToCommaString()
		if err := tx.Model(&order).Update("fabrics", newFabricsStr).Error; err != nil {
			return err
		}

		// 7. 构建响应
		response = &models.RemoveFabricFromOrderResponse{
			Success: true,
			Message: "布料已从订单中移除",
			Order:   &order,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

// RemoveFileFromOrder 从订单移除文件
func (s *OrderService) RemoveFileFromOrder(orderID uint, req *models.RemoveFileFromOrderRequest) (*models.RemoveFileFromOrderResponse, error) {
	var response *models.RemoveFileFromOrderResponse

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 检查订单是否存在
		var order models.Order
		if err := tx.First(&order, orderID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("订单不存在")
			}
			return err
		}

		// 2. 检查文件是否存在
		var file models.File
		if err := tx.Where("id = ?", req.FileID).First(&file).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("文件不存在")
			}
			return err
		}

		// 3. 根据文件类型更新订单的相应字段
		var existingFiles []string
		var fieldName string

		switch req.FileType {
		case "image":
			if order.Images != nil {
				json.Unmarshal(*order.Images, &existingFiles)
			}
			fieldName = "images"
		case "attachment":
			if order.Attachments != nil {
				json.Unmarshal(*order.Attachments, &existingFiles)
			}
			fieldName = "attachments"
		case "model":
			if order.Models != nil {
				json.Unmarshal(*order.Models, &existingFiles)
			}
			fieldName = "models"
		case "video":
			if order.Videos != nil {
				json.Unmarshal(*order.Videos, &existingFiles)
			}
			fieldName = "videos"
		default:
			return fmt.Errorf("不支持的文件类型: %s", req.FileType)
		}

		// 4. 从数组中移除指定的文件ID
		var newFiles []string
		found := false
		for _, fileID := range existingFiles {
			if fileID == req.FileID {
				found = true
			} else {
				newFiles = append(newFiles, fileID)
			}
		}

		if !found {
			return fmt.Errorf("订单中不包含该文件")
		}

		// 5. 更新订单的相应字段
		var jsonData datatypes.JSON
		if len(newFiles) > 0 {
			jsonBytes, _ := json.Marshal(newFiles)
			jsonData = datatypes.JSON(jsonBytes)
		} else {
			// 如果数组为空，设置为空数组
			emptyArray := []string{}
			jsonBytes, _ := json.Marshal(emptyArray)
			jsonData = datatypes.JSON(jsonBytes)
		}

		// 6. 更新订单
		updateData := map[string]interface{}{
			fieldName: jsonData,
		}
		if err := tx.Model(&order).Updates(updateData).Error; err != nil {
			return err
		}

		// 7. 重新获取更新后的订单
		if err := tx.First(&order, orderID).Error; err != nil {
			return err
		}

		// 8. 构建响应
		response = &models.RemoveFileFromOrderResponse{
			Success: true,
			Message: fmt.Sprintf("文件已从订单的%s中移除", req.FileType),
			Order:   &order,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetDB 获取数据库连接
func (s *OrderService) GetDB() *gorm.DB {
	return s.db
} 