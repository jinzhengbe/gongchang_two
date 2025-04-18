package controllers

import (
	"aneworder.com/backend/services"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"strconv"
)

type FileController struct {
	fileService *services.FileService
	uploadPath  string
}

func NewFileController(fileService *services.FileService, uploadPath string) *FileController {
	return &FileController{
		fileService: fileService,
		uploadPath:  uploadPath,
	}
}

// UploadFile handles file upload
func (c *FileController) UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get order ID from form data
	orderID, err := strconv.ParseUint(ctx.PostForm("order_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// Get file type from form data
	fileType := ctx.PostForm("file_type")
	if fileType == "" {
		fileType = "general" // Default file type
	}

	// Upload file
	fileIDs, err := c.fileService.UploadFiles([]*multipart.FileHeader{file}, uint(orderID), fileType, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"file_id": fileIDs[0],
	})
}

// GetOrderFiles gets all files associated with an order
func (c *FileController) GetOrderFiles(ctx *gin.Context) {
	orderID, err := strconv.ParseUint(ctx.Param("orderId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	files, err := c.fileService.GetFilesByOrderID(uint(orderID), "")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, files)
}

// GetFileInfo gets information about a specific file
func (c *FileController) GetFileInfo(ctx *gin.Context) {
	fileID, err := strconv.ParseUint(ctx.Param("fileId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	file, err := c.fileService.GetFileByID(uint(fileID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, file)
}

// DownloadFile handles file download
func (c *FileController) DownloadFile(ctx *gin.Context) {
	fileID, err := strconv.ParseUint(ctx.Param("fileId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	file, err := c.fileService.GetFileByID(uint(fileID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.File(file.FilePath)
}

// DeleteFile handles file deletion
func (c *FileController) DeleteFile(ctx *gin.Context) {
	fileID, err := strconv.ParseUint(ctx.Param("fileId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := c.fileService.DeleteFile(uint(fileID), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
} 