# 订单进度管理 API 文档

## 概述

订单进度管理模块用于跟踪和管理订单的生产进度。工厂可以为订单创建多个进度记录，记录不同阶段的完成情况。

## 数据库表结构

### order_progress 表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| `id` | bigint unsigned | 主键，自增 |
| `order_id` | bigint unsigned | 订单ID，外键关联orders表 |
| `factory_id` | varchar(191) | 工厂ID |
| `progress_type` | varchar(50) | 进度类型 |
| `percentage` | int | 完成百分比(0-100) |
| `status` | varchar(50) | 进度状态 |
| `description` | text | 进度描述 |
| `estimated_completion_time` | datetime(3) | 预计完成时间 |
| `actual_completion_time` | datetime(3) | 实际完成时间 |
| `creator_id` | varchar(191) | 创建者ID |
| `created_at` | datetime(3) | 创建时间 |
| `updated_at` | datetime(3) | 更新时间 |
| `deleted_at` | datetime(3) | 删除时间（软删除） |

## 进度类型说明

### ProgressType

- `design`: 设计阶段
- `material`: 材料准备
- `production`: 生产阶段
- `quality`: 质检阶段
- `packaging`: 包装阶段
- `shipping`: 发货阶段
- `custom`: 自定义阶段

## 进度状态说明

### ProgressStatus

- `not_started`: 未开始
- `in_progress`: 进行中
- `completed`: 已完成
- `delayed`: 延期
- `on_hold`: 暂停

## API 接口

### 1. 创建进度记录

**接口地址：** `POST /api/orders/{orderId}/progress`

**请求头：**
```
Content-Type: application/json
Authorization: Bearer {token}
```

**请求体：**
```json
{
  "order_id": 123,
  "factory_id": "factory_user_id",
  "progress_type": "production",
  "percentage": 50,
  "status": "in_progress",
  "description": "生产进度过半，质量良好",
  "estimated_completion_time": "2025-07-15T10:00:00Z",
  "actual_completion_time": null,
  "creator_id": "factory_user_id"
}
```

**响应：**
```json
{
  "id": 1,
  "order_id": 123,
  "factory_id": "factory_user_id",
  "progress_type": "production",
  "percentage": 50,
  "status": "in_progress",
  "description": "生产进度过半，质量良好",
  "estimated_completion_time": "2025-07-15T10:00:00Z",
  "actual_completion_time": null,
  "creator_id": "factory_user_id",
  "created_at": "2025-06-29T10:30:00Z",
  "updated_at": "2025-06-29T10:30:00Z"
}
```

**权限要求：** 只有工厂用户可以创建进度记录

### 2. 获取订单进度列表

**接口地址：** `GET /api/orders/{orderId}/progress`

**请求头：**
```
Authorization: Bearer {token}
```

**响应：**
```json
[
  {
    "id": 1,
    "order_id": 123,
    "factory_id": "factory_user_id",
    "progress_type": "production",
    "percentage": 50,
    "status": "in_progress",
    "description": "生产进度过半，质量良好",
    "estimated_completion_time": "2025-07-15T10:00:00Z",
    "actual_completion_time": null,
    "creator_id": "factory_user_id",
    "created_at": "2025-06-29T10:30:00Z",
    "updated_at": "2025-06-29T10:30:00Z",
    "order": {
      "id": 123,
      "title": "订单标题",
      "description": "订单描述"
    },
    "factory": {
      "user_id": "factory_user_id",
      "company_name": "工厂名称"
    }
  }
]
```

### 3. 更新进度记录

**接口地址：** `PUT /api/orders/{orderId}/progress/{progressId}`

**请求头：**
```
Content-Type: application/json
Authorization: Bearer {token}
```

**请求体：**
```json
{
  "progress_type": "production",
  "percentage": 75,
  "status": "in_progress",
  "description": "生产进度75%，即将完成",
  "estimated_completion_time": "2025-07-10T10:00:00Z",
  "actual_completion_time": null
}
```

