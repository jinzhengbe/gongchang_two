#!/bin/bash

# 接单管理 API 测试脚本
# 测试接单相关的所有API接口

BASE_URL="https://aneworders.com/api"
TOKEN="your_auth_token_here"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试计数器
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 测试结果记录
declare -a TEST_RESULTS

# 打印测试结果
print_result() {
    local test_name="$1"
    local status="$2"
    local message="$3"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if [ "$status" = "PASS" ]; then
        echo -e "${GREEN}✓ PASS${NC} - $test_name"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        TEST_RESULTS+=("PASS: $test_name")
    else
        echo -e "${RED}✗ FAIL${NC} - $test_name: $message"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        TEST_RESULTS+=("FAIL: $test_name - $message")
    fi
}

# 打印测试标题
print_title() {
    echo -e "\n${YELLOW}=== $1 ===${NC}"
}

# 发送HTTP请求并检查响应
send_request() {
    local method="$1"
    local url="$2"
    local data="$3"
    local expected_status="$4"
    local test_name="$5"
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$url" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $TOKEN" \
            -d "$data")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$url" \
            -H "Authorization: Bearer $TOKEN")
    fi
    
    # 分离响应体和状态码
    http_code=$(echo "$response" | tail -n1)
    response_body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" = "$expected_status" ]; then
        print_result "$test_name" "PASS" ""
        echo "Response: $response_body"
    else
        print_result "$test_name" "FAIL" "Expected $expected_status, got $http_code. Response: $response_body"
    fi
}

# 测试数据
ORDER_ID=1
FACTORY_ID="test_factory_001"
JIEDAN_ID=""

print_title "接单管理 API 测试"

echo "使用测试数据:"
echo "  - 订单ID: $ORDER_ID"
echo "  - 工厂ID: $FACTORY_ID"
echo "  - Token: $TOKEN"
echo ""

# 1. 创建接单记录
print_title "1. 创建接单记录测试"

send_request "POST" "$BASE_URL/jiedan" \
    "{\"order_id\": $ORDER_ID, \"factory_id\": \"$FACTORY_ID\"}" \
    "201" \
    "创建接单记录"

# 从响应中提取接单ID（如果创建成功）
if [ "$http_code" = "201" ]; then
    JIEDAN_ID=$(echo "$response_body" | grep -o '"id":[0-9]*' | cut -d':' -f2)
    echo "创建的接单ID: $JIEDAN_ID"
fi

# 2. 获取接单记录详情
print_title "2. 获取接单记录详情测试"

if [ -n "$JIEDAN_ID" ]; then
    send_request "GET" "$BASE_URL/jiedan/$JIEDAN_ID" \
        "" \
        "200" \
        "获取接单记录详情"
else
    print_result "获取接单记录详情" "FAIL" "没有可用的接单ID"
fi

# 3. 获取订单的接单记录列表
print_title "3. 获取订单接单记录列表测试"

send_request "GET" "$BASE_URL/orders/$ORDER_ID/jiedans" \
    "" \
    "200" \
    "获取订单的接单记录列表"

# 4. 获取工厂的接单记录列表
print_title "4. 获取工厂接单记录列表测试"

send_request "GET" "$BASE_URL/factories/$FACTORY_ID/jiedans?page=1&pageSize=10" \
    "" \
    "200" \
    "获取工厂的接单记录列表"

# 5. 同意接单
print_title "5. 同意接单测试"

if [ -n "$JIEDAN_ID" ]; then
    send_request "POST" "$BASE_URL/jiedan/$JIEDAN_ID/accept" \
        "{\"agree_user_id\": \"test_user_001\"}" \
        "200" \
        "同意接单"
else
    print_result "同意接单" "FAIL" "没有可用的接单ID"
fi

# 6. 拒绝接单（使用新的接单记录）
print_title "6. 拒绝接单测试"

# 先创建另一个接单记录用于拒绝测试
send_request "POST" "$BASE_URL/jiedan" \
    "{\"order_id\": $ORDER_ID, \"factory_id\": \"test_factory_002\"}" \
    "201" \
    "创建第二个接单记录用于拒绝测试"

# 从响应中提取接单ID
REJECT_JIEDAN_ID=$(echo "$response_body" | grep -o '"id":[0-9]*' | cut -d':' -f2)

if [ -n "$REJECT_JIEDAN_ID" ]; then
    send_request "POST" "$BASE_URL/jiedan/$REJECT_JIEDAN_ID/reject" \
        "{\"reason\": \"产能不足，无法承接此订单\"}" \
        "200" \
        "拒绝接单"
else
    print_result "拒绝接单" "FAIL" "没有可用的接单ID"
fi

# 7. 更新接单记录
print_title "7. 更新接单记录测试"

if [ -n "$JIEDAN_ID" ]; then
    send_request "PUT" "$BASE_URL/jiedan/$JIEDAN_ID" \
        "{\"status\": \"accepted\", \"agree_user_id\": \"updated_user_001\"}" \
        "200" \
        "更新接单记录"
else
    print_result "更新接单记录" "FAIL" "没有可用的接单ID"
fi

# 8. 获取接单统计信息
print_title "8. 获取接单统计信息测试"

send_request "GET" "$BASE_URL/factories/$FACTORY_ID/jiedan-statistics" \
    "" \
    "200" \
    "获取工厂接单统计信息"

# 9. 错误情况测试
print_title "9. 错误情况测试"

# 9.1 创建重复接单（应该失败）
send_request "POST" "$BASE_URL/jiedan" \
    "{\"order_id\": $ORDER_ID, \"factory_id\": \"$FACTORY_ID\"}" \
    "409" \
    "创建重复接单（应该失败）"

# 9.2 获取不存在的接单记录
send_request "GET" "$BASE_URL/jiedan/99999" \
    "" \
    "404" \
    "获取不存在的接单记录"

# 9.3 使用无效的订单ID
send_request "POST" "$BASE_URL/jiedan" \
    "{\"order_id\": 99999, \"factory_id\": \"$FACTORY_ID\"}" \
    "500" \
    "使用不存在的订单ID创建接单"

# 9.4 缺少必要参数
send_request "POST" "$BASE_URL/jiedan" \
    "{\"order_id\": $ORDER_ID}" \
    "400" \
    "缺少工厂ID参数"

# 10. 删除接单记录
print_title "10. 删除接单记录测试"

if [ -n "$JIEDAN_ID" ]; then
    send_request "DELETE" "$BASE_URL/jiedan/$JIEDAN_ID" \
        "" \
        "200" \
        "删除接单记录"
else
    print_result "删除接单记录" "FAIL" "没有可用的接单ID"
fi

# 打印测试总结
print_title "测试总结"

echo "总测试数: $TOTAL_TESTS"
echo -e "通过: ${GREEN}$PASSED_TESTS${NC}"
echo -e "失败: ${RED}$FAILED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "\n${GREEN}所有测试通过！${NC}"
    exit 0
else
    echo -e "\n${RED}有 $FAILED_TESTS 个测试失败${NC}"
    echo -e "\n失败的测试详情:"
    for result in "${TEST_RESULTS[@]}"; do
        if [[ $result == FAIL* ]]; then
            echo -e "${RED}$result${NC}"
        fi
    done
    exit 1
fi 