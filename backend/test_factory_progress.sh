#!/bin/bash

# 测试工厂用户创建进度记录

# 设置颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 配置
BASE_URL="http://localhost:8008"
ORDER_ID="39"
FACTORY_ID="3af8e32a-e267-45f1-8959-faf3f0787bfa"

echo -e "${YELLOW}开始测试工厂用户创建进度记录...${NC}"

# 1. 尝试使用工厂token（从用户提供的信息中提取）
echo -e "${GREEN}1. 使用工厂token测试...${NC}"

# 使用用户提供的工厂token
FACTORY_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiM2FmOGUzMmEtZTI2Ny00NWYxLTg5NTktZmFmM2YwNzg3YmZhIiwicm9sZSI6ImZhY3RvcnkiLCJleHAiOjE3NTE5Mjg5NDksImlhdCI6MTc1MTg0MjU0OX0.F9vdbHrRtuUIFr288adRjr27VwafLtigvH6RESHlKsM"

# 测试兼容路径 /api/orders/39/progresses
echo -e "${GREEN}2. 测试兼容路径 /api/orders/39/progresses...${NC}"
PROGRESS_RESPONSE=$(curl -s -X POST "$BASE_URL/api/orders/$ORDER_ID/progresses" \
  -H "Authorization: Bearer $FACTORY_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"order_id\": $ORDER_ID,
    \"factory_id\": \"$FACTORY_ID\",
    \"type\": \"production\",
    \"status\": \"in_progress\",
    \"description\": \"生产进度测试\",
    \"start_time\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
    \"completed_time\": null,
    \"images\": []
  }")

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 请求成功发送${NC}"
    echo "响应: $PROGRESS_RESPONSE"
    
    # 检查是否成功创建
    if echo "$PROGRESS_RESPONSE" | grep -q "success.*true"; then
        echo -e "${GREEN}✓ 进度记录创建成功！${NC}"
    elif echo "$PROGRESS_RESPONSE" | grep -q "该工厂未接此订单"; then
        echo -e "${YELLOW}⚠ 工厂未接此订单，需要先接单${NC}"
    elif echo "$PROGRESS_RESPONSE" | grep -q "只有已接单的工厂"; then
        echo -e "${YELLOW}⚠ 接单状态不是accepted，需要先同意接单${NC}"
    else
        echo -e "${YELLOW}⚠ 其他错误${NC}"
    fi
else
    echo -e "${RED}✗ 请求失败${NC}"
fi

# 3. 测试标准路径 /api/orders/39/progress
echo -e "${GREEN}3. 测试标准路径 /api/orders/39/progress...${NC}"
PROGRESS_RESPONSE_2=$(curl -s -X POST "$BASE_URL/api/orders/$ORDER_ID/progress" \
  -H "Authorization: Bearer $FACTORY_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"order_id\": $ORDER_ID,
    \"factory_id\": \"$FACTORY_ID\",
    \"type\": \"production\",
    \"status\": \"in_progress\",
    \"description\": \"生产进度测试2\",
    \"start_time\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
    \"completed_time\": null,
    \"images\": []
  }")

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 请求成功发送${NC}"
    echo "响应: $PROGRESS_RESPONSE_2"
else
    echo -e "${RED}✗ 请求失败${NC}"
fi

# 4. 获取进度记录列表
echo -e "${GREEN}4. 获取进度记录列表...${NC}"
GET_PROGRESS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/orders/$ORDER_ID/progresses" \
  -H "Authorization: Bearer $FACTORY_TOKEN" \
  -H "Accept: application/json")

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 获取进度记录成功${NC}"
    echo "响应: $GET_PROGRESS_RESPONSE"
else
    echo -e "${RED}✗ 获取进度记录失败${NC}"
fi

echo -e "${GREEN}测试完成！${NC}" 