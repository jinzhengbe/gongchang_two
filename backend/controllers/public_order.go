package controllers

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"aneworder.com/backend/models"
)

type PublicOrderController struct {
	db *gorm.DB
}

func NewPublicOrderController(db *gorm.DB) *PublicOrderController {
	return &PublicOrderController{db: db}
}

// GetPublicOrders 获取公开订单列表
// @Summary 获取公开订单列表
// @Description 获取公开的订单列表，无需认证
// @Tags 公开订单
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Param status query string false "订单状态" default(published)
// @Success 200 {object} models.PublicOrderResponse
// @Router /public/orders [get]
func (c *PublicOrderController) GetPublicOrders(ctx *gin.Context) {
	// 获取查询参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	status := ctx.DefaultQuery("status", "published")

	// 确保页码和每页数量合理
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// 查询订单
	var orders []models.Order
	var total int64

	// 构建查询条件
	query := c.db.Model(&models.Order{}).
		Where("status = ?", status).
		Preload("Factory")

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "获取订单总数失败"})
		return
	}

	// 分页查询
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "获取订单列表失败"})
		return
	}

	// 转换为公开订单格式
	publicOrders := make([]models.PublicOrder, len(orders))
	for i, order := range orders {
		publicOrders[i] = models.PublicOrder{
			ID:          order.ID,
			Title:       order.Title,
			Description: order.Description,
			Fabric:      order.Fabric,
			Quantity:    order.Quantity,
			Factory:     order.Factory.CompanyName,
			Status:      string(order.Status),
			CreateTime:  order.CreatedAt,
		}
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// 返回响应
	ctx.JSON(200, models.PublicOrderResponse{
		Orders:     publicOrders,
		Total:      int(total),
		Page:       page,
		PageSize:   limit,
		TotalPages: totalPages,
	})
} 