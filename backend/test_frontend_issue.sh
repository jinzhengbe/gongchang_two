#!/bin/bash

# 测试前端图片显示问题的后端API
# 验证images字段是否正确返回

# 配置
BASE_URL="http://localhost:8008"
LOGIN_USERNAME="gongchang"
LOGIN_PASSWORD="123456"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== 前端图片显示问题 - 后端API测试 ===${NC}"

# 1. 登录获取token
echo -e "\n${YELLOW}1. 登录获取token...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$LOGIN_USERNAME\",
    \"password\": \"$LOGIN_PASSWORD\"
  }")

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 登录请求成功${NC}"
    echo "$LOGIN_RESPONSE" | jq '.'
    
    # 提取token
    TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token // empty')
    
    if [ "$TOKEN" != "null" ] && [ "$TOKEN" != "" ]; then
        echo -e "${GREEN}✓ 获取到token: ${TOKEN:0:20}...${NC}"
    else
        echo -e "${RED}✗ 未获取到token${NC}"
        exit 1
    fi
else
    echo -e "${RED}✗ 登录失败${NC}"
    exit 1
fi

# 2. 测试工厂详情API
echo -e "\n${YELLOW}2. 测试工厂详情API...${NC}"
FACTORY_PROFILE_RESPONSE=$(curl -s -X GET "${BASE_URL}/api/factories/profile" \
  -H "Authorization: Bearer $TOKEN")

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 工厂详情API请求成功${NC}"
    echo "$FACTORY_PROFILE_RESPONSE" | jq '.'
    
    # 检查images字段
    HAS_IMAGES=$(echo "$FACTORY_PROFILE_RESPONSE" | jq -r '.data.images // empty')
    if [ "$HAS_IMAGES" != "null" ] && [ "$HAS_IMAGES" != "" ]; then
        echo -e "${GREEN}✓ images字段存在${NC}"
        
        # 检查images是否为数组
        IMAGES_COUNT=$(echo "$FACTORY_PROFILE_RESPONSE" | jq '.data.images | length // 0')
        if [ "$IMAGES_COUNT" -gt 0 ]; then
            echo -e "${GREEN}✓ images字段是数组类型，包含 $IMAGES_COUNT 个元素${NC}"
            
            # 显示所有图片URL
                echo -e "${GREEN}✓ 图片URL列表:${NC}"
                echo "$FACTORY_PROFILE_RESPONSE" | jq -r '.data.images[].url'
                
                # 测试第一个图片URL是否可访问
                FIRST_IMAGE_URL=$(echo "$FACTORY_PROFILE_RESPONSE" | jq -r '.data.images[0].url // empty')
                if [ "$FIRST_IMAGE_URL" != "null" ] && [ "$FIRST_IMAGE_URL" != "" ]; then
                    echo -e "${YELLOW}测试第一个图片URL: $FIRST_IMAGE_URL${NC}"
                    
                    # 构建完整URL
                    FULL_IMAGE_URL="${BASE_URL}${FIRST_IMAGE_URL}"
                    echo -e "${YELLOW}完整URL: $FULL_IMAGE_URL${NC}"
                    
                    # 测试图片是否可访问
                    IMAGE_RESPONSE=$(curl -s -I "$FULL_IMAGE_URL")
                    if [ $? -eq 0 ]; then
                        HTTP_STATUS=$(echo "$IMAGE_RESPONSE" | head -n 1 | cut -d' ' -f2)
                        if [ "$HTTP_STATUS" = "200" ]; then
                            echo -e "${GREEN}✓ 图片可正常访问 (HTTP $HTTP_STATUS)${NC}"
                        else
                            echo -e "${RED}✗ 图片访问失败 (HTTP $HTTP_STATUS)${NC}"
                        fi
                    else
                        echo -e "${RED}✗ 图片URL无法访问${NC}"
                    fi
                fi
            else
                echo -e "${YELLOW}⚠ 没有图片数据${NC}"
            fi
        else
            echo -e "${YELLOW}⚠ images字段为空或不是数组${NC}"
        fi
    else
        echo -e "${YELLOW}⚠ images字段为空或不存在${NC}"
    fi
    
    # 检查photos字段（兼容性）
    HAS_PHOTOS=$(echo "$FACTORY_PROFILE_RESPONSE" | jq -r '.data.photos // empty')
    if [ "$HAS_PHOTOS" != "null" ] && [ "$HAS_PHOTOS" != "" ]; then
        echo -e "${GREEN}✓ photos字段存在（兼容性）${NC}"
        echo -e "${YELLOW}photos字段内容: $HAS_PHOTOS${NC}"
    else
        echo -e "${YELLOW}⚠ photos字段为空${NC}"
    fi
