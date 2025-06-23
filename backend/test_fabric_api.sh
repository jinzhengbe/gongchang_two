#!/bin/bash

# 布料API测试脚本
BASE_URL="http://localhost:8008/api"

echo "=== 布料API测试 ==="
echo "基础URL: $BASE_URL"
echo ""

# 测试1: 获取所有布料
echo "1. 测试获取所有布料..."
curl -s -X GET "$BASE_URL/fabrics/all" | jq '.' || echo "获取所有布料失败"
echo ""

# 测试2: 获取布料分类
echo "2. 测试获取布料分类..."
curl -s -X GET "$BASE_URL/fabrics/categories" | jq '.' || echo "获取布料分类失败"
echo ""

# 测试3: 搜索布料
echo "3. 测试搜索布料..."
curl -s -X GET "$BASE_URL/fabrics/search?q=棉布&page=1&page_size=5" | jq '.' || echo "搜索布料失败"
echo ""

# 测试4: 根据分类获取布料
echo "4. 测试根据分类获取布料..."
curl -s -X GET "$BASE_URL/fabrics/category/棉布?page=1&page_size=5" | jq '.' || echo "根据分类获取布料失败"
echo ""

# 测试5: 获取布料统计信息
echo "5. 测试获取布料统计信息..."
curl -s -X GET "$BASE_URL/fabrics/statistics" | jq '.' || echo "获取布料统计信息失败"
echo ""

# 测试6: 获取特定布料详情
echo "6. 测试获取布料详情..."
curl -s -X GET "$BASE_URL/fabrics/1" | jq '.' || echo "获取布料详情失败"
echo ""

echo "=== 测试完成 ==="
echo ""
echo "注意: 如果需要测试需要认证的接口，请先获取JWT Token"
echo "认证接口包括:"
echo "- POST /api/fabrics (创建布料)"
echo "- PUT /api/fabrics/{id} (更新布料)"
echo "- DELETE /api/fabrics/{id} (删除布料)"
echo "- PUT /api/fabrics/{id}/stock (更新库存)" 