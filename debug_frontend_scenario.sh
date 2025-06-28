#!/bin/bash

# 模拟前端实际使用场景的调试脚本
echo "=== 模拟前端实际使用场景 ==="

# 设置测试环境
BASE_URL="http://localhost:8008"
LOG_FILE="frontend_scenario_$(date +%s).log"

# 记录日志的函数
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a $LOG_FILE
}

log "开始模拟前端实际使用场景..."

# 使用现有用户
TEST_USERNAME="frontend_user_$(date +%s)"
TEST_PASSWORD="frontend123456"

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

log "3. 创建订单（模拟前端创建新订单）..."
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"title\": \"前端测试订单\",
    \"description\": \"测试前端图片上传功能\",
    \"quantity\": 3,
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

log "5. 模拟前端上传第一张图片后自动调用 updateOrder（只传新图片ID）..."
UPDATE1_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"images\": [\"uploaded_file_001\"]
  }")

log "第一次更新响应: $UPDATE1_RESPONSE"

log "6. 查看第一次更新后的订单详情..."
DETAIL1=$(curl -s -X GET "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN")

log "第一次更新后详情: $DETAIL1"

log "7. 模拟前端上传第二张图片后自动调用 updateOrder（只传新图片ID）..."
UPDATE2_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"images\": [\"uploaded_file_002\"]
  }")

log "第二次更新响应: $UPDATE2_RESPONSE"

log "8. 查看第二次更新后的订单详情..."
DETAIL2=$(curl -s -X GET "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN")

log "第二次更新后详情: $DETAIL2"

log "9. 模拟前端点击保存按钮（可能不传images字段或传空）..."
SAVE_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"title\": \"前端测试订单 - 已保存\",
    \"description\": \"更新后的描述\"
  }")

log "保存响应: $SAVE_RESPONSE"

log "10. 查看保存后的订单详情..."
SAVE_DETAIL=$(curl -s -X GET "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN")

log "保存后详情: $SAVE_DETAIL"

log "11. 模拟前端刷新页面..."
REFRESH_DETAIL=$(curl -s -X GET "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN")

log "刷新后详情: $REFRESH_DETAIL"

log "12. 检查数据库中的实际数据..."
DB_DATA=$(mysql -h 192.168.0.10 -u gongchang -pgongchang gongchang -e "SELECT id, title, images FROM orders WHERE id = $ORDER_ID;" 2>/dev/null)
log "数据库中的数据: $DB_DATA"

echo ""
echo "=== 前端场景测试结果分析 ==="
echo "日志文件: $LOG_FILE"
echo ""
echo "关键检查点:"
echo "1. 初始订单 images 字段: $(echo $INITIAL_DETAIL | grep -o '"images":\[[^]]*\]')"
echo "2. 第一次上传后 images 字段: $(echo $DETAIL1 | grep -o '"images":\[[^]]*\]')"
echo "3. 第二次上传后 images 字段: $(echo $DETAIL2 | grep -o '"images":\[[^]]*\]')"
echo "4. 保存后 images 字段: $(echo $SAVE_DETAIL | grep -o '"images":\[[^]]*\]')"
echo "5. 刷新后 images 字段: $(echo $REFRESH_DETAIL | grep -o '"images":\[[^]]*\]')"
echo ""
echo "预期结果（修复后）:"
echo "- 初始: []"
echo "- 第一次上传后: [\"uploaded_file_001\"]"
echo "- 第二次上传后: [\"uploaded_file_001\", \"uploaded_file_002\"]"
echo "- 保存后: [\"uploaded_file_001\", \"uploaded_file_002\"] (保持不变)"
echo "- 刷新后: [\"uploaded_file_001\", \"uploaded_file_002\"] (保持不变)"
echo ""
echo "=== 前端场景测试完成 ===" 