else
    echo -e "${RED}✗ 工厂详情API请求失败${NC}"
    echo "$FACTORY_PROFILE_RESPONSE"
    exit 1
fi

# 3. 测试图片列表API
echo -e "\n${YELLOW}3. 测试图片列表API...${NC}"
FACTORY_ID=$(echo "$FACTORY_PROFILE_RESPONSE" | jq -r '.data.id // empty')
if [ "$FACTORY_ID" != "null" ] && [ "$FACTORY_ID" != "" ]; then
    PHOTOS_LIST_RESPONSE=$(curl -s -X GET "${BASE_URL}/api/factories/${FACTORY_ID}/photos" \
      -H "Authorization: Bearer $TOKEN")
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ 图片列表API请求成功${NC}"
        echo "$PHOTOS_LIST_RESPONSE" | jq '.'
        
        # 检查图片列表
        PHOTOS_COUNT=$(echo "$PHOTOS_LIST_RESPONSE" | jq '.photos | length // 0')
        echo -e "${GREEN}✓ 图片列表包含 $PHOTOS_COUNT 张图片${NC}"
        
        if [ "$PHOTOS_COUNT" -gt 0 ]; then
            echo -e "${GREEN}✓ 图片列表URL:${NC}"
            echo "$PHOTOS_LIST_RESPONSE" | jq -r '.photos[].url'
        fi
    else
        echo -e "${RED}✗ 图片列表API请求失败${NC}"
        echo "$PHOTOS_LIST_RESPONSE"
    fi
else
    echo -e "${RED}✗ 无法获取工厂ID${NC}"
fi

# 4. 总结
echo -e "\n${BLUE}=== 测试总结 ===${NC}"

# 检查后端API状态
if [ "$HAS_IMAGES" != "null" ] && [ "$HAS_IMAGES" != "" ]; then
    echo -e "${GREEN}✅ 后端API images字段正常${NC}"
else
    echo -e "${RED}❌ 后端API images字段异常${NC}"
fi

# 检查图片数据
if [ "$IMAGES_COUNT" -gt 0 ]; then
    echo -e "${GREEN}✅ 后端有图片数据 (${IMAGES_COUNT}张)${NC}"
else
    echo -e "${YELLOW}⚠ 后端没有图片数据${NC}"
fi

# 检查图片可访问性
if [ "$FIRST_IMAGE_URL" != "null" ] && [ "$FIRST_IMAGE_URL" != "" ]; then
    echo -e "${GREEN}✅ 图片URL格式正确${NC}"
else
    echo -e "${RED}❌ 图片URL格式异常${NC}"
fi

echo -e "\n${BLUE}=== 前端问题诊断 ===${NC}"
echo -e "${YELLOW}如果后端API正常但前端不显示图片，可能的原因：${NC}"
echo -e "1. 前端数据模型未更新，仍在使用旧的photos字段"
echo -e "2. 前端未在图片上传后刷新数据"
echo -e "3. 前端图片URL处理不正确"
echo -e "4. 前端缓存问题"
echo -e "5. 前端图片渲染组件问题"

echo -e "\n${BLUE}建议修复方案：${NC}"
echo -e "1. 更新前端FactoryProfile数据模型，添加images字段支持"
echo -e "2. 修改图片显示组件，使用images字段的url"
echo -e "3. 确保图片上传后调用刷新API"
echo -e "4. 添加图片URL处理函数，处理相对路径"
echo -e "5. 添加缓存控制，避免缓存问题"

echo -e "\n${GREEN}测试完成！${NC}" 