# 布料管理 API 文档

## 概述

布料管理API提供了完整的布料信息管理功能，包括布料的增删改查、搜索、分类管理等功能。

## 基础信息

- **基础URL**: `https://aneworders.com/api`
- **认证方式**: JWT Token (部分接口需要认证)
- **数据格式**: JSON

## API 接口列表

### 1. 获取所有布料 (公开接口)

获取所有可用的布料列表，用于前端下拉选择。

**接口地址**: `GET /api/fabrics/all`

**请求参数**: 无

**响应示例**:
```json
[
  {
    "id": 1,
    "name": "纯棉平纹布",
    "category": "棉布",
    "material": "100%棉",
    "color": "白色",
    "pattern": "平纹",
    "weight": 120.00,
    "width": 150.00,
    "price": 15.50,
    "unit": "米",
    "stock": 100,
    "min_order": 1,
    "description": "优质纯棉平纹布，透气性好，适合制作衬衫、T恤等",
    "image_url": "/uploads/fabrics/cotton_plain_white.jpg",
    "thumbnail_url": "/uploads/fabrics/thumbnails/cotton_plain_white.jpg",
    "tags": "棉布,白色,平纹,透气",
    "status": 1,
    "supplier_id": null,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

### 2. 搜索布料 (公开接口)

根据条件搜索布料列表。

**接口地址**: `GET /api/fabrics/search`

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| q | string | 否 | 搜索关键词 |
| category | string | 否 | 分类筛选 |
| material | string | 否 | 材质筛选 |
| color | string | 否 | 颜色筛选 |
| min_price | number | 否 | 最低价格 |
| max_price | number | 否 | 最高价格 |
| min_stock | int | 否 | 最低库存 |
| status | int | 否 | 状态筛选 (1:可用, 0:停用) |
| page | int | 否 | 页码 (默认1) |
| page_size | int | 否 | 每页数量 (默认10) |

**响应示例**:
```json
{
  "total": 10,
  "page": 1,
  "page_size": 10,
  "fabrics": [
    {
      "id": 1,
      "name": "纯棉平纹布",
      "category": "棉布",
      "material": "100%棉",
      "color": "白色",
      "pattern": "平纹",
      "weight": 120.00,
      "width": 150.00,
      "price": 15.50,
      "unit": "米",
      "stock": 100,
      "min_order": 1,
      "description": "优质纯棉平纹布，透气性好，适合制作衬衫、T恤等",
      "image_url": "/uploads/fabrics/cotton_plain_white.jpg",
      "thumbnail_url": "/uploads/fabrics/thumbnails/cotton_plain_white.jpg",
      "tags": "棉布,白色,平纹,透气",
      "status": 1,
      "supplier_id": null,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### 3. 获取布料详情 (公开接口)

根据ID获取布料的详细信息。

**接口地址**: `GET /api/fabrics/{id}`

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 布料ID |

**响应示例**:
```json
{
  "id": 1,
  "name": "纯棉平纹布",
  "category": "棉布",
  "material": "100%棉",
  "color": "白色",
  "pattern": "平纹",
  "weight": 120.00,
  "width": 150.00,
  "price": 15.50,
  "unit": "米",
  "stock": 100,
  "min_order": 1,
  "description": "优质纯棉平纹布，透气性好，适合制作衬衫、T恤等",
  "image_url": "/uploads/fabrics/cotton_plain_white.jpg",
  "thumbnail_url": "/uploads/fabrics/thumbnails/cotton_plain_white.jpg",
  "tags": "棉布,白色,平纹,透气",
  "status": 1,
  "supplier_id": null,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### 4. 获取布料分类 (公开接口)

获取所有布料分类列表。

**接口地址**: `GET /api/fabrics/categories`

**响应示例**:
```json
[
  {
    "id": 1,
    "name": "棉布",
    "description": "天然棉纤维制成的布料，透气性好，适合制作夏季服装",
    "icon": "cotton",
    "sort": 1,
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

### 5. 根据分类获取布料 (公开接口)

根据分类获取布料列表。

**接口地址**: `GET /api/fabrics/category/{category}`

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| category | string | 是 | 分类名称 |

**查询参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码 (默认1) |
| page_size | int | 否 | 每页数量 (默认10) |

### 6. 根据材质获取布料 (公开接口)

根据材质获取布料列表。

**接口地址**: `GET /api/fabrics/material/{material}`

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| material | string | 是 | 材质名称 |

**查询参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码 (默认1) |
| page_size | int | 否 | 每页数量 (默认10) |

### 7. 获取布料统计信息 (公开接口)

获取布料的统计信息。

**接口地址**: `GET /api/fabrics/statistics`

**响应示例**:
```json
{
  "total_fabrics": 10,
  "available_fabrics": 8,
  "low_stock_fabrics": 2,
  "category_stats": [
    {
      "category": "棉布",
      "count": 4
    },
    {
      "category": "丝绸",
      "count": 3
    }
  ]
}
```

## 需要认证的接口

以下接口需要JWT Token认证，请在请求头中添加：
```
Authorization: Bearer <your_jwt_token>
```

### 8. 创建布料 (需要认证)

创建新的布料记录。

**接口地址**: `POST /api/fabrics`

**请求体**:
```json
{
  "name": "新布料名称",
  "category": "棉布",
  "material": "100%棉",
  "color": "白色",
  "pattern": "平纹",
  "weight": 120.00,
  "width": 150.00,
  "price": 15.50,
  "unit": "米",
  "stock": 100,
  "min_order": 1,
  "description": "布料描述",
  "image_url": "/uploads/fabrics/image.jpg",
  "thumbnail_url": "/uploads/fabrics/thumbnail.jpg",
  "tags": "棉布,白色,平纹",
  "supplier_id": "supplier_123"
}
```

### 9. 更新布料 (需要认证)

更新布料的详细信息。

**接口地址**: `PUT /api/fabrics/{id}`

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 布料ID |

**请求体**: 同创建布料，但所有字段都是可选的

### 10. 删除布料 (需要认证)

删除指定的布料记录。

**接口地址**: `DELETE /api/fabrics/{id}`

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 布料ID |

### 11. 更新布料库存 (需要认证)

更新指定布料的库存数量。

**接口地址**: `PUT /api/fabrics/{id}/stock`

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 布料ID |

**请求体**:
```json
{
  "quantity": 10
}
```

## 错误响应

当请求出现错误时，API会返回相应的HTTP状态码和错误信息：

```json
{
  "error": "错误描述信息"
}
```

常见HTTP状态码：
- `400 Bad Request`: 请求参数错误
- `401 Unauthorized`: 未认证或认证失败
- `404 Not Found`: 资源不存在
- `500 Internal Server Error`: 服务器内部错误

## 前端集成示例

### Flutter/Dart 示例

```dart
import 'package:http/http.dart' as http;
import 'dart:convert';

class FabricService {
  static const String baseUrl = 'https://aneworders.com/api';
  
  // 获取所有布料
  static Future<List<Fabric>> getAllFabrics() async {
    try {
      final response = await http.get(Uri.parse('$baseUrl/fabrics/all'));
      
      if (response.statusCode == 200) {
        final List<dynamic> data = json.decode(response.body);
        return data.map((json) => Fabric.fromJson(json)).toList();
      } else {
        throw Exception('Failed to load fabrics');
      }
    } catch (e) {
      throw Exception('Network error: $e');
    }
  }
  
  // 搜索布料
  static Future<FabricListResponse> searchFabrics({
    String? query,
    String? category,
    String? material,
    String? color,
    double? minPrice,
    double? maxPrice,
    int? minStock,
    int page = 1,
    int pageSize = 10,
  }) async {
    try {
      final queryParams = <String, String>{};
      if (query != null) queryParams['q'] = query;
      if (category != null) queryParams['category'] = category;
      if (material != null) queryParams['material'] = material;
      if (color != null) queryParams['color'] = color;
      if (minPrice != null) queryParams['min_price'] = minPrice.toString();
      if (maxPrice != null) queryParams['max_price'] = maxPrice.toString();
      if (minStock != null) queryParams['min_stock'] = minStock.toString();
      queryParams['page'] = page.toString();
      queryParams['page_size'] = pageSize.toString();
      
      final uri = Uri.parse('$baseUrl/fabrics/search').replace(queryParameters: queryParams);
      final response = await http.get(uri);
      
      if (response.statusCode == 200) {
        return FabricListResponse.fromJson(json.decode(response.body));
      } else {
        throw Exception('Failed to search fabrics');
      }
    } catch (e) {
      throw Exception('Network error: $e');
    }
  }
}

// 布料模型
class Fabric {
  final int id;
  final String name;
  final String category;
  final String material;
  final String color;
  final String pattern;
  final double weight;
  final double width;
  final double price;
  final String unit;
  final int stock;
  final int minOrder;
  final String description;
  final String imageUrl;
  final String thumbnailUrl;
  final String tags;
  final int status;
  final String? supplierId;
  final DateTime createdAt;
  final DateTime updatedAt;
  
  Fabric({
    required this.id,
    required this.name,
    required this.category,
    required this.material,
    required this.color,
    required this.pattern,
    required this.weight,
    required this.width,
    required this.price,
    required this.unit,
    required this.stock,
    required this.minOrder,
    required this.description,
    required this.imageUrl,
    required this.thumbnailUrl,
    required this.tags,
    required this.status,
    this.supplierId,
    required this.createdAt,
    required this.updatedAt,
  });
  
  factory Fabric.fromJson(Map<String, dynamic> json) {
    return Fabric(
      id: json['id'],
      name: json['name'],
      category: json['category'],
      material: json['material'],
      color: json['color'],
      pattern: json['pattern'],
      weight: json['weight'].toDouble(),
      width: json['width'].toDouble(),
      price: json['price'].toDouble(),
      unit: json['unit'],
      stock: json['stock'],
      minOrder: json['min_order'],
      description: json['description'],
      imageUrl: json['image_url'],
      thumbnailUrl: json['thumbnail_url'],
      tags: json['tags'],
      status: json['status'],
      supplierId: json['supplier_id'],
      createdAt: DateTime.parse(json['created_at']),
      updatedAt: DateTime.parse(json['updated_at']),
    );
  }
}

// 布料列表响应模型
class FabricListResponse {
  final int total;
  final int page;
  final int pageSize;
  final List<Fabric> fabrics;
  
  FabricListResponse({
    required this.total,
    required this.page,
    required this.pageSize,
    required this.fabrics,
  });
  
  factory FabricListResponse.fromJson(Map<String, dynamic> json) {
    return FabricListResponse(
      total: json['total'],
      page: json['page'],
      pageSize: json['page_size'],
      fabrics: (json['fabrics'] as List)
          .map((fabricJson) => Fabric.fromJson(fabricJson))
          .toList(),
    );
  }
}
```

## 注意事项

1. **图片URL**: 布料图片的URL是相对路径，需要与域名拼接成完整URL
2. **价格单位**: 价格以元为单位，精确到分
3. **库存管理**: 库存数量不能为负数
4. **认证权限**: 只有管理员和供应商可以管理布料信息
5. **分页限制**: 每页最大数量限制为100条记录
6. **搜索功能**: 支持模糊搜索，会在名称、描述、材质、颜色、图案、标签等字段中查找 