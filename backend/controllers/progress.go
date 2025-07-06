package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"gongChang/models"
	"gongChang/services"

	"github.com/gin-gonic/gin"
)

type ProgressController struct {
	progressService *services.ProgressService
}

func NewProgressController(progressService *services.ProgressService) *ProgressController {
	return &ProgressController{
		progressService: progressService,
	}
}

// CreateProgress 创建进度记录
// @Summary 创建进度记录
// @Description 为指定订单创建新的进度记录
// @Tags 进度管理
// @Accept json
// @Produce json
// @Param orderId path int true "订单ID"
// @Param request body models.CreateProgressRequest true "创建进度请求"
// @Success 201 {object} gin.H
// @Router /api/orders/{orderId}/progress [post]
func (c *ProgressController) CreateProgress(ctx *gin.Context) {
	orderID, err := strconv.ParseUint(ctx.Param("orderId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的订单ID"})
		return
	}

	var req models.CreateProgressRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证订单ID一致性
	if uint(orderID) != req.OrderID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "URL中的订单ID与请求体中的订单ID不一致"})
		return
	}

	// 获取当前用户ID
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 验证用户角色
	userRole := ctx.GetString("user_role")
	if userRole != "factory" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "只有工厂用户可以创建进度记录"})
		return
	}

	// 验证工厂ID是否与当前用户匹配
	if req.FactoryID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "只能以自己的工厂身份创建进度记录"})
		return
	}

	progress, err := c.progressService.CreateProgress(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 处理图片数组
	var images []string
	if progress.Images != "" {
		if err := json.Unmarshal([]byte(progress.Images), &images); err != nil {
			images = []string{}
		}
	}

	// 返回符合要求文档的格式
	responseData := gin.H{
		"id":            progress.ID,
		"order_id":      progress.OrderID,
		"factory_id":    progress.FactoryID,
		"type":          progress.Type,
		"status":        progress.Status,
		"description":   progress.Description,
		"start_time":    progress.StartTime,
		"completed_time": progress.CompletedTime,
		"images":        images,
		"created_at":    progress.CreatedAt,
		"updated_at":    progress.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    responseData,
	})
}

// GetProgressByOrderID 获取订单进度列表
// @Summary 获取订单进度列表
// @Description 获取指定订单的所有进度记录
// @Tags 进度管理
// @Accept json
// @Produce json
// @Param orderId path int true "订单ID"
// @Success 200 {object} gin.H
// @Router /api/orders/{orderId}/progress [get]
func (c *ProgressController) GetProgressByOrderID(ctx *gin.Context) {
	orderID, err := strconv.ParseUint(ctx.Param("orderId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的订单ID"})
		return
	}

	progressList, err := c.progressService.GetProgressByOrderID(uint(orderID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为符合要求文档的格式
	var responseData []gin.H
	for _, progress := range progressList {
		// 处理图片数组
		var images []string
		if progress.Images != "" {
			if err := json.Unmarshal([]byte(progress.Images), &images); err != nil {
				images = []string{}
			}
		}

		responseData = append(responseData, gin.H{
			"id":            progress.ID,
			"order_id":      progress.OrderID,
			"factory_id":    progress.FactoryID,
			"type":          progress.Type,
			"status":        progress.Status,
			"description":   progress.Description,
			"start_time":    progress.StartTime,
			"completed_time": progress.CompletedTime,
			"images":        images,
			"created_at":    progress.CreatedAt,
			"updated_at":    progress.UpdatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    responseData,
	})
}

// UpdateProgress 更新进度记录
// @Summary 更新进度记录
// @Description 更新指定的进度记录
// @Tags 进度管理
// @Accept json
// @Produce json
// @Param orderId path int true "订单ID"
// @Param progressId path int true "进度记录ID"
// @Param request body models.UpdateProgressRequest true "更新进度请求"
// @Success 200 {object} models.OrderProgress
// @Router /api/orders/{orderId}/progress/{progressId} [put]
func (c *ProgressController) UpdateProgress(ctx *gin.Context) {
	orderID, err := strconv.ParseUint(ctx.Param("orderId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的订单ID"})
		return
	}

	progressID, err := strconv.ParseUint(ctx.Param("progressId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的进度记录ID"})
		return
	}

	var req models.UpdateProgressRequest
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

	// 验证用户角色
	userRole := ctx.GetString("user_role")
	if userRole != "factory" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "只有工厂用户可以更新进度记录"})
		return
	}

	// 验证权限：只能更新自己工厂的进度记录
	progress, err := c.progressService.GetProgressByID(uint(progressID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "进度记录不存在"})
		return
	}

	if progress.FactoryID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "只能更新自己工厂的进度记录"})
		return
	}

	// 验证订单ID一致性
	if progress.OrderID != uint(orderID) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "进度记录不属于指定的订单"})
		return
	}

	updatedProgress, err := c.progressService.UpdateProgress(uint(progressID), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 处理图片数组
	var images []string
	if updatedProgress.Images != "" {
		if err := json.Unmarshal([]byte(updatedProgress.Images), &images); err != nil {
			images = []string{}
		}
	}

	// 返回符合要求文档的格式
	responseData := gin.H{
		"id":            updatedProgress.ID,
		"order_id":      updatedProgress.OrderID,
		"factory_id":    updatedProgress.FactoryID,
		"type":          updatedProgress.Type,
		"status":        updatedProgress.Status,
		"description":   updatedProgress.Description,
		"start_time":    updatedProgress.StartTime,
		"completed_time": updatedProgress.CompletedTime,
		"images":        images,
		"created_at":    updatedProgress.CreatedAt,
		"updated_at":    updatedProgress.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    responseData,
	})
}

