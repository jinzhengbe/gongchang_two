# 工厂图片API文档

## 概述

本文档描述了工厂图片批量上传和管理API的完整功能，包括批量上传、图片列表获取、删除等功能。

## API端点

### 1. 批量上传工厂图片

**接口地址**: `POST /api/factories/{factoryId}/photos/batch`

**请求头**:
```
Content-Type: multipart/form-data
Authorization: Bearer {token}
```

**路径参数**:
- `factoryId`: 工厂ID (必填)

**表单参数**:
- `files`: 图片文件数组 (必填)
- `category`: 图片分类 (可选)

**请求示例**:
```bash
curl -X POST "http://localhost:8008/api/factories/gongchang/photos/batch" \
  -H "Authorization: Bearer {token}" \
  -F "files=@image1.jpg" \
  -F "files=@image2.jpg" \
  -F "category=workshop"
```

**成功响应** (200):
```json
{
  "success": true,
  "message": "批量上传成功",
  "uploaded_count": 2,
  "failed_count": 0,
  "photos": [
    {
      "id": "photo_id_1",
      "name": "image1.jpg",
      "url": "/uploads/photo_id_1.jpg",
      "thumbnail_url": "/uploads/thumbnails/photo_id_1_thumb.jpg",
      "category": "workshop",
      "size": 1024000,
      "factory_id": "gongchang",
      "status": "success",
      "created_at": "2025-01-15T10:30:00Z"
    }
  ],
  "failed_files": []
}
```

**错误响应**:
- 400: 请求参数错误
- 403: 权限不足
- 413: 文件过大
- 500: 服务器错误

### 2. 获取工厂图片列表

**接口地址**: `GET /api/factories/{factoryId}/photos`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `factoryId`: 工厂ID (必填)

**查询参数**:
- `category`: 图片分类筛选 (可选)
- `page`: 页码，默认1 (可选)
- `page_size`: 每页数量，默认20 (可选)

**请求示例**:
```bash
curl -X GET "http://localhost:8008/api/factories/gongchang/photos?category=workshop&page=1&page_size=10" \
  -H "Authorization: Bearer {token}"
```

**成功响应** (200):
```json
{
  "success": true,
  "total": 15,
  "photos": [
    {
      "id": "photo_id_1",
      "name": "image1.jpg",
      "url": "/uploads/photo_id_1.jpg",
      "thumbnail_url": "/uploads/thumbnails/photo_id_1_thumb.jpg",
      "category": "workshop",
      "size": 1024000,
      "factory_id": "gongchang",
      "status": "success",
      "created_at": "2025-01-15T10:30:00Z"
    }
  ],
  "categories": [
    {
      "id": 1,
      "factory_id": "gongchang",
      "name": "workshop",
      "color": "#FF5733",
      "count": 8
    }
  ]
}
```

### 3. 删除单张工厂图片

**接口地址**: `DELETE /api/factories/{factoryId}/photos/{photoId}`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `factoryId`: 工厂ID (必填)
- `photoId`: 图片ID (必填)

**请求示例**:
```bash
curl -X DELETE "http://localhost:8008/api/factories/gongchang/photos/photo_id_1" \
  -H "Authorization: Bearer {token}"
```

**成功响应** (200):
```json
{
  "success": true,
  "message": "图片删除成功"
}
```

### 4. 批量删除工厂图片

**接口地址**: `DELETE /api/factories/{factoryId}/photos/batch`

**请求头**:
```
Content-Type: application/json
Authorization: Bearer {token}
```

**路径参数**:
- `factoryId`: 工厂ID (必填)

**请求体**:
```json
{
  "photo_ids": ["photo_id_1", "photo_id_2", "photo_id_3"]
}
```

**请求示例**:
```bash
curl -X DELETE "http://localhost:8008/api/factories/gongchang/photos/batch" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "photo_ids": ["photo_id_1", "photo_id_2"]
  }'
```

**成功响应** (200):
```json
{
  "success": true,
  "message": "批量删除成功",
  "deleted_count": 2,
  "failed_count": 0,
  "failed_photo_ids": []
}
```

## 功能特性

### 1. 文件验证
- **文件格式**: 支持 JPG, JPEG, PNG, WebP
- **文件大小**: 每个文件最大 10MB
- **文件类型验证**: 检查文件头确保是真实图片

### 2. 安全特性
- **权限验证**: 只能操作自己的工厂图片
- **文件名安全**: 使用UUID避免文件名冲突
- **路径安全**: 限制文件保存路径，防止目录遍历攻击

### 3. 图片处理
- **自动缩略图**: 大于1MB的图片自动生成缩略图
- **分类管理**: 支持按分类组织图片
- **批量操作**: 支持批量上传和删除

### 4. 错误处理
- **错误隔离**: 单个文件失败不影响其他文件
- **详细错误信息**: 返回具体的失败原因
- **事务处理**: 确保数据一致性

