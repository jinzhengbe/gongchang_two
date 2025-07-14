# 设计师头像上传API文档

## 概述

本文档介绍设计师头像上传和设计师信息管理的API接口，包括头像上传、设计师信息获取和更新功能。

## API端点

### 1. 上传头像

**POST** `/api/designers/avatar`

上传设计师头像文件，返回头像URL。

#### 请求头
```
Content-Type: multipart/form-data
Authorization: Bearer {token}
```

#### 请求参数
- `avatar`: 头像文件 (必填，支持 JPG, PNG, WebP 格式，最大 5MB)

#### 响应格式

**成功响应 (200)**
```json
{
  "success": true,
  "message": "头像上传成功",
  "data": {
    "url": "/uploads/avatars/{filename}"
  }
}
```

**错误响应**

- **400** - 文件格式不支持
```json
{
  "error": "不支持的文件格式，支持: JPG, PNG, WebP"
}
```

- **400** - 文件过大
```json
{
  "error": "文件大小超过限制 (最大 5MB)"
}
```

- **401** - 未授权
```json
{
  "error": "未授权"
}
```

#### 使用示例

```bash
# 上传头像
curl -X POST "http://localhost:8008/api/designers/avatar" \
  -H "Authorization: Bearer {token}" \
  -F "avatar=@/path/to/avatar.jpg"
```

### 2. 获取设计师信息

**GET** `/api/designers/profile`

获取当前登录设计师的详细信息。

#### 请求头
```
Authorization: Bearer {token}
Content-Type: application/json
```

#### 响应格式

**成功响应 (200)**
```json
{
  "success": true,
  "data": {
    "ID": 1,
    "UserID": "designer-user-id",
    "CompanyName": "设计工作室",
    "Address": "北京市朝阳区",
    "Website": "http://design-studio.com",
    "Bio": "专业服装设计工作室",
    "Avatar": "/uploads/avatars/avatar-123.jpg",
    "Rating": 4.5,
    "Status": 1,
    "CreatedAt": "2025-01-01T00:00:00Z",
    "UpdatedAt": "2025-01-01T00:00:00Z",
    "User": {
      "id": "designer-user-id",
      "username": "designer1",
      "email": "designer@test.com",
      "role": "designer"
    }
  }
}
```

**错误响应**

- **401** - 未授权
```json
{
  "error": "未授权"
}
```

- **404** - 设计师档案不存在
```json
{
  "error": "设计师档案不存在"
}
```

#### 使用示例

```bash
# 获取设计师信息
curl -X GET "http://localhost:8008/api/designers/profile" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json"
```

### 3. 更新设计师信息

**PUT** `/api/designers/profile`

更新设计师的详细信息，包括头像URL。

#### 请求头
```
Authorization: Bearer {token}
Content-Type: application/json
```

#### 请求参数

```json
{
  "company_name": "新设计工作室名称",
  "address": "新地址",
  "website": "http://new-website.com",
  "bio": "新的个人简介",
  "avatar": "/uploads/avatars/new-avatar.jpg"
}
```

**参数说明：**
- `company_name`: 公司名称 (可选)
- `address`: 地址 (可选)
- `website`: 网站 (可选)
- `bio`: 个人简介 (可选)
- `avatar`: 头像URL (可选)

#### 响应格式

**成功响应 (200)**
```json
{
  "success": true,
  "message": "设计师信息更新成功",
  "data": {
    "ID": 1,
    "UserID": "designer-user-id",
    "CompanyName": "新设计工作室名称",
    "Address": "新地址",
    "Website": "http://new-website.com",
    "Bio": "新的个人简介",
    "Avatar": "/uploads/avatars/new-avatar.jpg",
    "Rating": 4.5,
    "Status": 1,
    "CreatedAt": "2025-01-01T00:00:00Z",
    "UpdatedAt": "2025-01-01T00:00:00Z"
  }
}
```

**错误响应**

- **400** - 请求参数错误
```json
{
  "error": "请求参数错误: {具体错误信息}"
}
```

- **401** - 未授权
```json
{
  "error": "未授权"
}
```

