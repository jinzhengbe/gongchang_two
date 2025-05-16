package controllers

import (
	"aneworder.com/backend/config"
	"aneworder.com/backend/middleware"
	"aneworder.com/backend/models"
	"aneworder.com/backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"log"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) Register(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建用户对象
	user := &models.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role:     models.UserRole(req.Role),
	}

	// 使用 UserService 注册用户
	if err := c.userService.Register(user); err != nil {
		if err == services.ErrUsernameExists {
			ctx.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func (c *UserController) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Login request binding error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Login attempt for user: %s, type: %s", req.Username, req.UserType)

	// 验证用户类型
	if req.UserType != string(models.RoleDesigner) && 
	   req.UserType != string(models.RoleFactory) && 
	   req.UserType != string(models.RoleSupplier) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user type"})
		return
	}

	user, err := c.userService.Login(req.Username, req.Password)
	if err != nil {
		log.Printf("Login failed for user %s: %v", req.Username, err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 验证用户类型是否匹配
	if string(user.Data.User.Role) != req.UserType {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user type mismatch"})
		return
	}

	// 生成JWT token
	cfg, err := config.LoadConfig()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load configuration"})
		return
	}

	token, err := middleware.GenerateToken(user.Data.User.ID, user.Data.User.Role, cfg.JWT.Secret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.Data.User.ID,
			"username": user.Data.User.Username,
			"email":    user.Data.User.Email,
			"role":     user.Data.User.Role,
		},
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