package controllers

import (
	"fmt"
	"net/http"
	"gongChang/models"
	"gongChang/services"

	"github.com/gin-gonic/gin"
)

type FactorySearchController struct {
	factorySearchService *services.FactorySearchService
}

func NewFactorySearchController(factorySearchService *services.FactorySearchService) *FactorySearchController {
	return &FactorySearchController{
		factorySearchService: factorySearchService,
	}
}

// SearchFactories 搜索工厂
// @Summary 搜索工厂
// @Description 搜索工厂，支持关键词搜索、地区筛选、专业领域筛选等
// @Tags 工厂搜索
// @Accept json
// @Produce json
// @Param query query string false "搜索关键词"
// @Param region query string false "地区筛选"
// @Param specialties query []string false "专业领域数组"
// @Param cooperation_status query string false "合作状态"
// @Param min_rating query number false "最低评分"
// @Param max_rating query number false "最高评分"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param sort_by query string false "排序字段" default(rating)
// @Param sort_order query string false "排序方向" default(desc)
// @Success 200 {object} models.FactorySearchResponse
// @Router /api/factories/search [get]
func (c *FactorySearchController) SearchFactories(ctx *gin.Context) {
	var req models.FactorySearchRequest
	
	// 绑定查询参数
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "参数绑定失败: " + err.Error(),
		})
		return
	}

	// 处理专业领域数组参数
	if specialties := ctx.QueryArray("specialties"); len(specialties) > 0 {
		req.Specialties = specialties
	}

	// 调用服务层搜索工厂
	result, err := c.factorySearchService.SearchFactories(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "搜索失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GetSearchSuggestions 获取搜索建议
// @Summary 获取工厂搜索建议
// @Description 获取工厂搜索建议，提供智能提示
// @Tags 工厂搜索
// @Accept json
// @Produce json
// @Param query query string true "搜索关键词"
// @Param limit query int false "建议数量" default(10)
// @Success 200 {object} models.FactorySearchSuggestionResponse
// @Router /api/factories/search/suggestions [get]
func (c *FactorySearchController) GetSearchSuggestions(ctx *gin.Context) {
	var req models.FactorySearchSuggestionRequest
	
	// 绑定查询参数
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "参数绑定失败: " + err.Error(),
		})
		return
	}

	// 验证必需参数
	if req.Query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "搜索关键词不能为空",
		})
		return
	}

	// 调用服务层获取搜索建议
	result, err := c.factorySearchService.GetSearchSuggestions(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "获取搜索建议失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// CreateFactorySpecialty 创建工厂专业领域
// @Summary 创建工厂专业领域
// @Description 为工厂添加专业领域标签
// @Tags 工厂搜索
// @Accept json
// @Produce json
// @Param factory_id path int true "工厂ID"
// @Param specialty body map[string]string true "专业领域信息"
// @Success 201 {object} map[string]interface{}
// @Router /api/factories/{factory_id}/specialties [post]
func (c *FactorySearchController) CreateFactorySpecialty(ctx *gin.Context) {
	// 获取工厂ID
	factoryIDStr := ctx.Param("factory_id")
	if factoryIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "工厂ID不能为空",
		})
		return
	}

	// 解析工厂ID
	var factoryID uint
	if _, err := fmt.Sscanf(factoryIDStr, "%d", &factoryID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的工厂ID",
		})
		return
	}

	// 获取专业领域信息
	var req struct {
		Specialty string `json:"specialty" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "参数绑定失败: " + err.Error(),
		})
		return
	}

	// 调用服务层创建专业领域
	err := c.factorySearchService.CreateFactorySpecialty(factoryID, req.Specialty)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "创建专业领域失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "专业领域创建成功",
	})
}

// CreateFactoryRating 创建工厂评分
// @Summary 创建工厂评分
// @Description 为工厂添加评分和评价
// @Tags 工厂搜索
// @Accept json
// @Produce json
// @Param factory_id path int true "工厂ID"
// @Param rating body map[string]interface{} true "评分信息"
// @Success 201 {object} map[string]interface{}
// @Router /api/factories/{factory_id}/ratings [post]
func (c *FactorySearchController) CreateFactoryRating(ctx *gin.Context) {
	// 获取工厂ID
	factoryIDStr := ctx.Param("factory_id")
	if factoryIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "工厂ID不能为空",
		})
		return
	}

	// 解析工厂ID
	var factoryID uint
	if _, err := fmt.Sscanf(factoryIDStr, "%d", &factoryID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的工厂ID",
		})
		return
	}

	// 获取评分信息
	var req struct {
		Rating  float64 `json:"rating" binding:"required,min=0,max=5"`
		Comment string  `json:"comment"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "参数绑定失败: " + err.Error(),
		})
		return
	}

	// 获取当前用户ID
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "用户未认证",
		})
		return
	}

	// 调用服务层创建评分
	err := c.factorySearchService.CreateFactoryRating(factoryID, req.Rating, req.Comment, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "创建评分失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "评分创建成功",
	})
} 