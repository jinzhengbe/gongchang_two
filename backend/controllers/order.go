package controllers

import (
	"gongChang/models"
	"gongChang/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OrderController struct {
	orderService *services.OrderService
	DB           *gorm.DB
}

func NewOrderController(orderService *services.OrderService, db *gorm.DB) *OrderController {
	return &OrderController{
		orderService: orderService,
		DB:           db,
	}
}

func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var req models.OrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证必要字段
	if req.Quantity <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "数量必须大于 0"})
		return
	}

	if req.Title == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "订单标题不能为空"})
		return
	}

	// 创建订单对象
	order := &models.Order{
		Title:             req.Title,
		Description:       req.Description,
		Fabric:            req.Fabric,
		Quantity:          req.Quantity,
		DesignerID:        req.DesignerID,
		CustomerID:        req.CustomerID,
		UnitPrice:         req.UnitPrice,
		TotalPrice:        req.TotalPrice,
		PaymentStatus:     req.PaymentStatus,
		ShippingAddress:   req.ShippingAddress,
		OrderType:         req.OrderType,
		Fabrics:           req.Fabrics,
		SpecialRequirements: req.SpecialRequirements,
	}

	// 处理文件关联
	if req.Attachments != nil && len(req.Attachments) > 0 {
		attachmentsJSON, err := json.Marshal(req.Attachments)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attachments format"})
			return
		}
		jsonData := datatypes.JSON(attachmentsJSON)
		order.Attachments = &jsonData
	}

	if req.Models != nil && len(req.Models) > 0 {
		modelsJSON, err := json.Marshal(req.Models)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid models format"})
			return
		}
		jsonData := datatypes.JSON(modelsJSON)
		order.Models = &jsonData
	}

	if req.Images != nil && len(req.Images) > 0 {
		imagesJSON, err := json.Marshal(req.Images)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid images format"})
			return
		}
		jsonData := datatypes.JSON(imagesJSON)
		order.Images = &jsonData
	}

	if req.Videos != nil && len(req.Videos) > 0 {
		videosJSON, err := json.Marshal(req.Videos)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid videos format"})
			return
		}
		jsonData := datatypes.JSON(videosJSON)
		order.Videos = &jsonData
	}

	// 设置默认值
	if req.Status == "" {
		order.Status = models.OrderStatusDraft
	} else {
		order.Status = models.OrderStatus(req.Status)
	}

	// 判断时间字段是否为 nil
	if req.DeliveryDate != nil {
		order.DeliveryDate = req.DeliveryDate
	}
	if req.OrderDate != nil {
		order.OrderDate = req.OrderDate
	}

	// 创建订单
	if err := c.orderService.CreateOrder(order); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回响应
	ctx.JSON(http.StatusCreated, gin.H{
		"id": order.ID,
		"title": order.Title,
		"description": order.Description,
		"fabric": order.Fabric,
		"quantity": order.Quantity,
		"designer_id": order.DesignerID,
		"customer_id": order.CustomerID,
		"status": order.Status,
		"attachments": req.Attachments,
		"models": req.Models,
		"images": req.Images,
		"videos": req.Videos,
		"created_at": order.CreatedAt,
		"updated_at": order.UpdatedAt,
	})
}

