package controllers

import (
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
// @Success 201 {object} models.OrderProgress
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

	// 设置创建者ID
	req.CreatorID = userID

	progress, err := c.progressService.CreateProgress(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, progress)
}

// GetProgressByOrderID 获取订单进度列表
// @Summary 获取订单进度列表
// @Description 获取指定订单的所有进度记录
// @Tags 进度管理
// @Accept json
// @Produce json
// @Param orderId path int true "订单ID"
// @Success 200 {array} models.OrderProgress
// @Router /api/orders/{orderId}/progress [get]
func (c *ProgressController) GetProgressByOrderID(ctx *gin.Context) {
	orderID, err := strconv.ParseUint(ctx.Param("orderId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的订单ID"})
		return
	}

	progress, err := c.progressService.GetProgressByOrderID(uint(orderID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, progress)
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

	ctx.JSON(http.StatusOK, updatedProgress)
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
		progressResponses = append(progressResponses, models.ProgressResponse{
			ID:                      p.ID,
			OrderID:                 p.OrderID,
			FactoryID:               p.FactoryID,
			ProgressType:            p.ProgressType,
			Percentage:              p.Percentage,
			Status:                  p.Status,
			Description:             p.Description,
			EstimatedCompletionTime: p.EstimatedCompletionTime,
			ActualCompletionTime:    p.ActualCompletionTime,
			CreatorID:               p.CreatorID,
			CreatedAt:               p.CreatedAt,
			UpdatedAt:               p.UpdatedAt,
			Order:                   &p.Order,
			Factory:                 &p.Factory,
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
// @Description 获取工厂的进度统计信息
// @Tags 进度管理
// @Accept json
// @Produce json
// @Param factory_id path string true "工厂ID"
// @Success 200 {object} map[string]int64
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

	ctx.JSON(http.StatusOK, stats)
} 