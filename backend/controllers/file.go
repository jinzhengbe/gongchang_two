package controllers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileController struct {
	uploadDir string
}

func NewFileController(uploadDir string) *FileController {
	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(err)
	}
	return &FileController{
		uploadDir: uploadDir,
	}
}

// UploadFile 处理文件上传
func (c *FileController) UploadFile(ctx *gin.Context) {
	// 获取订单ID
	orderID := ctx.PostForm("orderId")
	if orderID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "订单ID不能为空"})
		return
	}

	// 获取上传的文件
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
		return
	}
	defer file.Close()

	// 生成唯一的文件名
	ext := filepath.Ext(header.Filename)
	filename := uuid.New().String() + ext

	// 创建订单目录
	orderDir := filepath.Join(c.uploadDir, orderID)
	if err := os.MkdirAll(orderDir, 0755); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建目录失败"})
		return
	}

	// 保存文件
	filePath := filepath.Join(orderDir, filename)
	out, err := os.Create(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建文件失败"})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	// 返回文件信息
	ctx.JSON(http.StatusOK, gin.H{
		"id":         filename,
		"name":       header.Filename,
		"size":       header.Size,
		"uploadTime": time.Now().Unix(),
	})
}

// DownloadFile 处理文件下载
func (c *FileController) DownloadFile(ctx *gin.Context) {
	orderID := ctx.Param("orderId")
	fileID := ctx.Param("fileId")

	filePath := filepath.Join(c.uploadDir, orderID, fileID)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	ctx.File(filePath)
}

// DeleteFile 处理文件删除
func (c *FileController) DeleteFile(ctx *gin.Context) {
	orderID := ctx.Param("orderId")
	fileID := ctx.Param("fileId")

	filePath := filepath.Join(c.uploadDir, orderID, fileID)
	if err := os.Remove(filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除文件失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "文件删除成功"})
}

// GetFileInfo 获取文件信息
func (c *FileController) GetFileInfo(ctx *gin.Context) {
	orderID := ctx.Param("orderId")
	fileID := ctx.Param("fileId")

	filePath := filepath.Join(c.uploadDir, orderID, fileID)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         fileID,
		"name":       fileInfo.Name(),
		"size":       fileInfo.Size(),
		"uploadTime": fileInfo.ModTime().Unix(),
	})
}

// GetOrderFiles 获取订单的所有文件
func (c *FileController) GetOrderFiles(ctx *gin.Context) {
	orderID := ctx.Param("orderId")
	orderDir := filepath.Join(c.uploadDir, orderID)

	files, err := os.ReadDir(orderDir)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "读取文件列表失败"})
		return
	}

	var fileList []gin.H
	for _, file := range files {
		if !file.IsDir() {
			info, err := file.Info()
			if err != nil {
				continue
			}

			fileList = append(fileList, gin.H{
				"id":         info.Name(),
				"name":       info.Name(),
				"size":       info.Size(),
				"uploadTime": info.ModTime().Unix(),
			})
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"files": fileList})
} 