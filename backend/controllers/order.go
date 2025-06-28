package controllers

import (
	"gongChang/models"
	"gongChang/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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

	// 查询布料详细信息
	var fabrics []map[string]interface{}
	fabricsIDs := order.Fabrics
	if fabricsIDs != "" {
		// 解析布料ID字符串
		parts := strings.Split(fabricsIDs, ",")
		if len(parts) > 0 {
			// 构建ID查询条件
			var fabricIDs []string
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if part != "" {
					fabricIDs = append(fabricIDs, part)
				}
			}
			
			if len(fabricIDs) > 0 {
				// 查询布料信息
				rows, err := c.DB.Raw("SELECT id, name, category, material, color, pattern, weight, width, price, unit, stock, min_order, description, image_url, thumbnail_url, tags, status, designer_id, supplier_id, factory_id, created_at, updated_at FROM fabrics WHERE id IN (?)", fabricIDs).Rows()
				if err == nil {
					defer rows.Close()
					for rows.Next() {
						var fabric map[string]interface{}
						fabric = make(map[string]interface{})
						var id uint
						var name, category, material, color, pattern, unit, description, imageURL, thumbnailURL, tags string
						var weight, width, price float64
						var stock, minOrder, status int
						var designerID, supplierID, factoryID *string
						var createdAt, updatedAt time.Time
						
						rows.Scan(&id, &name, &category, &material, &color, &pattern, &weight, &width, &price, &unit, &stock, &minOrder, &description, &imageURL, &thumbnailURL, &tags, &status, &designerID, &supplierID, &factoryID, &createdAt, &updatedAt)
						
						fabric["id"] = id
						fabric["name"] = name
						fabric["category"] = category
						fabric["material"] = material
						fabric["color"] = color
						fabric["pattern"] = pattern
						fabric["weight"] = weight
						fabric["width"] = width
						fabric["price"] = price
						fabric["unit"] = unit
						fabric["stock"] = stock
						fabric["min_order"] = minOrder
						fabric["description"] = description
						fabric["image_url"] = imageURL
						fabric["thumbnail_url"] = thumbnailURL
						fabric["tags"] = tags
						fabric["status"] = status
						fabric["designer_id"] = designerID
						fabric["supplier_id"] = supplierID
						fabric["factory_id"] = factoryID
						fabric["created_at"] = createdAt
						fabric["updated_at"] = updatedAt
						
						fabrics = append(fabrics, fabric)
					}
				}
			}
		}
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
		"fabrics": fabrics,
		"fabrics_ids": fabricsIDs,
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
	// 从 JWT token 中获取设计师 ID
	designerID := ctx.GetString("user_id")
	if designerID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	log.Printf("Getting orders for designer ID: %s", designerID)

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

	// 获取设计师的订单列表
	orders, err := c.orderService.GetOrdersByUserID(designerID, status, page, pageSize)
	if err != nil {
		log.Printf("Error getting designer orders: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取总数
	total, err := c.orderService.GetOrdersCount(designerID, status)
	if err != nil {
		log.Printf("Error getting designer orders count: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Found %d orders for designer %s", total, designerID)

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
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"orders": orderList,
		"total": total,
		"page": page,
		"pageSize": pageSize,
	})
}

// UpdateOrder 更新订单
func (c *OrderController) UpdateOrder(ctx *gin.Context) {
	orderID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var req models.OrderUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户ID
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 更新订单
	if err := c.orderService.UpdateOrder(uint(orderID), &req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order updated successfully"})
}

// DeleteOrder 删除订单
func (c *OrderController) DeleteOrder(ctx *gin.Context) {
	orderID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// 获取当前用户ID
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 删除订单
	if err := c.orderService.DeleteOrder(uint(orderID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

// GetPublicOrders 获取公开订单列表（无需认证）
func (c *OrderController) GetPublicOrders(ctx *gin.Context) {
	// 获取查询参数
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

	// 获取公开订单列表
	orders, err := c.orderService.GetPublicOrders(page, pageSize)
	if err != nil {
		log.Printf("Error getting public orders: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取总数
	total, err := c.orderService.GetPublicOrdersCount()
	if err != nil {
		log.Printf("Error getting public orders count: %v", err)
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
			"factory": order.Factory,
			"status": order.Status,
			"createTime": order.CreateTime,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"orders": orderList,
		"total": total,
		"page": page,
		"pageSize": pageSize,
	})
}

// AddFabricToOrder 添加布料到订单
// @Summary 添加布料到订单
// @Description 创建新布料并关联到指定订单
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Param request body models.AddFabricToOrderRequest true "布料信息"
// @Success 201 {object} models.AddFabricToOrderResponse
// @Router /api/orders/{id}/add-fabric [post]
func (c *OrderController) AddFabricToOrder(ctx *gin.Context) {
	// 获取订单ID
	orderID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的订单ID"})
		return
	}

	// 绑定请求数据
	var req models.AddFabricToOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证订单ID一致性
	if uint(orderID) != req.OrderID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "URL中的订单ID与请求体中的订单ID不一致"})
		return
	}

	// 获取当前用户角色
	userRole, exists := ctx.Get("user_role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户角色未找到"})
		return
	}

	// 验证权限：只有设计师和供应商可以添加布料
	switch userRole.(string) {
	case "designer", "supplier":
		// 允许操作
	default:
		ctx.JSON(http.StatusForbidden, gin.H{"error": "用户角色不允许添加布料"})
		return
	}

	// 创建布料服务实例
	fabricService := services.NewFabricService(c.DB)

	// 根据用户角色设置相应的ID字段
	switch userRole.(string) {
	case "designer":
		// 设计师角色，在服务层处理设计师ID设置
	case "supplier":
		// 供应商角色，在服务层处理供应商ID设置
	}

	// 调用服务层方法
	response, err := c.orderService.AddFabricToOrder(uint(orderID), &req, fabricService)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// RemoveFabricFromOrder 从订单移除布料
// @Summary 从订单移除布料
// @Description 从指定订单中移除指定的布料
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Param request body models.RemoveFabricFromOrderRequest true "移除布料请求"
// @Success 200 {object} models.RemoveFabricFromOrderResponse
// @Router /api/orders/{id}/remove-fabric [delete]
func (c *OrderController) RemoveFabricFromOrder(ctx *gin.Context) {
	// 获取订单ID
	orderID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的订单ID"})
		return
	}

	// 绑定请求参数
	var req models.RemoveFabricFromOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 验证布料ID
	if req.FabricID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "布料ID不能为空"})
		return
	}

	// 创建布料服务实例
	fabricService := services.NewFabricService(c.orderService.GetDB())

	// 调用服务层方法
	response, err := c.orderService.RemoveFabricFromOrder(uint(orderID), &req, fabricService)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// RemoveFileFromOrder 从订单移除文件
// @Summary 从订单移除文件
// @Description 从指定订单中移除指定的文件（图片、附件、模型或视频）
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Param request body models.RemoveFileFromOrderRequest true "移除文件请求"
// @Success 200 {object} models.RemoveFileFromOrderResponse
// @Router /api/orders/{id}/remove-file [delete]
func (c *OrderController) RemoveFileFromOrder(ctx *gin.Context) {
	// 获取订单ID
	orderID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的订单ID"})
		return
	}

	// 绑定请求参数
	var req models.RemoveFileFromOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 验证文件ID
	if req.FileID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "文件ID不能为空"})
		return
	}

	// 验证文件类型
	if req.FileType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "文件类型不能为空"})
		return
	}

	// 验证文件类型是否有效
	validTypes := map[string]bool{
		"image":       true,
		"attachment":  true,
		"model":       true,
		"video":       true,
	}
	if !validTypes[req.FileType] {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件类型，支持的类型：image, attachment, model, video"})
		return
	}

	// 调用服务层方法
	response, err := c.orderService.RemoveFileFromOrder(uint(orderID), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
} 