package routes

import (
	"aneworder.com/backend/controllers"
	"aneworder.com/backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"aneworder.com/backend/middleware"
	"net/http"
)

func SetupRouter(router *gin.Engine, db *gorm.DB) {
	// 添加健康检查路由
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 添加 CORS 中间件
	router.Use(middleware.CORSMiddleware())

	// 创建服务实例
	userService := services.NewUserService(db)
	productService := services.NewProductService(db)
	orderService := services.NewOrderService(db)
	fileService := services.NewFileService(db, "uploads")

	// 创建控制器实例
	userController := controllers.NewUserController(userService)
	productController := controllers.NewProductController(productService)
	orderController := controllers.NewOrderController(orderService)
	fileController := controllers.NewFileController(fileService, "uploads")

	// 注册路由
	RegisterAuthRoutes(router, db)
	SetupUserRoutes(router, userController)
	SetupProductRoutes(router, productController)
	SetupOrderRoutes(router, orderController)
	RegisterFileRoutes(router, fileController)
} 