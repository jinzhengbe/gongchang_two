package services

import (
	"aneworder.com/backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io"
	"os"
	"path/filepath"
)

type FileService struct {
	db         *gorm.DB
	uploadPath string
}

func NewFileService(db *gorm.DB, uploadPath string) *FileService {
	return &FileService{
		db:         db,
		uploadPath: uploadPath,
	}
}

func (s *FileService) SaveFile(file io.Reader, filename string, orderID uint) (*models.File, error) {
	// 生成唯一文件名
	fileID := uuid.New().String()
	ext := filepath.Ext(filename)
	newFilename := fileID + ext

	// 确保上传目录存在
	if err := os.MkdirAll(s.uploadPath, 0755); err != nil {
		return nil, err
	}

	// 创建文件
	dst, err := os.Create(filepath.Join(s.uploadPath, newFilename))
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// 复制文件内容
	if _, err = io.Copy(dst, file); err != nil {
		return nil, err
	}

	// 创建文件记录
	fileRecord := &models.File{
		ID:       fileID,
		Name:     filename,
		Path:     newFilename,
		OrderID:  orderID,
	}

	if err := s.db.Create(fileRecord).Error; err != nil {
		// 如果数据库创建失败，删除已上传的文件
		os.Remove(filepath.Join(s.uploadPath, newFilename))
		return nil, err
	}

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