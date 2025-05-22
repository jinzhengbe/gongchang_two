package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupOrderRoutes(router *gin.Engine, orderController *controllers.OrderController) {
	orderGroup := router.Group("/api/orders")
	{
		// 所有订单路由都需要认证
		orderGroup.Use(middleware.AuthMiddleware())
		
		// 创建订单
		orderGroup.POST("/", orderController.CreateOrder)
		
		// 获取用户的所有订单
		orderGroup.GET("/user/:userID", orderController.GetOrdersByUserID)
		
		// 获取特定订单
		orderGroup.GET("/:id", orderController.GetOrderByID)
		
		// 更新订单状态
		orderGroup.PUT("/:id/status", orderController.UpdateOrderStatus)
		
		// 搜索订单
		orderGroup.GET("/search", orderController.SearchOrders)
		
		// 获取订单统计信息
		orderGroup.GET("/statistics", orderController.GetOrderStatistics)
		
		// 获取最近订单
		orderGroup.GET("/recent", orderController.GetRecentOrders)
	}
} 