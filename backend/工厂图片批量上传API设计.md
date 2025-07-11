# 工厂图片批量上传 API 设计文档

## 📋 API 概述

### 功能描述
批量上传图片到指定工厂，并自动关联到工厂信息中。

### API 端点
```
POST /api/factories/{factoryId}/photos/batch
```

### 请求头
```
Content-Type: multipart/form-data
Authorization: Bearer {token}
```

## 🔧 请求参数

### 路径参数
- `factoryId`: 工厂ID (必填)

### 表单参数
```json
{
  "files": [二进制文件1, 二进制文件2, 二进制文件3, ...],
  "type": "image"
}
```

**参数说明：**
- `files`: 图片文件数组 (必填)
- `type`: 文件类型，固定为 "image" (必填)

## 📤 响应格式

### 成功响应 (200)
```json
{
  "success": true,
  "message": "批量上传成功",
  "data": {
    "uploaded_count": 3,
    "photos": [
      {
        "id": "image_id_1",
        "name": "photo1.png",
        "url": "/uploads/image_id_1.png",
        "factory_id": "3af8e32a-e267-45f1-8959-faf3f0787bfa"
      },
      {
        "id": "image_id_2", 
        "name": "photo2.png",
        "url": "/uploads/image_id_2.png",
        "factory_id": "3af8e32a-e267-45f1-8959-faf3f0787bfa"
      },
      {
        "id": "image_id_3",
        "name": "photo3.png", 
        "url": "/uploads/image_id_3.png",
        "factory_id": "3af8e32a-e267-45f1-8959-faf3f0787bfa"
      }
    ]
  }
}
```

### 错误响应

#### 400 - 请求参数错误
```json
{
  "success": false,
  "error": "请选择要上传的图片"
}
```

#### 403 - 权限不足
```json
{
  "success": false,
  "error": "无权限操作此工厂"
}
```

#### 500 - 服务器错误
```json
{
  "success": false,
  "error": "服务器内部错误"
}
```

## 🗄️ 数据库设计

### 文件表结构
```sql
CREATE TABLE files (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    path VARCHAR(500) NOT NULL,
    url VARCHAR(500) NOT NULL,
    type VARCHAR(50) NOT NULL,
    factory_id VARCHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (factory_id) REFERENCES users(id)
);
```

### 工厂信息表更新
```sql
-- 更新工厂信息中的Photos字段
UPDATE users 
SET photos = CONCAT(COALESCE(photos, ''), ',', new_photo_urls)
WHERE id = ?;
```

## 🔧 后台实现逻辑

### 1. 路由定义
```go
// 批量上传工厂图片
router.POST("/api/factories/:factoryId/photos/batch", middleware.AuthMiddleware(), handlers.BatchUploadFactoryPhotos)
```

### 2. 处理函数
```go
func BatchUploadFactoryPhotos(c *gin.Context) {
    // 1. 获取工厂ID
    factoryId := c.Param("factoryId")
    
    // 2. 验证用户权限（只能给自己的工厂上传图片）
    userID := getUserIDFromToken(c)
    if !canManageFactory(userID, factoryId) {
        c.JSON(403, gin.H{"error": "无权限操作此工厂"})
        return
    }
    
    // 3. 获取上传的文件
    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(400, gin.H{"error": "文件格式错误"})
        return
    }
    
    files := form.File["files"]
    if len(files) == 0 {
        c.JSON(400, gin.H{"error": "请选择要上传的图片"})
        return
    }
    
    // 4. 批量处理文件
    var uploadedPhotos []PhotoInfo
    for _, file := range files {
        // 生成唯一文件名
        fileID := generateUUID()
        fileName := fileID + getFileExtension(file.Filename)
        
        // 保存文件到服务器
        filePath := "uploads/" + fileName
        if err := c.SaveUploadedFile(file, filePath); err != nil {
            continue // 跳过失败的文件
        }
        
        // 保存文件信息到数据库
        photoInfo := PhotoInfo{
            ID:        fileID,
            Name:      file.Filename,
            Path:      filePath,
            URL:       "/uploads/" + fileName,
            FactoryID: factoryId,
            CreatedAt: time.Now(),
        }
        
        // 插入数据库
        if err := db.Create(&photoInfo).Error; err != nil {
            continue
        }
        
        uploadedPhotos = append(uploadedPhotos, photoInfo)
    }
    
    // 5. 更新工厂信息中的Photos字段
    var photoURLs []string
    for _, photo := range uploadedPhotos {
        photoURLs = append(photoURLs, photo.URL)
    }
    
    // 获取现有图片URL并合并
    var factory Factory
    db.Where("user_id = ?", factoryId).First(&factory)
    
    existingPhotos := strings.Split(factory.Photos, ",")
    if factory.Photos == "" {
        existingPhotos = []string{}
    }
    
    // 合并新旧图片URL
    allPhotos := append(existingPhotos, photoURLs...)
    factory.Photos = strings.Join(allPhotos, ",")
    
    // 更新工厂信息
    db.Save(&factory)
    
    // 6. 返回结果
    c.JSON(200, gin.H{
        "success": true,
        "message": "批量上传成功",
        "data": gin.H{
            "uploaded_count": len(uploadedPhotos),
            "photos": uploadedPhotos,
        },
    })
}
```

