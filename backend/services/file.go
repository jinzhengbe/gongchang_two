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
	"mime/multipart"
	"encoding/json"
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

// 工厂图片相关方法

// BatchUploadFactoryPhotos 批量上传工厂图片
func (s *FileService) BatchUploadFactoryPhotos(files []*multipart.FileHeader, factoryID string, category string) (*models.BatchUploadFactoryPhotosResponse, error) {
	log.Printf("Starting batch upload for factory: %s, category: %s, files count: %d", factoryID, category, len(files))
	
	response := &models.BatchUploadFactoryPhotosResponse{
		Success:      true,
		Message:      "批量上传成功",
		UploadedCount: 0,
		FailedCount:  0,
		Photos:       make([]*models.FactoryPhotoInfo, 0),
		FailedFiles:  make([]*models.FailedFileInfo, 0),
	}

	// 验证工厂是否存在
	var factory models.FactoryProfile
	if err := s.db.Where("id = ?", factoryID).First(&factory).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("工厂不存在")
		}
		return nil, fmt.Errorf("验证工厂失败: %v", err)
	}

	// 批量处理文件
	for _, fileHeader := range files {
		photoInfo, err := s.processFactoryPhoto(fileHeader, factoryID, category)
		if err != nil {
			log.Printf("Failed to process file %s: %v", fileHeader.Filename, err)
			response.FailedCount++
			response.FailedFiles = append(response.FailedFiles, &models.FailedFileInfo{
				Name:  fileHeader.Filename,
				Error: err.Error(),
			})
			continue
		}
		
		response.UploadedCount++
		response.Photos = append(response.Photos, photoInfo)
	}

	// 更新工厂信息中的photos字段
	if err := s.updateFactoryPhotos(factoryID, response.Photos); err != nil {
		log.Printf("Failed to update factory photos: %v", err)
		// 不返回错误，因为文件已经上传成功
	}

	return response, nil
}

