package controllers

import (
	"aneworder.com/backend/services"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FileController struct {
	fileService *services.FileService
	uploadDir   string
}

func NewFileController(fileService *services.FileService, uploadDir string) *FileController {
	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("Failed to create upload directory: %v", err)
		panic(err)
	}
	log.Printf("Upload directory created/verified at: %s", uploadDir)
	return &FileController{
		fileService: fileService,
		uploadDir:   uploadDir,
	}
}

// UploadFile 处理文件上传
func (c *FileController) UploadFile(ctx *gin.Context) {
	// 获取并记录所有表单数据
	form, err := ctx.MultipartForm()
	if err != nil {
		log.Printf("Failed to parse multipart form: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无法解析表单数据"})
		return
	}

	// 记录所有表单字段
	log.Printf("Form fields received: %v", form.Value)

	// 获取订单ID
	orderIDStr := ctx.PostForm("orderId")
	log.Printf("Received orderId: %s", orderIDStr)

	if orderIDStr == "" {
		log.Printf("OrderId is empty")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "订单ID不能为空"})
		return
	}

	// 尝试解析订单ID
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		log.Printf("Failed to parse orderId '%s': %v", orderIDStr, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("无效的订单ID格式: %s", orderIDStr)})
		return
	}

	// 获取上传的文件
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		log.Printf("Failed to get uploaded file: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
		return
	}
	defer file.Close()

	log.Printf("Processing file upload: %s (size: %d bytes) for order %d", 
		header.Filename, header.Size, orderID)

	// 保存文件
	fileRecord, err := c.fileService.SaveFile(file, header.Filename, uint(orderID))
	if err != nil {
		log.Printf("Failed to save file: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("File uploaded successfully: %s", fileRecord.ID)
	ctx.JSON(http.StatusOK, fileRecord)
}

// GetOrderFiles 获取订单的所有文件
func (c *FileController) GetOrderFiles(ctx *gin.Context) {
	orderID, err := strconv.ParseUint(ctx.Param("orderId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的订单ID"})
		return
	}

	files, err := c.fileService.GetOrderFiles(uint(orderID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"files": files})
}

// GetFileInfo 获取文件信息
func (c *FileController) GetFileInfo(ctx *gin.Context) {
	fileID := ctx.Param("fileId")

	file, err := c.fileService.GetFileByID(fileID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	ctx.JSON(http.StatusOK, file)
}

// DownloadFile 处理文件下载
func (c *FileController) DownloadFile(ctx *gin.Context) {
	fileID := ctx.Param("fileId")

	filePath, err := c.fileService.GetFilePath(fileID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	ctx.File(filePath)
}

// DeleteFile 处理文件删除
func (c *FileController) DeleteFile(ctx *gin.Context) {
	fileID := ctx.Param("fileId")

	if err := c.fileService.DeleteFile(fileID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "文件删除成功"})
} 