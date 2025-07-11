#!/bin/bash

# 直接测试工厂信息API
echo "=== 直接测试工厂信息API ==="

# 配置
API_BASE="http://localhost:8008/api"

# 1. 注册一个新用户
echo "1. 注册新用户..."
REGISTER_RESPONSE=$(curl -s -X POST "$API_BASE/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newfactory",
    "password": "test123",
    "email": "newfactory@test.com",
    "role": "factory",
    "company_name": "新工厂",
    "address": "广东省深圳市南山区"
  }')

echo "注册响应: $REGISTER_RESPONSE"

# 2. 登录获取token
echo ""
echo "2. 登录获取token..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "gongchang",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "获取到的token: $TOKEN"

if [ -z "$TOKEN" ]; then
    echo "❌ 登录失败，无法获取token"
    exit 1
fi

echo "✅ 登录成功"

# 3. 获取工厂详细信息
echo ""
echo "3. 获取工厂详细信息..."
GET_PROFILE_RESPONSE=$(curl -s -X GET "$API_BASE/factories/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "获取工厂信息响应: $GET_PROFILE_RESPONSE"

# 4. 更新工厂详细信息
echo ""
echo "4. 更新工厂详细信息..."
UPDATE_RESPONSE=$(curl -s -X PUT "$API_BASE/factories/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "company_name": "更新后的工厂",
    "address": "广东省深圳市南山区科技园",
    "capacity": 2000,
    "equipment": "先进缝纫设备、裁剪设备、熨烫设备、激光切割机",
    "certificates": "ISO9001认证、OHSAS18001认证、环保认证",
    "photos": [
      "https://example.com/factory1.jpg",
      "https://example.com/factory2.jpg",
      "https://example.com/factory3.jpg"
    ],
    "videos": [
      "https://example.com/factory1.mp4",
      "https://example.com/factory2.mp4"
    ],
    "employee_count": 100
  }')

echo "更新工厂信息响应: $UPDATE_RESPONSE"

# 5. 再次获取工厂详细信息验证更新
echo ""
echo "5. 验证更新结果..."
GET_UPDATED_RESPONSE=$(curl -s -X GET "$API_BASE/factories/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "更新后的工厂信息: $GET_UPDATED_RESPONSE"

echo ""
echo "=== 测试完成 ===" 