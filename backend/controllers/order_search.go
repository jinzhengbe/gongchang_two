package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gongChang/models"
	"gongChang/services"
)

type OrderSearchController struct {
	orderSearchService *services.OrderSearchService
}

func NewOrderSearchController(orderSearchService *services.OrderSearchService) *OrderSearchController {
	return &OrderSearchController{
		orderSearchService: orderSearchService,
	}
}

// SearchOrders 高级订单搜索
// @Summary 高级订单搜索
// @Description 支持关键词搜索、状态筛选、时间筛选、排序和分页的订单搜索
// @Tags 订单搜索
// @Accept json
// @Produce json
// @Param query query string false "搜索关键词"
// @Param status query string false "订单状态筛选"
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param sort_by query string false "排序字段" default(created_at)
// @Param sort_order query string false "排序方向" default(desc)
// @Success 200 {object} models.OrderSearchResponse
// @Router /api/orders/search [get]
func (c *OrderSearchController) SearchOrders(ctx *gin.Context) {
	// 获取用户信息
	userID := ctx.GetString("user_id")
	userRole := ctx.GetString("user_role")

	// 构建搜索请求
	req := &models.OrderSearchRequest{
		UserID:   userID,
		UserRole: userRole,
	}

	// 解析查询参数
	if query := ctx.Query("query"); query != "" {
		req.Query = query
	}
	if status := ctx.Query("status"); status != "" {
		req.Status = status
	}
	if startDate := ctx.Query("start_date"); startDate != "" {
		req.StartDate = startDate
	}
	if endDate := ctx.Query("end_date"); endDate != "" {
		req.EndDate = endDate
	}

	// 解析分页参数
	if pageStr := ctx.DefaultQuery("page", "1"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			req.Page = page
		}
	}
	if pageSizeStr := ctx.DefaultQuery("page_size", "20"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			req.PageSize = pageSize
		}
	}

	// 解析排序参数
	if sortBy := ctx.DefaultQuery("sort_by", "created_at"); sortBy != "" {
		req.SortBy = sortBy
	}
	if sortOrder := ctx.DefaultQuery("sort_order", "desc"); sortOrder != "" {
		req.SortOrder = sortOrder
	}

	// 执行搜索
	result, err := c.orderSearchService.SearchOrders(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GetSearchSuggestions 获取搜索建议
// @Summary 获取搜索建议
// @Description 根据输入的关键词获取智能搜索建议
// @Tags 订单搜索
// @Accept json
// @Produce json
// @Param query query string true "搜索关键词"
// @Param limit query int false "建议数量" default(10)
// @Success 200 {object} models.SearchSuggestionResponse
// @Router /api/orders/search/suggestions [get]
func (c *OrderSearchController) GetSearchSuggestions(ctx *gin.Context) {
	// 构建建议请求
	req := &models.SearchSuggestionRequest{}

	// 解析查询参数
	if query := ctx.Query("query"); query != "" {
		req.Query = query
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "查询关键词不能为空",
		})
		return
	}

	// 解析限制参数
	if limitStr := ctx.DefaultQuery("limit", "10"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			req.Limit = limit
		}
	}

	// 获取搜索建议
	result, err := c.orderSearchService.GetSearchSuggestions(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GetSearchStatistics 获取搜索统计信息
// @Summary 获取搜索统计信息
// @Description 获取订单搜索相关的统计信息
// @Tags 订单搜索
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/orders/search/statistics [get]
func (c *OrderSearchController) GetSearchStatistics(ctx *gin.Context) {
	// 获取用户信息
	userID := ctx.GetString("user_id")
	userRole := ctx.GetString("user_role")

	// 这里可以添加搜索统计逻辑，比如热门搜索词、搜索趋势等
	// 目前返回基础统计信息
	stats := gin.H{
		"success": true,
		"data": gin.H{
			"total_orders": 0, // 这里应该从数据库获取实际数据
			"user_id":      userID,
			"user_role":    userRole,
			"hot_keywords": []string{
				"连衣裙",
				"真丝面料",
				"春季新款",
			},
		},
	}

	ctx.JSON(http.StatusOK, stats)
} 