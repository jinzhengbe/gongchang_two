# 订单搜索API文档

## 1. 订单搜索API

### 接口信息
- **请求方法：** GET
- **接口路径：** `/api/orders/search`
- **功能描述：** 搜索订单，支持关键词搜索、状态筛选、时间筛选等

### 请求参数
```json
{
  "query": "搜索关键词",
  "status": "订单状态筛选",
  "start_date": "开始日期",
  "end_date": "结束日期",
  "page": 1,
  "page_size": 20,
  "sort_by": "created_at",
  "sort_order": "desc"
}
```

### 参数说明
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| query | string | 否 | 搜索关键词，支持订单标题、订单号、面料、工厂名称 |
| status | string | 否 | 订单状态：all, published, in_progress, completed, cancelled |
| start_date | string | 否 | 开始日期，格式：YYYY-MM-DD |
| end_date | string | 否 | 结束日期，格式：YYYY-MM-DD |
| page | integer | 否 | 页码，默认1 |
| page_size | integer | 否 | 每页数量，默认20，最大100 |
| sort_by | string | 否 | 排序字段：created_at, updated_at, order_no |
| sort_order | string | 否 | 排序方向：asc, desc |

### 响应数据
```json
{
  "success": true,
  "data": {
    "orders": [
      {
        "id": "订单ID",
        "title": "订单标题",
        "order_no": "订单号",
        "status": "订单状态",
        "fabrics": ["面料1", "面料2"],
        "factory": {
          "id": "工厂ID",
          "name": "工厂名称"
        },
        "created_at": "创建时间",
        "updated_at": "更新时间"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 响应字段说明
| 字段名 | 类型 | 说明 |
|--------|------|------|
| success | boolean | 请求是否成功 |
| data.orders | array | 订单列表 |
| data.total | integer | 总记录数 |
| data.page | integer | 当前页码 |
| data.page_size | integer | 每页数量 |

## 2. 订单搜索建议API

### 接口信息
- **请求方法：** GET
- **接口路径：** `/api/orders/search/suggestions`
- **功能描述：** 获取搜索建议，提供智能提示

### 请求参数
```json
{
  "query": "搜索关键词",
  "limit": 10
}
```

### 参数说明
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| query | string | 是 | 搜索关键词 |
| limit | integer | 否 | 建议数量，默认10，最大20 |

### 响应数据
```json
{
  "success": true,
  "data": {
    "suggestions": [
      {
        "type": "order_title",
        "text": "2024春季新款连衣裙",
        "highlight": "2024春季新款<em>连衣裙</em>"
      },
      {
        "type": "fabric_name",
        "text": "真丝面料",
        "highlight": "真丝<em>面料</em>"
      },
      {
        "type": "factory_name",
        "text": "上海服装厂",
        "highlight": "上海<em>服装厂</em>"
      }
    ]
  }
}
```

### 响应字段说明
| 字段名 | 类型 | 说明 |
|--------|------|------|
| success | boolean | 请求是否成功 |
| data.suggestions | array | 搜索建议列表 |
| suggestions.type | string | 建议类型：order_title, fabric_name, factory_name |
| suggestions.text | string | 建议文本 |
| suggestions.highlight | string | 高亮显示的文本 |

## 3. 数据库索引设计

### 订单表索引
```sql
-- 订单搜索索引
CREATE INDEX idx_orders_search ON orders (
  title, 
  order_no, 
  status, 
  created_at
);

-- 订单面料关联索引
CREATE INDEX idx_order_fabrics ON order_fabrics (
  order_id, 
  fabric_name
);
```

## 4. 搜索算法设计

### 搜索策略
1. **全文搜索**：使用数据库全文搜索功能，支持中文分词
2. **模糊匹配**：支持部分关键词匹配
3. **权重计算**：
   - 订单标题：权重 1.0
   - 订单号：权重 0.8
   - 面料名称：权重 0.6
   - 工厂名称：权重 0.4
   - 订单描述：权重 0.3

### 相关性排序
- 精确匹配优先
- 关键词出现频率
- 最近更新时间
- 用户行为数据（点击率）

## 5. 性能优化

### 缓存机制
- 热门搜索词缓存（1小时）
- 搜索结果缓存（5分钟）
- 搜索建议缓存（10分钟）

### 分页优化
- 默认每页20条记录
- 支持游标分页
- 索引优化查询性能 