- **404** - 设计师档案不存在
```json
{
  "error": "设计师档案不存在"
}
```

#### 使用示例

```bash
# 更新设计师信息
curl -X PUT "http://localhost:8008/api/designers/profile" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "company_name": "新设计工作室",
    "address": "北京市朝阳区",
    "website": "http://new-design.com",
    "bio": "专业服装设计工作室",
    "avatar": "/uploads/avatars/avatar-123.jpg"
  }'
```

## 前端使用流程

### 1. 头像上传流程

```dart
Future<String?> uploadAvatar(File avatarFile) async {
  try {
    // 创建multipart request
    var request = http.MultipartRequest(
      'POST',
      Uri.parse('$baseUrl/api/designers/avatar'),
    );

    // 添加认证头
    request.headers['Authorization'] = 'Bearer $token';

    // 添加文件
    var stream = http.ByteStream(avatarFile.openRead());
    var length = await avatarFile.length();
    var multipartFile = http.MultipartFile(
      'avatar',
      stream,
      length,
      filename: avatarFile.path.split('/').last,
    );
    request.files.add(multipartFile);

    // 发送请求
    var response = await request.send();
    var responseData = await response.stream.bytesToString();
    var jsonData = jsonDecode(responseData);

    if (response.statusCode == 200 && jsonData['success'] == true) {
      return jsonData['data']['url']; // 返回头像URL
    } else {
      throw Exception(jsonData['error'] ?? '上传失败');
    }
  } catch (e) {
    print('头像上传失败: $e');
    return null;
  }
}
```

### 2. 保存设计师信息流程

```dart
Future<void> saveDesignerInfo({
  String? companyName,
  String? address,
  String? website,
  String? bio,
  String? avatarUrl,
}) async {
  try {
    final response = await http.put(
      Uri.parse('$baseUrl/api/designers/profile'),
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer $token',
      },
      body: jsonEncode({
        if (companyName != null) 'company_name': companyName,
        if (address != null) 'address': address,
        if (website != null) 'website': website,
        if (bio != null) 'bio': bio,
        if (avatarUrl != null) 'avatar': avatarUrl,
      }),
    );

    if (response.statusCode == 200) {
      print('设计师信息更新成功');
    } else {
      throw Exception('更新失败');
    }
  } catch (e) {
    print('保存设计师信息失败: $e');
    rethrow;
  }
}
```

### 3. 完整使用示例

```dart
Future<void> updateDesignerWithAvatar(File avatarFile) async {
  // 1. 上传头像
  String? avatarUrl = await uploadAvatar(avatarFile);
  
  if (avatarUrl != null) {
    // 2. 更新设计师信息
    await saveDesignerInfo(
      companyName: '新设计工作室',
      address: '北京市朝阳区',
      website: 'http://new-design.com',
      bio: '专业服装设计工作室',
      avatarUrl: avatarUrl,
    );
    
    print('设计师信息和头像更新成功');
  } else {
    print('头像上传失败');
  }
}
```

## 数据库结构

### designer_profiles 表

```sql
CREATE TABLE designer_profiles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id VARCHAR(191) UNIQUE,
    company_name VARCHAR(255),
    address VARCHAR(255),
    website VARCHAR(255),
    bio TEXT,
    avatar VARCHAR(500) DEFAULT "",  -- 新增：头像URL
    rating DECIMAL(3,2) DEFAULT 0,
    status INT DEFAULT 1,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

## 文件存储

头像文件存储在 `./uploads/avatars/` 目录下，通过 `/uploads/avatars/{filename}` 路径访问。

## 注意事项

1. **文件格式**: 只支持 JPG, PNG, WebP 格式
2. **文件大小**: 最大 5MB
3. **认证**: 所有API都需要有效的JWT token
4. **权限**: 只能操作自己的设计师档案
5. **URL访问**: 头像URL通过nginx静态文件服务提供访问

## 测试

运行测试脚本验证功能：

```bash
cd backend
chmod +x test_avatar_upload.sh
./test_avatar_upload.sh
``` 