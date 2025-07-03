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

// GetFactoryByUserID 根据用户ID获取工厂信息
func (fc *FactoryController) GetFactoryByUserID(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "用户ID不能为空",
		})
		return
	}

	var factory models.FactoryProfile
	err := fc.DB.Where("user_id = ?", userID).First(&factory).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "工厂不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "服务器内部错误",
		})
		return
	}

	// 获取关联的用户信息
	var user models.User
	fc.DB.Where("id = ?", factory.UserID).First(&user)
	factory.User = user

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": factory,
	})
}

// GetFactoryByID 根据工厂ID获取工厂详情
func (fc *FactoryController) GetFactoryByID(c *gin.Context) {
	factoryID := c.Param("id")
	if factoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "工厂ID不能为空",
		})
		return
	}

	// 将字符串ID转换为uint
	id, err := strconv.ParseUint(factoryID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的工厂ID",
		})
		return
	}

	var factory models.FactoryProfile
	err = fc.DB.Where("id = ?", uint(id)).First(&factory).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "工厂不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "服务器内部错误",
		})
		return
	}

	// 获取关联的用户信息
	var user models.User
	fc.DB.Where("id = ?", factory.UserID).First(&user)
	factory.User = user

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": factory,
	})
} 