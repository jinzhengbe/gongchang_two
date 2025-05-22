package routes

import (
	"github.com/gin-gonic/gin"
	"backend/controllers"
	"backend/services"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(router *gin.Engine, db *gorm.DB) {
	userService := services.NewUserService(db)
	authController := controllers.NewAuthController(userService)
	
	auth := router.Group("/api/auth")
	{
		auth.POST("/login", authController.Login)
	}
} 