package factory

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/your-project/models"
	"github.com/your-project/services"
)

type Handler struct {
	service *Service
	userService *UserService
}

func NewHandler(service *Service, userService *UserService) *Handler {
	return &Handler{service: service, userService: userService}
}

// Register 工厂注册
// @Summary 工厂注册
// @Description 注册新的工厂账号
// @Tags 工厂管理
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册信息"
// @Success 200 {object} Response
// @Router /api/factory/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建用户对象，注册到 users 表
	user := &models.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role:     models.RoleFactory,
	}

	// 使用 UserService 注册用户
	if err := h.userService.Register(user); err != nil {
		if err == services.ErrUsernameExists {
			c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// Login 工厂登录
// @Summary 工厂登录
// @Description 工厂账号登录
// @Tags 工厂管理
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} LoginResponse
// @Router /api/factory/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Login(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
} 