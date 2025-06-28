# PUT /api/orders/{id} API 规范

## 接口概述
更新指定订单的信息，支持部分字段更新。

## 请求格式

### URL
```
PUT /api/orders/{id}
```

### 请求头
```
Content-Type: application/json
Authorization: Bearer {token}
```

### 路径参数
- `id` (integer, required): 订单ID

### 请求体 (JSON)
```json
{
  "title": "string",
  "description": "string", 
  "fabric": "string",
  "quantity": "integer",
  "status": "string",
  "payment_status": "string",
  "shipping_address": "string",
  "orderType": "string",
  "fabrics": "string",
  "deliveryDate": "string (ISO 8601)",
  "order_date": "string (ISO 8601)",
  "specialRequirements": "string",
  "attachments": ["string"],
  "models": ["string"],
  "images": ["string"],
  "videos": ["string"]
}
```

## 字段说明

### 基础字段
| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| `title` | string | 否 | 订单标题 |
| `description` | string | 否 | 订单描述 |
| `fabric` | string | 否 | 面料信息 |
| `quantity` | integer | 否 | 数量 |
| `status` | string | 否 | 订单状态 (draft/published/completed/cancelled) |
| `payment_status` | string | 否 | 支付状态 |
| `shipping_address` | string | 否 | 收货地址 |
| `orderType` | string | 否 | 订单类型 |
| `fabrics` | string | 否 | 布料信息 |
| `deliveryDate` | string | 否 | 交货日期 (ISO 8601格式) |
| `order_date` | string | 否 | 订单日期 (ISO 8601格式) |
| `specialRequirements` | string | 否 | 特殊要求 |

### 文件字段
| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| `attachments` | array[string] | 否 | 附件文件ID列表 |
| `models` | array[string] | 否 | 模型文件ID列表 |
| `images` | array[string] | 否 | 图片文件ID列表 |
| `videos` | array[string] | 否 | 视频文件ID列表 |

## 响应格式

### 成功响应 (200 OK)
```json
{
  "message": "Order updated successfully"
}
```

### 错误响应

#### 400 Bad Request
```json
{
  "error": "Invalid order ID"
}
```
或
```json
{
  "error": "Invalid JSON format"
}
```

#### 401 Unauthorized
```json
{
  "error": "未授权"
}
```

#### 500 Internal Server Error
```json
{
  "error": "Database error message"
}
```

## 使用示例

### 1. 更新订单标题
```bash
curl -X PUT "http://localhost:8008/api/orders/123" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_token_here" \
  -d '{
    "title": "更新后的订单标题"
  }'
```

### 2. 添加图片到订单
```bash
curl -X PUT "http://localhost:8008/api/orders/123" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_token_here" \
  -d '{
    "images": ["file_001", "file_002", "file_003"]
  }'
```

### 3. 更新多个字段
```bash
curl -X PUT "http://localhost:8008/api/orders/123" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_token_here" \
  -d '{
    "title": "新标题",
    "description": "新描述",
    "quantity": 10,
    "status": "published",
    "images": ["img_001", "img_002"]
  }'
```

### 4. 只更新图片（不传其他字段）
```bash
curl -X PUT "http://localhost:8008/api/orders/123" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_token_here" \
  -d '{
    "images": ["new_image_001"]
  }'
```

## 重要注意事项

### 1. 图片字段合并逻辑
- 当传递 `images` 字段时，新图片ID会与现有图片ID合并
- 重复的图片ID会自动去重
- 当不传递 `images` 字段时，现有图片保持不变

### 2. 部分更新支持
- 可以只传递需要更新的字段
- 未传递的字段保持原值不变
- 所有字段都是可选的

### 3. 认证要求
- 必须提供有效的 JWT token
- Token 必须在 Authorization 头中以 `Bearer {token}` 格式传递

### 4. 数据验证
- 订单ID必须是有效的整数
- JSON 格式必须正确
- 日期字段必须是 ISO 8601 格式

### 5. 错误处理
- 如果订单不存在，返回 404 错误
- 如果用户无权限，返回 401 错误
- 如果数据格式错误，返回 400 错误

## 前端集成建议

### JavaScript/TypeScript 示例
```javascript
// 更新订单
async function updateOrder(orderId, updateData) {
  try {
    const response = await fetch(`/api/orders/${orderId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getToken()}`
      },
      body: JSON.stringify(updateData)
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const result = await response.json();
    return result;
  } catch (error) {
    console.error('更新订单失败:', error);
    throw error;
  }
}

// 使用示例
// 1. 更新标题
await updateOrder(123, { title: '新标题' });

// 2. 添加图片
await updateOrder(123, { images: ['file_001', 'file_002'] });

// 3. 更新多个字段
await updateOrder(123, {
  title: '新标题',
  description: '新描述',
  images: ['img_001', 'img_002']
});
```

### Flutter/Dart 示例
```dart
Future<Map<String, dynamic>> updateOrder(int orderId, Map<String, dynamic> updateData) async {
  try {
    final response = await http.put(
      Uri.parse('http://localhost:8008/api/orders/$orderId'),
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer $token',
      },
      body: jsonEncode(updateData),
    );
    
    if (response.statusCode == 200) {
      return jsonDecode(response.body);
    } else {
      throw Exception('更新失败: ${response.statusCode}');
    }
  } catch (e) {
    print('更新订单错误: $e');
    rethrow;
  }
}

// 使用示例
await updateOrder(123, {
  'title': '新标题',
  'images': ['file_001', 'file_002']
});
```

## 常见问题

### Q: 为什么图片会丢失？
A: 确保每次调用时传递完整的图片ID列表，或者不传递 `images` 字段让后端保持现有图片。

### Q: 如何只更新部分字段？
A: 只传递需要更新的字段，其他字段会自动保持原值。

### Q: 日期格式要求是什么？
A: 使用 ISO 8601 格式，例如：`"2025-06-27T15:30:00Z"`

### Q: 如何处理认证错误？
A: 检查 token 是否有效，确保在 Authorization 头中正确传递。 