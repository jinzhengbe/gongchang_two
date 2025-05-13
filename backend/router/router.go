package router

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	userGroup := r.Group("/api/users")
	userGroup.GET("/:id", userController.GetUser)
} 