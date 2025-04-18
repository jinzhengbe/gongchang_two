package routes

import (
	"aneworder.com/backend/controllers"
	"aneworder.com/backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"aneworder.com/backend/middleware"
)

func SetupRouter(router *gin.Engine, db *gorm.DB) {
	// 添加 CORS 中间件
	router.Use(middleware.CORSMiddleware())

	// 创建服务实例
	userService := services.NewUserService(db)
	productService := services.NewProductService(db)
	orderService := services.NewOrderService(db)
	fileService := services.NewFileService(db)

	// 创建控制器实例
	userController := controllers.NewUserController(userService)
	productController := controllers.NewProductController(productService)
	orderController := controllers.NewOrderController(orderService, fileService)
	fileController := controllers.NewFileController(fileService, "uploads")

	// 注册路由
	RegisterAuthRoutes(router)
	SetupUserRoutes(router, userController)
	SetupProductRoutes(router, productController)
	SetupOrderRoutes(router, orderController)
	RegisterFileRoutes(router, fileController)
} 