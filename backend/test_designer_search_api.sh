#!/bin/bash

# 设计师搜索API测试脚本
# 测试设计师搜索、建议、专业领域和评分功能

echo "开始测试设计师搜索API..."

# 基础URL
BASE_URL="http://localhost:8008"

# 测试账号信息
DESIGNER_USERNAME="sdf"
DESIGNER_PASSWORD="123456"

# 1. 测试设计师搜索接口
echo ""
echo "1. 测试设计师搜索接口"
echo "GET $BASE_URL/api/designers/search"

echo ""
echo "基础搜索测试："
curl -s -X GET "$BASE_URL/api/designers/search?page=1&page_size=5" | jq '.'

echo ""
echo "地区筛选搜索测试："
curl -s -X GET "$BASE_URL/api/designers/search?region=深圳&page=1&page_size=5" | jq '.'

echo ""
echo "评分筛选搜索测试："
curl -s -X GET "$BASE_URL/api/designers/search?min_rating=4.0&page=1&page_size=5" | jq '.'

echo ""
echo "按名称排序测试："
curl -s -X GET "$BASE_URL/api/designers/search?sort_by=name&sort_order=asc&page=1&page_size=5" | jq '.'

# 2. 测试搜索建议接口
echo ""
echo ""
echo "2. 测试搜索建议接口"
echo "GET $BASE_URL/api/designers/search/suggestions"

echo ""
echo "设计师名称建议测试："
curl -s -X GET "$BASE_URL/api/designers/search/suggestions?query=设计&limit=5" | jq '.'

echo ""
echo "地址建议测试："
curl -s -X GET "$BASE_URL/api/designers/search/suggestions?query=深圳&limit=5" | jq '.'

echo ""
echo "专业领域建议测试："
curl -s -X GET "$BASE_URL/api/designers/search/suggestions?query=服装&limit=5" | jq '.'

# 3. 测试设计师专业领域创建接口
echo ""
echo ""
echo "3. 测试设计师专业领域创建接口"
echo "POST $BASE_URL/api/designers/1/specialties"

echo ""
echo "登录获取token："
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\": \"$DESIGNER_USERNAME\", \"password\": \"$DESIGNER_PASSWORD\"}")

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
echo "Token: $TOKEN"

if [ "$TOKEN" != "null" ] && [ "$TOKEN" != "" ]; then
    echo ""
    echo "创建专业领域测试："
    curl -s -X POST "$BASE_URL/api/designers/1/specialties" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{"specialty": "服装设计"}' | jq '.'
else
    echo "登录失败，跳过需要认证的测试"
fi

# 4. 测试设计师评分创建接口
echo ""
echo ""
echo "4. 测试设计师评分创建接口"
echo "POST $BASE_URL/api/designers/1/ratings"

if [ "$TOKEN" != "null" ] && [ "$TOKEN" != "" ]; then
    echo ""
    echo "创建评分测试："
    curl -s -X POST "$BASE_URL/api/designers/1/ratings" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{"rating": 4.5, "comment": "设计很专业，沟通顺畅"}' | jq '.'
else
    echo "登录失败，跳过需要认证的测试"
fi

# 5. 测试错误情况
echo ""
echo ""
echo "5. 测试错误情况"

echo ""
echo "无效的搜索参数测试："
curl -s -X GET "$BASE_URL/api/designers/search?page=0&page_size=1000" | jq '.'

echo ""
echo "空的搜索建议测试："
curl -s -X GET "$BASE_URL/api/designers/search/suggestions" | jq '.'

echo ""
echo "无效的评分测试："
if [ "$TOKEN" != "null" ] && [ "$TOKEN" != "" ]; then
    curl -s -X POST "$BASE_URL/api/designers/1/ratings" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{"rating": 6.0, "comment": "无效评分"}' | jq '.'
else
    echo "登录失败，跳过需要认证的测试"
fi

echo ""
echo "设计师搜索API测试完成！"

# 6. 性能测试
echo ""
echo ""
echo "6. 性能测试"

echo ""
echo "测试搜索响应时间："
time curl -s -X GET "$BASE_URL/api/designers/search?page=1&page_size=10" > /dev/null

echo ""
echo "测试建议响应时间："
time curl -s -X GET "$BASE_URL/api/designers/search/suggestions?query=设计&limit=10" > /dev/null

echo ""
echo "性能测试完成！" 