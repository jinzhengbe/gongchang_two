#!/bin/bash

# 测试新的文件上传并关联到订单的接口
echo "=== 测试 POST /api/orders/{orderId}/add-file 接口 ==="

# 设置测试环境
BASE_URL="http://localhost:8008"
LOG_FILE="test_add_file_$(date +%s).log"

# 记录日志的函数
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a $LOG_FILE
}

log "开始测试新的文件上传接口..."

# 创建测试用户
TEST_USERNAME="file_test_user_$(date +%s)"
TEST_PASSWORD="file123456"

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
    \"title\": \"文件上传测试订单\",
    \"description\": \"测试新的文件上传接口\",
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

log "4. 查看初始订单详情..."
INITIAL_DETAIL=$(curl -s -X GET "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN")

log "初始订单详情: $INITIAL_DETAIL"

log "5. 创建测试图片文件..."
echo "This is a test image content" > test_image.jpg
echo "This is a test attachment content" > test_attachment.pdf
echo "This is a test model content" > test_model.stl
echo "This is a test video content" > test_video.mp4

log "6. 测试上传图片文件..."
IMAGE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/orders/$ORDER_ID/add-file" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_image.jpg" \
  -F "type=image" \
  -F "description=测试图片")

log "图片上传响应: $IMAGE_RESPONSE"

log "7. 查看上传图片后的订单详情..."
AFTER_IMAGE_DETAIL=$(curl -s -X GET "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN")

log "上传图片后订单详情: $AFTER_IMAGE_DETAIL"

log "8. 测试上传附件文件..."
ATTACHMENT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/orders/$ORDER_ID/add-file" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_attachment.pdf" \
  -F "type=attachment" \
  -F "description=测试附件")

log "附件上传响应: $ATTACHMENT_RESPONSE"

log "9. 测试上传模型文件..."
MODEL_RESPONSE=$(curl -s -X POST "$BASE_URL/api/orders/$ORDER_ID/add-file" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_model.stl" \
  -F "type=model" \
  -F "description=测试模型")

log "模型上传响应: $MODEL_RESPONSE"

log "10. 测试上传视频文件..."
VIDEO_RESPONSE=$(curl -s -X POST "$BASE_URL/api/orders/$ORDER_ID/add-file" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_video.mp4" \
  -F "type=video" \
  -F "description=测试视频")

log "视频上传响应: $VIDEO_RESPONSE"

log "11. 查看最终订单详情..."
FINAL_DETAIL=$(curl -s -X GET "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN")

log "最终订单详情: $FINAL_DETAIL"

log "12. 清理测试文件..."
rm -f test_image.jpg test_attachment.pdf test_model.stl test_video.mp4

echo ""
echo "=== 测试结果分析 ==="
echo "日志文件: $LOG_FILE"
echo ""
echo "关键检查点:"
echo "1. 初始订单文件字段: $(echo $INITIAL_DETAIL | grep -o '"images":\[[^]]*\]')"
echo "2. 上传图片后: $(echo $AFTER_IMAGE_DETAIL | grep -o '"images":\[[^]]*\]')"
echo "3. 最终订单详情:"
echo "   - images: $(echo $FINAL_DETAIL | grep -o '"images":\[[^]]*\]')"
echo "   - attachments: $(echo $FINAL_DETAIL | grep -o '"attachments":\[[^]]*\]')"
echo "   - models: $(echo $FINAL_DETAIL | grep -o '"models":\[[^]]*\]')"
echo "   - videos: $(echo $FINAL_DETAIL | grep -o '"videos":\[[^]]*\]')"
echo ""
echo "预期结果:"
echo "- 初始: 所有文件字段为空数组"
echo "- 上传图片后: images 字段包含1个文件ID"
echo "- 最终: 每个类型字段包含对应的文件ID"
echo ""
echo "=== 测试完成 ===" 