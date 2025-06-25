#!/bin/bash

# 添加布料到订单测试脚本
BASE_URL="http://localhost:8008/api"

echo "=== 添加布料到订单测试 ==="
echo "基础URL: $BASE_URL"
echo ""

# 测试数据
FABRIC_DATA='{
  "order_id": 1,
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

echo "1. 测试添加布料到订单..."
echo "请求数据: $FABRIC_DATA"
echo ""

response=$(curl -s -X POST "$BASE_URL/orders/1/add-fabric" \
  -H "Content-Type: application/json" \
  -d "$FABRIC_DATA")

echo "响应结果:"
echo "$response" | jq '.' 2>/dev/null || echo "$response"
echo ""

# 测试获取订单详情，查看布料是否已关联
echo "2. 测试获取订单详情（验证布料是否已关联）..."
curl -s -X GET "$BASE_URL/orders/1" | jq '.' 2>/dev/null || echo "获取订单详情失败"
echo ""

echo "=== 测试完成 ===" 