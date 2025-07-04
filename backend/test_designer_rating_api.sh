#!/bin/bash

# 设计师评分API扩展测试脚本
# 测试设计师评分的完整功能：创建、列表、统计

echo "开始测试设计师评分API扩展功能..."

# 基础URL
BASE_URL="http://localhost:8008"

# 测试账号信息
DESIGNER_USERNAME="sdf"
DESIGNER_PASSWORD="123456"

# 登录获取token
echo ""
echo "登录获取token："
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\": \"$DESIGNER_USERNAME\", \"password\": \"$DESIGNER_PASSWORD\"}")

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
echo "Token: $TOKEN"

if [ "$TOKEN" = "null" ] || [ "$TOKEN" = "" ]; then
    echo "登录失败，退出测试"
    exit 1
fi

# 设计师ID（使用之前测试过的设计师）
DESIGNER_ID=1

echo ""
echo "使用设计师ID: $DESIGNER_ID"

# 1. 创建多个评分测试
echo ""
echo "1. 创建多个评分测试"

echo ""
echo "创建评分1（5分）："
curl -s -X POST "$BASE_URL/api/designers/$DESIGNER_ID/ratings" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"rating": 5.0, "comment": "设计非常优秀，创意十足，沟通顺畅"}' | jq '.'

echo ""
echo "创建评分2（4分）："
curl -s -X POST "$BASE_URL/api/designers/$DESIGNER_ID/ratings" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"rating": 4.0, "comment": "设计不错，风格统一"}' | jq '.'

echo ""
echo "创建评分3（3分）："
curl -s -X POST "$BASE_URL/api/designers/$DESIGNER_ID/ratings" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"rating": 3.0, "comment": "设计一般，需要改进"}' | jq '.'

# 2. 获取评分列表
echo ""
echo ""
echo "2. 获取评分列表测试"

echo ""
echo "获取评分列表（第1页）："
curl -s -X GET "$BASE_URL/api/designers/$DESIGNER_ID/ratings?page=1&page_size=5" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""
echo "获取评分列表（第2页）："
curl -s -X GET "$BASE_URL/api/designers/$DESIGNER_ID/ratings?page=2&page_size=5" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# 3. 获取评分统计
echo ""
echo ""
echo "3. 获取评分统计测试"

echo ""
echo "获取评分统计："
curl -s -X GET "$BASE_URL/api/designers/$DESIGNER_ID/ratings/stats" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# 4. 测试错误情况
echo ""
echo ""
echo "4. 测试错误情况"

echo ""
echo "获取不存在的设计师评分："
curl -s -X GET "$BASE_URL/api/designers/999/ratings" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""
echo "获取不存在的设计师评分统计："
curl -s -X GET "$BASE_URL/api/designers/999/ratings/stats" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""
echo "无效的分页参数："
curl -s -X GET "$BASE_URL/api/designers/$DESIGNER_ID/ratings?page=0&page_size=1000" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# 5. 性能测试
echo ""
echo ""
echo "5. 性能测试"

echo ""
echo "测试评分列表响应时间："
time curl -s -X GET "$BASE_URL/api/designers/$DESIGNER_ID/ratings?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN" > /dev/null

echo ""
echo "测试评分统计响应时间："
time curl -s -X GET "$BASE_URL/api/designers/$DESIGNER_ID/ratings/stats" \
  -H "Authorization: Bearer $TOKEN" > /dev/null

echo ""
echo "设计师评分API扩展功能测试完成！"

# 6. 验证搜索功能中的评分
echo ""
echo ""
echo "6. 验证搜索功能中的评分"

echo ""
echo "搜索评分≥4.0的设计师："
curl -s -X GET "$BASE_URL/api/designers/search?min_rating=4.0&page=1&page_size=10" | jq '.'

echo ""
echo "搜索评分≥3.0的设计师："
curl -s -X GET "$BASE_URL/api/designers/search?min_rating=3.0&page=1&page_size=10" | jq '.' 