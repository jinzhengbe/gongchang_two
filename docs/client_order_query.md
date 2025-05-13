# 客户端订单查询说明

## 1. 查询单个订单

### 请求方式
```
GET /api/orders/{orderId}
```

### 请求参数
- orderId: 订单ID（路径参数）

### 请求示例
```bash
curl -X GET "http://localhost:8080/api/orders/1" \
     -H "Authorization: Bearer your_jwt_token"
```

### 响应示例
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "id": 1,
        "orderNo": "202402200001",
        "userId": "ttrr",
        "status": "PENDING",
        "totalAmount": 100.00,
        "createTime": "2024-02-20T10:00:00",
        "updateTime": "2024-02-20T10:00:00"
    }
}
```

## 2. 查询用户订单列表

### 请求方式
```
GET /api/orders/user/{userId}
```

### 请求参数
- userId: 用户ID（路径参数）
- page: 页码（可选，默认为1）
- size: 每页大小（可选，默认为10）

### 请求示例
```bash
curl -X GET "http://localhost:8080/api/orders/user/ttrr?page=1&size=10" \
     -H "Authorization: Bearer your_jwt_token"
```

### 响应示例
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "content": [
            {
                "id": 1,
                "orderNo": "202402200001",
                "userId": "ttrr",
                "status": "PENDING",
                "totalAmount": 100.00,
                "createTime": "2024-02-20T10:00:00",
                "updateTime": "2024-02-20T10:00:00"
            }
        ],
        "pageable": {
            "pageNumber": 0,
            "pageSize": 10,
            "sort": {
                "sorted": false
            }
        },
        "totalElements": 1,
        "totalPages": 1,
        "last": true,
        "first": true,
        "empty": false
    }
}
```

## 3. 查询所有订单（管理员）

### 请求方式
```
GET /api/orders
```

### 请求参数
- page: 页码（可选，默认为1）
- size: 每页大小（可选，默认为10）
- status: 订单状态（可选，用于筛选）

### 请求示例
```bash
curl -X GET "http://localhost:8080/api/orders?page=1&size=10&status=PENDING" \
     -H "Authorization: Bearer your_jwt_token"
```

### 响应示例
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "content": [
            {
                "id": 1,
                "orderNo": "202402200001",
                "userId": "ttrr",
                "status": "PENDING",
                "totalAmount": 100.00,
                "createTime": "2024-02-20T10:00:00",
                "updateTime": "2024-02-20T10:00:00"
            }
        ],
        "pageable": {
            "pageNumber": 0,
            "pageSize": 10,
            "sort": {
                "sorted": false
            }
        },
        "totalElements": 1,
        "totalPages": 1,
        "last": true,
        "first": true,
        "empty": false
    }
}
```

## 注意事项

1. 所有请求都需要在请求头中包含有效的 JWT token
2. 查询用户订单列表和查询所有订单接口支持分页
3. 订单状态包括：
   - PENDING: 待处理
   - PROCESSING: 处理中
   - COMPLETED: 已完成
   - CANCELLED: 已取消
4. 响应中的时间格式为 ISO 8601 标准格式
5. 分页参数从0开始计数

## 错误处理

如果请求失败，将返回如下格式的错误信息：

```json
{
    "code": 错误码,
    "message": "错误信息",
    "data": null
}
```

常见错误码：
- 401: 未授权
- 403: 禁止访问
- 404: 资源不存在
- 500: 服务器内部错误 