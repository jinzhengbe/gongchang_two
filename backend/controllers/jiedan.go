package controllers

import (
	"net/http"
	"strconv"
	"gongChang/models"
	"gongChang/services"

	"github.com/gin-gonic/gin"
)

type JiedanController struct {
	jiedanService *services.JiedanService
}

func NewJiedanController(jiedanService *services.JiedanService) *JiedanController {
	return &JiedanController{
		jiedanService: jiedanService,
	}
}

// CreateJiedan 创建接单记录
// @Summary 创建接单记录
// @Description 工厂对订单进行接单操作
// @Tags 接单管理
// @Accept json
// @Produce json
// @Param request body models.CreateJiedanRequest true "创建接单请求"
// @Success 201 {object} models.Jiedan
// @Router /api/jiedan [post]
func (c *JiedanController) CreateJiedan(ctx *gin.Context) {
	var req models.CreateJiedanRequest
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
		ctx.JSON(http.StatusForbidden, gin.H{"error": "只有工厂用户可以进行接单操作"})
		return
	}

	// 验证工厂ID是否与当前用户匹配
	if req.FactoryID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "只能以自己的工厂身份进行接单"})
		return
	}

	jiedan, err := c.jiedanService.CreateJiedan(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, jiedan)
}

// GetJiedanByID 根据ID获取接单记录
// @Summary 获取接单记录详情
// @Description 根据ID获取接单记录的详细信息
// @Tags 接单管理
// @Accept json
// @Produce json
// @Param id path int true "接单记录ID"
// @Success 200 {object} models.Jiedan
// @Router /api/jiedan/{id} [get]
func (c *JiedanController) GetJiedanByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的接单记录ID"})
		return
	}

	jiedan, err := c.jiedanService.GetJiedanByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "接单记录不存在"})
		return
	}

	ctx.JSON(http.StatusOK, jiedan)
}

// GetJiedansByOrderID 根据订单ID获取接单记录列表
// @Summary 获取订单的接单记录列表
// @Description 根据订单ID获取所有相关的接单记录
// @Tags 接单管理
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Success 200 {array} models.Jiedan
// @Router /api/orders/{id}/jiedans [get]
func (c *JiedanController) GetJiedansByOrderID(ctx *gin.Context) {
	orderID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的订单ID"})
		return
	}

	jiedans, err := c.jiedanService.GetJiedansByOrderID(uint(orderID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, jiedans)
}

// GetJiedansByFactoryID 根据工厂ID获取接单记录列表
// @Summary 获取工厂的接单记录列表
// @Description 根据工厂ID获取该工厂的所有接单记录
// @Tags 接单管理
// @Accept json
// @Produce json
// @Param factory_id path string true "工厂ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} models.JiedanListResponse
// @Router /api/factories/{factory_id}/jiedans [get]
func (c *JiedanController) GetJiedansByFactoryID(ctx *gin.Context) {
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

	jiedans, total, err := c.jiedanService.GetJiedansByFactoryID(factoryID, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为响应格式
	var jiedanResponses []models.JiedanResponse
	for _, jiedan := range jiedans {
		jiedanResponses = append(jiedanResponses, models.JiedanResponse{
			ID:          jiedan.ID,
			OrderID:     jiedan.OrderID,
			FactoryID:   jiedan.FactoryID,
			Status:      jiedan.Status,
			Price:       jiedan.Price,
			JiedanTime:  jiedan.JiedanTime,
			AgreeTime:   jiedan.AgreeTime,
			AgreeUserID: jiedan.AgreeUserID,
			CreatedAt:   jiedan.CreatedAt,
			UpdatedAt:   jiedan.UpdatedAt,
			Order:       &jiedan.Order,
			Factory:     &jiedan.Factory,
		})
	}

	response := models.JiedanListResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Jiedans:  jiedanResponses,
	}

	ctx.JSON(http.StatusOK, response)
}

// AcceptJiedan 同意接单
// @Summary 同意接单
// @Description 同意接单操作
// @Tags 接单管理
// @Accept json
// @Produce json
// @Param id path int true "接单记录ID"
// @Param request body models.AcceptJiedanRequest true "同意接单请求"
// @Success 200 {object} models.Jiedan
// @Router /api/jiedan/{id}/accept [post]
func (c *JiedanController) AcceptJiedan(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的接单记录ID"})
		return
	}

	var req models.AcceptJiedanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jiedan, err := c.jiedanService.AcceptJiedan(uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, jiedan)
}

// RejectJiedan 拒绝接单
// @Summary 拒绝接单
// @Description 拒绝接单操作
// @Tags 接单管理
// @Accept json
// @Produce json
// @Param id path int true "接单记录ID"
// @Param request body models.RejectJiedanRequest true "拒绝接单请求"
// @Success 200 {object} models.Jiedan
// @Router /api/jiedan/{id}/reject [post]
func (c *JiedanController) RejectJiedan(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的接单记录ID"})
		return
	}

	var req models.RejectJiedanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jiedan, err := c.jiedanService.RejectJiedan(uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, jiedan)
}

// UpdateJiedan 更新接单记录
// @Summary 更新接单记录
// @Description 更新接单记录信息
// @Tags 接单管理
// @Accept json
// @Produce json
// @Param id path int true "接单记录ID"
// @Param request body models.UpdateJiedanRequest true "更新接单请求"
// @Success 200 {object} models.Jiedan
// @Router /api/jiedan/{id} [put]
func (c *JiedanController) UpdateJiedan(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的接单记录ID"})
		return
	}

	var req models.UpdateJiedanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jiedan, err := c.jiedanService.UpdateJiedan(uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, jiedan)
}

// DeleteJiedan 删除接单记录
// @Summary 删除接单记录
// @Description 删除接单记录
// @Tags 接单管理
// @Accept json
// @Produce json
// @Param id path int true "接单记录ID"
// @Success 200 {object} gin.H
// @Router /api/jiedan/{id} [delete]
func (c *JiedanController) DeleteJiedan(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的接单记录ID"})
		return
	}

	if err := c.jiedanService.DeleteJiedan(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "接单记录删除成功"})
}

// GetJiedanStatistics 获取接单统计信息
// @Summary 获取接单统计信息
// @Description 获取工厂的接单统计信息
// @Tags 接单管理
// @Accept json
// @Produce json
// @Param factory_id path string true "工厂ID"
// @Success 200 {object} map[string]int64
// @Router /api/factories/{factory_id}/jiedan-statistics [get]
func (c *JiedanController) GetJiedanStatistics(ctx *gin.Context) {
	factoryID := ctx.Param("factory_id")

	stats, err := c.jiedanService.GetJiedanStatistics(factoryID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stats)
} 