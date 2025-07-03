# 设计师搜索API文档

## 概述

设计师搜索API提供了完整的设计师搜索、筛选、建议和评分功能，支持按名称、地区、专业领域、评分等条件进行搜索。

## API端点

### 1. 设计师搜索

**GET** `/api/designers/search`

搜索设计师，支持多种筛选条件。

#### 查询参数

| 参数 | 类型 | 必填 | 默认值 | 描述 |
|------|------|------|--------|------|
| query | string | 否 | - | 搜索关键词（设计师名称、地址） |
| region | string | 否 | - | 地区筛选 |
| specialties | string[] | 否 | - | 专业领域数组 |
| min_rating | number | 否 | - | 最低评分（0-5） |
| max_rating | number | 否 | - | 最高评分（0-5） |
| page | number | 否 | 1 | 页码 |
| page_size | number | 否 | 20 | 每页数量（最大100） |
| sort_by | string | 否 | rating | 排序字段（name, rating, created_at） |
| sort_order | string | 否 | desc | 排序方向（asc, desc） |

#### 响应格式

```json
{
  "success": true,
  "data": {
    "designers": [
      {
        "id": 1,
        "name": "设计师工作室",
        "address": "广东省深圳市",
        "specialties": ["服装设计", "时尚设计"],
        "rating": 4.5,
        "description": "专业服装设计工作室",
        "contact_info": {
          "phone": "designer@test.com",
          "email": "designer@test.com"
        },
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-01-01T00:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20
  }
}
```

#### 使用示例

```bash
# 基础搜索
curl -X GET "http://localhost:8008/api/designers/search?query=服装&page=1&page_size=20"

# 地区筛选
curl -X GET "http://localhost:8008/api/designers/search?region=深圳&page=1&page_size=20"

# 评分筛选
curl -X GET "http://localhost:8008/api/designers/search?min_rating=4.0&page=1&page_size=20"

# 专业领域筛选
curl -X GET "http://localhost:8008/api/designers/search?specialties=服装设计&page=1&page_size=20"

# 组合搜索
curl -X GET "http://localhost:8008/api/designers/search?query=设计&region=深圳&min_rating=4.0&specialties=服装设计&page=1&page_size=20"

# 按名称排序
curl -X GET "http://localhost:8008/api/designers/search?sort_by=name&sort_order=asc&page=1&page_size=20"
```

### 2. 搜索建议

**GET** `/api/designers/search/suggestions`

获取搜索建议，支持设计师名称、地址、专业领域建议。

#### 查询参数

| 参数 | 类型 | 必填 | 默认值 | 描述 |
|------|------|------|--------|------|
| query | string | 是 | - | 搜索关键词 |
| limit | number | 否 | 10 | 建议数量（最大20） |

#### 响应格式

```json
{
  "success": true,
  "data": {
    "suggestions": [
      {
        "type": "designer_name",
        "text": "设计师工作室",
        "highlight": "<em>设计</em>师工作室"
      },
      {
        "type": "designer_address",
        "text": "广东省深圳市",
        "highlight": "广东省<em>深圳</em>市"
      },
      {
        "type": "specialty",
        "text": "服装设计",
        "highlight": "<em>服装</em>设计"
      }
    ]
  }
}
```

#### 使用示例

```bash
# 获取设计师名称建议
curl -X GET "http://localhost:8008/api/designers/search/suggestions?query=设计&limit=10"

# 获取地址建议
curl -X GET "http://localhost:8008/api/designers/search/suggestions?query=深圳&limit=10"

# 获取专业领域建议
curl -X GET "http://localhost:8008/api/designers/search/suggestions?query=服装&limit=10"
```

### 3. 创建专业领域

**POST** `/api/designers/{designer_id}/specialties`

为设计师添加专业领域标签（需要认证）。

#### 路径参数

| 参数 | 类型 | 必填 | 描述 |
|------|------|------|------|
| designer_id | number | 是 | 设计师ID |

#### 请求体

```json
{
  "specialty": "服装设计"
}
```

#### 请求头

```
Authorization: Bearer <token>
Content-Type: application/json
```

#### 响应格式

```json
{
  "success": true,
  "message": "专业领域创建成功"
}
```

#### 使用示例

```bash
curl -X POST "http://localhost:8008/api/designers/1/specialties" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"specialty": "服装设计"}'
```

### 4. 创建评分

**POST** `/api/designers/{designer_id}/ratings`

为设计师添加评分和评价（需要认证）。

#### 路径参数

| 参数 | 类型 | 必填 | 描述 |
|------|------|------|------|
| designer_id | number | 是 | 设计师ID |

#### 请求体

