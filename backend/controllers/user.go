package controllers

import (
	"gongChang/config"
	"gongChang/middleware"
	"gongChang/models"
	"gongChang/services"
	"net/http"

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
	// 从请求头获取当前token
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "缺少认证token"})
		return
	}

	// 移除Bearer前缀
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// 验证token并获取用户信息
	claims, err := middleware.ValidateToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
		return
	}

	// 加载配置以获取JWT密钥
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Failed to load configuration for JWT: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "无法生成新token"})
		return
	}

	// 生成新的JWT token
	newToken, err := middleware.GenerateToken(claims.UserID, claims.Role, cfg.JWT.Secret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "无法生成新token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Token刷新成功",
		"token":   newToken,
	})
} 