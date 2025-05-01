package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"aneworder.com/backend/controllers"
)

func RegisterPublicRoutes(r *gin.Engine, db *gorm.DB) {
	public := r.Group("/public")
	{
		orderController := controllers.NewPublicOrderController(db)
		public.GET("/orders", orderController.GetPublicOrders)
	}
} 