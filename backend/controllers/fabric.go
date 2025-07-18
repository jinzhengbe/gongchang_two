package controllers

import (
	"gongChang/models"
	"gongChang/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FabricController struct {
	fabricService *services.FabricService
}

func NewFabricController(fabricService *services.FabricService) *FabricController {
	return &FabricController{
		fabricService: fabricService,
	}
}

// CreateFabric 创建布料
// @Summary 创建布料
// @Description 创建新的布料记录（仅设计师和供应商可创建）
// @Tags 布料管理
// @Accept json
// @Produce json
// @Param fabric body models.FabricRequest true "布料信息"
// @Success 201 {object} models.Fabric
// @Router /api/fabrics [post]
func (fc *FabricController) CreateFabric(c *gin.Context) {
	var req models.FabricRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 获取用户角色
	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户角色未找到"})
		return
	}

	// 根据用户角色设置相应的ID字段
	switch userRole.(string) {
	case "designer":
		userIDStr := userID.(string)
		req.DesignerID = userIDStr
	case "supplier":
		userIDStr := userID.(string)
		req.SupplierID = userIDStr
	case "factory":
		c.JSON(http.StatusForbidden, gin.H{"error": "工厂账号不允许创建布料"})
		return
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "用户角色不允许创建布料"})
		return
	}

	// 创建布料
	fabric, err := fc.fabricService.CreateFabric(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "布料创建成功",
		"fabric":  fabric,
	})
}

// GetFabricByID 根据ID获取布料
// @Summary 获取布料详情
// @Description 根据ID获取布料的详细信息
// @Tags 布料管理
// @Accept json
// @Produce json
// @Param id path int true "布料ID"
// @Success 200 {object} models.Fabric
// @Router /api/fabrics/{id} [get]
func (c *FabricController) GetFabricByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fabric ID"})
		return
	}

	fabric, err := c.fabricService.GetFabricByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Fabric not found"})
		return
	}

	ctx.JSON(http.StatusOK, fabric)
}

// UpdateFabric 更新布料
// @Summary 更新布料
// @Description 更新布料的详细信息
// @Tags 布料管理
// @Accept json
// @Produce json
// @Param id path int true "布料ID"
// @Param fabric body models.FabricUpdateRequest true "布料更新信息"
// @Success 200 {object} models.Fabric
// @Router /api/fabrics/{id} [put]
func (c *FabricController) UpdateFabric(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fabric ID"})
		return
	}

	var req models.FabricUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fabric, err := c.fabricService.UpdateFabric(uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, fabric)
}

// DeleteFabric 删除布料
// @Summary 删除布料
// @Description 删除指定的布料记录
// @Tags 布料管理
// @Accept json
// @Produce json
// @Param id path int true "布料ID"
// @Success 200 {object} gin.H
// @Router /api/fabrics/{id} [delete]
func (c *FabricController) DeleteFabric(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fabric ID"})
		return
	}

	if err := c.fabricService.DeleteFabric(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Fabric deleted successfully"})
}

// SearchFabrics 搜索布料
// @Summary 搜索布料
// @Description 根据条件搜索布料列表
// @Tags 布料管理
// @Accept json
// @Produce json
// @Param q query string false "搜索关键词"
// @Param category query string false "分类"
// @Param material query string false "材质"
// @Param color query string false "颜色"
// @Param min_price query number false "最低价格"
// @Param max_price query number false "最高价格"
// @Param min_stock query int false "最低库存"
// @Param status query int false "状态"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} models.FabricListResponse
// @Router /api/fabrics/search [get]
func (c *FabricController) SearchFabrics(ctx *gin.Context) {
	var req models.FabricSearchRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	result, err := c.fabricService.SearchFabrics(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GetAllFabrics 获取所有布料（用于前端下拉选择）
// @Summary 获取所有布料
// @Description 获取所有可用的布料列表，用于前端下拉选择
// @Tags 布料管理
// @Accept json
// @Produce json
// @Success 200 {array} models.FabricResponse
// @Router /api/fabrics/all [get]
func (c *FabricController) GetAllFabrics(ctx *gin.Context) {
	fabrics, err := c.fabricService.GetAllFabrics()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, fabrics)
}

// GetFabricCategories 获取布料分类
// @Summary 获取布料分类
// @Description 获取所有布料分类列表
// @Tags 布料管理
// @Accept json
// @Produce json
// @Success 200 {array} models.FabricCategory
// @Router /api/fabrics/categories [get]
func (c *FabricController) GetFabricCategories(ctx *gin.Context) {
	categories, err := c.fabricService.GetFabricCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

// GetFabricsByCategory 根据分类获取布料
// @Summary 根据分类获取布料
// @Description 根据分类获取布料列表
// @Tags 布料管理
// @Accept json
// @Produce json
// @Param category path string true "分类名称"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} models.FabricListResponse
// @Router /api/fabrics/category/{category} [get]
func (c *FabricController) GetFabricsByCategory(ctx *gin.Context) {
	category := ctx.Param("category")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	result, err := c.fabricService.GetFabricsByCategory(category, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GetFabricsByMaterial 根据材质获取布料
// @Summary 根据材质获取布料
// @Description 根据材质获取布料列表
// @Tags 布料管理
// @Accept json
// @Produce json
// @Param material path string true "材质名称"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} models.FabricListResponse
// @Router /api/fabrics/material/{material} [get]
func (c *FabricController) GetFabricsByMaterial(ctx *gin.Context) {
	material := ctx.Param("material")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	result, err := c.fabricService.GetFabricsByMaterial(material, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// UpdateFabricStock 更新布料库存
// @Summary 更新布料库存
// @Description 更新指定布料的库存数量
// @Tags 布料管理
// @Accept json
// @Produce json
// @Param id path int true "布料ID"
// @Param quantity body int true "库存变化量"
// @Success 200 {object} gin.H
// @Router /api/fabrics/{id}/stock [put]
func (c *FabricController) UpdateFabricStock(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fabric ID"})
		return
	}

	var req struct {
		Quantity int `json:"quantity" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.fabricService.UpdateFabricStock(uint(id), req.Quantity); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Stock updated successfully"})
}

// GetFabricStatistics 获取布料统计信息
// @Summary 获取布料统计
// @Description 获取布料的统计信息
// @Tags 布料管理
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Router /api/fabrics/statistics [get]
func (c *FabricController) GetFabricStatistics(ctx *gin.Context) {
	stats, err := c.fabricService.GetFabricStatistics()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stats)
} 