#!/bin/bash

# 修改密码API测试脚本
BASE_URL="http://localhost:8008/api"

echo "=== 修改密码API测试 ==="

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

# 测试修改密码API
echo -e "\n2. 测试修改密码API..."
echo "API: POST /api/users/change-password"

# 测试成功修改密码
echo "测试成功修改密码:"
RESPONSE=$(curl -s -X POST "$BASE_URL/users/change-password" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "123456",
    "new_password": "newpassword123"
  }')

echo "响应: $RESPONSE"

# 测试旧密码错误
echo -e "\n测试旧密码错误:"
RESPONSE2=$(curl -s -X POST "$BASE_URL/users/change-password" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "wrongpassword",
    "new_password": "newpassword123"
  }')

echo "响应: $RESPONSE2"

# 测试新密码太短
echo -e "\n测试新密码太短:"
RESPONSE3=$(curl -s -X POST "$BASE_URL/users/change-password" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "123456",
    "new_password": "123"
  }')

echo "响应: $RESPONSE3"

# 测试缺少参数
echo -e "\n测试缺少参数:"
RESPONSE4=$(curl -s -X POST "$BASE_URL/users/change-password" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "123456"
  }')

echo "响应: $RESPONSE4"

# 恢复原密码
echo -e "\n3. 恢复原密码..."
RESPONSE5=$(curl -s -X POST "$BASE_URL/users/change-password" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "newpassword123",
    "new_password": "123456"
  }')

echo "恢复原密码响应: $RESPONSE5"

echo -e "\n=== 测试完成 ===" 