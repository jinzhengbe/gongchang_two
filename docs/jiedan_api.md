# 接单管理 API 文档

## 概述

接单管理模块用于处理工厂对订单的接单、同意、拒绝等操作。接单表 `jiedan` 记录了工厂对订单的接单状态和相关信息。

## 数据库表结构

### jiedan 表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| `id` | bigint unsigned | 主键，自增 |
| `order_id` | bigint unsigned | 订单ID，外键关联orders表 |
| `factory_id` | varchar(191) | 工厂ID |
| `status` | varchar(50) | 状态：pending-待处理, accepted-已同意, rejected-已拒绝 |
| `price` | decimal(10,2) | 接单价格 |
| `jiedan_time` | datetime(3) | 接单时间 |
| `agree_time` | datetime(3) | 同意时间 |
| `agree_user_id` | varchar(191) | 同意的用户ID |
| `created_at` | datetime(3) | 创建时间 |
| `updated_at` | datetime(3) | 更新时间 |
| `deleted_at` | datetime(3) | 删除时间（软删除） |

## API 接口

### 1. 创建接单记录

**接口地址：** `POST /api/jiedan`

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
  "price": 1500.50
}
```

**响应：**
```json
{
  "id": 1,
  "order_id": 123,
  "factory_id": "factory_user_id",
  "status": "pending",
  "price": 1500.50,
  "jiedan_time": "2025-06-28T10:30:00Z",
  "agree_time": null,
  "agree_user_id": null,
  "created_at": "2025-06-28T10:30:00Z",
  "updated_at": "2025-06-28T10:30:00Z"
}
```

**权限要求：** 只有工厂用户可以进行接单操作

### 2. 获取接单记录详情

**接口地址：** `GET /api/jiedan/{id}`

**请求头：**
```
Authorization: Bearer {token}
```

**响应：**
```json
{
  "id": 1,
  "order_id": 123,
  "factory_id": "factory_user_id",
  "status": "pending",
  "price": 1500.50,
  "jiedan_time": "2025-06-28T10:30:00Z",
  "agree_time": null,
  "agree_user_id": null,
  "created_at": "2025-06-28T10:30:00Z",
  "updated_at": "2025-06-28T10:30:00Z",
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
```

### 3. 获取订单的接单记录列表

**接口地址：** `GET /api/orders/{orderId}/jiedans`

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
    "status": "pending",
    "price": 1500.50,
    "jiedan_time": "2025-06-28T10:30:00Z",
    "agree_time": null,
    "agree_user_id": null,
    "created_at": "2025-06-28T10:30:00Z",
    "updated_at": "2025-06-28T10:30:00Z",
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
```

### 4. 获取工厂的接单记录列表

**接口地址：** `GET /api/factories/{factoryId}/jiedans?page=1&pageSize=10`

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
  "jiedans": [
    {
      "id": 1,
      "order_id": 123,
      "factory_id": "factory_user_id",
      "status": "pending",
      "price": 1500.50,
      "jiedan_time": "2025-06-28T10:30:00Z",
      "agree_time": null,
      "agree_user_id": null,
      "created_at": "2025-06-28T10:30:00Z",
      "updated_at": "2025-06-28T10:30:00Z",
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

### 5. 同意接单

**接口地址：** `POST /api/jiedan/{id}/accept`

**请求头：**
```
Content-Type: application/json
Authorization: Bearer {token}
```

**请求体：**
```json
{
  "agree_user_id": "user_id"
}
```

**响应：**
```json
{
  "id": 1,
  "order_id": 123,
  "factory_id": "factory_user_id",
  "status": "accepted",
  "price": 1500.50,
  "jiedan_time": "2025-06-28T10:30:00Z",
  "agree_time": "2025-06-28T11:00:00Z",
  "agree_user_id": "user_id",
  "created_at": "2025-06-28T10:30:00Z",
  "updated_at": "2025-06-28T11:00:00Z"
}
```

### 6. 拒绝接单

**接口地址：** `POST /api/jiedan/{id}/reject`

**请求头：**
```
Content-Type: application/json
Authorization: Bearer {token}
```

**请求体：**
```json
{
  "reason": "拒绝原因"
}
```

**响应：**
```json
{
  "id": 1,
  "order_id": 123,
  "factory_id": "factory_user_id",
  "status": "rejected",
  "price": 1500.50,
  "jiedan_time": "2025-06-28T10:30:00Z",
  "agree_time": null,
  "agree_user_id": null,
  "created_at": "2025-06-28T10:30:00Z",
  "updated_at": "2025-06-28T11:00:00Z"
}
```

### 7. 更新接单记录

**接口地址：** `PUT /api/jiedan/{id}`

**请求头：**
```
Content-Type: application/json
Authorization: Bearer {token}
```

**请求体：**
```json
{
  "status": "accepted",
  "agree_user_id": "user_id"
}
```

**响应：**
```json
{
  "id": 1,
  "order_id": 123,
  "factory_id": "factory_user_id",
  "status": "accepted",
  "price": 1500.50,
  "jiedan_time": "2025-06-28T10:30:00Z",
  "agree_time": null,
  "agree_user_id": "user_id",
  "created_at": "2025-06-28T10:30:00Z",
  "updated_at": "2025-06-28T11:00:00Z"
}
```

### 8. 删除接单记录

**接口地址：** `DELETE /api/jiedan/{id}`

**请求头：**
```
Authorization: Bearer {token}
```

**响应：**
```json
{
  "message": "接单记录删除成功"
}
```

### 9. 获取接单统计信息

**接口地址：** `GET /api/factories/{factoryId}/jiedan-statistics`

**请求头：**
```
Authorization: Bearer {token}
```

**响应：**
```json
{
  "pending": 5,
  "accepted": 10,
  "rejected": 2
}
```

## 状态说明

### 接单状态 (JiedanStatus)

- `pending`: 待处理 - 工厂已接单，等待处理
- `accepted`: 已同意 - 接单已被同意
- `rejected`: 已拒绝 - 接单已被拒绝

## 业务规则

1. **接单限制**：同一工厂对同一订单只能接单一次
2. **状态流转**：只能对待处理状态的接单进行同意或拒绝操作
3. **权限控制**：只有工厂用户可以进行接单操作
4. **身份验证**：工厂只能以自己的身份进行接单

## 错误码

| 错误码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 权限不足 |
| 404 | 接单记录不存在 |
| 409 | 该工厂已对该订单进行过接单操作 |
| 422 | 只能对待处理的接单进行同意/拒绝操作 |
| 500 | 服务器内部错误 |

## 使用示例

### 工厂接单流程

1. **工厂接单**
```bash
curl -X POST "https://aneworders.com/api/jiedan" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "order_id": 123,
    "factory_id": "factory_user_id",
    "price": 1500.50
  }'
```

2. **查看接单记录**
```bash
curl -X GET "https://aneworders.com/api/jiedan/1" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

3. **同意接单**
```bash
curl -X POST "https://aneworders.com/api/jiedan/1/accept" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "agree_user_id": "user_id"
  }'
```

4. **查看工厂接单统计**
```bash
curl -X GET "https://aneworders.com/api/factories/factory_user_id/jiedan-statistics" \
  -H "Authorization: Bearer YOUR_TOKEN"
``` 