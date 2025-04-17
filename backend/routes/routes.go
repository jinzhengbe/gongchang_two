package routes

import (
	"backend/controllers"
	"backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// 初始化服务
	userService := services.NewUserService(db)

	// 初始化控制器
	userController := controllers.NewUserController(userService)

	// 健康检查
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "API is running",
		})
	})

	// 用户认证路由
	auth := router.Group("/api/auth")
	{
		auth.POST("/login", userController.Login)
		auth.POST("/register", userController.Register)
	}

	return router
} 