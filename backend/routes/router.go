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
	fileService := services.NewFileService(db, "./uploads")

	// 创建控制器实例
	userController := controllers.NewUserController(userService)
	productController := controllers.NewProductController(productService)
	orderController := controllers.NewOrderController(orderService)
	fileController := controllers.NewFileController(fileService, "./uploads")

	// API 路由组
	api := r.Group("/api")
	{
		// 用户路由
		userGroup := api.Group("/users")
		{
			userGroup.POST("/register", userController.Register)
			userGroup.POST("/login", userController.Login)
		}

		// 需要认证的路由
		authGroup := api.Group("")
		authGroup.Use(middleware.AuthMiddleware())
		{
			// 用户路由
			userGroup := authGroup.Group("/users")
			{
				userGroup.GET("/:id", userController.GetUser)
				userGroup.PUT("/:id", userController.UpdateUser)
				userGroup.DELETE("/:id", userController.DeleteUser)
			}

			// 产品路由
			productGroup := authGroup.Group("/products")
			{
				productGroup.POST("", productController.CreateProduct)
				productGroup.GET("", productController.GetProducts)
				productGroup.GET("/:id", productController.GetProduct)
				productGroup.PUT("/:id", productController.UpdateProduct)
				productGroup.DELETE("/:id", productController.DeleteProduct)
			}

			// 订单路由
			orderGroup := authGroup.Group("/orders")
			{
				orderGroup.POST("", orderController.CreateOrder)
				orderGroup.GET("", orderController.GetOrdersByUserID)
				orderGroup.GET("/:id", orderController.GetOrderByID)
				orderGroup.PUT("/:id/status", orderController.UpdateOrderStatus)
				orderGroup.GET("/search", orderController.SearchOrders)
				orderGroup.GET("/statistics", orderController.GetOrderStatistics)
				orderGroup.GET("/recent", orderController.GetRecentOrders)
			}

			// 文件路由
			fileGroup := authGroup.Group("/files")
			{
				fileGroup.POST("/upload", fileController.UploadFile)
				fileGroup.POST("/batch", fileController.GetBatchFileDetails)
				fileGroup.GET("/:id", fileController.GetFileDetails)
				fileGroup.GET("/download/:id", fileController.DownloadFile)
				fileGroup.DELETE("/:id", fileController.DeleteFile)
				fileGroup.GET("/order/:id", fileController.GetOrderFiles)
			}
		}
	}

	return r
} 