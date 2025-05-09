package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"aneworder.com/backend/controllers"
	"aneworder.com/backend/services"
)

func RegisterPublicRoutes(r *gin.Engine, db *gorm.DB) {
	// 用户相关路由
	userService := services.NewUserService(db)
	userController := controllers.NewUserController(userService)

	r.POST("/api/register", userController.Register)

	// 订单相关路由
	public := r.Group("/public")
	{
		orderController := controllers.NewPublicOrderController(db)
		public.GET("/orders", orderController.GetPublicOrders)
	}
} 