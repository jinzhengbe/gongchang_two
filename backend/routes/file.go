package routes

import (
	"github.com/gin-gonic/gin"
	"aneworder.com/backend/controllers"
	"aneworder.com/backend/middleware"
)

func RegisterFileRoutes(r *gin.Engine, fileController *controllers.FileController) {
	// 将文件路由注册到 API 组下
	api := r.Group("/api")
	{
		files := api.Group("/files")
		files.Use(middleware.AuthMiddleware())
		{
			files.POST("/upload", fileController.UploadFile)
			files.GET("/order/:orderId", fileController.GetOrderFiles)
			files.GET("/order/:orderId/file/:fileId", fileController.GetFileInfo)
			files.GET("/order/:orderId/file/:fileId/download", fileController.DownloadFile)
			files.DELETE("/order/:orderId/file/:fileId", fileController.DeleteFile)
		}
	}
} 