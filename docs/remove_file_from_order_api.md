# DELETE /api/orders/{id}/remove-file API 规范

## 接口概述
从指定订单中移除指定的文件（图片、附件、模型或视频）。

## 请求格式

### URL
```
DELETE /api/orders/{id}/remove-file
```

### 请求头
```
Content-Type: application/json
Authorization: Bearer {token}
```

### 路径参数
- `id` (integer, required): 订单ID（数字类型）

### 请求体 (JSON)
```json
{
  "fileId": "string",     // 文件ID（UUID格式）
  "fileType": "string"    // 文件类型：image, attachment, model, video
}
```

## 字段说明

| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| `fileId` | string | 是 | 要移除的文件ID（UUID格式） |
| `fileType` | string | 是 | 文件类型，支持：image, attachment, model, video |

## 响应格式

### 成功响应 (200 OK)
```json
{
  "success": true,
  "message": "文件已从订单的image中移除",
  "order": {
    "id": 29,
    "title": "订单标题",
    "description": "订单描述",
    "images": ["file1_id", "file2_id"],
    "attachments": [],
    "models": [],
    "videos": [],
    // ... 其他订单字段
  }
}
```

### 错误响应

#### 400 Bad Request
```json
{
  "error": "无效的订单ID"
}
```

```json
{
  "error": "文件ID不能为空"
}
```

```json
{
  "error": "文件类型不能为空"
}
```

```json
{
  "error": "无效的文件类型，支持的类型：image, attachment, model, video"
}
```

#### 404 Not Found
```json
{
  "error": "订单不存在"
}
```

```json
{
  "error": "文件不存在"
}
```

```json
{
  "error": "订单中不包含该文件"
}
```

#### 500 Internal Server Error
```json
{
  "error": "服务器内部错误"
}
```

## 使用示例

### 删除订单中的图片
```bash
curl -X DELETE "https://aneworders.com/api/orders/29/remove-file" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "fileId": "fbdd3f3e-0a2a-4180-8478-7e334e7d9fe7",
    "fileType": "image"
  }'
```

### 删除订单中的附件
```bash
curl -X DELETE "https://aneworders.com/api/orders/29/remove-file" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "fileId": "abc123-def456-ghi789",
    "fileType": "attachment"
  }'
```

## 注意事项

1. **订单ID必须是数字类型**：路径参数中的订单ID必须是数字，不能是UUID
2. **文件ID必须是UUID格式**：请求体中的fileId必须是有效的UUID格式
3. **文件类型必须明确指定**：fileType字段必须指定为image、attachment、model或video之一
4. **文件必须存在于订单中**：要删除的文件ID必须存在于订单的对应文件数组中
5. **事务安全**：操作在数据库事务中执行，确保数据一致性

## 与现有API的区别

- **DELETE /api/orders/{id}**：删除整个订单（订单ID必须是数字）
- **DELETE /api/orders/{id}/remove-file**：从订单中移除指定文件（订单ID是数字，文件ID是UUID）
- **DELETE /api/files/{fileId}**：删除文件记录（文件ID是UUID） 