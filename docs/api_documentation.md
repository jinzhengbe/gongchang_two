# API 文档

## 工厂管理

### 注册工厂
- **URL**: `/api/factory/register`
- **方法**: `POST`
- **请求体**:
```json
{
    "name": "工厂名称",
    "username": "用户名",
    "password": "密码",
    "address": "地址",
    "contact": "联系人",
    "phone": "联系电话",
    "email": "邮箱",
    "license": "营业执照号",
    "description": "工厂描述"
}
```
- **响应**:
```json
{
    "message": "注册成功"
}
```

### 工厂登录
- **URL**: `/api/factory/login`
- **方法**: `POST`
- **请求体**:
```json
{
    "username": "用户名",
    "password": "密码"
}
```
- **响应**:
```json
{
    "token": "JWT令牌",
    "factory": {
        "id": 1,
        "name": "工厂名称",
        "username": "用户名",
        "address": "地址",
        "contact": "联系人",
        "phone": "联系电话",
        "email": "邮箱",
        "license": "营业执照号",
        "description": "工厂描述",
        "status": 1,
        "created_at": "创建时间",
        "updated_at": "更新时间"
    }
}
```

## 订单管理

### 获取工厂订单列表
- **URL**: `/api/factory/orders`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer <token>`
- **响应**:
```json
{
    "orders": [
        {
            "id": 1,
            "title": "订单标题",
            "description": "订单描述",
            "status": "订单状态",
            "created_at": "创建时间"
        }
    ]
}
```

### 更新订单状态
- **URL**: `/api/orders/:id/status`
- **方法**: `PUT`
- **请求头**: `Authorization: Bearer <token>`
- **请求体**:
```json
{
    "status": "新状态"
}
```
- **响应**:
```json
{
    "message": "更新成功"
}
```

## 文件管理

### 上传文件
- **URL**: `/api/files/upload`
- **方法**: `POST`
- **请求头**: 
  - `Authorization: Bearer <token>`
  - `Content-Type: multipart/form-data`
- **请求体**:
  - `file`: 文件
  - `order_id`: 订单ID
- **响应**:
```json
{
    "file_id": "文件ID",
    "url": "文件访问URL"
}
```

### 获取订单文件列表
- **URL**: `/api/files/order/:id`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer <token>`
- **响应**:
```json
{
    "files": [
        {
            "id": "文件ID",
            "name": "文件名",
            "url": "文件URL",
            "created_at": "上传时间"
        }
    ]
}
``` 