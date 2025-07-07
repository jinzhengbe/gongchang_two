#!/bin/bash

# 工厂订单详情页面API测试脚本
# 测试 checkAcceptOrderStatus() 和 fetchOrderProgress() 对应的后端API

# 设置颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 配置
BASE_URL="http://localhost:8008"
TOKEN=""
ORDER_ID="39"

echo -e "${YELLOW}开始测试工厂订单详情页面API...${NC}"

# 1. 登录获取token
echo -e "${GREEN}1. 登录获取token...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "sdf",
    "password": "123456"
  }')

if [ $? -eq 0 ]; then
    TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    if [ -n "$TOKEN" ]; then
        echo -e "${GREEN}✓ 登录成功，获取到token${NC}"
    else
        echo -e "${RED}✗ 登录失败，无法获取token${NC}"
        echo "响应: $LOGIN_RESPONSE"
        exit 1
    fi
else
    echo -e "${RED}✗ 登录请求失败${NC}"
    exit 1
fi

# 2. 测试获取订单接单记录 (checkAcceptOrderStatus)
echo -e "${GREEN}2. 测试获取订单接单记录...${NC}"
JIEDAN_RESPONSE=$(curl -s -X GET "$BASE_URL/api/orders/$ORDER_ID/jiedans" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Accept: application/json")

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 获取接单记录成功${NC}"
    echo "响应: $JIEDAN_RESPONSE"
else
    echo -e "${RED}✗ 获取接单记录失败${NC}"
fi

# 3. 测试获取订单进度记录 (fetchOrderProgress)
echo -e "${GREEN}3. 测试获取订单进度记录...${NC}"
PROGRESS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/orders/$ORDER_ID/progress" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Accept: application/json")

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 获取进度记录成功${NC}"
    echo "响应: $PROGRESS_RESPONSE"
else
    echo -e "${RED}✗ 获取进度记录失败${NC}"
fi

# 4. 测试创建接单记录（需要工厂权限，这里会失败）
echo -e "${GREEN}4. 测试创建接单记录（预期失败，因为用户是设计师）...${NC}"
CREATE_JIEDAN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/jiedan" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"order_id\": $ORDER_ID,
    \"factory_id\": \"sdf\",
    \"price\": 5000.0
  }")

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 创建接单记录请求成功${NC}"
    echo "响应: $CREATE_JIEDAN_RESPONSE"
else
    echo -e "${YELLOW}⚠ 创建接单记录失败（符合预期，因为用户是设计师）${NC}"
fi

# 5. 测试工厂接受订单（需要工厂权限，这里会失败）
echo -e "${GREEN}5. 测试工厂接受订单（预期失败，因为用户是设计师）...${NC}"
ACCEPT_ORDER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/orders/$ORDER_ID/accept" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"order_id\": $ORDER_ID,
    \"factory_id\": \"sdf\",
    \"status\": \"accepted\",
    \"accepted_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
    \"action\": \"accept_order\",
    \"price_quote\": 5000.0,
    \"message\": \"我们接受这个订单\"
  }")

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 工厂接受订单请求成功${NC}"
    echo "响应: $ACCEPT_ORDER_RESPONSE"
else
    echo -e "${YELLOW}⚠ 工厂接受订单失败（符合预期，因为用户是设计师）${NC}"
fi

# 6. 测试创建进度记录（需要工厂权限，这里会失败）
echo -e "${GREEN}6. 测试创建进度记录（预期失败，因为用户是设计师）...${NC}"
CREATE_PROGRESS_RESPONSE=$(curl -s -X POST "$BASE_URL/api/orders/$ORDER_ID/progress" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"order_id\": $ORDER_ID,
    \"factory_id\": \"sdf\",
    \"type\": \"production\",
    \"status\": \"in_progress\",
    \"description\": \"开始生产阶段\",
    \"start_time\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
    \"completed_time\": null,
    \"images\": []
  }")

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 创建进度记录请求成功${NC}"
    echo "响应: $CREATE_PROGRESS_RESPONSE"
else
    echo -e "${YELLOW}⚠ 创建进度记录失败（符合预期，因为用户是设计师）${NC}"
fi

echo -e "${GREEN}测试完成！${NC}"
echo -e "${YELLOW}注意：部分测试失败是预期的，因为测试用户是设计师角色，而某些API需要工厂权限${NC}" 