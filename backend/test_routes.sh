#!/bin/bash

# 测试工厂图片批量上传和管理API
BASE_URL="http://localhost:8008"

echo "=== 测试工厂图片批量上传和管理API ==="

# 1. 登录获取token和factory_id
echo "1. 登录获取token..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "gongchang",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token和factory_id
TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
# 从profile中提取工厂ID，使用UserID作为factory_id
FACTORY_ID=$(echo "$LOGIN_RESPONSE" | grep -o '"UserID":"[^"]*"' | cut -d'"' -f4)

echo "Token: $TOKEN"
echo "Factory ID: $FACTORY_ID"

if [ -z "$TOKEN" ] || [ -z "$FACTORY_ID" ]; then
    echo "登录失败或无法获取token/factory_id"
    exit 1
fi

echo ""
echo "2. 测试批量上传图片..."
# 创建测试图片文件
echo "创建测试图片文件..."
mkdir -p test_images
for i in {1..3}; do
    echo "这是测试图片 $i 的内容" > "test_images/test_image_$i.jpg"
done

# 批量上传图片
UPLOAD_RESPONSE=$(curl -s -X POST "$BASE_URL/api/factories/$FACTORY_ID/photos/batch" \
  -H "Authorization: Bearer $TOKEN" \
  -F "files=@test_images/test_image_1.jpg" \
  -F "files=@test_images/test_image_2.jpg" \
  -F "files=@test_images/test_image_3.jpg")

echo "批量上传响应: $UPLOAD_RESPONSE"

echo ""
echo "3. 测试获取图片列表..."
LIST_RESPONSE=$(curl -s -X GET "$BASE_URL/api/factories/$FACTORY_ID/photos" \
  -H "Authorization: Bearer $TOKEN")

echo "图片列表响应: $LIST_RESPONSE"

# 提取第一个图片ID用于删除测试
FIRST_IMAGE_ID=$(echo "$LIST_RESPONSE" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ ! -z "$FIRST_IMAGE_ID" ]; then
    echo ""
    echo "4. 测试删除单个图片 (ID: $FIRST_IMAGE_ID)..."
    DELETE_RESPONSE=$(curl -s -X DELETE "$BASE_URL/api/factories/$FACTORY_ID/photos/$FIRST_IMAGE_ID" \
      -H "Authorization: Bearer $TOKEN")
    
    echo "删除单个图片响应: $DELETE_RESPONSE"
    
    echo ""
    echo "5. 测试批量删除图片..."
    # 获取剩余图片ID
    REMAINING_IMAGES=$(curl -s -X GET "$BASE_URL/api/factories/$FACTORY_ID/photos" \
      -H "Authorization: Bearer $TOKEN")
    
    IMAGE_IDS=$(echo "$REMAINING_IMAGES" | grep -o '"id":"[^"]*"' | cut -d'"' -f4 | tr '\n' ',' | sed 's/,$//')
    
    if [ ! -z "$IMAGE_IDS" ]; then
        BATCH_DELETE_RESPONSE=$(curl -s -X DELETE "$BASE_URL/api/factories/$FACTORY_ID/photos/batch" \
          -H "Authorization: Bearer $TOKEN" \
          -H "Content-Type: application/json" \
          -d "{\"photo_ids\":[\"$IMAGE_IDS\"]}")
        
        echo "批量删除响应: $BATCH_DELETE_RESPONSE"
    else
        echo "没有剩余图片可删除"
    fi
else
    echo "没有找到图片ID，跳过删除测试"
fi

echo ""
echo "6. 清理测试文件..."
rm -rf test_images

echo "=== 测试完成 ===" 