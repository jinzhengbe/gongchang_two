package controllers

import (
	"aneworder.com/backend/services"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	log.Printf("Starting file upload process...")
	
	// 获取上传的文件
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		log.Printf("Failed to get uploaded file: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
		return
	}
	defer file.Close()
	log.Printf("Successfully received file: %s", header.Filename)

	// 获取订单ID（可选）
	var orderID *uint
	orderIDStr := ctx.PostForm("orderId")
	log.Printf("Received orderId: %s", orderIDStr)
	
	if orderIDStr != "" {
		// 尝试解析订单ID
		parsedID, err := strconv.ParseUint(orderIDStr, 10, 32)
		if err != nil {
			log.Printf("Failed to parse orderId '%s': %v", orderIDStr, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("无效的订单ID格式: %s", orderIDStr)})
			return
		}
		uintID := uint(parsedID)
		orderID = &uintID
		log.Printf("Successfully parsed orderId: %d", *orderID)
	} else {
		log.Printf("No orderId provided, proceeding with file upload without order association")
	}

	// 保存文件
	fileRecord, err := c.fileService.SaveFile(file, header.Filename, orderID)
	if err != nil {
		log.Printf("Failed to save file: %v", err)
		if err.Error() == "订单不存在" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
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

// GetFileDetails 获取单个文件详情
func (c *FileController) GetFileDetails(ctx *gin.Context) {
	fileID := ctx.Param("id")

	file, err := c.fileService.GetFileByID(fileID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	// 构建文件URL
	fileURL := fmt.Sprintf("/api/files/%s/download", file.ID)

	ctx.JSON(http.StatusOK, gin.H{
		"id": file.ID,
		"name": file.Name,
		"url": fileURL,
		"type": filepath.Ext(file.Name),
	})
}

// GetBatchFileDetails 批量获取文件详情
func (c *FileController) GetBatchFileDetails(ctx *gin.Context) {
	var req struct {
		IDs []string `json:"ids" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	files, err := c.fileService.GetFilesByIDs(req.IDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 构建文件详情列表
	fileDetails := make([]gin.H, 0, len(files))
	for _, file := range files {
		fileURL := fmt.Sprintf("/api/files/%s/download", file.ID)
		fileDetails = append(fileDetails, gin.H{
			"id": file.ID,
			"name": file.Name,
			"url": fileURL,
			"type": filepath.Ext(file.Name),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"files": fileDetails})
}