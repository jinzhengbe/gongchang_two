package routes

import (
	"github.com/gin-gonic/gin"
	"backend/controllers"
)

func RegisterFileRoutes(r *gin.Engine, fileController *controllers.FileController) {
	files := r.Group("/api/files")
	{
		files.POST("/upload", fileController.UploadFile)
		files.GET("/order/:orderId", fileController.GetOrderFiles)
		files.GET("/order/:orderId/file/:fileId", fileController.GetFileInfo)
		files.GET("/order/:orderId/file/:fileId/download", fileController.DownloadFile)
		files.DELETE("/order/:orderId/file/:fileId", fileController.DeleteFile)
	}
} 