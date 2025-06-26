#!/bin/bash

# 测试订单文件更新API
echo "=== 测试订单文件更新API ==="

# 1. 登录获取token
echo "1. 登录获取token..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8008/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "sdf",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "Token: $TOKEN"

if [ -z "$TOKEN" ]; then
    echo "登录失败，无法获取token"
    exit 1
fi

# 2. 获取订单列表，选择一个订单进行测试
echo "2. 获取订单列表..."
ORDERS_RESPONSE=$(curl -s -X GET http://localhost:8008/api/orders \
  -H "Authorization: Bearer $TOKEN")

echo "订单列表响应: $ORDERS_RESPONSE"

# 提取第一个订单ID
ORDER_ID=$(echo $ORDERS_RESPONSE | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
echo "选择的订单ID: $ORDER_ID"

if [ -z "$ORDER_ID" ]; then
    echo "没有找到订单，无法进行测试"
    exit 1
fi

# 3. 获取订单详情（更新前）
echo "3. 获取订单详情（更新前）..."
BEFORE_RESPONSE=$(curl -s -X GET http://localhost:8008/api/orders/$ORDER_ID \
  -H "Authorization: Bearer $TOKEN")

echo "更新前的订单详情: $BEFORE_RESPONSE"

# 4. 更新订单文件信息
echo "4. 更新订单文件信息..."
UPDATE_RESPONSE=$(curl -s -X PUT http://localhost:8008/api/orders/$ORDER_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "测试订单文件更新",
    "description": "测试文件保存功能",
    "images": ["847a56c1-8329-4f62-98f5-50b0b701eb00", "test-image-2"],
    "attachments": ["test-attachment-1", "test-attachment-2"],
    "models": ["test-model-1"],
    "videos": ["test-video-1"]
  }')

echo "更新响应: $UPDATE_RESPONSE"

# 5. 获取订单详情（更新后）
echo "5. 获取订单详情（更新后）..."
AFTER_RESPONSE=$(curl -s -X GET http://localhost:8008/api/orders/$ORDER_ID \
  -H "Authorization: Bearer $TOKEN")

echo "更新后的订单详情: $AFTER_RESPONSE"

# 6. 测试不同的字段名组合
echo "6. 测试不同的字段名组合..."
UPDATE_RESPONSE2=$(curl -s -X PUT http://localhost:8008/api/orders/$ORDER_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "测试多种字段名",
    "image_ids": ["new-image-1", "new-image-2"],
    "attachment_ids": ["new-attachment-1"],
    "model_ids": ["new-model-1", "new-model-2"],
    "video_ids": ["new-video-1"]
  }')

echo "多种字段名更新响应: $UPDATE_RESPONSE2"

# 7. 最终获取订单详情
echo "7. 最终获取订单详情..."
FINAL_RESPONSE=$(curl -s -X GET http://localhost:8008/api/orders/$ORDER_ID \
  -H "Authorization: Bearer $TOKEN")

echo "最终订单详情: $FINAL_RESPONSE"

echo "=== 测试完成 ===" 