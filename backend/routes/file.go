package routes

import (
	"github.com/gin-gonic/gin"
	"gongChang/controllers"
	"gongChang/middleware"
)

func SetupFileRoutes(router *gin.Engine, fileController *controllers.FileController) {
	fileGroup := router.Group("/api/files")
	{
		// 公开路由
		fileGroup.GET("/:id", fileController.GetFileDetails)
		fileGroup.GET("/download/:id", fileController.DownloadFile)

		// 需要认证的路由
		authGroup := fileGroup.Group("")
		authGroup.Use(middleware.AuthMiddleware())
		{
			authGroup.POST("/upload", fileController.UploadFile)
			authGroup.POST("/batch", fileController.GetBatchFileDetails)
			authGroup.DELETE("/:id", fileController.DeleteFile)
			authGroup.GET("/order/:id", fileController.GetOrderFiles)
		}
	}
} 