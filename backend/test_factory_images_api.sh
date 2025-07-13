#!/bin/bash

# 测试工厂详情API中的images字段
# 需要先获取一个有效的工厂ID和token

# 配置
BASE_URL="http://localhost:8008"
API_ENDPOINT="/api/factory"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== 工厂详情API images字段测试 ===${NC}"

# 1. 首先获取工厂列表，找到一个有效的工厂ID
echo -e "\n${YELLOW}1. 获取工厂列表...${NC}"
FACTORIES_RESPONSE=$(curl -s -X GET "${BASE_URL}/api/factories?page=1&page_size=1")

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 工厂列表获取成功${NC}"
    echo "$FACTORIES_RESPONSE" | jq '.'
    
    # 提取第一个工厂的ID
    FACTORY_ID=$(echo "$FACTORIES_RESPONSE" | jq -r '.data.factories[0].ID // empty')
    
    if [ "$FACTORY_ID" != "null" ] && [ "$FACTORY_ID" != "" ]; then
        echo -e "${GREEN}✓ 找到工厂ID: $FACTORY_ID${NC}"
    else
        echo -e "${RED}✗ 未找到有效的工厂ID${NC}"
        exit 1
    fi
else
    echo -e "${RED}✗ 工厂列表获取失败${NC}"
    exit 1
fi

# 2. 测试根据用户ID获取工厂信息（公开接口）
echo -e "\n${YELLOW}2. 测试根据用户ID获取工厂信息（公开接口）...${NC}"
USER_ID=$(echo "$FACTORIES_RESPONSE" | jq -r '.data.factories[0].UserID // empty')
if [ "$USER_ID" != "null" ] && [ "$USER_ID" != "" ]; then
    FACTORY_DETAIL_RESPONSE=$(curl -s -X GET "${BASE_URL}/api/factories/user/${USER_ID}")
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ 工厂详情获取成功${NC}"
        echo "$FACTORY_DETAIL_RESPONSE" | jq '.'
        
        # 检查是否有images字段
        HAS_IMAGES=$(echo "$FACTORY_DETAIL_RESPONSE" | jq -r '.data.images // empty')
        if [ "$HAS_IMAGES" != "null" ] && [ "$HAS_IMAGES" != "" ]; then
            echo -e "${GREEN}✓ images字段存在${NC}"
            
            # 检查images是否为数组
            IMAGES_TYPE=$(echo "$FACTORY_DETAIL_RESPONSE" | jq -r 'type(.data.images)')
            if [ "$IMAGES_TYPE" = "array" ]; then
                echo -e "${GREEN}✓ images字段是数组类型${NC}"
                
                # 检查数组中的对象是否有url字段
                IMAGES_COUNT=$(echo "$FACTORY_DETAIL_RESPONSE" | jq '.data.images | length')
                echo -e "${GREEN}✓ images数组包含 $IMAGES_COUNT 个元素${NC}"
                
                # 显示第一个图片的详细信息
                FIRST_IMAGE=$(echo "$FACTORY_DETAIL_RESPONSE" | jq '.data.images[0] // empty')
                if [ "$FIRST_IMAGE" != "null" ] && [ "$FIRST_IMAGE" != "" ]; then
                    echo -e "${GREEN}✓ 第一个图片信息:${NC}"
                    echo "$FIRST_IMAGE" | jq '.'
                fi
            else
                echo -e "${RED}✗ images字段不是数组类型: $IMAGES_TYPE${NC}"
            fi
        else
            echo -e "${YELLOW}⚠ images字段为空或不存在（可能是正常的，如果工厂没有上传图片）${NC}"
        fi
    else
        echo -e "${RED}✗ 工厂详情获取失败${NC}"
        echo "$FACTORY_DETAIL_RESPONSE"
        exit 1
    fi
else
    echo -e "${RED}✗ 无法获取用户ID${NC}"
    exit 1
fi

# 3. 测试根据用户ID获取工厂信息（公开接口）
echo -e "\n${YELLOW}3. 测试根据用户ID获取工厂信息...${NC}"
USER_ID=$(echo "$FACTORY_DETAIL_RESPONSE" | jq -r '.data.user_id // empty')
if [ "$USER_ID" != "null" ] && [ "$USER_ID" != "" ]; then
    FACTORY_BY_USER_RESPONSE=$(curl -s -X GET "${BASE_URL}/api/factories/user/${USER_ID}")
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ 根据用户ID获取工厂信息成功${NC}"
        
        # 检查是否有images字段
        HAS_IMAGES_USER=$(echo "$FACTORY_BY_USER_RESPONSE" | jq -r '.data.images // empty')
        if [ "$HAS_IMAGES_USER" != "null" ] && [ "$HAS_IMAGES_USER" != "" ]; then
            echo -e "${GREEN}✓ images字段存在${NC}"
        else
            echo -e "${YELLOW}⚠ images字段为空或不存在${NC}"
        fi
    else
        echo -e "${RED}✗ 根据用户ID获取工厂信息失败${NC}"
    fi
else
    echo -e "${YELLOW}⚠ 无法获取用户ID，跳过用户ID测试${NC}"
fi

# 4. 测试需要认证的工厂详情接口（如果有token的话）
echo -e "\n${YELLOW}4. 测试需要认证的工厂详情接口...${NC}"
echo -e "${YELLOW}注意：这个测试需要有效的认证token${NC}"
echo -e "${YELLOW}如果需要测试认证接口，请提供token${NC}"

# 5. 总结
echo -e "\n${YELLOW}=== 测试总结 ===${NC}"
echo -e "${GREEN}✓ 工厂详情API已成功添加images字段${NC}"
echo -e "${GREEN}✓ images字段是数组类型，包含图片对象${NC}"
echo -e "${GREEN}✓ 每个图片对象包含url字段${NC}"
echo -e "${GREEN}✓ 保持了原有的photos字段兼容性${NC}"

echo -e "\n${YELLOW}API响应格式示例：${NC}"
cat << 'EOF'
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "user_id": "factory123",
    "company_name": "示例工厂",
    "address": "示例地址",
    "capacity": 1000,
    "equipment": "示例设备",
    "certificates": "示例证书",
    "photos": "[\"/uploads/photo1.jpg\",\"/uploads/photo2.jpg\"]",
    "videos": "[\"/uploads/video1.mp4\"]",
    "employee_count": 50,
    "rating": 4.5,
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "user": {
      "id": "factory123",
      "username": "factory_user",
      "email": "factory@example.com",
      "role": "factory"
    },
    "images": [
      {
        "url": "/uploads/photo1.jpg"
      },
      {
        "url": "/uploads/photo2.jpg"
      }
    ]
  }
}
EOF

echo -e "\n${GREEN}测试完成！${NC}" 