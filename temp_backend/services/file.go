package services

import (
	"aneworder.com/backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io"
	"log"
	"os"
	"path/filepath"
	"fmt"
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
	// 生成唯一文件名
	fileID := uuid.New().String()
	ext := filepath.Ext(filename)
	newFilename := fileID + ext

	// 确保上传目录存在
	if err := os.MkdirAll(s.uploadPath, 0755); err != nil {
		log.Printf("Failed to create upload directory: %v", err)
		return nil, err
	}

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

	// 复制文件内容并检查大小
	written, err := io.Copy(dst, file)
	if err != nil {
		log.Printf("Failed to copy file content: %v", err)
		return nil, err
	}

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

	// 创建文件记录
	fileRecord := &models.File{
		ID:      fileID,
		Name:    filename,
		Path:    newFilename,
		OrderID: orderID,
	}

	if err := s.db.Create(fileRecord).Error; err != nil {
		log.Printf("Failed to create file record in database: %v", err)
		// 如果数据库创建失败，删除已上传的文件
		os.Remove(finalPath)
		return nil, err
	}

	log.Printf("File saved successfully: %s (%d bytes)", newFilename, written)
	return fileRecord, nil
}

func (s *FileService) GetOrderFiles(orderID uint) ([]models.File, error) {
	var files []models.File
	err := s.db.Where("order_id = ?", orderID).Find(&files).Error
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
		return "", err
	}
	return filepath.Join(s.uploadPath, file.Path), nil
} 