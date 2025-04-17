package routes

import (
	"aneworder.com/backend/controllers"
	"aneworder.com/backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine, userController *controllers.UserController) {
	userGroup := router.Group("/api/users")
	{
		// 公开路由
		userGroup.POST("/register", userController.Register)
		userGroup.POST("/login", userController.Login)

		// 需要认证的路由
		authGroup := userGroup.Group("")
		authGroup.Use(middleware.AuthMiddleware())
		{
			authGroup.GET("/:id", userController.GetUser)
			authGroup.PUT("/:id", userController.UpdateUser)
			authGroup.DELETE("/:id", userController.DeleteUser)
		}
	}
} 