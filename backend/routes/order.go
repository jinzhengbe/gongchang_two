package routes

import (
	"gongChang/controllers"
	"gongChang/middleware"

	"github.com/gin-gonic/gin"
)

func SetupOrderRoutes(router *gin.Engine, orderController *controllers.OrderController) {
	orderGroup := router.Group("/api/orders")
	{
		// 公开路由（已由 router.go 注册，这里注释掉，避免重复和认证冲突）
		// orderGroup.GET("/recent", orderController.GetRecentOrders)

		// 需要认证的路由
		authGroup := orderGroup.Group("")
		authGroup.Use(middleware.AuthMiddleware())
		{
			// 创建订单
			authGroup.POST("/", orderController.CreateOrder)
			
			// 获取用户的所有订单
			authGroup.GET("/user/:userID", orderController.GetOrdersByUserID)
			
			// 获取特定订单
			authGroup.GET("/:id", orderController.GetOrderByID)
			
			// 更新订单状态
			authGroup.PUT("/:id/status", orderController.UpdateOrderStatus)
			
			// 搜索订单
			authGroup.GET("/search", orderController.SearchOrders)
			
			// 获取订单统计信息
			authGroup.GET("/statistics", orderController.GetOrderStatistics)
		}
	}
} 