**响应：**
```json
{
  "id": 1,
  "order_id": 123,
  "factory_id": "factory_user_id",
  "progress_type": "production",
  "percentage": 75,
  "status": "in_progress",
  "description": "生产进度75%，即将完成",
  "estimated_completion_time": "2025-07-10T10:00:00Z",
  "actual_completion_time": null,
  "creator_id": "factory_user_id",
  "created_at": "2025-06-29T10:30:00Z",
  "updated_at": "2025-06-29T11:00:00Z"
}
```

**权限要求：** 只有工厂用户可以更新自己创建的进度记录

### 4. 删除进度记录

**接口地址：** `DELETE /api/orders/{orderId}/progress/{progressId}`

**请求头：**
```
Authorization: Bearer {token}
```

**响应：**
```json
{
  "message": "进度记录删除成功"
}
```

**权限要求：** 只有工厂用户可以删除自己创建的进度记录

### 5. 获取工厂进度列表

**接口地址：** `GET /api/factories/{factoryId}/progress?page=1&pageSize=10`

**请求头：**
```
Authorization: Bearer {token}
```

**查询参数：**
- `page`: 页码（默认1）
- `pageSize`: 每页数量（默认10）

**响应：**
```json
{
  "total": 25,
  "page": 1,
  "page_size": 10,
  "progress": [
    {
      "id": 1,
      "order_id": 123,
      "factory_id": "factory_user_id",
      "progress_type": "production",
      "percentage": 75,
      "status": "in_progress",
      "description": "生产进度75%，即将完成",
      "estimated_completion_time": "2025-07-10T10:00:00Z",
      "actual_completion_time": null,
      "creator_id": "factory_user_id",
      "created_at": "2025-06-29T10:30:00Z",
      "updated_at": "2025-06-29T11:00:00Z",
      "order": {
        "id": 123,
        "title": "订单标题"
      },
      "factory": {
        "user_id": "factory_user_id",
        "company_name": "工厂名称"
      }
    }
  ]
}
```

**权限要求：** 只能查看自己工厂的进度记录

### 6. 获取进度统计信息

**接口地址：** `GET /api/factories/{factoryId}/progress-statistics`

**请求头：**
```
Authorization: Bearer {token}
```

**响应：**
```json
{
  "not_started": 5,
  "in_progress": 10,
  "completed": 8,
  "delayed": 2,
  "on_hold": 1
}
```

**权限要求：** 只能查看自己工厂的统计信息

## 业务规则

1. **权限控制**：只有工厂用户可以创建、更新、删除进度记录
2. **数据一致性**：进度记录必须属于指定的订单
3. **身份验证**：工厂只能以自己的身份进行进度管理
4. **百分比范围**：完成百分比必须在0-100之间
5. **时间逻辑**：实际完成时间不能早于创建时间

## 错误码

| 错误码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 权限不足 |
| 404 | 进度记录不存在 |
| 422 | 数据验证失败 |
| 500 | 服务器内部错误 |

## 使用示例

### 工厂进度管理流程

1. **创建进度记录**
```bash
curl -X POST "https://aneworders.com/api/orders/123/progress" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "order_id": 123,
    "factory_id": "factory_user_id",
    "progress_type": "production",
    "percentage": 50,
    "status": "in_progress",
    "description": "生产进度过半，质量良好",
    "estimated_completion_time": "2025-07-15T10:00:00Z",
    "creator_id": "factory_user_id"
  }'
```

2. **查看订单进度**
```bash
curl -X GET "https://aneworders.com/api/orders/123/progress" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

3. **更新进度**
```bash
curl -X PUT "https://aneworders.com/api/orders/123/progress/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "percentage": 75,
    "description": "生产进度75%，即将完成"
  }'
```

4. **查看工厂所有进度**
```bash
curl -X GET "https://aneworders.com/api/factories/factory_user_id/progress?page=1&pageSize=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

5. **查看进度统计**
```bash
curl -X GET "https://aneworders.com/api/factories/factory_user_id/progress-statistics" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

6. **删除进度记录**
```bash
curl -X DELETE "https://aneworders.com/api/orders/123/progress/1" \
  -H "Authorization: Bearer YOUR_TOKEN"
``` 