## 数据模型

### File 模型扩展
```go
type File struct {
    ID        string     `json:"id"`
    Name      string     `json:"name"`
    Path      string     `json:"path"`
    Type      string     `json:"type"`                    // 文件类型
    OrderID   *uint      `json:"order_id,omitempty"`     // 关联订单
    FactoryID string     `json:"factory_id,omitempty"`   // 关联工厂
    Category  string     `json:"category,omitempty"`     // 图片分类
    Size      int64      `json:"size,omitempty"`         // 文件大小
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
}
```

### 图片分类
系统预定义了以下图片分类：
- `workshop`: 车间照片
- `equipment`: 设备照片
- `products`: 产品照片
- `certificates`: 证书照片

## 使用示例

### Flutter/Dart 示例

```dart
// 批量上传图片
Future<Map<String, dynamic>> uploadFactoryPhotos({
  required String factoryId,
  required List<String> filePaths,
  String? category,
}) async {
  try {
    final request = http.MultipartRequest(
      'POST',
      Uri.parse('$baseUrl/api/factories/$factoryId/photos/batch'),
    );
    
    request.headers['Authorization'] = 'Bearer $token';
    
    for (String filePath in filePaths) {
      final file = File(filePath);
      final stream = http.ByteStream(file.openRead());
      final length = await file.length();
      
      request.files.add(http.MultipartFile(
        'files',
        stream,
        length,
        filename: path.basename(filePath),
      ));
    }
    
    if (category != null) {
      request.fields['category'] = category;
    }
    
    final response = await request.send();
    final responseBody = await response.stream.bytesToString();
    
    if (response.statusCode == 200) {
      return json.decode(responseBody);
    } else {
      throw Exception('上传失败');
    }
  } catch (e) {
    throw Exception('上传失败: $e');
  }
}

// 获取图片列表
Future<Map<String, dynamic>> getFactoryPhotos({
  required String factoryId,
  String? category,
  int page = 1,
  int pageSize = 20,
}) async {
  try {
    final queryParams = <String, String>{
      'page': page.toString(),
      'page_size': pageSize.toString(),
    };
    
    if (category != null) {
      queryParams['category'] = category;
    }
    
    final uri = Uri.parse('$baseUrl/api/factories/$factoryId/photos')
        .replace(queryParameters: queryParams);
    
    final response = await http.get(
      uri,
      headers: {
        'Authorization': 'Bearer $token',
      },
    );
    
    if (response.statusCode == 200) {
      return json.decode(response.body);
    } else {
      throw Exception('获取图片列表失败');
    }
  } catch (e) {
    throw Exception('获取图片列表失败: $e');
  }
}
```

### JavaScript 示例

```javascript
// 批量上传图片
async function uploadFactoryPhotos(factoryId, files, category = null) {
  const formData = new FormData();
  
  files.forEach(file => {
    formData.append('files', file);
  });
  
  if (category) {
    formData.append('category', category);
  }
  
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
    throw new Error(`上传失败: ${error.message}`);
  }
}

// 获取图片列表
async function getFactoryPhotos(factoryId, category = null, page = 1, pageSize = 20) {
  const params = new URLSearchParams({
    page: page.toString(),
    page_size: pageSize.toString()
  });
  
  if (category) {
    params.append('category', category);
  }
  
  try {
    const response = await fetch(`/api/factories/${factoryId}/photos?${params}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    const result = await response.json();
    return result;
  } catch (error) {
    throw new Error(`获取图片列表失败: ${error.message}`);
  }
}
```

## 部署说明

### 1. 数据库迁移
确保运行了数据库迁移脚本，添加了新的字段：
```sql
ALTER TABLE files ADD COLUMN type VARCHAR(50) NULL;
ALTER TABLE files ADD COLUMN factory_id VARCHAR(191) NULL;
ALTER TABLE files ADD COLUMN category VARCHAR(100) NULL;
ALTER TABLE files ADD COLUMN size BIGINT NULL;
```

### 2. 目录权限
确保上传目录有正确的权限：
```bash
mkdir -p uploads/thumbnails
chmod 755 uploads
chmod 755 uploads/thumbnails
```

### 3. 配置检查
检查配置文件中的上传路径设置：
```yaml
upload:
  path: "./uploads"
  max_size: 10
```

## 测试

运行测试脚本验证功能：
```bash
cd backend
chmod +x test_factory_photos_api.sh
./test_factory_photos_api.sh
```

## 注意事项

1. **文件大小限制**: 每个图片文件最大10MB
2. **权限控制**: 只能操作自己的工厂图片
3. **分类管理**: 建议使用预定义的分类名称
4. **缩略图**: 大于1MB的图片会自动生成缩略图
5. **错误处理**: 批量操作中单个文件失败不影响其他文件 