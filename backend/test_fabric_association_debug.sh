#!/bin/bash

# 配置
API_URL="http://127.0.0.1:8008/api"
USERNAME="sdf"
PASSWORD="123456"

echo "=== 布料关联调试测试 ==="

# 1. 登录获取token
echo "1. 登录获取token..."
LOGIN_RESP=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "'$USERNAME'", "password": "'$PASSWORD'"}')
TOKEN=$(echo "$LOGIN_RESP" | grep -o '"token":"[^"]*"' | cut -d '"' -f4)
USER_ID=$(echo "$LOGIN_RESP" | grep -o '"id":"[^"]*"' | cut -d '"' -f4)
echo "Token: $TOKEN"
echo "User ID: $USER_ID"

# 2. 创建订单
echo ""
echo "2. 创建订单..."
ORDER_RESP=$(curl -s -X POST "$API_URL/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "调试测试订单",
    "description": "用于调试布料关联的订单",
    "quantity": 5,
    "designer_id": "'$USER_ID'",
    "customer_id": "'$USER_ID'"
  }')
ORDER_ID=$(echo "$ORDER_RESP" | grep -o '"id":[0-9]*' | cut -d ':' -f2)
echo "Order ID: $ORDER_ID"

# 3. 查询订单详情（添加布料前）
echo ""
echo "3. 查询订单详情（添加布料前）..."
BEFORE_RESP=$(curl -s -X GET "$API_URL/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN")
echo "添加布料前的订单详情:"
echo "$BEFORE_RESP" | jq '.'

# 4. 添加布料到订单
echo ""
echo "4. 添加布料到订单..."
ADD_FABRIC_RESP=$(curl -s -X POST "$API_URL/orders/$ORDER_ID/add-fabric" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "order_id": '$ORDER_ID',
    "name": "调试测试布料",
    "composition": "100%棉",
    "color": "蓝色",
    "width": 160.0,
    "weight": 130.0,
    "price": 30.00
  }')
echo "添加布料响应:"
echo "$ADD_FABRIC_RESP" | jq '.'

# 5. 查询订单详情（添加布料后）
echo ""
echo "5. 查询订单详情（添加布料后）..."
AFTER_RESP=$(curl -s -X GET "$API_URL/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN")
echo "添加布料后的订单详情:"
echo "$AFTER_RESP" | jq '.'

# 6. 再次添加一个布料
echo ""
echo "6. 再次添加一个布料..."
ADD_FABRIC_RESP2=$(curl -s -X POST "$API_URL/orders/$ORDER_ID/add-fabric" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "order_id": '$ORDER_ID',
    "name": "调试测试布料2",
    "composition": "100%丝绸",
    "color": "红色",
    "width": 140.0,
    "weight": 110.0,
    "price": 50.00
  }')
echo "添加第二个布料响应:"
echo "$ADD_FABRIC_RESP2" | jq '.'

# 7. 最终查询订单详情
echo ""
echo "7. 最终查询订单详情..."
FINAL_RESP=$(curl -s -X GET "$API_URL/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN")
echo "最终订单详情:"
echo "$FINAL_RESP" | jq '.'

echo ""
echo "=== 调试测试完成 ===" 