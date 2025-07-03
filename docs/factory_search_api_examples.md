# 工厂搜索API使用示例

## 概述
工厂搜索API提供了强大的工厂搜索功能，支持关键词搜索、地区筛选、专业领域筛选、评分筛选等多种功能。

## API端点

### 1. 工厂搜索
- **URL**: `GET /api/factories/search`
- **描述**: 搜索工厂，支持多种筛选条件
- **认证**: 无需认证（公开接口）

### 2. 搜索建议
- **URL**: `GET /api/factories/search/suggestions`
- **描述**: 获取搜索建议，提供智能提示
- **认证**: 无需认证（公开接口）

### 3. 创建专业领域
- **URL**: `POST /api/factories/{factory_id}/specialties`
- **描述**: 为工厂添加专业领域标签
- **认证**: 需要认证

### 4. 创建评分
- **URL**: `POST /api/factories/{factory_id}/ratings`
- **描述**: 为工厂添加评分和评价
- **认证**: 需要认证

## 使用示例

### JavaScript/TypeScript 示例

#### 1. 基础搜索
```javascript
// 基础关键词搜索
async function searchFactories(query) {
  const response = await fetch(`/api/factories/search?query=${encodeURIComponent(query)}&page=1&page_size=20`);
  const data = await response.json();
  return data;
}

// 使用示例
searchFactories('服装工厂').then(result => {
  console.log('搜索结果:', result);
});
```

#### 2. 高级搜索
```javascript
// 多条件搜索
async function advancedSearch(params) {
  const queryParams = new URLSearchParams({
    query: params.query || '',
    region: params.region || '',
    min_rating: params.minRating || '',
    max_rating: params.maxRating || '',
    page: params.page || 1,
    page_size: params.pageSize || 20,
    sort_by: params.sortBy || 'rating',
    sort_order: params.sortOrder || 'desc'
  });
  
  // 处理专业领域数组
  if (params.specialties && params.specialties.length > 0) {
    params.specialties.forEach(specialty => {
      queryParams.append('specialties', specialty);
    });
  }
  
  const response = await fetch(`/api/factories/search?${queryParams}`);
  const data = await response.json();
  return data;
}

// 使用示例
advancedSearch({
  query: '服装',
  region: '上海',
  specialties: ['服装制造', '配饰制造'],
  minRating: 4.0,
  maxRating: 5.0,
  sortBy: 'rating',
  sortOrder: 'desc'
}).then(result => {
  console.log('高级搜索结果:', result);
});
```

#### 3. 搜索建议
```javascript
// 获取搜索建议
async function getSearchSuggestions(query, limit = 10) {
  const response = await fetch(`/api/factories/search/suggestions?query=${encodeURIComponent(query)}&limit=${limit}`);
  const data = await response.json();
  return data;
}

// 使用示例
getSearchSuggestions('服装', 5).then(result => {
  console.log('搜索建议:', result);
});
```

#### 4. 创建专业领域（需要认证）
```javascript
// 为工厂添加专业领域
async function addFactorySpecialty(factoryId, specialty, token) {
  const response = await fetch(`/api/factories/${factoryId}/specialties`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({ specialty })
  });
  const data = await response.json();
  return data;
}

// 使用示例
addFactorySpecialty(1, '服装制造', 'your-jwt-token').then(result => {
  console.log('专业领域创建结果:', result);
});
```

#### 5. 创建评分（需要认证）
```javascript
// 为工厂添加评分
async function addFactoryRating(factoryId, rating, comment, token) {
  const response = await fetch(`/api/factories/${factoryId}/ratings`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({ rating, comment })
  });
  const data = await response.json();
  return data;
}

// 使用示例
addFactoryRating(1, 4.5, '服务很好，质量不错', 'your-jwt-token').then(result => {
  console.log('评分创建结果:', result);
});
```

### React 组件示例

