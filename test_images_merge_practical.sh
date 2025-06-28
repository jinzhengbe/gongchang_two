#!/bin/bash

# 实用的 images 字段合并测试脚本
echo "=== 测试 images 字段合并修复 ==="

# 设置测试环境
BASE_URL="http://localhost:8008"

# 测试用户信息（需要先注册或使用现有用户）
TEST_USERNAME="test_user_$(date +%s)"
TEST_PASSWORD="test123456"

echo "1. 注册测试用户..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$TEST_USERNAME\",
    \"password\": \"$TEST_PASSWORD\",
    \"email\": \"$TEST_USERNAME@test.com\",
    \"role\": \"designer\"
  }")

echo "注册响应: $REGISTER_RESPONSE"

# 提取用户ID（假设注册成功）
USER_ID=$(echo $REGISTER_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
if [ -z "$USER_ID" ]; then
    echo "❌ 注册失败，使用默认用户ID: 1"
    USER_ID="1"
fi

echo "2. 登录获取token..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$TEST_USERNAME\",
    \"password\": \"$TEST_PASSWORD\"
  }")

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
if [ -z "$TOKEN" ]; then
    echo "❌ 登录失败，无法获取token"
    exit 1
fi

echo "Token: ${TOKEN:0:20}..."

echo "3. 创建测试订单（包含初始图片）..."
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"title\": \"测试图片合并订单\",
    \"description\": \"测试images字段合并功能\",
    \"quantity\": 10,
    \"designer_id\": \"$USER_ID\",
    \"customer_id\": \"$USER_ID\",
    \"images\": [\"img_001\", \"img_002\"]
  }")

echo "创建订单响应: $CREATE_RESPONSE"

# 提取订单ID
ORDER_ID=$(echo $CREATE_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
if [ -z "$ORDER_ID" ]; then
    echo "❌ 创建订单失败"
    exit 1
fi

echo "订单ID: $ORDER_ID"

echo "4. 查看初始订单详情..."
INITIAL_DETAIL=$(curl -s -X GET "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "初始订单详情: $INITIAL_DETAIL"

echo "5. 第一次更新 - 添加新图片..."
UPDATE1_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"title\": \"测试图片合并订单\",
    \"images\": [\"img_003\", \"img_004\"]
  }")

echo "第一次更新响应: $UPDATE1_RESPONSE"

echo "6. 查看第一次更新后的订单详情..."
DETAIL1=$(curl -s -X GET "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "第一次更新后详情: $DETAIL1"

echo "7. 第二次更新 - 只更新标题，不传images字段..."
UPDATE2_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"title\": \"更新后的测试订单标题\"
  }")

echo "第二次更新响应: $UPDATE2_RESPONSE"

echo "8. 查看第二次更新后的订单详情..."
DETAIL2=$(curl -s -X GET "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "第二次更新后详情: $DETAIL2"

echo "9. 第三次更新 - 添加重复图片ID..."
UPDATE3_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"images\": [\"img_001\", \"img_005\"]
  }")

echo "第三次更新响应: $UPDATE3_RESPONSE"

echo "10. 查看最终订单详情..."
FINAL_DETAIL=$(curl -s -X GET "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "最终订单详情: $FINAL_DETAIL"

echo "=== 测试结果分析 ==="
echo "预期结果:"
echo "- 初始图片: [img_001, img_002]"
echo "- 第一次更新后: [img_001, img_002, img_003, img_004]"
echo "- 第二次更新后: 图片应该保持不变"
echo "- 最终结果: [img_001, img_002, img_003, img_004, img_005] (去重后)"

echo "=== 测试完成 ===" 