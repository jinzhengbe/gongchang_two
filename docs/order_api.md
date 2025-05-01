# 订单管理API文档

## 重要更新
- 自 2024-05-01 起，API 路径已统一更新
  - 旧路径: `/api/v1/orders`
  - 新路径: `/api/orders`
  - 变更原因: 统一 API 路径格式，简化路由配置
  - 兼容性: 临时支持旧路径，建议尽快迁移到新路径

## 目录
1. [接口说明](#接口说明)
2. [认证说明](#认证说明)
3. [通用响应格式](#通用响应格式)
4. [错误码说明](#错误码说明)
5. [接口详情](#接口详情)
   - [创建订单](#创建订单)
   - [上传模型文件](#上传模型文件)
   - [上传详情图片](#上传详情图片)
   - [获取订单列表](#获取订单列表)
   - [获取订单详情](#获取订单详情)
   - [更新订单](#更新订单)
   - [删除订单](#删除订单)

## 接口说明

所有接口都需要在请求头中添加认证信息:
```
Authorization: Bearer <token>
```

## 认证说明

- 所有接口都需要JWT认证
- token通过登录接口获取
- token过期时间为24小时
- 认证失败返回401状态码

## 通用响应格式

### 成功响应
```json
{
    "code": 200,
    "message": "success",
    "data": {
        // 具体数据
    }
}
```

### 错误响应
```json
{
    "code": 400,
    "message": "错误信息",
    "data": null
}
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 409 | 资源冲突 |
| 500 | 服务器错误 |

## 接口详情

### 创建订单

**接口说明**: 创建新订单

**请求方式**: POST

**接口地址**: `/api/orders`

**请求头**:
```
Content-Type: application/json
Authorization: Bearer <token>
```

**请求参数**:
```json
{
    "designer_id": 1,          // 设计师ID
    "customer_id": 2,          // 客户ID
    "product_id": 3,           // 产品ID
    "quantity": 1,             // 数量
    "unit_price": 100.0,       // 单价
    "total_price": 100.0,      // 总价
    "status": "pending",       // 订单状态
    "payment_status": "unpaid",// 支付状态
    "shipping_address": "北京市朝阳区", // 收货地址
    "order_date": "2024-04-17T10:00:00Z" // 订单日期
}
```

**响应示例**:
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "id": 1,
        "designer_id": 1,
        "customer_id": 2,
        "product_id": 3,
        "quantity": 1,
        "unit_price": 100.0,
        "total_price": 100.0,
        "status": "pending",
        "payment_status": "unpaid",
        "shipping_address": "北京市朝阳区",
        "order_date": "2024-04-17T10:00:00Z",
        "created_at": "2024-04-17T10:00:00Z",
        "updated_at": "2024-04-17T10:00:00Z"
    }
}
```

### 上传模型文件

**接口说明**: 上传订单的模型文件

**请求方式**: POST

**接口地址**: `/api/orders/{orderId}/model-files`

**请求头**:
```
Content-Type: multipart/form-data
Authorization: Bearer <token>
```

**请求参数**:
- files: 文件数组(支持多个文件)

**文件限制**:
- 单个文件大小: ≤10MB
- 支持格式: .stl, .obj, .3ds等3D模型文件

**响应示例**:
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "message": "Files uploaded successfully",
        "file_ids": [1, 2, 3]
    }
}
```

### 上传详情图片

**接口说明**: 上传订单的详情图片

**请求方式**: POST

**接口地址**: `/api/orders/{orderId}/detail-images`

**请求头**:
```
Content-Type: multipart/form-data
Authorization: Bearer <token>
```

**请求参数**:
- files: 图片数组(支持多个图片)

**图片限制**:
- 单个图片大小: ≤10MB
- 支持格式: .jpg, .png, .gif

**响应示例**:
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "message": "Images uploaded successfully",
        "file_ids": [1, 2, 3]
    }
}
```

### 获取订单列表

**接口说明**: 获取订单列表，支持分页和筛选

**请求方式**: GET

**接口地址**: `/api/orders`

**请求头**:
```
Authorization: Bearer <token>
```

**查询参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认1 |
| pageSize | int | 否 | 每页数量，默认10 |
| status | string | 否 | 订单状态筛选 |
| startDate | string | 否 | 开始日期，格式: YYYY-MM-DD |
| endDate | string | 否 | 结束日期，格式: YYYY-MM-DD |

**响应示例**:
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "total": 100,
        "page": 1,
        "pageSize": 10,
        "orders": [
            {
                "id": 1,
                "designer_id": 1,
                "customer_id": 2,
                "product_id": 3,
                "quantity": 1,
                "unit_price": 100.0,
                "total_price": 100.0,
                "status": "pending",
                "payment_status": "unpaid",
                "shipping_address": "北京市朝阳区",
                "order_date": "2024-04-17T10:00:00Z",
                "model_files": [
                    {
                        "id": 1,
                        "file_name": "model1.stl",
                        "file_path": "uploads/model/20240417100000_model1.stl",
                        "file_type": "model",
                        "uploaded_by": 1
                    }
                ],
                "detail_images": [
                    {
                        "id": 1,
                        "file_name": "image1.jpg",
                        "file_path": "uploads/detail/20240417100000_image1.jpg",
                        "file_type": "detail",
                        "uploaded_by": 1
                    }
                ]
            }
        ]
    }
}
```

### 获取订单详情

**接口说明**: 获取单个订单的详细信息

**请求方式**: GET

**接口地址**: `/api/orders/{orderId}`

**请求头**:
```
Authorization: Bearer <token>
```

**响应示例**:
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "id": 1,
        "designer_id": 1,
        "customer_id": 2,
        "product_id": 3,
        "quantity": 1,
        "unit_price": 100.0,
        "total_price": 100.0,
        "status": "pending",
        "payment_status": "unpaid",
        "shipping_address": "北京市朝阳区",
        "order_date": "2024-04-17T10:00:00Z",
        "model_files": [
            {
                "id": 1,
                "file_name": "model1.stl",
                "file_path": "uploads/model/20240417100000_model1.stl",
                "file_type": "model",
                "uploaded_by": 1
            }
        ],
        "detail_images": [
            {
                "id": 1,
                "file_name": "image1.jpg",
                "file_path": "uploads/detail/20240417100000_image1.jpg",
                "file_type": "detail",
                "uploaded_by": 1
            }
        ]
    }
}
```

### 更新订单

**接口说明**: 更新订单信息

**请求方式**: PUT

**接口地址**: `/api/orders/{orderId}`

**请求头**:
```
Content-Type: application/json
Authorization: Bearer <token>
```

**请求参数**:
```json
{
    "quantity": 2,             // 数量
    "unit_price": 90.0,        // 单价
    "total_price": 180.0,      // 总价
    "status": "processing",    // 订单状态
    "payment_status": "paid",  // 支付状态
    "shipping_address": "北京市海淀区" // 收货地址
}
```

**响应示例**:
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "message": "Order updated successfully"
    }
}
```

### 删除订单

**接口说明**: 删除订单

**请求方式**: DELETE

**接口地址**: `/api/orders/{orderId}`

**请求头**:
```
Authorization: Bearer <token>
```

**响应示例**:
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "message": "Order deleted successfully"
    }
}
```

## 订单状态说明

| 状态 | 说明 |
|------|------|
| pending | 待处理 |
| accepted | 已接受 |
| in_progress | 进行中 |
| completed | 已完成 |
| cancelled | 已取消 |
| rejected | 已拒绝 |
| on_hold | 暂停中 |

## 支付状态说明

| 状态 | 说明 |
|------|------|
| unpaid | 未支付 |
| paid | 已支付 |
| refunded | 已退款 |
| failed | 支付失败 | 
| processing | 处理中 |
| partially_paid | 部分支付 | 