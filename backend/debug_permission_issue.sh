#!/bin/bash

# 调试权限问题的脚本
# 检查用户ID、工厂ID和权限校验逻辑

echo "=== 调试权限问题 ==="

# 1. 检查当前登录用户信息
echo "1. 检查当前登录用户信息..."
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

if [ -z "$TOKEN" ]; then
    echo "❌ 无法获取token"
    exit 1
fi

# 2. 检查用户信息
echo -e "\n2. 检查用户信息..."
USER_INFO=$(curl -s -X GET "http://localhost:8008/api/auth/profile" \
  -H "Authorization: Bearer $TOKEN")

echo "用户信息: $USER_INFO"

# 提取用户ID
USER_ID=$(echo "$USER_INFO" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "用户ID: $USER_ID"

# 3. 检查工厂信息
echo -e "\n3. 检查工厂信息..."
FACTORY_RESPONSE=$(curl -s -X GET "http://localhost:8008/api/factories/profile" \
  -H "Authorization: Bearer $TOKEN")

echo "工厂信息: $FACTORY_RESPONSE"

# 提取工厂ID
FACTORY_ID=$(echo "$FACTORY_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "工厂ID: $FACTORY_ID"

# 4. 检查数据库中的工厂记录
echo -e "\n4. 检查数据库中的工厂记录..."
echo "用户ID: $USER_ID"
echo "工厂ID: $FACTORY_ID"

# 5. 测试权限校验逻辑
echo -e "\n5. 测试权限校验逻辑..."

# 使用正确的工厂ID测试上传
if [ ! -z "$FACTORY_ID" ]; then
    echo "使用正确的工厂ID测试上传: $FACTORY_ID"
    
    # 创建测试图片
    echo "创建测试图片..."
    convert -size 100x100 xc:red test_permission.jpg 2>/dev/null || echo "red" > test_permission.jpg
    
    UPLOAD_RESPONSE=$(curl -s -X POST "http://localhost:8008/api/factories/$FACTORY_ID/photos/batch" \
      -H "Authorization: Bearer $TOKEN" \
      -F "files=@test_permission.jpg")
    
    echo "上传响应: $UPLOAD_RESPONSE"
    
    # 清理测试文件
    rm -f test_permission.jpg
else
    echo "❌ 无法获取工厂ID"
fi

# 6. 检查是否有其他工厂记录
echo -e "\n6. 检查是否有其他工厂记录..."
echo "尝试使用用户ID作为工厂ID..."

UPLOAD_WITH_USER_ID=$(curl -s -X POST "http://localhost:8008/api/factories/$USER_ID/photos/batch" \
  -H "Authorization: Bearer $TOKEN" \
  -F "files=@/dev/null")

echo "使用用户ID作为工厂ID的响应: $UPLOAD_WITH_USER_ID"

# 7. 分析问题
echo -e "\n=== 问题分析 ==="
echo "用户ID: $USER_ID"
echo "工厂ID: $FACTORY_ID"

if [ "$USER_ID" = "$FACTORY_ID" ]; then
    echo "✅ 用户ID和工厂ID匹配"
else
    echo "❌ 用户ID和工厂ID不匹配"
    echo "这是权限问题的根本原因"
fi

echo -e "\n=== 解决方案 ==="
echo "1. 确保工厂记录中的user_id字段与当前用户ID一致"
echo "2. 或者修改权限校验逻辑，允许用户操作自己的工厂"
echo "3. 检查数据库中users表的id和user_id字段"

echo -e "\n=== 调试完成 ===" 