package controllers

import (
	"aneworder.com/backend/models"
	"aneworder.com/backend/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type OrderController struct {
	orderService *services.OrderService
}

func NewOrderController(orderService *services.OrderService) *OrderController {
	return &OrderController{
		orderService: orderService,
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

	// 处理数组字段
	if req.Attachments != nil {
		// 确保数组不为空
		if len(req.Attachments) > 0 {
			attachmentsJSON, err := json.Marshal(req.Attachments)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attachments format"})
				return
			}
			jsonData := datatypes.JSON(attachmentsJSON)
			order.Attachments = &jsonData
		} else {
			// 如果是空数组，设置为 null
			order.Attachments = nil
		}
	}

	if req.Models != nil {
		if len(req.Models) > 0 {
			modelsJSON, err := json.Marshal(req.Models)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid models format"})
				return
			}
			jsonData := datatypes.JSON(modelsJSON)
			order.Models = &jsonData
		} else {
			order.Models = nil
		}
	}

	if req.Images != nil {
		if len(req.Images) > 0 {
			imagesJSON, err := json.Marshal(req.Images)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid images format"})
				return
			}
			jsonData := datatypes.JSON(imagesJSON)
			order.Images = &jsonData
		} else {
			order.Images = nil
		}
	}

	if req.Videos != nil {
		if len(req.Videos) > 0 {
			videosJSON, err := json.Marshal(req.Videos)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid videos format"})
				return
			}
			jsonData := datatypes.JSON(videosJSON)
			order.Videos = &jsonData
		} else {
			order.Videos = nil
		}
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
	factoryID := ctx.GetString("user_id")
	if factoryID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

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
	orders, err := c.orderService.GetOrdersByUserID(factoryID, status, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取总数
	total, err := c.orderService.GetOrdersCount(factoryID, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
			"attachments": nil,
			"models": nil,
			"images": nil,
			"videos": nil,
			"created_at": order.CreatedAt,
			"updated_at": order.UpdatedAt,
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

	ctx.JSON(http.StatusOK, gin.H{
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
		"files": order.Files,
		"created_at": order.CreatedAt,
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
			"created_at": order.CreatedAt,
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
			"created_at": order.CreatedAt,
			"updated_at": order.UpdatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"orders": orderList})
} 