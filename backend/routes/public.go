package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gongChang/controllers"
)

func RegisterPublicRoutes(r *gin.Engine, db *gorm.DB) {
	// 创建公开路由组
	public := r.Group("/public")
	{
		// 订单相关路由
		orderController := controllers.NewPublicOrderController(db)
		public.GET("/orders", orderController.GetPublicOrders)
		public.GET("/orders/:id", orderController.GetPublicOrderDetail)
		
		// 订单文件相关路由
		fileController := controllers.NewPublicFileController(db)
		public.POST("/orders/:id/files", fileController.UploadOrderFiles)
		public.GET("/orders/:id/files", fileController.GetOrderFiles)
		public.GET("/files/:fileId", fileController.GetFile)
		public.DELETE("/files/:fileId", fileController.DeleteFile)
		
		// 订单3D模型相关路由
		public.POST("/orders/:id/models", fileController.UploadOrderModels)
		public.GET("/orders/:id/models", fileController.GetOrderModels)
		
		// 订单图片相关路由
		public.POST("/orders/:id/images", fileController.UploadOrderImages)
		public.GET("/orders/:id/images", fileController.GetOrderImages)
		
		// 订单视频相关路由
		public.POST("/orders/:id/videos", fileController.UploadOrderVideos)
		public.GET("/orders/:id/videos", fileController.GetOrderVideos)
	}
} 