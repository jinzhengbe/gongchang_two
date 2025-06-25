package factory

import (
	"github.com/gin-gonic/gin"
	"gongChang/services"
	"net/http"
)

type Handler struct {
	service *Service
	userService *services.UserService
}

func NewHandler(service *Service, userService *services.UserService) *Handler {
	return &Handler{service: service, userService: userService}
}

/*
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
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 确保角色是 factory
	req.Role = "factory"

	// 使用 UserService 注册用户
	if err := h.userService.Register(req); err != nil {
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
*/

// GetFactoryOrders 获取工厂的订单列表
// @Summary 获取工厂的订单列表
// @Description 获取当前登录工厂的订单列表
// @Tags 工厂管理
// @Accept json
// @Produce json
// @Param request query OrderListRequest true "查询参数"
// @Success 200 {object} OrderListResponse
// @Router /api/factory/orders [get]
func (h *Handler) GetFactoryOrders(c *gin.Context) {
	// 获取当前登录的工厂ID
	factoryID, exists := c.Get("factory_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或登录已过期"})
		return
	}

	// 解析查询参数
	var req OrderListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取订单列表
	resp, err := h.service.GetFactoryOrders(factoryID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetDesignerOrders 获取设计师的订单列表
// @Summary 获取设计师的订单列表
// @Description 获取当前登录设计师的订单列表
// @Tags 设计师管理
// @Accept json
// @Produce json
// @Param request query OrderListRequest true "查询参数"
// @Success 200 {object} OrderListResponse
// @Router /api/designer/orders [get]
func (h *Handler) GetDesignerOrders(c *gin.Context) {
	// 获取当前登录的设计师ID
	designerID, exists := c.Get("designer_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或登录已过期"})
		return
	}

	// 解析查询参数
	var req OrderListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取订单列表
	resp, err := h.service.GetDesignerOrders(designerID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
} 