// DeleteProgress 删除进度记录
// @Summary 删除进度记录
// @Description 删除指定的进度记录
// @Tags 进度管理
// @Accept json
// @Produce json
// @Param orderId path int true "订单ID"
// @Param progressId path int true "进度记录ID"
// @Success 200 {object} gin.H
// @Router /api/orders/{orderId}/progress/{progressId} [delete]
func (c *ProgressController) DeleteProgress(ctx *gin.Context) {
	orderID, err := strconv.ParseUint(ctx.Param("orderId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的订单ID"})
		return
	}

	progressID, err := strconv.ParseUint(ctx.Param("progressId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的进度记录ID"})
		return
	}

	// 获取当前用户ID
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 验证用户角色
	userRole := ctx.GetString("user_role")
	if userRole != "factory" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "只有工厂用户可以删除进度记录"})
		return
	}

	// 验证权限：只能删除自己工厂的进度记录
	progress, err := c.progressService.GetProgressByID(uint(progressID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "进度记录不存在"})
		return
	}

	if progress.FactoryID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "只能删除自己工厂的进度记录"})
		return
	}

	// 验证订单ID一致性
	if progress.OrderID != uint(orderID) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "进度记录不属于指定的订单"})
		return
	}

	if err := c.progressService.DeleteProgress(uint(progressID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "进度记录删除成功"})
}

// GetProgressByFactoryID 获取工厂进度列表
// @Summary 获取工厂进度列表
// @Description 获取指定工厂的所有进度记录
// @Tags 进度管理
// @Accept json
// @Produce json
// @Param factory_id path string true "工厂ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} models.ProgressListResponse
// @Router /api/factories/{factory_id}/progress [get]
func (c *ProgressController) GetProgressByFactoryID(ctx *gin.Context) {
	factoryID := ctx.Param("factory_id")
	
	// 获取分页参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")
	
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// 获取当前用户ID
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 验证权限：只能查看自己工厂的进度记录
	if factoryID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "只能查看自己工厂的进度记录"})
		return
	}

	progress, total, err := c.progressService.GetProgressByFactoryID(factoryID, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为响应格式
	var progressResponses []models.ProgressResponse
	for _, p := range progress {
		// 处理图片数组
		var images []string
		if p.Images != "" {
			if err := json.Unmarshal([]byte(p.Images), &images); err != nil {
				images = []string{}
			}
		}

		progressResponses = append(progressResponses, models.ProgressResponse{
			ID:            p.ID,
			OrderID:       p.OrderID,
			FactoryID:     p.FactoryID,
			Type:          p.Type,
			Status:        p.Status,
			Description:   p.Description,
			StartTime:     p.StartTime,
			CompletedTime: p.CompletedTime,
			Images:        images,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
			Order:         &p.Order,
			Factory:       &p.Factory,
		})
	}

	response := models.ProgressListResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Progress: progressResponses,
	}

	ctx.JSON(http.StatusOK, response)
}

// GetProgressStatistics 获取进度统计信息
// @Summary 获取进度统计信息
// @Description 获取指定工厂的进度统计信息
// @Tags 进度管理
// @Accept json
// @Produce json
// @Param factory_id path string true "工厂ID"
// @Success 200 {object} models.ProgressStatistics
// @Router /api/factories/{factory_id}/progress-statistics [get]
func (c *ProgressController) GetProgressStatistics(ctx *gin.Context) {
	factoryID := ctx.Param("factory_id")

	// 获取当前用户ID
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 验证权限：只能查看自己工厂的统计信息
	if factoryID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "只能查看自己工厂的统计信息"})
		return
	}

	stats, err := c.progressService.GetProgressStatistics(factoryID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为统计响应格式
	response := models.ProgressStatistics{
		NotStarted: stats["not_started"],
		InProgress: stats["in_progress"],
		Completed:  stats["completed"],
		Delayed:    stats["delayed"],
		OnHold:     stats["on_hold"],
		Total:      stats["not_started"] + stats["in_progress"] + stats["completed"] + stats["delayed"] + stats["on_hold"],
	}

	ctx.JSON(http.StatusOK, response)
} 