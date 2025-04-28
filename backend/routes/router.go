package routes

import (
	"aneworder.com/backend/controllers"
	"aneworder.com/backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"aneworder.com/backend/middleware"
	"net/http"
	"aneworder.com/backend/config"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// 设置受信任的代理
	r.SetTrustedProxies(cfg.Server.TrustedProxies)

	// 添加健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 添加 CORS 中间件
	r.Use(middleware.CORSMiddleware())

	// 创建服务实例
	userService := services.NewUserService(db)
	productService := services.NewProductService(db)
	orderService := services.NewOrderService(db)

	// 创建控制器实例
	userController := controllers.NewUserController(userService)
	productController := controllers.NewProductController(productService)
	orderController := controllers.NewOrderController(orderService)

	// API 路由组
	api := r.Group("/api")
	{
		// 用户相关路由
		api.POST("/users/register", userController.Register)
		api.POST("/users/login", userController.Login)
		api.GET("/users/:id", userController.GetUser)
		api.PUT("/users/:id", userController.UpdateUser)
		api.DELETE("/users/:id", userController.DeleteUser)

		// 产品相关路由
		api.GET("/products", productController.GetProducts)
		api.GET("/products/:id", productController.GetProduct)
		api.POST("/products", productController.CreateProduct)
		api.PUT("/products/:id", productController.UpdateProduct)
		api.DELETE("/products/:id", productController.DeleteProduct)

		// 订单相关路由
		api.POST("/orders", middleware.AuthMiddleware(), orderController.CreateOrder)
		api.GET("/orders/recent", middleware.AuthMiddleware(), orderController.GetRecentOrders)
		api.GET("/orders/latest", middleware.AuthMiddleware(), orderController.GetLatestOrders)
		api.GET("/orders/hot", middleware.AuthMiddleware(), orderController.GetHotOrders)
		api.GET("/orders/:id", middleware.AuthMiddleware(), orderController.GetOrderByID)
		api.PUT("/orders/:id/status", middleware.AuthMiddleware(), orderController.UpdateOrderStatus)
		api.GET("/orders", middleware.AuthMiddleware(), orderController.GetOrdersByUserID)
	}

	return r
} 