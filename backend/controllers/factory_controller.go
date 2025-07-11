package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gongChang/models"
	"gongChang/services"
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

// GetFactoryProfile 获取当前用户的工厂详细信息
func (fc *FactoryController) GetFactoryProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未授权",
		})
		return
	}

	var factory models.FactoryProfile
	err := fc.DB.Where("user_id = ?", userID).First(&factory).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "工厂信息不存在",
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

// UpdateFactoryProfile 更新工厂详细信息
func (fc *FactoryController) UpdateFactoryProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未授权",
		})
		return
	}

	var req models.UpdateFactoryProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	var factory models.FactoryProfile
	err := fc.DB.Where("user_id = ?", userID).First(&factory).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "工厂信息不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "服务器内部错误",
		})
		return
	}

	// 更新工厂信息
	updates := make(map[string]interface{})

	if req.CompanyName != "" {
		updates["company_name"] = req.CompanyName
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.Capacity != nil {
		updates["capacity"] = *req.Capacity
	}
	if req.Equipment != "" {
		updates["equipment"] = req.Equipment
	}
	if req.Certificates != "" {
		updates["certificates"] = req.Certificates
	}
	if req.EmployeeCount != nil {
		updates["employee_count"] = *req.EmployeeCount
	}
	if req.Photos != nil {
		// 将照片数组转换为JSON字符串存储
		photosJSON, err := json.Marshal(req.Photos)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "照片数据格式错误",
			})
			return
		}
		updates["photos"] = string(photosJSON)
	}
	
	if req.Videos != nil {
		// 将视频数组转换为JSON字符串存储
		videosJSON, err := json.Marshal(req.Videos)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "视频数据格式错误",
			})
			return
		}
		updates["videos"] = string(videosJSON)
	}

	updates["updated_at"] = time.Now()

	if err := fc.DB.Model(&factory).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新失败: " + err.Error(),
		})
		return
	}

	// 重新查询获取更新后的数据
	fc.DB.Where("user_id = ?", userID).First(&factory)
	var user models.User
	fc.DB.Where("id = ?", factory.UserID).First(&user)
	factory.User = user

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "工厂信息更新成功",
		"data": factory,
	})
} 

// 工厂图片管理相关方法

// BatchUploadPhotos 批量上传工厂图片
func (fc *FactoryController) BatchUploadPhotos(c *gin.Context) {
	// 获取工厂ID
	factoryID := c.Param("factory_id")
	if factoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "工厂ID不能为空",
		})
		return
	}

	// 验证用户权限（只能给自己的工厂上传图片）
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "未授权",
		})
		return
	}

	// 验证用户是否有权限操作此工厂
	if userID != factoryID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "无权限操作此工厂",
		})
		return
	}

	// 获取上传的文件
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "文件格式错误",
		})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "请选择要上传的图片",
		})
		return
	}

	// 获取分类参数
	category := c.PostForm("category")

	// 调用服务层处理批量上传
	fileService := services.NewFileService(fc.DB, "./uploads")
	response, err := fileService.BatchUploadFactoryPhotos(files, factoryID, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetFactoryPhotos 获取工厂图片列表
func (fc *FactoryController) GetFactoryPhotos(c *gin.Context) {
	// 获取工厂ID
	factoryID := c.Param("factory_id")
	if factoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "工厂ID不能为空",
		})
		return
	}

	// 验证用户权限
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "未授权",
		})
		return
	}

	// 验证用户是否有权限查看此工厂
	if userID != factoryID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "无权限查看此工厂",
		})
		return
	}

	// 获取查询参数
	category := c.Query("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 调用服务层获取图片列表
	fileService := services.NewFileService(fc.DB, "./uploads")
	response, err := fileService.GetFactoryPhotos(factoryID, category, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteFactoryPhoto 删除单张工厂图片
func (fc *FactoryController) DeleteFactoryPhoto(c *gin.Context) {
	// 获取工厂ID和图片ID
	factoryID := c.Param("factory_id")
	photoID := c.Param("photoId")
	
	if factoryID == "" || photoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "工厂ID和图片ID不能为空",
		})
		return
	}

	// 验证用户权限
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "未授权",
		})
		return
	}

	// 验证用户是否有权限操作此工厂
	if userID != factoryID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "无权限操作此工厂",
		})
		return
	}

	// 调用服务层删除图片
	fileService := services.NewFileService(fc.DB, "./uploads")
	err := fileService.DeleteFactoryPhoto(photoID, factoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "图片删除成功",
	})
}

// BatchDeletePhotos 批量删除工厂图片
func (fc *FactoryController) BatchDeletePhotos(c *gin.Context) {
	// 获取工厂ID
	factoryID := c.Param("factory_id")
	if factoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "工厂ID不能为空",
		})
		return
	}

	// 验证用户权限
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "未授权",
		})
		return
	}

	// 验证用户是否有权限操作此工厂
	if userID != factoryID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "无权限操作此工厂",
		})
		return
	}

	// 解析请求体
	var req models.BatchDeleteFactoryPhotosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "请求参数错误: " + err.Error(),
		})
		return
	}

	// 调用服务层批量删除图片
	fileService := services.NewFileService(fc.DB, "./uploads")
	response, err := fileService.BatchDeleteFactoryPhotos(req.PhotoIDs, factoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
} 