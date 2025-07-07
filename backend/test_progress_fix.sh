#!/bin/bash

# 测试进度API修复脚本
# 验证工厂接单后可以创建进度记录

# 设置颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 配置
BASE_URL="http://localhost:8008"
TOKEN=""
ORDER_ID="39"
FACTORY_ID="3af8e32a-e267-45f1-8959-faf3f0787bfa"  # 订单39的接单工厂ID

echo -e "${YELLOW}开始测试进度API修复...${NC}"

# 1. 登录获取token（使用工厂用户）
echo -e "${GREEN}1. 登录获取token...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"factory_user\",
    \"password\": \"123456\"
  }")

if [ $? -eq 0 ]; then
    TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    if [ -n "$TOKEN" ]; then
        echo -e "${GREEN}✓ 登录成功，获取到token${NC}"
    else
        echo -e "${RED}✗ 登录失败，无法获取token${NC}"
        echo "响应: $LOGIN_RESPONSE"
        echo -e "${YELLOW}尝试使用其他用户...${NC}"
        
        # 尝试使用其他用户
        LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
          -H "Content-Type: application/json" \
          -d '{
            "username": "sdf",
            "password": "123456"
          }')
        
        TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
        if [ -n "$TOKEN" ]; then
            echo -e "${GREEN}✓ 使用备用用户登录成功${NC}"
        else
            echo -e "${RED}✗ 所有用户登录都失败${NC}"
            exit 1
        fi
    fi
else
    echo -e "${RED}✗ 登录请求失败${NC}"
    exit 1
fi

# 2. 检查订单39的接单记录
echo -e "${GREEN}2. 检查订单39的接单记录...${NC}"
JIEDAN_RESPONSE=$(curl -s -X GET "$BASE_URL/api/orders/$ORDER_ID/jiedans" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Accept: application/json")

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 获取接单记录成功${NC}"
    echo "响应: $JIEDAN_RESPONSE"
    
    # 检查是否有接单记录
    if echo "$JIEDAN_RESPONSE" | grep -q "factory_id.*$FACTORY_ID"; then
        echo -e "${GREEN}✓ 找到工厂 $FACTORY_ID 的接单记录${NC}"
    else
        echo -e "${YELLOW}⚠ 未找到工厂 $FACTORY_ID 的接单记录${NC}"
    fi
else
    echo -e "${RED}✗ 获取接单记录失败${NC}"
fi

# 3. 测试创建进度记录（使用正确的工厂ID）
echo -e "${GREEN}3. 测试创建进度记录（使用接单工厂ID）...${NC}"
CREATE_PROGRESS_RESPONSE=$(curl -s -X POST "$BASE_URL/api/orders/$ORDER_ID/progress" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"order_id\": $ORDER_ID,
    \"factory_id\": \"$FACTORY_ID\",
    \"type\": \"material\",
    \"status\": \"in_progress\",
    \"description\": \"材料采购进行中\",
    \"start_time\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
    \"completed_time\": null,
    \"images\": []
  }")

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 创建进度记录请求成功${NC}"
    echo "响应: $CREATE_PROGRESS_RESPONSE"
    
    # 检查是否成功创建
    if echo "$CREATE_PROGRESS_RESPONSE" | grep -q "success.*true"; then
        echo -e "${GREEN}✓ 进度记录创建成功！${NC}"
    else
        echo -e "${YELLOW}⚠ 进度记录创建失败，但请求成功发送${NC}"
    fi
else
    echo -e "${RED}✗ 创建进度记录请求失败${NC}"
fi

# 4. 再次获取进度记录，验证是否创建成功
echo -e "${GREEN}4. 获取订单进度记录...${NC}"
PROGRESS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/orders/$ORDER_ID/progress" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Accept: application/json")

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 获取进度记录成功${NC}"
    echo "响应: $PROGRESS_RESPONSE"
else
    echo -e "${RED}✗ 获取进度记录失败${NC}"
fi

echo -e "${GREEN}测试完成！${NC}" 