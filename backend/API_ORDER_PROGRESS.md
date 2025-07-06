# 订单进度相关接口文档

## 1. 创建订单进度（Add Progress）

- **接口路径**：`POST /api/orders/{orderId}/progress`
- **请求参数（JSON）**：

```json
{
  "order_id": "39",                // 订单ID（必填）
  "factory_id": "888",              // 工厂ID（必填）
  "type": "design",                 // 阶段类型: design/material/production/quality/packaging/shipping/custom
  "status": "in_progress",          // 状态: not_started/in_progress/completed/delayed/on_hold
  "description": "进度说明",         // 进度描述
  "start_time": "2025-07-29T00:00:00.000Z",      // 可选，开始时间
  "completed_time": "2025-07-30T00:00:00.000Z",  // 可选，完成时间
  "images": [
    "https://xxx.com/upload/xxx.jpg"   // 可选，图片URL数组
  ]
}
```

- **返回值（JSON）**：

```json
{
  "success": true,
  "data": {
    "id": "123",
    "order_id": "39",
    "factory_id": "888",
    "type": "design",
    "status": "in_progress",
    "description": "进度说明",
    "start_time": "2025-07-29T00:00:00.000Z",
    "completed_time": "2025-07-30T00:00:00.000Z",
    "images": [
      "https://xxx.com/upload/xxx.jpg"
    ],
    "created_at": "2025-07-29T12:00:00.000Z",
    "updated_at": "2025-07-29T12:00:00.000Z"
  }
}
```

---

## 2. 获取订单进度列表（Get Progress List）

- **接口路径**：`GET /api/orders/{orderId}/progress`
- **返回值（JSON）**：

```json
{
  "success": true,
  "data": [
    {
      "id": "123",
      "order_id": "39",
      "factory_id": "888",
      "type": "design",
      "status": "completed",
      "description": "设计已完成",
      "start_time": "2025-07-29T00:00:00.000Z",
      "completed_time": "2025-07-30T00:00:00.000Z",
      "images": [
        "https://xxx.com/upload/xxx.jpg"
      ],
      "created_at": "2025-07-29T12:00:00.000Z",
      "updated_at": "2025-07-29T12:00:00.000Z"
    }
    // ...更多进度
  ]
}
```

---

## 3. 图片上传（Upload Image）

- **接口路径**：`POST /api/files/upload`
- **请求类型**：`multipart/form-data`
- **参数**：
  - `file`：图片文件
  - `type`：可选，文件类型（如 image）

- **返回值（JSON）**：

```json
{
  "success": true,
  "files": [
    {
      "url": "https://xxx.com/upload/xxx.jpg",
      "name": "xxx.jpg"
    }
  ]
}
```

---

## 4. 字段说明

- `order_id`：订单唯一标识，必填。
- `factory_id`：工厂唯一标识，必填。
- `type`：进度阶段类型。
- `status`：进度状态。
- `description`：进度描述。
- `start_time`/`completed_time`：时间字段，ISO8601格式。
- `images`：图片URL数组，由图片上传接口返回。

---

## 5. 说明

- 图片上传接口可复用 `/api/files/upload`，前端已统一调用，无需新开发。
- 进度相关接口建议返回标准 success/error 字段，便于前端判断。
- 进度的 images 字段为图片 URL 数组，前端会展示缩略图和大图预览。
- 如需扩展其它字段，可与前端协商。 