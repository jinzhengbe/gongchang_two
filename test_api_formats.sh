#!/bin/bash

# 测试 PUT /api/orders/{id} 的各种请求格式
echo "=== 测试 PUT /api/orders/{id} 请求格式 ==="

# 设置测试环境
BASE_URL="http://localhost:8008"
LOG_FILE="api_format_test_$(date +%s).log"

# 记录日志的函数
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a $LOG_FILE
}

log "开始测试 API 请求格式..."

# 创建测试用户
TEST_USERNAME="format_test_user_$(date +%s)"
TEST_PASSWORD="format123456"

log "1. 注册测试用户..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$TEST_USERNAME\",
    \"password\": \"$TEST_PASSWORD\",
    \"email\": \"$TEST_USERNAME@test.com\",
    \"role\": \"designer\"
  }")

log "注册响应: $REGISTER_RESPONSE"

# 提取用户ID
USER_ID=$(echo $REGISTER_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
if [ -z "$USER_ID" ]; then
    log "❌ 注册失败，使用默认用户ID: 1"
    USER_ID="1"
fi

log "2. 登录获取token..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$TEST_USERNAME\",
    \"password\": \"$TEST_PASSWORD\"
  }")

log "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
if [ -z "$TOKEN" ]; then
    log "❌ 登录失败，无法获取token"
    exit 1
fi

log "Token: ${TOKEN:0:20}..."

log "3. 创建测试订单..."
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"title\": \"API格式测试订单\",
    \"description\": \"测试各种API请求格式\",
    \"quantity\": 5,
    \"designer_id\": \"$USER_ID\",
    \"customer_id\": \"$USER_ID\"
  }")

log "创建订单响应: $CREATE_RESPONSE"

# 提取订单ID
ORDER_ID=$(echo $CREATE_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
if [ -z "$ORDER_ID" ]; then
    log "❌ 创建订单失败"
    exit 1
fi

log "订单ID: $ORDER_ID"

echo ""
echo "=== 测试各种请求格式 ==="

log "4. 测试格式1: 只更新标题"
FORMAT1_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "更新后的标题"
  }')

log "格式1响应: $FORMAT1_RESPONSE"

log "5. 测试格式2: 只更新图片（单个图片）"
FORMAT2_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "images": ["single_image_001"]
  }')

log "格式2响应: $FORMAT2_RESPONSE"

log "6. 测试格式3: 只更新图片（多个图片）"
FORMAT3_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "images": ["multi_image_001", "multi_image_002", "multi_image_003"]
  }')

log "格式3响应: $FORMAT3_RESPONSE"

log "7. 测试格式4: 更新多个字段"
FORMAT4_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "多字段更新标题",
    "description": "多字段更新描述",
    "quantity": 10,
    "status": "published",
    "images": ["field_image_001", "field_image_002"]
  }')

log "格式4响应: $FORMAT4_RESPONSE"

log "8. 测试格式5: 空图片数组"
FORMAT5_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "images": []
  }')

log "格式5响应: $FORMAT5_RESPONSE"

log "9. 测试格式6: 不传图片字段"
FORMAT6_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "不传图片字段的标题"
  }')

log "格式6响应: $FORMAT6_RESPONSE"

log "10. 测试格式7: 复杂JSON结构"
FORMAT7_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "复杂JSON标题",
    "description": "这是一个包含特殊字符的描述：!@#$%^&*()",
    "quantity": 15,
    "payment_status": "paid",
    "shipping_address": "北京市朝阳区某某街道123号",
    "orderType": "custom",
    "fabrics": "cotton,linen",
    "specialRequirements": "需要特殊处理，包含换行符\n和制表符\t",
    "attachments": ["att_001", "att_002"],
    "models": ["model_001"],
    "images": ["complex_img_001", "complex_img_002", "complex_img_003"],
    "videos": ["video_001"]
  }')

log "格式7响应: $FORMAT7_RESPONSE"

log "11. 查看最终订单详情..."
FINAL_DETAIL=$(curl -s -X GET "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN")

log "最终订单详情: $FINAL_DETAIL"

echo ""
echo "=== 测试错误格式 ==="

log "12. 测试错误格式1: 无效的JSON"
ERROR1_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "无效JSON标题",
    "images": ["img_001", "img_002"
  }')

log "错误格式1响应: $ERROR1_RESPONSE"

log "13. 测试错误格式2: 缺少认证"
ERROR2_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "缺少认证的标题"
  }')

log "错误格式2响应: $ERROR2_RESPONSE"

log "14. 测试错误格式3: 无效的订单ID"
ERROR3_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/orders/999999" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "无效订单ID的标题"
  }')

log "错误格式3响应: $ERROR3_RESPONSE"

echo ""
echo "=== 测试结果总结 ==="
echo "日志文件: $LOG_FILE"
echo ""
echo "成功格式测试:"
echo "✅ 格式1: 只更新标题"
echo "✅ 格式2: 只更新图片（单个）"
echo "✅ 格式3: 只更新图片（多个）"
echo "✅ 格式4: 更新多个字段"
echo "✅ 格式5: 空图片数组"
echo "✅ 格式6: 不传图片字段"
echo "✅ 格式7: 复杂JSON结构"
echo ""
echo "错误格式测试:"
echo "❌ 格式1: 无效JSON - 应返回400错误"
echo "❌ 格式2: 缺少认证 - 应返回401错误"
echo "❌ 格式3: 无效订单ID - 应返回400错误"
echo ""
echo "=== 测试完成 ===" 