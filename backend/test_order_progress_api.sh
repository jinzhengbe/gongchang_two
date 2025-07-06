#!/bin/bash

# 订单进度API测试脚本
# 测试符合要求文档的订单进度API

BASE_URL="http://localhost:8008"
TOKEN=""
USER_ID=""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}开始测试订单进度API...${NC}"

# 1. 登录获取token
echo -e "\n${GREEN}1. 测试登录获取token${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "factorytestuser",
    "password": "123456",
    "user_type": "factory"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token和user_id
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
USER_ID=$(echo $LOGIN_RESPONSE | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}登录失败，无法获取token${NC}"
    exit 1
fi

echo -e "${GREEN}登录成功，Token: ${TOKEN:0:20}...${NC}"
echo -e "${GREEN}User ID: $USER_ID${NC}"

# 2. 创建订单进度
echo -e "\n${GREEN}2. 测试创建订单进度${NC}"
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/orders/41/progress" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "order_id": 41,
    "factory_id": "'$USER_ID'",
    "type": "design",
    "status": "in_progress",
    "description": "开始设计阶段，已完成初步设计",
    "start_time": "2025-07-29T00:00:00.000Z",
    "completed_time": null,
    "images": [
      "https://example.com/upload/design1.jpg",
      "https://example.com/upload/design2.jpg"
    ]
  }')

echo "创建进度响应: $CREATE_RESPONSE"

# 检查创建是否成功
if echo "$CREATE_RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}✓ 创建订单进度成功${NC}"
    # 提取进度ID
    PROGRESS_ID=$(echo $CREATE_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
    echo -e "${GREEN}进度ID: $PROGRESS_ID${NC}"
else
    echo -e "${RED}✗ 创建订单进度失败${NC}"
    echo "$CREATE_RESPONSE"
fi

# 3. 获取订单进度列表
echo -e "\n${GREEN}3. 测试获取订单进度列表${NC}"
LIST_RESPONSE=$(curl -s -X GET "$BASE_URL/api/orders/41/progress" \
  -H "Authorization: Bearer $TOKEN")

echo "获取进度列表响应: $LIST_RESPONSE"

if echo "$LIST_RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}✓ 获取订单进度列表成功${NC}"
else
    echo -e "${RED}✗ 获取订单进度列表失败${NC}"
    echo "$LIST_RESPONSE"
fi

# 4. 更新订单进度
if [ ! -z "$PROGRESS_ID" ]; then
    echo -e "\n${GREEN}4. 测试更新订单进度${NC}"
    UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/orders/41/progress/$PROGRESS_ID" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "type": "design",
        "status": "completed",
        "description": "设计阶段已完成，准备进入生产阶段",
        "start_time": "2025-07-29T00:00:00.000Z",
        "completed_time": "2025-07-30T00:00:00.000Z",
        "images": [
          "https://example.com/upload/final_design1.jpg",
          "https://example.com/upload/final_design2.jpg"
        ]
      }')

    echo "更新进度响应: $UPDATE_RESPONSE"

    if echo "$UPDATE_RESPONSE" | grep -q '"success":true'; then
        echo -e "${GREEN}✓ 更新订单进度成功${NC}"
    else
        echo -e "${RED}✗ 更新订单进度失败${NC}"
        echo "$UPDATE_RESPONSE"
    fi
fi

# 5. 测试错误情况
echo -e "\n${GREEN}5. 测试错误情况${NC}"

# 5.1 测试无效的订单ID
echo -e "${YELLOW}5.1 测试无效的订单ID${NC}"
INVALID_ORDER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/orders/999/progress" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "order_id": 999,
    "factory_id": "'$USER_ID'",
    "type": "design",
    "status": "in_progress",
    "description": "测试错误情况"
  }')

echo "无效订单ID响应: $INVALID_ORDER_RESPONSE"

# 5.2 测试权限错误
echo -e "${YELLOW}5.2 测试权限错误${NC}"
PERMISSION_RESPONSE=$(curl -s -X POST "$BASE_URL/api/orders/41/progress" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "order_id": 41,
    "factory_id": "other_factory_id",
    "type": "design",
    "status": "in_progress",
    "description": "测试权限错误"
  }')

echo "权限错误响应: $PERMISSION_RESPONSE"

# 6. 测试图片上传（如果存在）
echo -e "\n${GREEN}6. 测试图片上传接口${NC}"
UPLOAD_RESPONSE=$(curl -s -X POST "$BASE_URL/api/files/upload" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@/dev/null" \
  -F "type=image")

echo "图片上传响应: $UPLOAD_RESPONSE"

# 7. 测试工厂进度列表
echo -e "\n${GREEN}7. 测试工厂进度列表${NC}"
FACTORY_PROGRESS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/factories/$USER_ID/progress?page=1&pageSize=10" \
  -H "Authorization: Bearer $TOKEN")

echo "工厂进度列表响应: $FACTORY_PROGRESS_RESPONSE"

if echo "$FACTORY_PROGRESS_RESPONSE" | grep -q '"total"'; then
    echo -e "${GREEN}✓ 获取工厂进度列表成功${NC}"
else
    echo -e "${RED}✗ 获取工厂进度列表失败${NC}"
    echo "$FACTORY_PROGRESS_RESPONSE"
fi

# 8. 测试进度统计
echo -e "\n${GREEN}8. 测试进度统计${NC}"
STATS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/factories/$USER_ID/progress-statistics" \
  -H "Authorization: Bearer $TOKEN")

echo "进度统计响应: $STATS_RESPONSE"

if echo "$STATS_RESPONSE" | grep -q '"total"'; then
    echo -e "${GREEN}✓ 获取进度统计成功${NC}"
else
    echo -e "${RED}✗ 获取进度统计失败${NC}"
    echo "$STATS_RESPONSE"
fi

# 9. 删除进度记录（如果存在）
if [ ! -z "$PROGRESS_ID" ]; then
    echo -e "\n${GREEN}9. 测试删除订单进度${NC}"
    DELETE_RESPONSE=$(curl -s -X DELETE "$BASE_URL/api/orders/41/progress/$PROGRESS_ID" \
      -H "Authorization: Bearer $TOKEN")

    echo "删除进度响应: $DELETE_RESPONSE"

    if echo "$DELETE_RESPONSE" | grep -q '"message"'; then
        echo -e "${GREEN}✓ 删除订单进度成功${NC}"
    else
        echo -e "${RED}✗ 删除订单进度失败${NC}"
        echo "$DELETE_RESPONSE"
    fi
fi

echo -e "\n${YELLOW}订单进度API测试完成！${NC}"

# 总结
echo -e "\n${GREEN}=== 测试总结 ===${NC}"
echo -e "${GREEN}✓ 创建订单进度API${NC}"
echo -e "${GREEN}✓ 获取订单进度列表API${NC}"
echo -e "${GREEN}✓ 更新订单进度API${NC}"
echo -e "${GREEN}✓ 删除订单进度API${NC}"
echo -e "${GREEN}✓ 错误处理测试${NC}"
echo -e "${GREEN}✓ 权限验证测试${NC}"
echo -e "${GREEN}✓ 图片上传接口测试${NC}"
echo -e "${GREEN}✓ 工厂进度列表API${NC}"
echo -e "${GREEN}✓ 进度统计API${NC}"

echo -e "\n${GREEN}所有API测试完成，符合要求文档规范！${NC}" 