// processFactoryPhoto 处理单个工厂图片
func (s *FileService) processFactoryPhoto(fileHeader *multipart.FileHeader, factoryID string, category string) (*models.FactoryPhotoInfo, error) {
	// 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 检查文件大小
	if fileHeader.Size > 10*1024*1024 { // 10MB限制
		return nil, fmt.Errorf("文件大小超过限制 (最大 10MB)")
	}

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	supportedImageExts := []string{".jpg", ".jpeg", ".png", ".webp"}
	supported := false
	for _, supportedExt := range supportedImageExts {
		if ext == supportedExt {
			supported = true
			break
		}
	}
	if !supported {
		return nil, fmt.Errorf("不支持的文件格式: %s (支持: JPG, PNG, WebP)", ext)
	}

	// 生成唯一文件名
	fileID := uuid.New().String()
	newFilename := fileID + ext

	// 保存文件
	finalPath := filepath.Join(s.uploadPath, newFilename)
	dst, err := os.Create(finalPath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %v", err)
	}
	defer dst.Close()

	// 复制文件内容
	written, err := io.Copy(dst, file)
	if err != nil {
		os.Remove(finalPath) // 清理失败的文件
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}

	// 创建文件记录
	fileRecord := &models.File{
		ID:        fileID,
		Name:      fileHeader.Filename,
		Path:      newFilename,
		Type:      "image",
		FactoryID: factoryID,
		Category:  category,
		Size:      written,
	}

	// 保存到数据库
	if err := s.db.Create(fileRecord).Error; err != nil {
		os.Remove(finalPath) // 清理失败的文件
		return nil, fmt.Errorf("保存文件记录失败: %v", err)
	}

	// 生成缩略图（可选）
	thumbnailURL := ""
	if written > 1024*1024 { // 大于1MB的图片生成缩略图
		thumbnailPath, err := s.generateThumbnail(finalPath, fileID)
		if err == nil {
			thumbnailURL = "/uploads/thumbnails/" + filepath.Base(thumbnailPath)
		}
	}

	return &models.FactoryPhotoInfo{
		ID:           fileID,
		Name:         fileHeader.Filename,
		URL:          "/uploads/" + newFilename,
		ThumbnailURL: thumbnailURL,
		Category:     category,
		Size:         written,
		FactoryID:    factoryID,
		Status:       "success",
		CreatedAt:    fileRecord.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// GetFactoryPhotos 获取工厂图片列表
func (s *FileService) GetFactoryPhotos(factoryID string, category string, page, pageSize int) (*models.GetFactoryPhotosResponse, error) {
	query := s.db.Model(&models.File{}).Where("factory_id = ? AND type = ?", factoryID, "image")
	
	// 按分类筛选
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var files []models.File
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&files)

	// 转换为响应格式
	photos := make([]*models.FactoryPhotoInfo, 0, len(files))
	for _, file := range files {
		thumbnailURL := ""
		if file.Size > 1024*1024 {
			thumbnailURL = "/uploads/thumbnails/" + strings.TrimSuffix(file.Path, filepath.Ext(file.Path)) + "_thumb" + filepath.Ext(file.Path)
		}

		photos = append(photos, &models.FactoryPhotoInfo{
			ID:           file.ID,
			Name:         file.Name,
			URL:          "/uploads/" + file.Path,
			ThumbnailURL: thumbnailURL,
			Category:     file.Category,
			Size:         file.Size,
			FactoryID:    file.FactoryID,
			Status:       "success",
			CreatedAt:    file.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	// 获取分类统计
	categories, err := s.getPhotoCategories(factoryID)
	if err != nil {
		log.Printf("Failed to get photo categories: %v", err)
	}

	return &models.GetFactoryPhotosResponse{
		Success:    true,
		Total:      total,
		Photos:     photos,
		Categories: categories,
	}, nil
}

// DeleteFactoryPhoto 删除单张工厂图片
func (s *FileService) DeleteFactoryPhoto(photoID string, factoryID string) error {
	var file models.File
	if err := s.db.Where("id = ? AND factory_id = ? AND type = ?", photoID, factoryID, "image").First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("图片不存在")
		}
		return err
	}

	// 删除物理文件
	filePath := filepath.Join(s.uploadPath, file.Path)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		log.Printf("Failed to remove file %s: %v", filePath, err)
	}

	// 删除缩略图（如果存在）
	thumbnailPath := filepath.Join(s.uploadPath, "thumbnails", strings.TrimSuffix(file.Path, filepath.Ext(file.Path))+"_thumb"+filepath.Ext(file.Path))
	os.Remove(thumbnailPath)

	// 删除数据库记录
	return s.db.Delete(&file).Error
}

// BatchDeleteFactoryPhotos 批量删除工厂图片
func (s *FileService) BatchDeleteFactoryPhotos(photoIDs []string, factoryID string) (*models.BatchDeleteFactoryPhotosResponse, error) {
	response := &models.BatchDeleteFactoryPhotosResponse{
		Success:        true,
		Message:        "批量删除成功",
		DeletedCount:   0,
		FailedCount:    0,
		FailedPhotoIDs: make([]string, 0),
	}

	for _, photoID := range photoIDs {
		if err := s.DeleteFactoryPhoto(photoID, factoryID); err != nil {
			response.FailedCount++
			response.FailedPhotoIDs = append(response.FailedPhotoIDs, photoID)
			log.Printf("Failed to delete photo %s: %v", photoID, err)
		} else {
			response.DeletedCount++
		}
	}

	return response, nil
}

// updateFactoryPhotos 更新工厂信息中的photos字段
func (s *FileService) updateFactoryPhotos(factoryID string, photos []*models.FactoryPhotoInfo) error {
	// 获取现有的photos字段
	var factory models.FactoryProfile
	if err := s.db.Where("id = ?", factoryID).First(&factory).Error; err != nil {
		return err
	}

	// 解析现有的photos
	var existingPhotos []string
	if factory.Photos != "" {
		if err := json.Unmarshal([]byte(factory.Photos), &existingPhotos); err != nil {
			log.Printf("Failed to unmarshal existing photos: %v", err)
			existingPhotos = []string{}
		}
	}

	// 添加新的图片URL
	for _, photo := range photos {
		existingPhotos = append(existingPhotos, photo.URL)
	}

	// 去重
	uniquePhotos := make([]string, 0)
	seen := make(map[string]bool)
	for _, photo := range existingPhotos {
		if !seen[photo] {
			seen[photo] = true
			uniquePhotos = append(uniquePhotos, photo)
		}
	}

	// 更新工厂信息
	photosJSON, err := json.Marshal(uniquePhotos)
	if err != nil {
		return err
	}

	return s.db.Model(&factory).Update("photos", string(photosJSON)).Error
}

// getPhotoCategories 获取图片分类统计
func (s *FileService) getPhotoCategories(factoryID string) ([]*models.PhotoCategory, error) {
	// 这里可以扩展为从数据库查询分类，目前返回默认分类
	defaultCategories := []*models.PhotoCategory{
		{ID: 1, FactoryID: factoryID, Name: "workshop", Color: "#FF5733", Count: 0},
		{ID: 2, FactoryID: factoryID, Name: "equipment", Color: "#33FF57", Count: 0},
		{ID: 3, FactoryID: factoryID, Name: "products", Color: "#3357FF", Count: 0},
		{ID: 4, FactoryID: factoryID, Name: "certificates", Color: "#F3FF33", Count: 0},
	}

	// 统计每个分类的图片数量
	for _, category := range defaultCategories {
		var count int64
		s.db.Model(&models.File{}).Where("factory_id = ? AND type = ? AND category = ?", factoryID, "image", category.Name).Count(&count)
		category.Count = int(count)
	}

	return defaultCategories, nil
}

// generateThumbnail 生成缩略图
func (s *FileService) generateThumbnail(imagePath, fileID string) (string, error) {
	// 这里可以集成图片处理库（如imaging）来生成缩略图
	// 目前返回原图路径
	return imagePath, nil
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