```json
{
  "rating": 4.5,
  "comment": "设计很专业，沟通顺畅"
}
```

#### 请求头

```
Authorization: Bearer <token>
Content-Type: application/json
```

#### 响应格式

```json
{
  "success": true,
  "message": "评分创建成功"
}
```

#### 使用示例

```bash
curl -X POST "http://localhost:8008/api/designers/1/ratings" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"rating": 4.5, "comment": "设计很专业，沟通顺畅"}'
```

## 错误处理

### 常见错误响应

#### 400 Bad Request
```json
{
  "success": false,
  "error": "参数绑定失败: Key: 'DesignerSearchRequest.Page' Error:Field validation for 'Page' failed on the 'min' tag"
}
```

#### 401 Unauthorized
```json
{
  "success": false,
  "error": "Authorization header is required"
}
```

#### 500 Internal Server Error
```json
{
  "success": false,
  "error": "搜索设计师失败: database connection error"
}
```

## 数据库表结构

### designer_profiles 表
```sql
CREATE TABLE designer_profiles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id VARCHAR(191) UNIQUE,
    company_name VARCHAR(255),
    address VARCHAR(255),
    website VARCHAR(255),
    bio TEXT,
    rating DECIMAL(3,2) DEFAULT 0,
    status INT DEFAULT 1,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

### designer_specialties 表
```sql
CREATE TABLE designer_specialties (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    designer_id BIGINT NOT NULL,
    specialty VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### designer_ratings 表
```sql
CREATE TABLE designer_ratings (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    designer_id BIGINT NOT NULL,
    rating DECIMAL(3,2) NOT NULL,
    comment TEXT,
    rater_id VARCHAR(191) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## 性能优化

### 数据库索引
- `idx_designer_profiles_search` - 设计师名称和地址复合索引
- `idx_designer_specialties_designer_id` - 专业领域设计师ID索引
- `idx_designer_ratings_designer_id` - 评分设计师ID索引
- `idx_designer_ratings_rating` - 评分索引
- `idx_designer_profiles_status_rating` - 状态和评分复合索引

### 查询优化
- 使用子查询计算平均评分
- 支持分页查询
- 索引覆盖查询
- 连接查询优化

## 客户端调用示例

#### JavaScript
```javascript
// 基础搜索
async function searchDesigners(query) {
  const response = await fetch(`/api/designers/search?query=${encodeURIComponent(query)}&page=1&page_size=20`);
  const data = await response.json();
  return data;
}

// 获取搜索建议
async function getDesignerSuggestions(query) {
  const response = await fetch(`/api/designers/search/suggestions?query=${encodeURIComponent(query)}&limit=10`);
  const data = await response.json();
  return data;
}

// 创建专业领域
async function createDesignerSpecialty(designerId, specialty, token) {
  const response = await fetch(`/api/designers/${designerId}/specialties`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({ specialty })
  });
  return await response.json();
}

// 创建评分
async function createDesignerRating(designerId, rating, comment, token) {
  const response = await fetch(`/api/designers/${designerId}/ratings`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({ rating, comment })
  });
  return await response.json();
}
```

#### Python
```python
import requests

# 基础搜索
def search_designers(query):
    response = requests.get(f"http://localhost:8008/api/designers/search", 
                          params={"query": query, "page": 1, "page_size": 20})
    return response.json()

# 获取搜索建议
def get_designer_suggestions(query):
    response = requests.get(f"http://localhost:8008/api/designers/search/suggestions", 
                          params={"query": query, "limit": 10})
    return response.json()

# 创建专业领域
def create_designer_specialty(designer_id, specialty, token):
    response = requests.post(f"http://localhost:8008/api/designers/{designer_id}/specialties",
                           headers={"Authorization": f"Bearer {token}"},
                           json={"specialty": specialty})
    return response.json()

# 创建评分
def create_designer_rating(designer_id, rating, comment, token):
    response = requests.post(f"http://localhost:8008/api/designers/{designer_id}/ratings",
                           headers={"Authorization": f"Bearer {token}"},
                           json={"rating": rating, "comment": comment})
    return response.json()
```

## 测试账号

### 设计师账号
- 用户名：`testuser1`
- 密码：`test123`

### 使用说明
1. 使用设计师账号登录获取token
2. 使用token进行需要认证的操作
3. 搜索和建议接口无需认证

## 注意事项

1. **评分范围**：评分必须在0-5之间
2. **分页限制**：每页最大100条记录
3. **认证要求**：专业领域和评分创建需要认证
4. **数据验证**：所有输入都会进行验证
5. **性能考虑**：大量数据时建议使用分页
6. **索引优化**：已为常用查询创建索引 