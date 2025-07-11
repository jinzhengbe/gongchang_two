#!/bin/bash

# 工厂信息API测试脚本
echo "=== 工厂信息API测试 ==="

# 配置
API_BASE="http://localhost:8080/api"
TOKEN=""

# 登录获取token
echo "1. 登录获取token..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "factory1",
    "password": "test123"
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

# 2. 获取工厂详细信息
echo ""
echo "2. 获取工厂详细信息..."
GET_PROFILE_RESPONSE=$(curl -s -X GET "$API_BASE/factories/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "获取工厂信息响应: $GET_PROFILE_RESPONSE"

# 3. 更新工厂详细信息
echo ""
echo "3. 更新工厂详细信息..."
UPDATE_RESPONSE=$(curl -s -X PUT "$API_BASE/factories/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "company_name": "测试工厂更新",
    "address": "广东省深圳市南山区科技园",
    "capacity": 1000,
    "equipment": "先进缝纫设备、裁剪设备、熨烫设备",
    "certificates": "ISO9001认证、OHSAS18001认证",
    "photos": [
      "https://example.com/factory1.jpg",
      "https://example.com/factory2.jpg"
    ],
    "videos": [
      "https://example.com/factory1.mp4",
      "https://example.com/factory2.mp4"
    ],
    "employee_count": 50
  }')

echo "更新工厂信息响应: $UPDATE_RESPONSE"

# 4. 再次获取工厂详细信息验证更新
echo ""
echo "4. 验证更新结果..."
GET_UPDATED_RESPONSE=$(curl -s -X GET "$API_BASE/factories/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "更新后的工厂信息: $GET_UPDATED_RESPONSE"

echo ""
echo "=== 测试完成 ===" 