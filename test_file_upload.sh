#!/bin/bash

# 设置颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}开始测试文件上传和访问功能...${NC}"

# 测试文件上传
echo -e "${GREEN}1. 测试文件上传...${NC}"
UPLOAD_RESPONSE=$(curl -s -X POST \
  -F "file=@test.txt" \
  -F "orderId=1" \
  http://localhost:8008/api/files/upload)

echo "上传响应: $UPLOAD_RESPONSE"

# 提取文件ID
FILE_ID=$(echo $UPLOAD_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo -e "${GREEN}文件ID: $FILE_ID${NC}"

if [ -z "$FILE_ID" ]; then
    echo -e "${RED}错误: 无法获取文件ID${NC}"
    exit 1
fi

# 测试获取文件详情
echo -e "${GREEN}2. 测试获取文件详情...${NC}"
DETAIL_RESPONSE=$(curl -s -X GET \
  http://localhost:8008/api/files/$FILE_ID)

echo "文件详情响应: $DETAIL_RESPONSE"

# 测试文件下载
echo -e "${GREEN}3. 测试文件下载...${NC}"
DOWNLOAD_RESPONSE=$(curl -s -I \
  http://localhost:8008/api/files/download/$FILE_ID)

echo "下载响应头: $DOWNLOAD_RESPONSE"

# 测试CORS头
echo -e "${GREEN}4. 测试CORS头...${NC}"
CORS_RESPONSE=$(curl -s -I \
  -H "Origin: https://aneworders.com" \
  http://localhost:8008/api/files/download/$FILE_ID)

echo "CORS响应头: $CORS_RESPONSE"

# 检查是否包含CORS头
if echo "$CORS_RESPONSE" | grep -q "Access-Control-Allow-Origin"; then
    echo -e "${GREEN}✓ CORS头设置正确${NC}"
else
    echo -e "${RED}✗ CORS头设置有问题${NC}"
fi

echo -e "${GREEN}文件上传和访问功能测试完成！${NC}" 