### 3. 数据结构
```go
type PhotoInfo struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Path      string    `json:"path"`
    URL       string    `json:"url"`
    FactoryID string    `json:"factory_id"`
    CreatedAt time.Time `json:"created_at"`
}

type Factory struct {
    ID       uint   `json:"id"`
    UserID   string `json:"user_id"`
    Photos   string `json:"photos"` // 存储为逗号分隔的URL字符串
    // ... 其他字段
}
```

### 4. 辅助函数
```go
func generateUUID() string {
    return uuid.New().String()
}

func getFileExtension(filename string) string {
    return filepath.Ext(filename)
}

func canManageFactory(userID, factoryID string) bool {
    // 检查用户是否有权限管理此工厂
    return userID == factoryID
}
```

## 📱 前端调用示例

### Dart/Flutter 实现
```dart
// 批量上传工厂图片
Future<Map<String, dynamic>> uploadFactoryPhotosBatch({
  required String factoryId,
  required List<String> filePaths,
  required List<String> fileNames,
}) async {
  try {
    final request = http.MultipartRequest(
      'POST',
      Uri.parse('https://aneworders.com/api/factories/$factoryId/photos/batch'),
    );
    
    // 添加认证头
    request.headers['Authorization'] = 'Bearer $token';
    
    // 添加多个文件
    for (int i = 0; i < filePaths.length; i++) {
      final file = File(filePaths[i]);
      final stream = http.ByteStream(file.openRead());
      final length = await file.length();
      
      request.files.add(http.MultipartFile(
        'files',
        stream,
        length,
        filename: fileNames[i],
      ));
    }
    
    final response = await request.send();
    final responseBody = await response.stream.bytesToString();
    
    if (response.statusCode == 200) {
      return json.decode(responseBody);
    } else {
      return {'success': false, 'error': '上传失败'};
    }
  } catch (e) {
    return {'success': false, 'error': e.toString()};
  }
}
```

### JavaScript 实现
```javascript
async function uploadFactoryPhotosBatch(factoryId, files) {
  const formData = new FormData();
  
  // 添加文件
  files.forEach(file => {
    formData.append('files', file);
  });
  
  try {
    const response = await fetch(`/api/factories/${factoryId}/photos/batch`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`
      },
      body: formData
    });
    
    const result = await response.json();
    return result;
  } catch (error) {
    return { success: false, error: error.message };
  }
}
```

## ✅ 功能优势

1. **一步完成**：上传图片的同时直接关联到工厂信息
2. **批量处理**：支持一次上传多张图片
3. **权限控制**：只能给自己的工厂上传图片
4. **错误处理**：单个文件失败不影响其他文件
5. **数据一致性**：自动更新工厂信息中的Photos字段
6. **性能优化**：批量处理减少数据库操作次数

## 🔒 安全考虑

1. **文件类型验证**：只允许图片文件上传
2. **文件大小限制**：设置合理的文件大小上限
3. **权限验证**：确保用户只能操作自己的工厂
4. **文件名安全**：使用UUID避免文件名冲突
5. **路径安全**：限制文件保存路径，防止目录遍历攻击

## 📊 使用场景

1. **工厂信息编辑页面**：用户选择多张图片后批量上传
2. **工厂展示页面**：展示工厂的所有图片
3. **图片管理**：支持删除和重新排序图片

## 🚀 部署注意事项

1. **文件存储**：确保uploads目录有写入权限
2. **磁盘空间**：监控磁盘使用情况
3. **备份策略**：定期备份上传的图片文件
4. **CDN配置**：考虑使用CDN加速图片访问
5. **清理策略**：定期清理未关联的图片文件 