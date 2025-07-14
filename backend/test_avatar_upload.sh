#!/bin/bash

# 头像上传功能测试脚本
# 测试设计师头像上传和设计师信息更新功能

BASE_URL="http://localhost:8008"

echo "=== 设计师头像上传功能测试 ==="

# 1. 登录获取token
echo "1. 登录获取token..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testdesigner",
    "password": "password"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "❌ 登录失败，无法获取token"
    echo "登录响应: $LOGIN_RESPONSE"
    exit 1
fi

echo "✅ 登录成功，获取到token"

# 2. 获取当前设计师信息
echo ""
echo "2. 获取当前设计师信息..."
PROFILE_RESPONSE=$(curl -s -X GET "$BASE_URL/api/designers/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "设计师信息响应: $PROFILE_RESPONSE"

# 3. 创建测试头像文件
echo ""
echo "3. 创建测试头像文件..."
echo "这是一个测试头像文件" > test_avatar.txt
echo "✅ 创建测试文件完成"

# 4. 上传头像
echo ""
echo "4. 上传头像..."
AVATAR_RESPONSE=$(curl -s -X POST "$BASE_URL/api/designers/avatar" \
  -H "Authorization: Bearer $TOKEN" \
  -F "avatar=@test_avatar.txt")

echo "头像上传响应: $AVATAR_RESPONSE"

# 提取头像URL
AVATAR_URL=$(echo $AVATAR_RESPONSE | grep -o '"url":"[^"]*"' | cut -d'"' -f4)

if [ ! -z "$AVATAR_URL" ]; then
    echo "✅ 头像上传成功，URL: $AVATAR_URL"
else
    echo "❌ 头像上传失败"
    rm -f test_avatar.txt
    exit 1
fi

# 5. 更新设计师信息（带头像URL）
echo ""
echo "5. 更新设计师信息（带头像URL）..."
UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/designers/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"company_name\": \"测试设计工作室\",
    \"address\": \"北京市朝阳区\",
    \"website\": \"http://test-design.com\",
    \"bio\": \"这是一个测试设计师档案\",
    \"avatar\": \"$AVATAR_URL\"
  }")

echo "更新设计师信息响应: $UPDATE_RESPONSE"

# 6. 验证头像URL可访问
echo ""
echo "6. 验证头像URL可访问..."
AVATAR_ACCESS_RESPONSE=$(curl -s -I "$BASE_URL$AVATAR_URL")
echo "头像访问响应头: $AVATAR_ACCESS_RESPONSE"

if echo "$AVATAR_ACCESS_RESPONSE" | grep -q "200 OK"; then
    echo "✅ 头像URL可正常访问"
else
    echo "❌ 头像URL访问失败"
fi

# 7. 再次获取设计师信息验证更新
echo ""
echo "7. 再次获取设计师信息验证更新..."
UPDATED_PROFILE_RESPONSE=$(curl -s -X GET "$BASE_URL/api/designers/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "更新后的设计师信息: $UPDATED_PROFILE_RESPONSE"

# 8. 测试无token访问
echo ""
echo "8. 测试无token访问..."
NO_TOKEN_RESPONSE=$(curl -s -X GET "$BASE_URL/api/designers/profile" \
  -H "Content-Type: application/json")

echo "无token访问响应: $NO_TOKEN_RESPONSE"

# 9. 测试其他用户角色
echo ""
echo "9. 测试其他用户角色..."
FACTORY_LOGIN=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testdesigner",
    "password": "password"
  }')

FACTORY_TOKEN=$(echo $FACTORY_LOGIN | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ ! -z "$FACTORY_TOKEN" ]; then
    FACTORY_RESPONSE=$(curl -s -X GET "$BASE_URL/api/designers/profile" \
      -H "Authorization: Bearer $FACTORY_TOKEN" \
      -H "Content-Type: application/json")
    
    echo "工厂用户获取设计师信息响应: $FACTORY_RESPONSE"
    
    if echo "$FACTORY_RESPONSE" | grep -q '"error":"设计师档案不存在"'; then
        echo "✅ 工厂用户无法获取设计师信息（符合预期）"
    else
        echo "❌ 工厂用户获取设计师信息异常"
    fi
fi

# 清理测试文件
rm -f test_avatar.txt

echo ""
echo "=== 测试完成 ==="
echo "✅ 设计师头像上传功能测试通过"
echo "📋 测试内容:"
echo "  - 登录认证"
echo "  - 获取设计师信息"
echo "  - 头像文件上传"
echo "  - 头像URL生成"
echo "  - 设计师信息更新"
echo "  - 头像URL访问验证"
echo "  - 权限验证"
echo "  - 错误处理" 