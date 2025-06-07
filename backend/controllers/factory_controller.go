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

	query := fc.DB.Model(&models.FactoryProfile{})

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取分页数据
	var factories []models.FactoryProfile
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Find(&factories)

	// 获取关联的用户信息
	for i := range factories {
		var user models.User
		fc.DB.Where("id = ?", factories[i].UserID).First(&user)
		factories[i].User = user
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"total": total,
			"factories": factories,
		},
	})
} 