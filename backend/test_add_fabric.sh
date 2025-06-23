#!/bin/bash

# 布料添加测试脚本
BASE_URL="http://localhost:8008/api"

echo "=== 布料添加测试 ==="
echo "基础URL: $BASE_URL"
echo ""

# 测试数据
FABRIC_DATA='{
  "name": "测试棉布",
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
  "description": "测试用棉布，透气性好",
  "image_url": "/uploads/fabrics/test_cotton.jpg",
  "thumbnail_url": "/uploads/fabrics/thumbnails/test_cotton.jpg",
  "tags": "测试,棉布,白色,平纹"
}'

echo "1. 测试添加布料..."
echo "请求数据: $FABRIC_DATA"
echo ""

response=$(curl -s -X POST "$BASE_URL/fabrics/create" \
  -H "Content-Type: application/json" \
  -d "$FABRIC_DATA")

echo "响应结果:"
echo "$response" | jq '.' 2>/dev/null || echo "$response"
echo ""

# 测试获取所有布料
echo "2. 测试获取所有布料（验证是否添加成功）..."
curl -s -X GET "$BASE_URL/fabrics/all" | jq '.' 2>/dev/null || echo "获取布料列表失败"
echo ""

echo "=== 测试完成 ===" 