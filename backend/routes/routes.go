package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// 注册公开路由
	RegisterPublicRoutes(r, db)
	
	// 注册其他路由...
	// RegisterAuthRoutes(r, db)
	// RegisterOrderRoutes(r, db)
	// ...
} 