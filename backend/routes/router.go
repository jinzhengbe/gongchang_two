package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gongChang/controllers"
	"gongChang/services"
	"gongChang/middleware"
	"gongChang/config"
	"net/http"
	"strings"
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

	// 添加静态文件服务，专门用于提供上传的文件
	r.Static("/uploads", "./uploads")
	
	// 为静态文件添加CORS头
	r.Use(func(c *gin.Context) {
		if c.Request.URL.Path == "/uploads" || strings.HasPrefix(c.Request.URL.Path, "/uploads/") {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Max-Age", "86400")
		}
		c.Next()
	})

	// 创建服务实例
	userService := services.NewUserService(db)
	productService := services.NewProductService(db)
	orderService := services.NewOrderService(db)
	fileService := services.NewFileService(db, "./uploads")
	fabricService := services.NewFabricService(db)

	// 创建控制器实例
	userController := controllers.NewUserController(userService)
	productController := controllers.NewProductController(productService)
	orderController := controllers.NewOrderController(orderService, db)
	fileController := controllers.NewFileController(fileService, "./uploads", cfg)
	factoryController := &controllers.FactoryController{DB: db}
	fabricController := controllers.NewFabricController(fabricService)

	// API 路由组
	api := r.Group("/api")
	{
		// 认证相关路由（无需认证）
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", userController.Login)
			authGroup.POST("/register", userController.Register)
			authGroup.POST("/refresh", userController.RefreshToken)
		}

		// 公开路由（无需认证）
		publicGroup := api.Group("/public")
		{
			publicGroup.GET("/orders", orderController.GetPublicOrders)
		}

		// 布料公开路由（无需认证）
		fabricPublicGroup := api.Group("/fabrics")
		{
			fabricPublicGroup.GET("/all", fabricController.GetAllFabrics)
			fabricPublicGroup.GET("/categories", fabricController.GetFabricCategories)
			fabricPublicGroup.GET("/search", fabricController.SearchFabrics)
			fabricPublicGroup.GET("/category/:category", fabricController.GetFabricsByCategory)
			fabricPublicGroup.GET("/material/:material", fabricController.GetFabricsByMaterial)
			fabricPublicGroup.GET("/:id", fabricController.GetFabricByID)
			fabricPublicGroup.GET("/statistics", fabricController.GetFabricStatistics)
		}

		// 工厂列表路由（公开）
		api.GET("/factories", factoryController.GetFactoryList)

		// 获取最近订单（公开路由）
		api.GET("/orders/recent", orderController.GetRecentOrders)

		// 需要认证的路由
		authRequiredGroup := api.Group("")
		authRequiredGroup.Use(middleware.AuthMiddleware())
		{
			// 用户管理路由
			userGroup := authRequiredGroup.Group("/users")
			{
				userGroup.GET("/profile", userController.GetUserProfile)
				userGroup.PUT("/profile", userController.UpdateUserProfile)
				userGroup.GET("/:id", userController.GetUser)
				userGroup.PUT("/:id", userController.UpdateUser)
				userGroup.DELETE("/:id", userController.DeleteUser)
			}

			// 产品路由
			productGroup := authRequiredGroup.Group("/products")
			{
				productGroup.POST("", productController.CreateProduct)
				productGroup.GET("", productController.GetProducts)
				productGroup.GET("/:id", productController.GetProduct)
				productGroup.PUT("/:id", productController.UpdateProduct)
				productGroup.DELETE("/:id", productController.DeleteProduct)
			}

			// 订单路由（完整的CRUD操作）
			orderGroup := authRequiredGroup.Group("/orders")
			{
				orderGroup.POST("", orderController.CreateOrder)
				orderGroup.GET("", orderController.GetOrdersByUserID)
				orderGroup.GET("/:id", orderController.GetOrderByID)
				orderGroup.PUT("/:id", orderController.UpdateOrder)
				orderGroup.DELETE("/:id", orderController.DeleteOrder)
				orderGroup.PUT("/:id/status", orderController.UpdateOrderStatus)
				orderGroup.GET("/search", orderController.SearchOrders)
				orderGroup.GET("/statistics", orderController.GetOrderStatistics)
			}

			// 工厂订单路由
			factoryGroup := authRequiredGroup.Group("/factory")
			{
				factoryGroup.GET("/orders", orderController.GetOrdersByUserID)
				factoryGroup.PUT("/orders/:id", orderController.UpdateOrderStatus)
			}

			// 设计师订单路由
			designerGroup := authRequiredGroup.Group("/designer")
			{
				designerGroup.GET("/orders", orderController.GetOrdersByDesignerID)
				designerGroup.POST("/orders", orderController.CreateOrder)
			}

			// 文件路由
			fileGroup := authRequiredGroup.Group("/files")
			{
				fileGroup.POST("/upload", fileController.UploadFile)
				fileGroup.POST("/batch", fileController.GetBatchFileDetails)
				fileGroup.GET("/:id", fileController.GetFileDetails)
				fileGroup.GET("/download/:id", fileController.DownloadFile)
				fileGroup.DELETE("/:id", fileController.DeleteFile)
				fileGroup.GET("/order/:id", fileController.GetOrderFiles)
			}

			// 布料管理路由（需要认证）
			fabricGroup := authRequiredGroup.Group("/fabrics")
			{
				fabricGroup.POST("", fabricController.CreateFabric)
				fabricGroup.PUT("/:id", fabricController.UpdateFabric)
				fabricGroup.DELETE("/:id", fabricController.DeleteFabric)
				fabricGroup.PUT("/:id/stock", fabricController.UpdateFabricStock)
			}
		}
	}

	// 注册公开路由
	RegisterPublicRoutes(r, db)

	return r
} 