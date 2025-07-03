# 工厂搜索API文档

## 1. 工厂搜索API

### 接口信息
- **请求方法：** GET
- **接口路径：** `/api/factories/search`
- **功能描述：** 搜索工厂，支持关键词搜索、地区筛选、专业领域筛选等

### 请求参数
```json
{
  "query": "搜索关键词",
  "region": "地区筛选",
  "specialties": ["专业领域1", "专业领域2"],
  "cooperation_status": "合作状态",
  "min_rating": 4.0,
  "max_rating": 5.0,
  "page": 1,
  "page_size": 20,
  "sort_by": "rating",
  "sort_order": "desc"
}
```

### 参数说明
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| query | string | 否 | 搜索关键词，支持工厂名称、地址、专业领域 |
| region | string | 否 | 地区筛选，如：上海市、浙江省 |
| specialties | array | 否 | 专业领域数组：["服装", "配饰", "鞋类"] |
| cooperation_status | string | 否 | 合作状态：all, cooperating, not_cooperating |
| min_rating | float | 否 | 最低评分，范围0.0-5.0 |
| max_rating | float | 否 | 最高评分，范围0.0-5.0 |
| page | integer | 否 | 页码，默认1 |
| page_size | integer | 否 | 每页数量，默认20，最大100 |
| sort_by | string | 否 | 排序字段：rating, name, created_at |
| sort_order | string | 否 | 排序方向：asc, desc |

### 响应数据
```json
{
  "success": true,
  "data": {
    "factories": [
      {
        "id": "工厂ID",
        "name": "工厂名称",
        "address": "工厂地址",
        "specialties": ["服装", "配饰"],
        "rating": 4.5,
        "cooperation_status": "合作中",
        "description": "工厂描述",
        "contact_info": {
          "phone": "联系电话",
          "email": "邮箱"
        },
        "capacity": {
          "monthly_orders": 100,
          "max_order_size": 10000
        },
        "created_at": "创建时间",
        "updated_at": "更新时间"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20
  }
}
```

### 响应字段说明
| 字段名 | 类型 | 说明 |
|--------|------|------|
| success | boolean | 请求是否成功 |
| data.factories | array | 工厂列表 |
| data.total | integer | 总记录数 |
| data.page | integer | 当前页码 |
| data.page_size | integer | 每页数量 |

## 2. 工厂搜索建议API

### 接口信息
- **请求方法：** GET
- **接口路径：** `/api/factories/search/suggestions`
- **功能描述：** 获取工厂搜索建议，提供智能提示

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
        "type": "factory_name",
        "text": "上海服装厂",
        "highlight": "上海<em>服装厂</em>"
      },
      {
        "type": "factory_address",
        "text": "上海市浦东新区",
        "highlight": "上海市<em>浦东新区</em>"
      },
      {
        "type": "specialty",
        "text": "服装制造",
        "highlight": "<em>服装</em>制造"
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
| suggestions.type | string | 建议类型：factory_name, factory_address, specialty |
| suggestions.text | string | 建议文本 |
| suggestions.highlight | string | 高亮显示的文本 |

## 3. 数据库索引设计

### 工厂表索引
```sql
-- 工厂搜索索引
CREATE INDEX idx_factories_search ON factories (
  name, 
  address, 
  rating,
  cooperation_status
);

-- 工厂专业领域索引
CREATE INDEX idx_factory_specialties ON factory_specialties (
  factory_id, 
  specialty
);

-- 工厂地区索引
CREATE INDEX idx_factories_region ON factories (
  province,
  city
);
```

## 4. 搜索算法设计

### 搜索策略
1. **全文搜索**：使用数据库全文搜索功能，支持中文分词
2. **模糊匹配**：支持部分关键词匹配
3. **权重计算**：
   - 工厂名称：权重 1.0
   - 工厂地址：权重 0.8
   - 专业领域：权重 0.6
   - 工厂描述：权重 0.4

### 相关性排序
- 精确匹配优先
- 评分高低
- 合作状态（合作中优先）
- 最近更新时间

## 5. 性能优化

### 缓存机制
- 热门搜索词缓存（1小时）
- 搜索结果缓存（5分钟）
- 搜索建议缓存（10分钟）
- 地区数据缓存（24小时）

### 分页优化
- 默认每页20条记录
- 支持游标分页
- 索引优化查询性能

## 6. 筛选功能

### 地区筛选
- 支持省份筛选
- 支持城市筛选
- 支持多地区选择

### 专业领域筛选
- 服装制造
- 配饰制造
- 鞋类制造
- 箱包制造
- 其他

### 评分筛选
- 支持评分范围选择
- 支持最低评分筛选
- 支持最高评分筛选 