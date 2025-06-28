#!/bin/bash

echo "=== 测试文件URL格式修复 ==="

# 设置测试环境
BASE_URL="http://localhost:8008"
LOG_FILE="test_url_fix_$(date +%s).log"

# 记录日志的函数
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a $LOG_FILE
}

log "开始测试文件URL格式..."

# 创建测试用户
TEST_USERNAME="url_test_user_$(date +%s)"
TEST_PASSWORD="url123456"

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

log "2. 登录获取token..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$TEST_USERNAME\",
    \"password\": \"$TEST_PASSWORD\"
  }")

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
    \"title\": \"URL格式测试订单\",
    \"description\": \"测试文件URL格式修复\",
    \"quantity\": 1,
    \"designer_id\": \"1\",
    \"customer_id\": \"1\"
  }")

# 提取订单ID
ORDER_ID=$(echo $CREATE_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
if [ -z "$ORDER_ID" ]; then
    log "❌ 创建订单失败"
    exit 1
fi

log "订单ID: $ORDER_ID"

log "4. 创建测试图片文件..."
echo "This is a test image for URL format" > test_url_image.jpg

log "5. 上传图片文件..."
UPLOAD_RESPONSE=$(curl -s -X POST "$BASE_URL/api/orders/$ORDER_ID/add-file" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_url_image.jpg" \
  -F "type=image" \
  -F "description=URL格式测试图片")

log "上传响应: $UPLOAD_RESPONSE"

log "6. 检查URL格式..."
# 提取返回的URL
RETURNED_URL=$(echo $UPLOAD_RESPONSE | grep -o '"url":"[^"]*"' | cut -d'"' -f4)

if [ -n "$RETURNED_URL" ]; then
    log "返回的URL: $RETURNED_URL"
    
    # 检查URL格式
    if [[ $RETURNED_URL == /uploads/* ]]; then
        log "✅ URL格式正确: 以 /uploads/ 开头"
    else
        log "❌ URL格式错误: 应该以 /uploads/ 开头"
    fi
    
    # 检查是否包含完整域名
    if [[ $RETURNED_URL == https://* ]]; then
        log "❌ URL格式错误: 不应该包含完整域名"
    else
        log "✅ URL格式正确: 不包含完整域名"
    fi
else
    log "❌ 无法提取URL"
fi

log "7. 测试直接访问文件..."
if [ -n "$RETURNED_URL" ]; then
    # 提取文件名
    FILENAME=$(echo $RETURNED_URL | sed 's|/uploads/||')
    log "尝试访问文件: $FILENAME"
    
    # 测试直接访问
    DIRECT_ACCESS=$(curl -s -I "$BASE_URL$RETURNED_URL")
    if [[ $DIRECT_ACCESS == *"200 OK"* ]]; then
        log "✅ 文件可以直接访问"
    else
        log "❌ 文件无法直接访问"
        log "响应: $DIRECT_ACCESS"
    fi
fi

log "8. 清理测试文件..."
rm -f test_url_image.jpg

echo ""
echo "=== 测试结果总结 ==="
echo "日志文件: $LOG_FILE"
echo ""
echo "预期结果:"
echo "- URL格式应该是: /uploads/filename.ext"
echo "- 不应该包含完整域名"
echo "- 文件应该可以直接访问"
echo ""
echo "=== 测试完成 ===" 