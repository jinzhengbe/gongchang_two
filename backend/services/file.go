package services

import (
	"gongChang/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io"
	"log"
	"os"
	"path/filepath"
	"fmt"
	"strings"
)

const (
	MaxFileSize = 100 * 1024 * 1024 // 100MB
)

type FileService struct {
	db         *gorm.DB
	uploadPath string
}

func NewFileService(db *gorm.DB, uploadPath string) *FileService {
	// 确保上传目录存在
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		log.Printf("Failed to create upload directory: %v", err)
		panic(err)
	}
	log.Printf("File service initialized with upload path: %s", uploadPath)
	return &FileService{
		db:         db,
		uploadPath: uploadPath,
	}
}

func (s *FileService) SaveFile(file io.Reader, filename string, orderID *uint) (*models.File, error) {
	log.Printf("Starting SaveFile process for file: %s", filename)
	if orderID != nil {
		log.Printf("OrderID provided: %d", *orderID)
	} else {
		log.Printf("No OrderID provided")
	}

	// 生成唯一文件名
	fileID := uuid.New().String()
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".txt" // 如果没有扩展名，默认使用.txt
	}
	newFilename := fileID + ext
	log.Printf("Generated new filename: %s", newFilename)

	// 确保上传目录存在
	if err := os.MkdirAll(s.uploadPath, 0755); err != nil {
		log.Printf("Failed to create upload directory: %v", err)
		return nil, err
	}
	log.Printf("Upload directory verified: %s", s.uploadPath)

	// 创建临时文件
	tempFile := filepath.Join(s.uploadPath, "temp_"+newFilename)
	dst, err := os.Create(tempFile)
	if err != nil {
		log.Printf("Failed to create temporary file: %v", err)
		return nil, err
	}
	defer func() {
		dst.Close()
		os.Remove(tempFile) // 清理临时文件
	}()
	log.Printf("Created temporary file: %s", tempFile)

	// 复制文件内容并检查大小
	written, err := io.Copy(dst, file)
	if err != nil {
		log.Printf("Failed to copy file content: %v", err)
		return nil, err
	}
	log.Printf("File content copied, size: %d bytes", written)

	// 检查文件大小
	if written > MaxFileSize {
		log.Printf("File too large: %d bytes (max: %d bytes)", written, MaxFileSize)
		return nil, fmt.Errorf("文件大小超过限制 (最大 %d MB)", MaxFileSize/1024/1024)
	}

	// 重命名临时文件为最终文件
	finalPath := filepath.Join(s.uploadPath, newFilename)
	if err := os.Rename(tempFile, finalPath); err != nil {
		log.Printf("Failed to rename temporary file: %v", err)
		return nil, err
	}
	log.Printf("File renamed to final path: %s", finalPath)

	// 创建文件记录
	fileRecord := &models.File{
		ID:      fileID,
		Name:    filename,
		Path:    newFilename,  // 确保包含扩展名
		OrderID: orderID,
	}
	log.Printf("Created file record: %+v", fileRecord)

	// 使用事务来确保数据一致性
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 如果提供了订单ID，检查订单是否存在
		if orderID != nil {
			var count int64
			if err := tx.Model(&models.Order{}).Where("id = ?", *orderID).Count(&count).Error; err != nil {
				log.Printf("Failed to check order existence: %v", err)
				return err
			}
			if count == 0 {
				log.Printf("Order not found: %d", *orderID)
				return fmt.Errorf("订单不存在")
			}
			log.Printf("Order exists: %d", *orderID)
		}

		// 创建文件记录
		if err := tx.Create(fileRecord).Error; err != nil {
			log.Printf("Failed to create file record: %v", err)
			return err
		}
		log.Printf("File record created successfully")

		return nil
	})

	if err != nil {
		log.Printf("Transaction failed: %v", err)
		// 如果数据库操作失败，删除已上传的文件
		os.Remove(finalPath)
		return nil, err
	}

	log.Printf("File saved successfully: %s (%d bytes)", newFilename, written)
	if orderID != nil {
		log.Printf("Associated with order: %d", *orderID)
	} else {
		log.Printf("No order association")
	}
	return fileRecord, nil
}

func (s *FileService) GetOrderFiles(orderID uint) ([]models.File, error) {
	var files []models.File
	err := s.db.Where("orderID = ?", orderID).Find(&files).Error
	return files, err
}

func (s *FileService) GetFileByID(fileID string) (*models.File, error) {
	var file models.File
	err := s.db.First(&file, "id = ?", fileID).Error
	return &file, err
}

func (s *FileService) DeleteFile(fileID string) error {
	var file models.File
	if err := s.db.First(&file, "id = ?", fileID).Error; err != nil {
		return err
	}

	// 删除物理文件
	if err := os.Remove(filepath.Join(s.uploadPath, file.Path)); err != nil && !os.IsNotExist(err) {
		return err
	}

	// 删除数据库记录
	return s.db.Delete(&file).Error
}

func (s *FileService) GetFilePath(fileID string) (string, error) {
	var file models.File
	if err := s.db.First(&file, "id = ?", fileID).Error; err != nil {
		log.Printf("Failed to find file with ID %s: %v", fileID, err)
		return "", err
	}

	// 确保文件路径是绝对路径
	filePath := filepath.Join(s.uploadPath, file.Path)
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		log.Printf("Failed to get absolute path for %s: %v", filePath, err)
		return "", err
	}

	// 验证文件路径是否在上传目录内
	if !strings.HasPrefix(absPath, s.uploadPath) {
		log.Printf("Invalid file path: %s (outside upload directory)", absPath)
		return "", fmt.Errorf("invalid file path")
	}

	return absPath, nil
}

// GetFilesByIDs 批量获取文件信息
func (s *FileService) GetFilesByIDs(fileIDs []string) ([]models.File, error) {
	var files []models.File
	err := s.db.Where("id IN ?", fileIDs).Find(&files).Error
	return files, err
}