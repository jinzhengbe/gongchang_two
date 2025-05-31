package routes

import (
	"gongChang/controllers"
	"gongChang/middleware"

	"github.com/gin-gonic/gin"
)

func SetupProductRoutes(router *gin.Engine, productController *controllers.ProductController) {
	productGroup := router.Group("/api/products")
	{
		// 公开路由
		productGroup.GET("", productController.GetProducts)
		productGroup.GET("/:id", productController.GetProduct)

		// 需要认证的路由
		authGroup := productGroup.Group("")
		authGroup.Use(middleware.AuthMiddleware())
		{
			authGroup.POST("", productController.CreateProduct)
			authGroup.PUT("/:id", productController.UpdateProduct)
			authGroup.DELETE("/:id", productController.DeleteProduct)
		}
	}
} 