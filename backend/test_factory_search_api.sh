#!/bin/bash

# 工厂搜索API测试脚本
# 用于测试工厂搜索功能的各个接口

BASE_URL="http://localhost:8008"
API_BASE="$BASE_URL/api"

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}开始测试工厂搜索API...${NC}"

# 1. 测试工厂搜索接口
echo -e "\n${GREEN}1. 测试工厂搜索接口${NC}"
echo "GET $API_BASE/factories/search"

# 基础搜索
echo -e "\n${YELLOW}基础搜索测试：${NC}"
curl -s -X GET "$API_BASE/factories/search?query=工厂&page=1&page_size=5" | jq '.'

# 地区筛选搜索
echo -e "\n${YELLOW}地区筛选搜索测试：${NC}"
curl -s -X GET "$API_BASE/factories/search?region=上海&page=1&page_size=5" | jq '.'

# 评分筛选搜索
echo -e "\n${YELLOW}评分筛选搜索测试：${NC}"
curl -s -X GET "$API_BASE/factories/search?min_rating=4.0&page=1&page_size=5" | jq '.'

# 排序测试
echo -e "\n${YELLOW}按名称排序测试：${NC}"
curl -s -X GET "$API_BASE/factories/search?sort_by=name&sort_order=asc&page=1&page_size=5" | jq '.'

# 2. 测试搜索建议接口
echo -e "\n${GREEN}2. 测试搜索建议接口${NC}"
echo "GET $API_BASE/factories/search/suggestions"

# 工厂名称建议
echo -e "\n${YELLOW}工厂名称建议测试：${NC}"
curl -s -X GET "$API_BASE/factories/search/suggestions?query=工厂&limit=5" | jq '.'

# 地址建议
echo -e "\n${YELLOW}地址建议测试：${NC}"
curl -s -X GET "$API_BASE/factories/search/suggestions?query=上海&limit=5" | jq '.'

# 专业领域建议
echo -e "\n${YELLOW}专业领域建议测试：${NC}"
curl -s -X GET "$API_BASE/factories/search/suggestions?query=服装&limit=5" | jq '.'

# 3. 测试工厂专业领域创建接口（需要认证）
echo -e "\n${GREEN}3. 测试工厂专业领域创建接口${NC}"
echo "POST $API_BASE/factories/1/specialties"

# 先登录获取token
echo -e "\n${YELLOW}登录获取token：${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testfactory",
    "password": "123456"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
echo "Token: $TOKEN"

if [ "$TOKEN" != "null" ] && [ "$TOKEN" != "" ]; then
    echo -e "\n${YELLOW}创建专业领域测试：${NC}"
    curl -s -X POST "$API_BASE/factories/1/specialties" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "specialty": "服装制造"
      }' | jq '.'
else
    echo -e "${RED}登录失败，跳过需要认证的测试${NC}"
fi

# 4. 测试工厂评分创建接口（需要认证）
echo -e "\n${GREEN}4. 测试工厂评分创建接口${NC}"
echo "POST $API_BASE/factories/1/ratings"

if [ "$TOKEN" != "null" ] && [ "$TOKEN" != "" ]; then
    echo -e "\n${YELLOW}创建评分测试：${NC}"
    curl -s -X POST "$API_BASE/factories/1/ratings" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "rating": 4.5,
        "comment": "服务很好，质量不错"
      }' | jq '.'
else
    echo -e "${RED}登录失败，跳过需要认证的测试${NC}"
fi

# 5. 测试错误情况
echo -e "\n${GREEN}5. 测试错误情况${NC}"

# 无效的搜索参数
echo -e "\n${YELLOW}无效的搜索参数测试：${NC}"
curl -s -X GET "$API_BASE/factories/search?page_size=1000" | jq '.'

# 空的搜索建议
echo -e "\n${YELLOW}空的搜索建议测试：${NC}"
curl -s -X GET "$API_BASE/factories/search/suggestions" | jq '.'

# 无效的评分
echo -e "\n${YELLOW}无效的评分测试：${NC}"
if [ "$TOKEN" != "null" ] && [ "$TOKEN" != "" ]; then
    curl -s -X POST "$API_BASE/factories/1/ratings" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "rating": 6.0,
        "comment": "无效评分"
      }' | jq '.'
fi

echo -e "\n${GREEN}工厂搜索API测试完成！${NC}"

# 6. 性能测试
echo -e "\n${GREEN}6. 性能测试${NC}"

echo -e "\n${YELLOW}测试搜索响应时间：${NC}"
time curl -s -X GET "$API_BASE/factories/search?page=1&page_size=20" > /dev/null

echo -e "\n${YELLOW}测试建议响应时间：${NC}"
time curl -s -X GET "$API_BASE/factories/search/suggestions?query=工厂&limit=10" > /dev/null

echo -e "\n${GREEN}性能测试完成！${NC}" 