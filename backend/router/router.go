package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gongChang/controllers"
	"gongChang/services"
	"gongChang/middleware"
	"gongChang/config"
	"gongChang/internal/factory"
)

func SetupRouter(r *gin.Engine, db *gorm.DB) {
	// 用户相关路由
	userGroup := r.Group("/api/users")
	userController := &controllers.UserController{DB: db}
	userGroup.GET("/:id", userController.GetUser)

	// 工厂相关路由
	factoryGroup := r.Group("/api")
	factoryController := &controllers.FactoryController{DB: db}
	factoryGroup.GET("/factories", factoryController.GetFactoryList)
} 