#### 工厂搜索组件
```jsx
import React, { useState, useEffect } from 'react';

const FactorySearch = () => {
  const [searchParams, setSearchParams] = useState({
    query: '',
    region: '',
    minRating: '',
    maxRating: '',
    specialties: []
  });
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(false);

  const searchFactories = async () => {
    setLoading(true);
    try {
      const queryParams = new URLSearchParams({
        ...searchParams,
        page: 1,
        page_size: 20
      });
      
      const response = await fetch(`/api/factories/search?${queryParams}`);
      const data = await response.json();
      setResults(data.data.factories);
    } catch (error) {
      console.error('搜索失败:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="factory-search">
      <div className="search-form">
        <input
          type="text"
          placeholder="搜索工厂名称或地址"
          value={searchParams.query}
          onChange={(e) => setSearchParams({...searchParams, query: e.target.value})}
        />
        <input
          type="text"
          placeholder="地区筛选"
          value={searchParams.region}
          onChange={(e) => setSearchParams({...searchParams, region: e.target.value})}
        />
        <input
          type="number"
          placeholder="最低评分"
          value={searchParams.minRating}
          onChange={(e) => setSearchParams({...searchParams, minRating: e.target.value})}
        />
        <button onClick={searchFactories} disabled={loading}>
          {loading ? '搜索中...' : '搜索'}
        </button>
      </div>
      
      <div className="search-results">
        {results.map(factory => (
          <div key={factory.id} className="factory-card">
            <h3>{factory.name}</h3>
            <p>地址: {factory.address}</p>
            <p>评分: {factory.rating}</p>
            <p>专业领域: {factory.specialties.join(', ')}</p>
            <p>合作状态: {factory.cooperation_status}</p>
          </div>
        ))}
      </div>
    </div>
  );
};

export default FactorySearch;
```

#### 搜索建议组件
```jsx
import React, { useState, useEffect } from 'react';

const SearchSuggestions = ({ query, onSuggestionClick }) => {
  const [suggestions, setSuggestions] = useState([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (query.length > 0) {
      fetchSuggestions();
    } else {
      setSuggestions([]);
    }
  }, [query]);

  const fetchSuggestions = async () => {
    setLoading(true);
    try {
      const response = await fetch(`/api/factories/search/suggestions?query=${encodeURIComponent(query)}&limit=10`);
      const data = await response.json();
      setSuggestions(data.data.suggestions);
    } catch (error) {
      console.error('获取建议失败:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="search-suggestions">
      {loading && <div>加载中...</div>}
      {suggestions.map((suggestion, index) => (
        <div
          key={index}
          className="suggestion-item"
          onClick={() => onSuggestionClick(suggestion.text)}
          dangerouslySetInnerHTML={{ __html: suggestion.highlight }}
        />
      ))}
    </div>
  );
};

export default SearchSuggestions;
```

### cURL 示例

#### 1. 基础搜索
```bash
curl -X GET "http://localhost:8008/api/factories/search?query=服装工厂&page=1&page_size=10"
```

#### 2. 地区筛选搜索
```bash
curl -X GET "http://localhost:8008/api/factories/search?region=上海&min_rating=4.0&page=1&page_size=10"
```

#### 3. 专业领域筛选搜索
```bash
curl -X GET "http://localhost:8008/api/factories/search?specialties=服装制造&specialties=配饰制造&page=1&page_size=10"
```

#### 4. 获取搜索建议
```bash
curl -X GET "http://localhost:8008/api/factories/search/suggestions?query=服装&limit=5"
```

#### 5. 创建专业领域（需要认证）
```bash
curl -X POST "http://localhost:8008/api/factories/1/specialties" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"specialty": "服装制造"}'
```

#### 6. 创建评分（需要认证）
```bash
curl -X POST "http://localhost:8008/api/factories/1/ratings" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"rating": 4.5, "comment": "服务很好，质量不错"}'
```

## 响应格式

### 搜索响应
```json
{
  "success": true,
  "data": {
    "factories": [
      {
        "id": 1,
        "name": "上海服装厂",
        "address": "上海市浦东新区",
        "specialties": ["服装制造", "配饰制造"],
        "rating": 4.5,
        "cooperation_status": "合作中",
        "description": "专业服装制造工厂",
        "contact_info": {
          "phone": "13800138000",
          "email": "contact@shanghai-factory.com"
        },
        "capacity": {
          "monthly_orders": 100,
          "max_order_size": 10000
        },
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20
  }
}
```

### 搜索建议响应
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

## 错误处理

### 常见错误响应
```json
{
  "success": false,
  "error": "错误描述"
}
```

### 错误码说明
- `400`: 请求参数错误
- `401`: 未认证（需要token的接口）
- `403`: 权限不足
- `500`: 服务器内部错误

## 性能优化建议

1. **使用分页**: 避免一次性获取大量数据
2. **合理使用筛选**: 减少不必要的查询条件
3. **缓存结果**: 对热门搜索词进行缓存
4. **异步加载**: 使用异步方式加载搜索建议
5. **防抖处理**: 对搜索输入进行防抖处理

## 注意事项

1. 搜索接口是公开的，无需认证
2. 专业领域和评分创建需要认证
3. 评分范围是0.0-5.0
4. 分页大小最大为100
5. 搜索建议数量最大为20
6. 所有时间字段使用ISO 8601格式 