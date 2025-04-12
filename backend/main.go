package main

import (
	"backend/config"
	"backend/controllers"
	"backend/middleware"
	"backend/models"
	"backend/routes"
	"backend/services"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// 初始化数据库连接
	dsn := cfg.DBUser + ":" + cfg.DBPassword + "@tcp(" + cfg.DBHost + ":" + cfg.DBPort + ")/" + cfg.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// 自动迁移数据库表
	err = db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.OrderProgress{},
		&models.OrderAttachment{},
	)
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	// 初始化服务
	userService := services.NewUserService(db)
	productService := services.NewProductService(db)
	orderService := services.NewOrderService(db)

	// 初始化控制器
	userController := controllers.NewUserController(userService)
	productController := controllers.NewProductController(productService)
	orderController := controllers.NewOrderController(orderService)

	// 初始化路由
	router := gin.Default()

	// 配置CORS
	router.Use(middleware.CORSMiddleware())

	// 设置路由
	routes.SetupUserRoutes(router, userController)
	routes.SetupProductRoutes(router, productController)
	routes.SetupOrderRoutes(router, orderController)

	// 创建控制器
	fileController := controllers.NewFileController("uploads")

	// 注册路由
	routes.RegisterFileRoutes(router, fileController)

	// 注册认证路由
	routes.RegisterAuthRoutes(router)

	// 启动服务器
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
} 