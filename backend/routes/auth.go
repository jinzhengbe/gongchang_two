package routes

import (
	"github.com/gin-gonic/gin"
	"aneworder.com/backend/controllers"
)

func RegisterAuthRoutes(router *gin.Engine) {
	authController := new(controllers.AuthController)
	
	auth := router.Group("/api/auth")
	{
		auth.POST("/login", authController.Login)
	}
} 