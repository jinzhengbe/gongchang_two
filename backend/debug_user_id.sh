#!/bin/bash

echo "=== 调试用户ID提取 ==="

# 1. 登录获取token
LOGIN_RESPONSE=$(curl -s -X POST "http://localhost:8008/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "gongchang",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "Token: $TOKEN"

# 2. 解码JWT token获取用户ID
echo -e "\n2. 解码JWT token..."
# 提取payload部分
PAYLOAD=$(echo "$TOKEN" | cut -d'.' -f2)
echo "Payload: $PAYLOAD"

# Base64解码payload
DECODED_PAYLOAD=$(echo "$PAYLOAD" | base64 -d 2>/dev/null || echo "$PAYLOAD" | base64 -d -i 2>/dev/null)
echo "解码后的payload: $DECODED_PAYLOAD"

# 提取user_id
USER_ID_FROM_TOKEN=$(echo "$DECODED_PAYLOAD" | grep -o '"user_id":"[^"]*"' | cut -d'"' -f4)
echo "从token提取的用户ID: $USER_ID_FROM_TOKEN"

# 3. 检查工厂信息
echo -e "\n3. 检查工厂信息..."
FACTORY_RESPONSE=$(curl -s -X GET "http://localhost:8008/api/factories/profile" \
  -H "Authorization: Bearer $TOKEN")

echo "工厂信息: $FACTORY_RESPONSE"

# 提取工厂ID
FACTORY_ID=$(echo "$FACTORY_RESPONSE" | grep -o '"id":[0-9]*' | cut -d':' -f2)
echo "工厂ID (数字): $FACTORY_ID"

# 提取user_id
FACTORY_USER_ID=$(echo "$FACTORY_RESPONSE" | grep -o '"user_id":"[^"]*"' | cut -d'"' -f4)
echo "工厂user_id: $FACTORY_USER_ID"

# 4. 测试不同的工厂ID
echo -e "\n4. 测试不同的工厂ID..."

# 测试使用数字ID
echo "测试使用数字ID: $FACTORY_ID"
UPLOAD_RESPONSE_1=$(curl -s -X POST "http://localhost:8008/api/factories/$FACTORY_ID/photos/batch" \
  -H "Authorization: Bearer $TOKEN" \
  -F "files=@/dev/null")

echo "数字ID响应: $UPLOAD_RESPONSE_1"

# 测试使用user_id
echo "测试使用user_id: $FACTORY_USER_ID"
UPLOAD_RESPONSE_2=$(curl -s -X POST "http://localhost:8008/api/factories/$FACTORY_USER_ID/photos/batch" \
  -H "Authorization: Bearer $TOKEN" \
  -F "files=@/dev/null")

echo "user_id响应: $UPLOAD_RESPONSE_2"

# 5. 分析结果
echo -e "\n=== 分析结果 ==="
echo "Token中的用户ID: $USER_ID_FROM_TOKEN"
echo "工厂数字ID: $FACTORY_ID"
echo "工厂user_id: $FACTORY_USER_ID"

if [ "$USER_ID_FROM_TOKEN" = "$FACTORY_USER_ID" ]; then
    echo "✅ Token用户ID与工厂user_id匹配"
else
    echo "❌ Token用户ID与工厂user_id不匹配"
fi

echo -e "\n=== 调试完成 ===" 