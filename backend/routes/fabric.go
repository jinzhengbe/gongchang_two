package routes

import (
	"github.com/gin-gonic/gin"
	"gongChang/controllers"
	"gongChang/middleware"
)

func SetupFabricRoutes(router *gin.Engine, fabricController *controllers.FabricController) {
	// 公开路由 - 无需认证
	publicFabricGroup := router.Group("/api/fabrics")
	{
		// 获取所有布料（用于前端下拉选择）
		publicFabricGroup.GET("/all", fabricController.GetAllFabrics)
		
		// 获取布料分类
		publicFabricGroup.GET("/categories", fabricController.GetFabricCategories)
		
		// 搜索布料
		publicFabricGroup.GET("/search", fabricController.SearchFabrics)
		
		// 根据分类获取布料
		publicFabricGroup.GET("/category/:category", fabricController.GetFabricsByCategory)
		
		// 根据材质获取布料
		publicFabricGroup.GET("/material/:material", fabricController.GetFabricsByMaterial)
		
		// 获取布料详情
		publicFabricGroup.GET("/:id", fabricController.GetFabricByID)
		
		// 获取布料统计信息
		publicFabricGroup.GET("/statistics", fabricController.GetFabricStatistics)
	}

	// 需要认证的路由 - 管理员和供应商可以管理布料
	authFabricGroup := router.Group("/api/fabrics")
	authFabricGroup.Use(middleware.AuthMiddleware())
	{
		// 创建布料
		authFabricGroup.POST("", fabricController.CreateFabric)
		
		// 更新布料
		authFabricGroup.PUT("/:id", fabricController.UpdateFabric)
		
		// 删除布料
		authFabricGroup.DELETE("/:id", fabricController.DeleteFabric)
		
		// 更新布料库存
		authFabricGroup.PUT("/:id/stock", fabricController.UpdateFabricStock)
	}
} 