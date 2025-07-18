package controllers

import (
	"gongChang/services"
	"gongChang/config"
	"gongChang/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type FileController struct {
	fileService *services.FileService
	uploadDir   string
	config      *config.Config
}

func NewFileController(fileService *services.FileService, uploadDir string, cfg *config.Config) *FileController {
	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("Failed to create upload directory: %v", err)
		panic(err)
	}
	log.Printf("Upload directory created/verified at: %s", uploadDir)
	return &FileController{
		fileService: fileService,
		uploadDir:   uploadDir,
		config:      cfg,
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
	fileRecord, err := c.fileService.SaveFile(file, header.Filename, orderID, "")
	if err != nil {
		log.Printf("Failed to save file: %v", err)
		if err.Error() == "订单不存在" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 构建文件访问URL - 直接访问uploads目录
	fileURL := fmt.Sprintf("/uploads/%s", fileRecord.Path)
	
	// 返回包含完整URL的响应
	response := gin.H{
		"id":        fileRecord.ID,
		"name":      fileRecord.Name,
		"path":      fileRecord.Path,
		"url":       fileURL,
		"order_id":  fileRecord.OrderID,
		"created_at": fileRecord.CreatedAt,
		"updated_at": fileRecord.UpdatedAt,
	}

	log.Printf("File uploaded successfully: %s", fileRecord.ID)
	ctx.JSON(http.StatusOK, response)
}

// GetOrderFiles 获取订单的所有文件
func (c *FileController) GetOrderFiles(ctx *gin.Context) {
	orderID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的订单ID"})
		return
	}

	files, err := c.fileService.GetOrderFiles(uint(orderID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 构建文件详情列表，包含完整的URL
	fileDetails := make([]gin.H, 0, len(files))
	baseURL := c.config.Server.BaseURL
	
	for _, file := range files {
		fileURL := fmt.Sprintf("%s/api/files/download/%s", baseURL, file.ID)
		fileDetails = append(fileDetails, gin.H{
			"id":        file.ID,
			"name":      file.Name,
			"path":      file.Path,
			"url":       fileURL,
			"order_id":  file.OrderID,
			"type":      filepath.Ext(file.Name),
			"created_at": file.CreatedAt,
			"updated_at": file.UpdatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"files": fileDetails})
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
	fileID := ctx.Param("id")
	log.Printf("Attempting to download file with ID: %s", fileID)

	filePath, err := c.fileService.GetFilePath(fileID)
	if err != nil {
		log.Printf("Failed to get file path for ID %s: %v", fileID, err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("File not found at path: %s", filePath)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Printf("Failed to get file info for %s: %v", filePath, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取文件信息"})
		return
	}

	// 设置CORS头
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	ctx.Header("Access-Control-Max-Age", "86400")

	// 设置响应头
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", filepath.Base(filePath)))
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// 发送文件
	ctx.File(filePath)
}

// DeleteFile 处理文件删除
func (c *FileController) DeleteFile(ctx *gin.Context) {
	fileID := ctx.Param("id")

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

	// 构建完整的文件访问URL
	fileURL := fmt.Sprintf("/uploads/%s", file.Path)

	ctx.JSON(http.StatusOK, gin.H{
		"id":        file.ID,
		"name":      file.Name,
		"path":      file.Path,
		"url":       fileURL,
		"order_id":  file.OrderID,
		"type":      filepath.Ext(file.Name),
		"created_at": file.CreatedAt,
		"updated_at": file.UpdatedAt,
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
		fileURL := fmt.Sprintf("/uploads/%s", file.Path)
		fileDetails = append(fileDetails, gin.H{
			"id":        file.ID,
			"name":      file.Name,
			"path":      file.Path,
			"url":       fileURL,
			"order_id":  file.OrderID,
			"type":      filepath.Ext(file.Name),
			"created_at": file.CreatedAt,
			"updated_at": file.UpdatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"files": fileDetails})
}

// AddFileToOrder 上传并关联文件到订单
// @Summary 上传并关联文件到订单
// @Description 上传文件并自动关联到指定订单，同时更新订单的文件字段
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "订单ID"
// @Param file formData file true "上传的文件"
// @Param type formData string true "文件类型" Enums(image,attachment,model,video)
// @Param description formData string false "文件描述"
// @Success 200 {object} models.AddFileToOrderResponse
// @Router /api/orders/{id}/add-file [post]
func (c *FileController) AddFileToOrder(ctx *gin.Context) {
	log.Printf("Starting add file to order process...")
	
	// 获取订单ID
	orderID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		log.Printf("Invalid order ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的订单ID"})
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
	log.Printf("Successfully received file: %s for order: %d", header.Filename, orderID)

	// 绑定表单数据
	var req models.AddFileToOrderRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("Failed to bind form data: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "表单数据绑定失败: " + err.Error()})
		return
	}
	
	log.Printf("File type: %s, description: %s", req.Type, req.Description)

	// 保存文件并关联到订单
	orderIDUint := uint(orderID)
	fileRecord, err := c.fileService.SaveFile(file, header.Filename, &orderIDUint, req.Type)
	if err != nil {
		log.Printf("Failed to save file: %v", err)
		if err.Error() == "订单不存在" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 更新订单的文件字段
	orderService := services.NewOrderService(c.fileService.GetDB())
	order, err := orderService.GetOrderByID(uint(orderID))
	if err != nil {
		log.Printf("Failed to get order: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取订单信息失败"})
		return
	}

	// 根据文件类型更新订单的相应字段
	var updateReq models.OrderUpdateRequest
	
	// 获取现有文件列表
	var existingFiles []string
	switch req.Type {
	case "image":
		if order.Images != nil {
			json.Unmarshal(*order.Images, &existingFiles)
		}
		existingFiles = append(existingFiles, fileRecord.ID)
		updateReq.Images = existingFiles
	case "attachment":
		if order.Attachments != nil {
			json.Unmarshal(*order.Attachments, &existingFiles)
		}
		existingFiles = append(existingFiles, fileRecord.ID)
		updateReq.Attachments = existingFiles
	case "model":
		if order.Models != nil {
			json.Unmarshal(*order.Models, &existingFiles)
		}
		existingFiles = append(existingFiles, fileRecord.ID)
		updateReq.Models = existingFiles
	case "video":
		if order.Videos != nil {
			json.Unmarshal(*order.Videos, &existingFiles)
		}
		existingFiles = append(existingFiles, fileRecord.ID)
		updateReq.Videos = existingFiles
	}

	// 更新订单
	if err := orderService.UpdateOrder(uint(orderID), &updateReq); err != nil {
		log.Printf("Failed to update order: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新订单失败"})
		return
	}

	// 重新获取更新后的订单
	updatedOrder, err := orderService.GetOrderByID(uint(orderID))
	if err != nil {
		log.Printf("Failed to get updated order: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取更新后的订单失败"})
		return
	}

	// 构建文件访问URL - 直接访问uploads目录
	fileURL := fmt.Sprintf("/uploads/%s", fileRecord.Path)
	
	// 构建响应
	fileInfo := &models.FileInfo{
		ID:          fileRecord.ID,
		URL:         fileURL,
		Type:        req.Type,
		Name:        fileRecord.Name,
		Description: req.Description,
	}

	response := models.AddFileToOrderResponse{
		Success: true,
		Order:   updatedOrder,
		File:    fileInfo,
	}

	log.Printf("File successfully added to order: %s", fileRecord.ID)
	ctx.JSON(http.StatusOK, response)
}