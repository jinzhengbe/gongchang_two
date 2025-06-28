#!/bin/bash

# 测试 images 字段合并修复的脚本
echo "=== 测试 images 字段合并修复 ==="

# 设置测试环境
BASE_URL="http://localhost:8008"
ORDER_ID="1"  # 假设订单ID为1

echo "1. 创建测试订单..."
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "title": "测试订单",
    "description": "测试图片合并功能",
    "quantity": 10,
    "designer_id": "test_designer",
    "customer_id": "test_customer",
    "images": ["img1", "img2"]
  }')

echo "创建订单响应: $CREATE_RESPONSE"

echo "2. 第一次更新 - 添加新图片..."
UPDATE1_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "title": "测试订单",
    "images": ["img3", "img4"]
  }')

echo "第一次更新响应: $UPDATE1_RESPONSE"

echo "3. 查看订单详情..."
DETAIL_RESPONSE=$(curl -s -X GET "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE")

echo "订单详情: $DETAIL_RESPONSE"

echo "4. 第二次更新 - 只更新标题，不传images..."
UPDATE2_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "title": "更新后的测试订单"
  }')

echo "第二次更新响应: $UPDATE2_RESPONSE"

echo "5. 再次查看订单详情..."
DETAIL2_RESPONSE=$(curl -s -X GET "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE")

echo "最终订单详情: $DETAIL2_RESPONSE"

echo "=== 测试完成 ==="
echo "预期结果:"
echo "- 第一次更新后，images 应该包含 [img1, img2, img3, img4]"
echo "- 第二次更新后，images 应该保持不变，不丢失" 