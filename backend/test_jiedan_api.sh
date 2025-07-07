#!/bin/bash

# 测试接单API接口
BASE_URL="http://localhost:8008/api"

echo "=== 测试接单API接口 ==="

# 获取认证token
echo "1. 获取认证token..."
TOKEN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "gongchang",
    "password": "123456"
  }')

TOKEN=$(echo $TOKEN_RESPONSE | jq -r '.token')
if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo "获取token失败: $TOKEN_RESPONSE"
    exit 1
fi
echo "Token获取成功"

# 测试获取订单下指定工厂的接单记录
echo -e "\n2. 测试获取订单下指定工厂的接单记录..."
echo "API: GET /api/orders/{orderId}/jiedan?factory_id={factoryId}"

# 测试订单39和工厂factory1
echo "测试订单39和工厂factory1:"
RESPONSE=$(curl -s -X GET "$BASE_URL/orders/39/jiedan?factory_id=factory1" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "响应: $RESPONSE"

# 测试不存在的工厂
echo -e "\n测试订单39和不存在的工厂nonexistent:"
RESPONSE2=$(curl -s -X GET "$BASE_URL/orders/39/jiedan?factory_id=nonexistent" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "响应: $RESPONSE2"

# 测试无效的订单ID
echo -e "\n测试无效的订单ID:"
RESPONSE3=$(curl -s -X GET "$BASE_URL/orders/999999/jiedan?factory_id=factory1" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "响应: $RESPONSE3"

# 测试缺少factory_id参数
echo -e "\n测试缺少factory_id参数:"
RESPONSE4=$(curl -s -X GET "$BASE_URL/orders/39/jiedan" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "响应: $RESPONSE4"

echo -e "\n=== 测试完成 ===" 