func (c *OrderController) GetOrdersByUserID(ctx *gin.Context) {
	// 从 JWT token 中获取用户 ID
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	log.Printf("Getting orders for user ID: %s", userID)

	// 获取查询参数
	status := ctx.Query("status")
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	// 转换分页参数
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// 获取订单列表
	orders, err := c.orderService.GetOrdersByUserID(userID, status, page, pageSize)
	if err != nil {
		log.Printf("Error getting orders: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取总数
	total, err := c.orderService.GetOrdersCount(userID, status)
	if err != nil {
		log.Printf("Error getting orders count: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Found %d orders for user %s", total, userID)

	// 组装返回格式
	orderList := make([]gin.H, 0, len(orders))
	for _, order := range orders {
		orderList = append(orderList, gin.H{
			"id": order.ID,
			"title": order.Title,
			"description": order.Description,
			"fabric": order.Fabric,
			"quantity": order.Quantity,
			"factory_id": order.FactoryID,
			"status": order.Status,
			"attachments": order.Attachments,
			"models": order.Models,
			"images": order.Images,
			"videos": order.Videos,
			"createTime": order.CreatedAt,
			"updated_at": order.UpdatedAt,
			"designer_id": order.DesignerID,
			"customer_id": order.CustomerID,
			"unit_price": order.UnitPrice,
			"total_price": order.TotalPrice,
			"payment_status": order.PaymentStatus,
			"shipping_address": order.ShippingAddress,
			"order_type": order.OrderType,
			"fabrics": order.Fabrics,
			"delivery_date": order.DeliveryDate,
			"order_date": order.OrderDate,
			"special_requirements": order.SpecialRequirements,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": total,
		"page": page,
		"pageSize": pageSize,
		"orders": orderList,
	})
}

func (c *OrderController) GetOrderByID(ctx *gin.Context) {
	orderID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := c.orderService.GetOrderByID(uint(orderID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 处理文件信息
	var attachments, models, images, videos []string
	if order.Attachments != nil {
		json.Unmarshal(*order.Attachments, &attachments)
	}
	if order.Models != nil {
		json.Unmarshal(*order.Models, &models)
	}
	if order.Images != nil {
		json.Unmarshal(*order.Images, &images)
	}
	if order.Videos != nil {
		json.Unmarshal(*order.Videos, &videos)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": order.ID,
		"title": order.Title,
		"description": order.Description,
		"fabric": order.Fabric,
		"quantity": order.Quantity,
		"factory_id": order.FactoryID,
		"status": order.Status,
		"attachments": attachments,
		"models": models,
		"images": images,
		"videos": videos,
		"files": order.Files,
		"createTime": order.CreatedAt,
		"updated_at": order.UpdatedAt,
		"designer_id": order.DesignerID,
		"customer_id": order.CustomerID,
		"unit_price": order.UnitPrice,
		"total_price": order.TotalPrice,
		"payment_status": order.PaymentStatus,
		"shipping_address": order.ShippingAddress,
		"order_type": order.OrderType,
		"fabrics": order.Fabrics,
		"delivery_date": order.DeliveryDate,
		"order_date": order.OrderDate,
		"special_requirements": order.SpecialRequirements,
	})
}

func (c *OrderController) UpdateOrderStatus(ctx *gin.Context) {
	orderID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var statusUpdate struct {
		Status models.OrderStatus `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&statusUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.orderService.UpdateOrderStatus(uint(orderID), statusUpdate.Status); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

func (c *OrderController) SearchOrders(ctx *gin.Context) {
	factoryID := ctx.GetString("user_id")
	if factoryID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	query := ctx.Query("q")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	orders, err := c.orderService.SearchOrders(query, factoryID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	orderList := make([]gin.H, 0, len(orders))
	for _, order := range orders {
		orderList = append(orderList, gin.H{
			"id": order.ID,
			"title": order.Title,
			"description": order.Description,
			"fabric": order.Fabric,
			"quantity": order.Quantity,
			"factory_id": order.FactoryID,
			"status": order.Status,
			"attachments": nil,
			"models": nil,
			"images": nil,
			"videos": nil,
			"createTime": order.CreatedAt,
			"updated_at": order.UpdatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"orders": orderList})
}

// GetOrderStatistics 获取订单统计信息
func (c *OrderController) GetOrderStatistics(ctx *gin.Context) {
	factoryID := ctx.GetString("user_id")
	if factoryID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	stats, err := c.orderService.GetOrderStatistics(factoryID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

// GetRecentOrders 获取最近订单
func (c *OrderController) GetRecentOrders(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}

	orders, err := c.orderService.GetRecentOrders(limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	orderList := make([]gin.H, 0, len(orders))
	for _, order := range orders {
		orderList = append(orderList, gin.H{
			"id": order.ID,
			"title": order.Title,
			"description": order.Description,
			"fabric": order.Fabric,
			"quantity": order.Quantity,
			"factory_id": order.FactoryID,
			"status": order.Status,
			"attachments": nil,
			"models": nil,
			"images": nil,
			"videos": nil,
			"createTime": order.CreatedAt,
			"updated_at": order.UpdatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"orders": orderList})
}

// GetOrdersByDesignerID 根据设计师ID获取订单列表
func (c *OrderController) GetOrdersByDesignerID(ctx *gin.Context) {
	designerID := ctx.Query("designer_id")
	if designerID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "设计师ID不能为空"})
		return
	}

	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	orders := []models.Order{}
	err = c.DB.Where("designer_id = ?", designerID).
		Order("id desc").
		Offset(offset).
		Limit(pageSize).
		Find(&orders).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	total := int64(0)
	err = c.DB.Model(&models.Order{}).Where("designer_id = ?", designerID).Count(&total).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": total,
		"page": page,
		"page_size": pageSize,
		"orders": orders,
	})
} 