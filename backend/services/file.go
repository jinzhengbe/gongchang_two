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
	"mime"
)

const (
	MaxFileSize = 100 * 1024 * 1024 // 100MB
)

// 支持的文件类型映射
var SupportedFileTypes = map[string][]string{
	"image": {
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg",
	},
	"attachment": {
		".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".zip", ".rar",
	},
	"model": {
		".stl", ".obj", ".fbx", ".dae", ".3ds", ".max", ".blend",
	},
	"video": {
		".mp4", ".avi", ".mov", ".wmv", ".flv", ".webm", ".mkv",
	},
}

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

func (s *FileService) SaveFile(file io.Reader, filename string, orderID *uint, fileType string) (*models.File, error) {
	log.Printf("Starting SaveFile process for file: %s, type: %s", filename, fileType)
	if orderID != nil {
		log.Printf("OrderID provided: %d", *orderID)
	} else {
		log.Printf("No OrderID provided")
	}

	// 获取原始扩展名 - 保持原始大小写
	originalExt := filepath.Ext(filename)
	log.Printf("Original extension: %s", originalExt)

	// 验证文件类型（使用小写进行比较）
	lowerExt := strings.ToLower(originalExt)
	if fileType != "" {
		if supportedExts, exists := SupportedFileTypes[fileType]; exists {
			extSupported := false
			for _, ext := range supportedExts {
				if lowerExt == ext {
					extSupported = true
					break
				}
			}
			if !extSupported {
				return nil, fmt.Errorf("不支持的文件类型: %s (支持的类型: %v)", lowerExt, supportedExts)
			}
		}
	}

	// 生成唯一文件名
	fileID := uuid.New().String()
	
	// 处理扩展名 - 保持客户端原始扩展名
	var finalExt string
	if originalExt == "" {
		// 只有在没有扩展名时才设置默认扩展名
		switch fileType {
		case "image":
			finalExt = ".jpg"
		case "attachment":
			finalExt = ".pdf"
		case "model":
			finalExt = ".stl"
		case "video":
			finalExt = ".mp4"
		default:
			finalExt = ".txt"
		}
		log.Printf("No extension found in filename, using default: %s", finalExt)
	} else {
		// 保持客户端传过来的原始扩展名
		finalExt = originalExt
		log.Printf("Using original extension from client: %s", finalExt)
	}
	
	newFilename := fileID + finalExt
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

	// 验证文件内容（可选：检查文件头）
	if err := s.validateFileContent(tempFile, fileType, finalExt); err != nil {
		log.Printf("File content validation failed: %v", err)
		return nil, err
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

// GetDB 获取数据库连接
func (s *FileService) GetDB() *gorm.DB {
	return s.db
}

// validateFileContent 验证文件内容
func (s *FileService) validateFileContent(filePath, fileType, extension string) error {
	// 读取文件头进行验证
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("无法打开文件进行验证: %v", err)
	}
	defer file.Close()

	// 读取前512字节用于文件头检测
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("读取文件头失败: %v", err)
	}

	// 检测MIME类型
	detectedMIME := mime.TypeByExtension(extension)
	log.Printf("Detected MIME type for %s: %s", extension, detectedMIME)

	// 根据文件类型进行基本验证
	switch fileType {
	case "image":
		return s.validateImageFile(buffer, extension)
	case "video":
		return s.validateVideoFile(buffer, extension)
	case "attachment":
		return s.validateAttachmentFile(buffer, extension)
	case "model":
		return s.validateModelFile(buffer, extension)
	default:
		log.Printf("No specific validation for file type: %s", fileType)
		return nil
	}
}

// validateImageFile 验证图片文件
func (s *FileService) validateImageFile(buffer []byte, extension string) error {
	// 检查常见图片文件头
	imageHeaders := map[string][]byte{
		".jpg":  {0xFF, 0xD8, 0xFF},
		".jpeg": {0xFF, 0xD8, 0xFF},
		".png":  {0x89, 0x50, 0x4E, 0x47},
		".gif":  {0x47, 0x49, 0x46},
		".bmp":  {0x42, 0x4D},
		".webp": {0x52, 0x49, 0x46, 0x46},
	}

	if header, exists := imageHeaders[extension]; exists {
		if len(buffer) >= len(header) {
			for i, b := range header {
				if buffer[i] != b {
					return fmt.Errorf("图片文件头验证失败: %s", extension)
				}
			}
			log.Printf("Image file header validated: %s", extension)
			return nil
		}
	}
	
	log.Printf("No specific header validation for image type: %s", extension)
	return nil
}

// validateVideoFile 验证视频文件
func (s *FileService) validateVideoFile(buffer []byte, extension string) error {
	// 检查常见视频文件头
	videoHeaders := map[string][]byte{
		".mp4": {0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70}, // MP4
		".avi": {0x52, 0x49, 0x46, 0x46}, // AVI
		".mov": {0x00, 0x00, 0x00, 0x14, 0x66, 0x74, 0x79, 0x70}, // MOV
	}

	if header, exists := videoHeaders[extension]; exists {
		if len(buffer) >= len(header) {
			for i, b := range header {
				if buffer[i] != b {
					return fmt.Errorf("视频文件头验证失败: %s", extension)
				}
			}
			log.Printf("Video file header validated: %s", extension)
			return nil
		}
	}
	
	log.Printf("No specific header validation for video type: %s", extension)
	return nil
}

// validateAttachmentFile 验证附件文件
func (s *FileService) validateAttachmentFile(buffer []byte, extension string) error {
	// 检查常见文档文件头
	docHeaders := map[string][]byte{
		".pdf": {0x25, 0x50, 0x44, 0x46}, // PDF
		".zip": {0x50, 0x4B, 0x03, 0x04}, // ZIP
		".rar": {0x52, 0x61, 0x72, 0x21}, // RAR
	}

	if header, exists := docHeaders[extension]; exists {
		if len(buffer) >= len(header) {
			for i, b := range header {
				if buffer[i] != b {
					return fmt.Errorf("附件文件头验证失败: %s", extension)
				}
			}
			log.Printf("Attachment file header validated: %s", extension)
			return nil
		}
	}
	
	log.Printf("No specific header validation for attachment type: %s", extension)
	return nil
}

// validateModelFile 验证模型文件
func (s *FileService) validateModelFile(buffer []byte, extension string) error {
	// 3D模型文件通常没有固定的文件头，这里只做基本检查
	log.Printf("Model file validation skipped for: %s", extension)
	return nil
}