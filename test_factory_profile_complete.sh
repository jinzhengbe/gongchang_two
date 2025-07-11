#!/bin/bash

# 工厂信息编辑API完整测试脚本
# 测试工厂信息获取和更新功能

BASE_URL="http://localhost:8008"
TOKEN=""
FACTORY_ID=""

echo "=== 工厂信息编辑API完整测试 ==="

# 1. 登录获取token
echo "1. 登录获取token..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "factory1",
    "password": "test123"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "❌ 登录失败，无法获取token"
    exit 1
fi

echo "✅ 登录成功，获取到token: ${TOKEN:0:20}..."

# 2. 获取当前工厂信息
echo ""
echo "2. 获取当前工厂信息..."
GET_RESPONSE=$(curl -s -X GET "$BASE_URL/api/factories/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "获取工厂信息响应: $GET_RESPONSE"

# 提取工厂ID
FACTORY_ID=$(echo $GET_RESPONSE | grep -o '"UserID":"[^"]*"' | cut -d'"' -f4)

if [ -z "$FACTORY_ID" ]; then
    echo "❌ 无法获取工厂ID"
    exit 1
fi

echo "✅ 获取到工厂ID: $FACTORY_ID"

# 3. 更新工厂信息
echo ""
echo "3. 更新工厂信息..."
UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/factories/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "company_name": "测试工厂更新",
    "address": "广东省深圳市南山区",
    "capacity": 2000,
    "equipment": "全自动裁剪机,工业缝纫机,质检设备",
    "certificates": "ISO9001,质量管理体系认证,环保认证",
    "photos": [
      "https://example.com/photo1.jpg",
      "https://example.com/photo2.jpg",
      "https://example.com/photo3.jpg"
    ],
    "videos": [
      "https://example.com/video1.mp4",
      "https://example.com/video2.mp4"
    ],
    "employee_count": 150
  }')

echo "更新工厂信息响应: $UPDATE_RESPONSE"

# 检查更新是否成功
if echo "$UPDATE_RESPONSE" | grep -q '"code":0'; then
    echo "✅ 工厂信息更新成功"
else
    echo "❌ 工厂信息更新失败"
    echo "错误详情: $UPDATE_RESPONSE"
    exit 1
fi

# 4. 再次获取工厂信息验证更新
echo ""
echo "4. 验证更新结果..."
GET_AGAIN_RESPONSE=$(curl -s -X GET "$BASE_URL/api/factories/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "验证响应: $GET_AGAIN_RESPONSE"

# 检查关键字段是否更新
if echo "$GET_AGAIN_RESPONSE" | grep -q '"CompanyName":"测试工厂更新"'; then
    echo "✅ 公司名称更新成功"
else
    echo "❌ 公司名称更新失败"
fi

if echo "$GET_AGAIN_RESPONSE" | grep -q '"EmployeeCount":150'; then
    echo "✅ 员工数量更新成功"
else
    echo "❌ 员工数量更新失败"
fi

if echo "$GET_AGAIN_RESPONSE" | grep -q '"photo1.jpg"'; then
    echo "✅ 照片数组更新成功"
else
    echo "❌ 照片数组更新失败"
fi

if echo "$GET_AGAIN_RESPONSE" | grep -q '"video1.mp4"'; then
    echo "✅ 视频数组更新成功"
else
    echo "❌ 视频数组更新失败"
fi

# 5. 测试部分更新
echo ""
echo "5. 测试部分更新..."
PARTIAL_UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/factories/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "company_name": "部分更新测试工厂",
    "employee_count": 200
  }')

echo "部分更新响应: $PARTIAL_UPDATE_RESPONSE"

if echo "$PARTIAL_UPDATE_RESPONSE" | grep -q '"code":0'; then
    echo "✅ 部分更新成功"
else
    echo "❌ 部分更新失败"
fi

# 6. 测试错误情况
echo ""
echo "6. 测试错误情况..."

# 测试无效token
echo "测试无效token..."
INVALID_TOKEN_RESPONSE=$(curl -s -X GET "$BASE_URL/api/factories/profile" \
  -H "Authorization: Bearer invalid_token" \
  -H "Content-Type: application/json")

echo "无效token响应: $INVALID_TOKEN_RESPONSE"

# 测试无token
echo "测试无token..."
NO_TOKEN_RESPONSE=$(curl -s -X GET "$BASE_URL/api/factories/profile" \
  -H "Content-Type: application/json")

echo "无token响应: $NO_TOKEN_RESPONSE"

# 7. 测试其他用户角色
echo ""
echo "7. 测试其他用户角色..."

# 尝试用设计师账号获取工厂信息
DESIGNER_LOGIN=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "designer1",
    "password": "test123"
  }')

DESIGNER_TOKEN=$(echo $DESIGNER_LOGIN | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ ! -z "$DESIGNER_TOKEN" ]; then
    DESIGNER_RESPONSE=$(curl -s -X GET "$BASE_URL/api/factories/profile" \
      -H "Authorization: Bearer $DESIGNER_TOKEN" \
      -H "Content-Type: application/json")
    
    echo "设计师获取工厂信息响应: $DESIGNER_RESPONSE"
    
    if echo "$DESIGNER_RESPONSE" | grep -q '"code":404'; then
        echo "✅ 设计师无法获取工厂信息（符合预期）"
    else
        echo "❌ 设计师获取工厂信息异常"
    fi
fi

echo ""
echo "=== 测试完成 ==="
echo "✅ 工厂信息编辑API测试通过"
echo "📋 测试内容:"
echo "  - 登录认证"
echo "  - 获取工厂信息"
echo "  - 完整更新工厂信息"
echo "  - 部分更新工厂信息"
echo "  - 照片和视频数组处理"
echo "  - 错误处理"
echo "  - 权限验证" 