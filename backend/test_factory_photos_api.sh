#!/bin/bash

# 工厂图片API测试脚本
echo "=== 工厂图片API测试 ==="

# 配置
BASE_URL="http://localhost:8008"
TOKEN=""
FACTORY_ID=""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. 登录获取token
echo -e "${YELLOW}1. 登录获取token...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "gongchang",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}❌ 登录失败，无法获取token${NC}"
    exit 1
fi

echo -e "${GREEN}✅ 登录成功，获取到token: ${TOKEN:0:20}...${NC}"

# 2. 获取工厂ID
echo -e "\n${YELLOW}2. 获取工厂ID...${NC}"
FACTORY_RESPONSE=$(curl -s -X GET "$BASE_URL/api/factories/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "工厂信息响应: $FACTORY_RESPONSE"

# 提取工厂ID（假设是用户ID）
FACTORY_ID=$(echo $LOGIN_RESPONSE | grep -o '"user_id":"[^"]*"' | cut -d'"' -f4)
if [ -z "$FACTORY_ID" ]; then
    FACTORY_ID="gongchang" # 使用用户名作为工厂ID
fi

echo -e "${GREEN}✅ 工厂ID: $FACTORY_ID${NC}"

# 3. 测试批量上传图片
echo -e "\n${YELLOW}3. 测试批量上传图片...${NC}"

# 创建测试图片文件
echo "创建测试图片文件..."
echo "test image content" > test_image1.jpg
echo "test image content 2" > test_image2.jpg

UPLOAD_RESPONSE=$(curl -s -X POST "$BASE_URL/api/factories/$FACTORY_ID/photos/batch" \
  -H "Authorization: Bearer $TOKEN" \
  -F "files=@test_image1.jpg" \
  -F "files=@test_image2.jpg" \
  -F "category=workshop")

echo "批量上传响应: $UPLOAD_RESPONSE"

if echo "$UPLOAD_RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}✅ 批量上传成功${NC}"
else
    echo -e "${RED}❌ 批量上传失败${NC}"
    echo "$UPLOAD_RESPONSE"
fi

# 4. 测试获取图片列表
echo -e "\n${YELLOW}4. 测试获取图片列表...${NC}"
LIST_RESPONSE=$(curl -s -X GET "$BASE_URL/api/factories/$FACTORY_ID/photos?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "获取图片列表响应: $LIST_RESPONSE"

if echo "$LIST_RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}✅ 获取图片列表成功${NC}"
    
    # 提取第一个图片ID用于删除测试
    FIRST_PHOTO_ID=$(echo $LIST_RESPONSE | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
    echo "第一个图片ID: $FIRST_PHOTO_ID"
else
    echo -e "${RED}❌ 获取图片列表失败${NC}"
    echo "$LIST_RESPONSE"
fi

# 5. 测试删除单张图片
if [ ! -z "$FIRST_PHOTO_ID" ]; then
    echo -e "\n${YELLOW}5. 测试删除单张图片...${NC}"
    DELETE_RESPONSE=$(curl -s -X DELETE "$BASE_URL/api/factories/$FACTORY_ID/photos/$FIRST_PHOTO_ID" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json")

    echo "删除图片响应: $DELETE_RESPONSE"

    if echo "$DELETE_RESPONSE" | grep -q '"success":true'; then
        echo -e "${GREEN}✅ 删除单张图片成功${NC}"
    else
        echo -e "${RED}❌ 删除单张图片失败${NC}"
        echo "$DELETE_RESPONSE"
    fi
else
    echo -e "${YELLOW}⚠️  跳过删除单张图片测试（没有图片ID）${NC}"
fi

# 6. 测试批量删除图片
echo -e "\n${YELLOW}6. 测试批量删除图片...${NC}"

# 先获取剩余的图片ID
REMAINING_PHOTOS=$(curl -s -X GET "$BASE_URL/api/factories/$FACTORY_ID/photos?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

# 提取所有图片ID
PHOTO_IDS=$(echo $REMAINING_PHOTOS | grep -o '"id":"[^"]*"' | cut -d'"' -f4 | tr '\n' ',' | sed 's/,$//')

if [ ! -z "$PHOTO_IDS" ]; then
    BATCH_DELETE_RESPONSE=$(curl -s -X DELETE "$BASE_URL/api/factories/$FACTORY_ID/photos/batch" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "{\"photo_ids\":[\"$PHOTO_IDS\"]}")

    echo "批量删除响应: $BATCH_DELETE_RESPONSE"

    if echo "$BATCH_DELETE_RESPONSE" | grep -q '"success":true'; then
        echo -e "${GREEN}✅ 批量删除图片成功${NC}"
    else
        echo -e "${RED}❌ 批量删除图片失败${NC}"
        echo "$BATCH_DELETE_RESPONSE"
    fi
else
    echo -e "${YELLOW}⚠️  跳过批量删除测试（没有图片）${NC}"
fi

# 7. 测试错误情况
echo -e "\n${YELLOW}7. 测试错误情况...${NC}"

# 7.1 测试无权限
echo -e "${YELLOW}7.1 测试无权限...${NC}"
NO_PERMISSION_RESPONSE=$(curl -s -X POST "$BASE_URL/api/factories/other_factory/photos/batch" \
  -H "Authorization: Bearer $TOKEN" \
  -F "files=@test_image1.jpg")

echo "无权限响应: $NO_PERMISSION_RESPONSE"

# 7.2 测试无文件上传
echo -e "${YELLOW}7.2 测试无文件上传...${NC}"
NO_FILES_RESPONSE=$(curl -s -X POST "$BASE_URL/api/factories/$FACTORY_ID/photos/batch" \
  -H "Authorization: Bearer $TOKEN")

echo "无文件响应: $NO_FILES_RESPONSE"

# 清理测试文件
echo -e "\n${YELLOW}清理测试文件...${NC}"
rm -f test_image1.jpg test_image2.jpg

echo -e "\n${GREEN}=== 测试完成 ===${NC}" 