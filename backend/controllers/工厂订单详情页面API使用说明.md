# 工厂订单详情页面API使用说明

## 概述

本文档详细说明了工厂订单详情页面中 `checkAcceptOrderStatus()` 和 `fetchOrderProgress()` 两个方法对应的后端API接口。

## 1. checkAcceptOrderStatus() 方法对应的API

### 1.1 主要API接口

#### 获取订单接单记录
- **接口路径**: `GET /api/orders/{orderId}/jiedans`
- **请求方法**: GET
- **认证**: 需要Bearer Token
- **参数**: 
  - `orderId` (路径参数): 订单ID (数字格式)

#### 请求示例
```bash
curl -X GET "https://aneworders.com/api/orders/39/jiedans" \
  -H "Accept: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### 响应格式
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "order_id": 39,
      "factory_id": 123,
      "status": "accepted",
      "message": "我们接受这个订单",
      "price_quote": 5000.0,
      "estimated_delivery_date": "2024-02-15",
      "jiedan_time": "2024-01-15T10:30:00Z",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

### 1.2 相关API接口

#### 创建接单记录
- **接口路径**: `POST /api/jiedan`
- **请求方法**: POST
- **认证**: 需要Bearer Token
- **请求体**:
```json
{
  "order_id": 39,
  "factory_id": 123,
  "status": "pending",
  "message": "我们愿意接受这个订单",
  "price_quote": 5000.0,
  "estimated_delivery_date": "2024-02-15"
}
```

#### 工厂接单API
- **接口路径**: `POST /api/orders/{orderId}/accept`
- **请求方法**: POST
- **认证**: 需要Bearer Token
- **请求体**:
```json
{
  "order_id": 39,
  "factory_id": 123,
  "status": "accepted",
  "accepted_at": "2024-01-15T10:30:00Z",
  "action": "accept_order"
}
```

## 2. fetchOrderProgress() 方法对应的API

### 2.1 主要API接口

#### 获取订单进度记录
- **接口路径**: `GET /api/orders/{orderId}/progress`
- **请求方法**: GET
- **认证**: 需要Bearer Token
- **参数**: 
  - `orderId` (路径参数): 订单ID (数字格式)

#### 请求示例
```bash
curl -X GET "https://aneworders.com/api/orders/39/progress" \
  -H "Accept: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### 响应格式
```json
{
  "success": true,
  "data": [
    {
      "id": "1",
      "order_id": "39",
      "type": "design",
      "status": "completed",
      "description": "设计阶段已完成",
      "start_time": "2024-01-10T09:00:00Z",
      "completed_time": "2024-01-12T17:30:00Z",
      "created_at": "2024-01-10T09:00:00Z",
      "updated_at": "2024-01-12T17:30:00Z"
    },
    {
      "id": "2",
      "order_id": "39",
      "type": "material",
      "status": "in_progress",
      "description": "材料采购进行中",
      "start_time": "2024-01-13T08:00:00Z",
      "completed_time": null,
      "created_at": "2024-01-13T08:00:00Z",
      "updated_at": "2024-01-14T16:45:00Z"
    }
  ]
}
```

### 2.2 相关API接口

#### 创建进度记录
- **接口路径**: `POST /api/progress`
- **请求方法**: POST
- **认证**: 需要Bearer Token
- **请求体**:
```json
{
  "order_id": 39,
  "type": "production",
  "status": "in_progress",
  "description": "开始生产阶段",
  "start_time": "2024-01-15T09:00:00Z",
  "completed_time": null
}
```

#### 更新进度记录
- **接口路径**: `PUT /api/progress/{progressId}`
- **请求方法**: PUT
- **认证**: 需要Bearer Token
- **请求体**:
```json
{
  "type": "production",
  "status": "completed",
  "description": "生产阶段已完成",
  "start_time": "2024-01-15T09:00:00Z",
  "completed_time": "2024-01-20T17:00:00Z"
}
```

#### 删除进度记录
- **接口路径**: `DELETE /api/progress/{progressId}`
- **请求方法**: DELETE
- **认证**: 需要Bearer Token

## 3. 前端实现逻辑

### 3.1 checkAcceptOrderStatus() 实现流程

1. **调用 `fetchJiedanInfo()` 方法**
2. **发送请求**: `GET /api/orders/{orderId}/jiedans`
3. **处理响应**:
   - 如果成功且有接单记录，设置 `hasAcceptedOrder.value = true`
   - 如果失败或无记录，设置 `hasAcceptedOrder.value = false`
4. **更新UI状态**: 根据接单状态显示不同的按钮和内容

### 3.2 fetchOrderProgress() 实现流程

1. **发送请求**: `GET /api/orders/{orderId}/progress`
2. **处理响应**:
   - 如果成功，解析进度记录并更新 `progressRecords`
   - 如果失败，使用本地模拟数据
3. **更新UI**: 显示进度列表和添加进度按钮

## 4. 错误处理

### 4.1 常见错误码

- **400**: 请求参数错误
- **401**: 未授权，需要重新登录
- **404**: 订单不存在或API路径错误
- **500**: 服务器内部错误

### 4.2 前端错误处理策略

1. **API失败时使用模拟数据**: 确保用户体验不受影响
2. **显示错误提示**: 通过snackbar通知用户
3. **重试机制**: 提供手动刷新功能

## 5. 认证要求

所有API调用都需要在请求头中包含有效的Bearer Token：

```bash
Authorization: Bearer YOUR_JWT_TOKEN
```

## 6. 数据格式要求

### 6.1 订单ID格式
- 必须是数字格式
- 前端会自动验证和转换

### 6.2 时间格式
- 使用ISO 8601格式: `YYYY-MM-DDTHH:mm:ssZ`
- 例如: `2024-01-15T10:30:00Z`

### 6.3 进度状态
- `not_started`: 未开始
- `in_progress`: 进行中
- `completed`: 已完成
- `cancelled`: 已取消

## 7. 测试建议

### 7.1 API测试
```bash
# 测试获取接单记录
curl -X GET "https://aneworders.com/api/orders/39/jiedans" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 测试获取进度记录
curl -X GET "https://aneworders.com/api/orders/39/progress" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 7.2 前端测试
1. 使用工厂账号登录
2. 访问订单详情页面
3. 检查接单状态显示
4. 测试添加进度功能

## 8. 注意事项

1. **订单ID验证**: 确保订单ID是数字格式
2. **工厂权限**: 只有已接单的工厂才能添加进度
3. **数据一致性**: 前端会缓存数据，需要及时刷新
4. **错误回退**: API失败时会使用模拟数据确保功能可用

## 9. 相关文件

- **前端控制器**: `lib/features/factory/controller/factory_order_detail_controller.dart`
- **API服务**: `lib/features/order/services/order_service.dart`
- **进度模型**: `lib/features/order/models/progress_record.dart` 