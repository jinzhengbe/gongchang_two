package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gongChang/models"
)

type FactoryController struct {
	DB *gorm.DB
}

// GetFactoryList 获取工厂列表
func (fc *FactoryController) GetFactoryList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))

	query := fc.DB.Model(&models.Factory{})

	// 如果指定了状态，添加状态过滤
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取分页数据
	var factories []models.Factory
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Find(&factories)

	response := models.FactoryListResponse{
		Total:     total,
		Factories: factories,
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": response,
	})
} 