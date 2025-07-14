package controllers

import (
	"fmt"
	"gongChang/config"
	"gongChang/middleware"
	"gongChang/models"
	"gongChang/services"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"log"
	"errors"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// Register a new user
func (uc *UserController) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// 统一注册服务
	err := uc.userService.Register(req)
	if err != nil {
		if errors.Is(err, services.ErrUsernameExists) || errors.Is(err, services.ErrEmailExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login a user
func (uc *UserController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// 统一登录服务
	user, profile, err := uc.userService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 加载配置以获取JWT密钥
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Failed to load configuration for JWT: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// 生成JWT token
	token, err := middleware.GenerateToken(user.ID, user.Role, cfg.JWT.Secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user":    user,
		"profile": profile,
	})
}

func (c *UserController) GetUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	user, err := c.userService.GetUserByID(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = userID
	if err := c.userService.UpdateUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	if err := c.userService.DeleteUser(userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetUserProfile 获取当前用户信息
func (c *UserController) GetUserProfile(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	user, err := c.userService.GetUserByID(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// UpdateUserProfile 更新当前用户信息
func (c *UserController) UpdateUserProfile(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var req models.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	user, err := c.userService.GetUserByID(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 更新用户信息
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := c.userService.UpdateUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "用户信息更新成功",
		"user":    user,
	})
}

// RefreshToken 刷新Token
func (c *UserController) RefreshToken(ctx *gin.Context) {
	// 从请求头获取token
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未提供token"})
		return
	}

	// 去掉Bearer前缀
	if strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}

	// 验证token
	claims, err := middleware.ValidateToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token无效"})
		return
	}

	// 生成新token
	cfg, err := config.LoadConfig()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "加载配置失败"})
		return
	}
	newToken, err := middleware.GenerateToken(claims.UserID, claims.Role, cfg.JWT.Secret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": newToken,
	})
}

// UploadAvatar 上传头像
func (c *UserController) UploadAvatar(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取上传的文件
	file, header, err := ctx.Request.FormFile("avatar")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
		return
	}
	defer file.Close()

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(header.Filename))
	supportedExts := []string{".jpg", ".jpeg", ".png", ".webp"}
	supported := false
	for _, supportedExt := range supportedExts {
		if ext == supportedExt {
			supported = true
			break
		}
	}
	if !supported {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件格式，支持: JPG, PNG, WebP"})
		return
	}

	// 检查文件大小（限制5MB）
	if header.Size > 5*1024*1024 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "文件大小超过限制 (最大 5MB)"})
		return
	}

	// 生成唯一文件名
	fileID := uuid.New().String()
	newFilename := fileID + ext

	// 确保上传目录存在
	uploadDir := "./uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建上传目录失败"})
		return
	}

	// 保存文件
	finalPath := filepath.Join(uploadDir, newFilename)
	dst, err := os.Create(finalPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建文件失败"})
		return
	}
	defer dst.Close()

	// 复制文件内容
	_, err = io.Copy(dst, file)
	if err != nil {
		os.Remove(finalPath) // 清理失败的文件
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	// 构建文件访问URL
	fileURL := fmt.Sprintf("/uploads/avatars/%s", newFilename)

	// 返回成功响应
	response := models.UploadAvatarResponse{
		Success: true,
		Message: "头像上传成功",
		Data: struct {
			URL string `json:"url"`
		}{
			URL: fileURL,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateDesignerProfile 更新设计师信息
func (c *UserController) UpdateDesignerProfile(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var req models.UpdateDesignerProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 获取设计师档案
	var designerProfile models.DesignerProfile
	err := c.userService.GetDB().Where("user_id = ?", userID).First(&designerProfile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "设计师档案不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	// 更新设计师信息
	updates := make(map[string]interface{})

	if req.CompanyName != "" {
		updates["company_name"] = req.CompanyName
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.Website != "" {
		updates["website"] = req.Website
	}
	if req.Bio != "" {
		updates["bio"] = req.Bio
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	// 执行更新
	if err := c.userService.GetDB().Model(&designerProfile).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败: " + err.Error()})
		return
	}

	// 获取更新后的完整信息
	var updatedProfile models.DesignerProfile
	c.userService.GetDB().Where("user_id = ?", userID).First(&updatedProfile)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "设计师信息更新成功",
		"data":    updatedProfile,
	})
}

// GetDesignerProfile 获取设计师信息
func (c *UserController) GetDesignerProfile(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var designerProfile models.DesignerProfile
	err := c.userService.GetDB().Where("user_id = ?", userID).First(&designerProfile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "设计师档案不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	// 获取关联的用户信息
	var user models.User
	c.userService.GetDB().Where("id = ?", userID).First(&user)
	designerProfile.User = user

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    designerProfile,
	})
}

// ChangePassword 修改密码
// @Summary 修改用户密码
// @Description 用户修改自己的密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body models.ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} gin.H
// @Router /api/users/change-password [post]
func (c *UserController) ChangePassword(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var req models.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 验证新密码长度
	if len(req.NewPassword) < 6 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "新密码长度不能少于6位"})
		return
	}

	// 调用服务层修改密码
	err := c.userService.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "密码修改成功",
	})
} 