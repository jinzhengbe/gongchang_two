package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gongChang/controllers"
	"gongChang/services"
	"gongChang/middleware"
	"gongChang/config"
	"net/http"
	"log"
	"gongChang/internal/factory"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// 设置受信任的代理
	r.SetTrustedProxies(cfg.Server.TrustedProxies)

	// 添加健康检查路由
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 添加 CORS 中间件
	r.Use(middleware.CORSMiddleware())

	// 创建服务实例
	userService := services.NewUserService(db)
	productService := services.NewProductService(db)
	orderService := services.NewOrderService(db)
	fileService := services.NewFileService(db, "./uploads")
	factoryService := factory.NewService(db)

	// 创建控制器实例
	userController := controllers.NewUserController(userService)
	productController := controllers.NewProductController(productService)
	orderController := controllers.NewOrderController(orderService)
	fileController := controllers.NewFileController(fileService, "./uploads")
	factoryHandler := factory.NewHandler(factoryService, nil)
	factoryController := factory.NewController(factoryService)

	// API 路由组
	api := r.Group("/api")
	{
		// 工厂路由
		factoryGroup := api.Group("/factory")
		{
			factoryGroup.POST("/register", factoryHandler.Register)
			factoryGroup.POST("/login", factoryHandler.Login)
		}

		// 工厂清单路由
		api.GET("/factories", factoryController.GetFactoryList)

		// 获取最近订单（公开路由）
		api.GET("/orders/recent", orderController.GetRecentOrders)

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
				orderGroup.GET("/:id", orderController.GetOrderByID)
				orderGroup.PUT("/:id/status", orderController.UpdateOrderStatus)
				orderGroup.GET("/search", orderController.SearchOrders)
				orderGroup.GET("/statistics", orderController.GetOrderStatistics)
			}

			// 工厂订单路由
			factoryGroup := authGroup.Group("/factory")
			{
				factoryGroup.GET("/orders", orderController.GetOrdersByUserID)
				factoryGroup.PUT("/orders/:id", orderController.UpdateOrderStatus)
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

		// 设计师订单路由（单独注册，确保一定生效）
		designerGroup := api.Group("/designer")
		designerGroup.Use(middleware.AuthMiddleware())
		{
			designerGroup.GET("/orders", orderController.GetOrdersByDesignerID)
			designerGroup.POST("/orders", orderController.CreateOrder)
		}
		log.Println("!!! DESIGNER ROUTE REGISTERED !!!")
	}

	// 注册公开路由
	RegisterPublicRoutes(r, db)

	// 打印所有已注册的路由
	for _, route := range r.Routes() {
		log.Printf("[ROUTE] %s %s", route.Method, route.Path)
	}

	return r
} 