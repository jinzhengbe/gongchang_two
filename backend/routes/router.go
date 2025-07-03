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
	jiedanService := services.NewJiedanService(db)
	progressService := services.NewProgressService(db)
	employeeService := services.NewEmployeeService(db)
	orderSearchService := services.NewOrderSearchService(db)

	// 创建控制器实例
	userController := controllers.NewUserController(userService)
	productController := controllers.NewProductController(productService)
	orderController := controllers.NewOrderController(orderService, db)
	fileController := controllers.NewFileController(fileService, "./uploads", cfg)
	factoryController := &controllers.FactoryController{DB: db}
	fabricController := controllers.NewFabricController(fabricService)
	jiedanController := controllers.NewJiedanController(jiedanService)
	progressController := controllers.NewProgressController(progressService)
	employeeController := controllers.NewEmployeeController(employeeService)
	orderSearchController := controllers.NewOrderSearchController(orderSearchService)

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
		// 根据用户ID获取单个工厂信息（公开）
		api.GET("/factories/user/:userId", factoryController.GetFactoryByUserID)

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

			// 订单搜索路由（高级搜索功能）- 使用不同的路径避免冲突
			authRequiredGroup.GET("/order-search", orderSearchController.SearchOrders)
			authRequiredGroup.GET("/order-search/suggestions", orderSearchController.GetSearchSuggestions)
			authRequiredGroup.GET("/order-search/statistics", orderSearchController.GetSearchStatistics)

			// 订单路由（完整的CRUD操作）
			orderGroup := authRequiredGroup.Group("/orders")
			{
				orderGroup.POST("", orderController.CreateOrder)
				orderGroup.GET("", orderController.GetOrdersByUserID)
				orderGroup.GET("/:id", orderController.GetOrderByID)
				orderGroup.PUT("/:id", orderController.UpdateOrder)
				orderGroup.DELETE("/:id", orderController.DeleteOrder)
				orderGroup.PUT("/:id/status", orderController.UpdateOrderStatus)
				orderGroup.GET("/statistics", orderController.GetOrderStatistics)
				orderGroup.POST("/:id/add-fabric", orderController.AddFabricToOrder)
				orderGroup.DELETE("/:id/remove-fabric", orderController.RemoveFabricFromOrder)
				orderGroup.POST("/:id/add-file", fileController.AddFileToOrder)
				orderGroup.DELETE("/:id/remove-file", orderController.RemoveFileFromOrder)
				orderGroup.GET("/:id/jiedans", jiedanController.GetJiedansByOrderID)
				
				// 进度管理路由
				orderGroup.POST("/:id/progress", progressController.CreateProgress)
				orderGroup.GET("/:id/progress", progressController.GetProgressByOrderID)
				orderGroup.PUT("/:id/progress/:progressId", progressController.UpdateProgress)
				orderGroup.DELETE("/:id/progress/:progressId", progressController.DeleteProgress)
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

			// 接单管理路由（需要认证）
			jiedanGroup := authRequiredGroup.Group("/jiedan")
			{
				jiedanGroup.POST("", jiedanController.CreateJiedan)
				jiedanGroup.GET("/:id", jiedanController.GetJiedanByID)
				jiedanGroup.PUT("/:id", jiedanController.UpdateJiedan)
				jiedanGroup.DELETE("/:id", jiedanController.DeleteJiedan)
				jiedanGroup.POST("/:id/accept", jiedanController.AcceptJiedan)
				jiedanGroup.POST("/:id/reject", jiedanController.RejectJiedan)
			}

			// 工厂接单相关路由
			authRequiredGroup.GET("/factories/:factoryId/jiedans", jiedanController.GetJiedansByFactoryID)
			authRequiredGroup.GET("/factories/:factoryId/jiedan-statistics", jiedanController.GetJiedanStatistics)
			
			// 工厂进度管理路由
			authRequiredGroup.GET("/factories/:factoryId/progress", progressController.GetProgressByFactoryID)
			authRequiredGroup.GET("/factories/:factoryId/progress-statistics", progressController.GetProgressStatistics)
			
			// 根据工厂ID获取工厂详情（需要认证）
			authRequiredGroup.GET("/factory/:id", factoryController.GetFactoryByID)
			
			// 职工管理路由（仅工厂角色）
			employeeGroup := authRequiredGroup.Group("/employees")
			employeeGroup.Use(middleware.FactoryRoleMiddleware())
			{
				employeeGroup.POST("", employeeController.CreateEmployee)
				employeeGroup.GET("", employeeController.GetEmployees)
				employeeGroup.GET("/statistics", employeeController.GetEmployeeStatistics)
				employeeGroup.GET("/search", employeeController.SearchEmployees)
				employeeGroup.GET("/:id", employeeController.GetEmployee)
				employeeGroup.PUT("/:id", employeeController.UpdateEmployee)
				employeeGroup.DELETE("/:id", employeeController.DeleteEmployee)
			}
		}
	}

	// 注册公开路由
	RegisterPublicRoutes(r, db)

	return r
} 