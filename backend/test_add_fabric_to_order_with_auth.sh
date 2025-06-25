#!/bin/bash

# 配置
API_URL="http://127.0.0.1:8008/api"
USERNAME="sdf"   # 请替换为你的设计师用户名
PASSWORD="123456"      # 请替换为你的设计师密码

# 1. 登录获取token
echo "== 登录获取token =="
LOGIN_RESP=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "'$USERNAME'", "password": "'$PASSWORD'"}')
TOKEN=$(echo "$LOGIN_RESP" | grep -o '"token":"[^"]*"' | cut -d '"' -f4)
USER_ID=$(echo "$LOGIN_RESP" | grep -o '"id":"[^"]*"' | cut -d '"' -f4)
if [ -z "$TOKEN" ]; then
  echo "登录失败: $LOGIN_RESP"
  exit 1
fi
echo "Token: $TOKEN"
echo "User ID: $USER_ID"

# 2. 创建订单
echo "== 创建订单 =="
ORDER_RESP=$(curl -s -X POST "$API_URL/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "测试订单",
    "description": "自动化测试订单",
    "quantity": 10,
    "designer_id": "'$USER_ID'",
    "customer_id": "'$USER_ID'"
  }')
ORDER_ID=$(echo "$ORDER_RESP" | grep -o '"id":[0-9]*' | cut -d ':' -f2)
if [ -z "$ORDER_ID" ]; then
  echo "创建订单失败: $ORDER_RESP"
  exit 1
fi
echo "Order ID: $ORDER_ID"

# 3. 添加布料到订单
echo "== 添加布料到订单 =="
ADD_FABRIC_RESP=$(curl -s -X POST "$API_URL/orders/$ORDER_ID/add-fabric" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "order_id": '$ORDER_ID',
    "name": "测试布料",
    "composition": "棉100%",
    "color": "白色",
    "width": 150.0,
    "weight": 120.0
  }')
echo "添加布料响应: $ADD_FABRIC_RESP"

# 4. 测试获取订单详情，查看布料是否已关联
echo "4. 测试获取订单详情（验证布料是否已关联）..."
curl -s -X GET "$API_URL/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN" | jq '.' 2>/dev/null || echo "获取订单详情失败"
echo ""

echo "=== 测试完成 ===" 