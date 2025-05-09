package controllers

import (
	"aneworder.com/backend/models"
	"aneworder.com/backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	var order models.Order
	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证必要字段
	if order.Quantity <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "数量必须大于 0"})
		return
	}

	if order.Title == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "订单标题不能为空"})
		return
	}

	// 设置默认值
	if order.Status == "" {
		order.Status = models.OrderStatusDraft
	}

	if err := c.orderService.CreateOrder(&order); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, order)
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

	// 组装返回格式，保证四个字段为数组
	orderList := make([]gin.H, 0, len(orders))
	for _, order := range orders {
		fileIDs := order.FileIDs
		if fileIDs == nil { fileIDs = []string{} }
		modelIDs := order.ModelIDs
		if modelIDs == nil { modelIDs = []string{} }
		imageIDs := order.ImageIDs
		if imageIDs == nil { imageIDs = []string{} }
		videoIDs := order.VideoIDs
		if videoIDs == nil { videoIDs = []string{} }
		orderList = append(orderList, gin.H{
			"id": order.ID,
			"title": order.Title,
			"description": order.Description,
			"fabric": order.Fabric,
			"quantity": order.Quantity,
			"factory_id": order.FactoryID,
			"status": order.Status,
			"file_ids": fileIDs,
			"model_ids": modelIDs,
			"image_ids": imageIDs,
			"video_ids": videoIDs,
			"created_at": order.CreatedAt,
			"updated_at": order.UpdatedAt,
			// 可按需补充其他字段
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

	// 保证四个字段为非nil数组
	fileIDs := order.FileIDs
	if fileIDs == nil { fileIDs = []string{} }
	modelIDs := order.ModelIDs
	if modelIDs == nil { modelIDs = []string{} }
	imageIDs := order.ImageIDs
	if imageIDs == nil { imageIDs = []string{} }
	videoIDs := order.VideoIDs
	if videoIDs == nil { videoIDs = []string{} }

	ctx.JSON(http.StatusOK, gin.H{
		"id": order.ID,
		"title": order.Title,
		"description": order.Description,
		"fabric": order.Fabric,
		"quantity": order.Quantity,
		"factory_id": order.FactoryID,
		"status": order.Status,
		"file_ids": fileIDs,
		"model_ids": modelIDs,
		"image_ids": imageIDs,
		"video_ids": videoIDs,
		"created_at": order.CreatedAt,
		"updated_at": order.UpdatedAt,
		// 可按需补充其他字段
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
		fileIDs := order.FileIDs
		if fileIDs == nil { fileIDs = []string{} }
		modelIDs := order.ModelIDs
		if modelIDs == nil { modelIDs = []string{} }
		imageIDs := order.ImageIDs
		if imageIDs == nil { imageIDs = []string{} }
		videoIDs := order.VideoIDs
		if videoIDs == nil { videoIDs = []string{} }
		orderList = append(orderList, gin.H{
			"id": order.ID,
			"title": order.Title,
			"description": order.Description,
			"fabric": order.Fabric,
			"quantity": order.Quantity,
			"factory_id": order.FactoryID,
			"status": order.Status,
			"file_ids": fileIDs,
			"model_ids": modelIDs,
			"image_ids": imageIDs,
			"video_ids": videoIDs,
			"created_at": order.CreatedAt,
			"updated_at": order.UpdatedAt,
			// 可按需补充其他字段
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
		fileIDs := order.FileIDs
		if fileIDs == nil { fileIDs = []string{} }
		modelIDs := order.ModelIDs
		if modelIDs == nil { modelIDs = []string{} }
		imageIDs := order.ImageIDs
		if imageIDs == nil { imageIDs = []string{} }
		videoIDs := order.VideoIDs
		if videoIDs == nil { videoIDs = []string{} }
		orderList = append(orderList, gin.H{
			"id": order.ID,
			"title": order.Title,
			"description": order.Description,
			"fabric": order.Fabric,
			"quantity": order.Quantity,
			"factory_id": order.FactoryID,
			"status": order.Status,
			"file_ids": fileIDs,
			"model_ids": modelIDs,
			"image_ids": imageIDs,
			"video_ids": videoIDs,
			"created_at": order.CreatedAt,
			"updated_at": order.UpdatedAt,
			// 可按需补充其他字段
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"orders": orderList})
} 