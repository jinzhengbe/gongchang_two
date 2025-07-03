# 订单搜索API使用示例

## 1. 基础搜索

### 关键词搜索
```bash
# 搜索包含"连衣裙"的订单
curl -X GET "http://localhost:8008/api/orders/search?query=连衣裙" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 状态筛选
```bash
# 搜索已发布的订单
curl -X GET "http://localhost:8008/api/orders/search?status=published" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 搜索所有状态的订单
curl -X GET "http://localhost:8008/api/orders/search?status=all" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 时间范围搜索
```bash
# 搜索2024年1月的订单
curl -X GET "http://localhost:8008/api/orders/search?start_date=2024-01-01&end_date=2024-01-31" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 2. 高级搜索

### 组合搜索
```bash
# 搜索2024年春季的连衣裙订单
curl -X GET "http://localhost:8008/api/orders/search?query=连衣裙&start_date=2024-03-01&end_date=2024-05-31&status=published" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 分页搜索
```bash
# 第2页，每页10条记录
curl -X GET "http://localhost:8008/api/orders/search?page=2&page_size=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 排序搜索
```bash
# 按标题升序排列
curl -X GET "http://localhost:8008/api/orders/search?sort_by=title&sort_order=asc" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 按更新时间倒序排列
curl -X GET "http://localhost:8008/api/orders/search?sort_by=updated_at&sort_order=desc" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 3. 搜索建议

### 获取搜索建议
```bash
# 获取"连衣裙"的搜索建议
curl -X GET "http://localhost:8008/api/orders/search/suggestions?query=连衣裙" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 限制建议数量
curl -X GET "http://localhost:8008/api/orders/search/suggestions?query=真丝&limit=5" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 搜索建议响应示例
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

## 4. 搜索统计

### 获取搜索统计信息
```bash
curl -X GET "http://localhost:8008/api/orders/search/statistics" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 统计信息响应示例
```json
{
  "success": true,
  "data": {
    "total_orders": 1250,
    "user_id": "user123",
    "user_role": "factory",
    "hot_keywords": [
      "连衣裙",
      "真丝面料",
      "春季新款"
    ]
  }
}
```

## 5. 完整搜索响应示例

### 搜索响应格式
```json
{
  "success": true,
  "data": {
    "orders": [
      {
        "id": 123,
        "title": "2024春季新款连衣裙",
        "order_no": "ORD000123",
        "status": "published",
        "fabrics": ["真丝面料", "蕾丝花边"],
        "factory": {
          "id": "factory123",
          "name": "上海服装厂"
        },
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
      }
    ],
    "total": 25,
    "page": 1,
    "page_size": 20
  }
}
```

## 6. 错误处理

### 常见错误响应
```json
{
  "success": false,
  "error": "查询关键词不能为空"
}
```

```json
{
  "success": false,
  "error": "获取总数失败: database connection error"
}
```

## 7. 权限控制

### 不同角色的搜索权限
- **工厂用户**: 只能搜索分配给自己的订单
- **设计师**: 只能搜索自己创建的订单
- **管理员**: 可以搜索所有订单
- **普通用户**: 只能搜索公开的订单

## 8. 性能优化建议

### 查询优化
1. 使用合适的索引
2. 避免使用 `%keyword%` 模式（如果可能）
3. 合理设置分页大小（建议不超过100）
4. 使用缓存热门搜索词

### 前端优化
1. 实现搜索防抖
2. 使用搜索建议减少无效请求
3. 合理缓存搜索结果
4. 实现虚拟滚动处理大量数据

## 9. 测试用例

### 功能测试
```bash
# 测试空查询
curl -X GET "http://localhost:8008/api/orders/search" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 测试无效日期格式
curl -X GET "http://localhost:8008/api/orders/search?start_date=invalid-date" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 测试超大数据量
curl -X GET "http://localhost:8008/api/orders/search?page_size=1000" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 性能测试
```bash
# 压力测试
ab -n 1000 -c 10 "http://localhost:8008/api/orders/search?query=test"

# 并发测试
curl -X GET "http://localhost:8008/api/orders/search?query=test" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -w "Time: %{time_total}s\n"
``` 