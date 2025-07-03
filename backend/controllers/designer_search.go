package controllers

import (
	"fmt"
	"net/http"
	"gongChang/models"
	"gongChang/services"

	"github.com/gin-gonic/gin"
)

type DesignerSearchController struct {
	designerSearchService *services.DesignerSearchService
}

func NewDesignerSearchController(designerSearchService *services.DesignerSearchService) *DesignerSearchController {
	return &DesignerSearchController{
		designerSearchService: designerSearchService,
	}
}

// SearchDesigners 搜索设计师
// @Summary 搜索设计师
// @Description 搜索设计师，支持关键词搜索、地区筛选、专业领域筛选等
// @Tags 设计师搜索
// @Accept json
// @Produce json
// @Param query query string false "搜索关键词"
// @Param region query string false "地区筛选"
// @Param specialties query []string false "专业领域数组"
// @Param min_rating query number false "最低评分"
// @Param max_rating query number false "最高评分"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param sort_by query string false "排序字段" default(rating)
// @Param sort_order query string false "排序方向" default(desc)
// @Success 200 {object} models.DesignerSearchResponse
// @Router /api/designers/search [get]
func (c *DesignerSearchController) SearchDesigners(ctx *gin.Context) {
	var req models.DesignerSearchRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "参数绑定失败: " + err.Error(),
		})
		return
	}

	response, err := c.designerSearchService.SearchDesigners(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "搜索设计师失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetSearchSuggestions 获取搜索建议
// @Summary 获取设计师搜索建议
// @Description 根据关键词获取设计师搜索建议
// @Tags 设计师搜索
// @Accept json
// @Produce json
// @Param query query string true "搜索关键词"
// @Param limit query int false "建议数量" default(10)
// @Success 200 {object} models.DesignerSearchSuggestionResponse
// @Router /api/designers/search/suggestions [get]
func (c *DesignerSearchController) GetSearchSuggestions(ctx *gin.Context) {
	var req models.DesignerSearchSuggestionRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "参数绑定失败: " + err.Error(),
		})
		return
	}

	response, err := c.designerSearchService.GetSearchSuggestions(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "获取搜索建议失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// CreateDesignerSpecialty 创建设计师专业领域
// @Summary 创建设计师专业领域
// @Description 为设计师添加专业领域标签
// @Tags 设计师搜索
// @Accept json
// @Produce json
// @Param designer_id path int true "设计师ID"
// @Param specialty body map[string]string true "专业领域信息"
// @Success 201 {object} map[string]interface{}
// @Router /api/designers/{designer_id}/specialties [post]
func (c *DesignerSearchController) CreateDesignerSpecialty(ctx *gin.Context) {
	// 获取设计师ID
	designerIDStr := ctx.Param("designer_id")
	if designerIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "设计师ID不能为空",
		})
		return
	}

	// 解析设计师ID
	var designerID uint
	if _, err := fmt.Sscanf(designerIDStr, "%d", &designerID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的设计师ID",
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
	err := c.designerSearchService.CreateDesignerSpecialty(designerID, req.Specialty)
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

// CreateDesignerRating 创建设计师评分
// @Summary 创建设计师评分
// @Description 为设计师添加评分和评价
// @Tags 设计师搜索
// @Accept json
// @Produce json
// @Param designer_id path int true "设计师ID"
// @Param rating body map[string]interface{} true "评分信息"
// @Success 201 {object} map[string]interface{}
// @Router /api/designers/{designer_id}/ratings [post]
func (c *DesignerSearchController) CreateDesignerRating(ctx *gin.Context) {
	// 获取设计师ID
	designerIDStr := ctx.Param("designer_id")
	if designerIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "设计师ID不能为空",
		})
		return
	}

	// 解析设计师ID
	var designerID uint
	if _, err := fmt.Sscanf(designerIDStr, "%d", &designerID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的设计师ID",
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
	err := c.designerSearchService.CreateDesignerRating(designerID, req.Rating, req.Comment, userID)
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