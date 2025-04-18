package services

import (
	"aneworder.com/backend/models"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type FileService struct {
	db *gorm.DB
}

func NewFileService(db *gorm.DB) *FileService {
	return &FileService{
		db: db,
	}
}

// UploadFiles handles file uploads and associates them with an order
func (s *FileService) UploadFiles(files []*multipart.FileHeader, orderID uint, fileType string, userID uint) ([]uint, error) {
	var fileIDs []uint

	// Create upload directory if it doesn't exist
	uploadDir := filepath.Join("uploads", fileType)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}

	for _, file := range files {
		// Generate unique filename
		filename := time.Now().Format("20060102150405") + "_" + file.Filename
		filePath := filepath.Join(uploadDir, filename)

		// Save file
		src, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return nil, err
		}

		// Create file record in database
		fileRecord := models.File{
			FileName:   file.Filename,
			FilePath:   filePath,
			FileType:   fileType,
			UploadedBy: userID,
			OrderID:    orderID,
		}

		if err := s.db.Create(&fileRecord).Error; err != nil {
			return nil, err
		}

		// Associate file with order
		if fileType == "model" {
			if err := s.db.Model(&models.Order{}).Where("id = ?", orderID).Association("ModelFiles").Append(&fileRecord); err != nil {
				return nil, err
			}
		} else if fileType == "detail" {
			if err := s.db.Model(&models.Order{}).Where("id = ?", orderID).Association("DetailImages").Append(&fileRecord); err != nil {
				return nil, err
			}
		}

		fileIDs = append(fileIDs, fileRecord.ID)
	}

	return fileIDs, nil
}

// GetFileByID retrieves a file by its ID
func (s *FileService) GetFileByID(fileID uint) (*models.File, error) {
	var file models.File
	if err := s.db.First(&file, fileID).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

// DeleteFile deletes a file and its record
func (s *FileService) DeleteFile(fileID uint, userID uint) error {
	var file models.File
	if err := s.db.First(&file, fileID).Error; err != nil {
		return err
	}

	// Check if user has permission to delete the file
	if file.UploadedBy != userID {
		return errors.New("unauthorized to delete this file")
	}

	// Delete file from storage
	if err := os.Remove(file.FilePath); err != nil {
		return err
	}

	// Delete file record from database
	return s.db.Delete(&file).Error
}

// GetFilesByOrderID retrieves all files associated with an order
func (s *FileService) GetFilesByOrderID(orderID uint, fileType string) ([]models.File, error) {
	var files []models.File
	query := s.db.Where("order_id = ?", orderID)
	if fileType != "" {
		query = query.Where("file_type = ?", fileType)
	}
	err := query.Find(&files).Error
	return files, err
} 