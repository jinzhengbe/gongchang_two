#!/bin/bash

# 测试API路由配置
echo "=== 测试API路由配置 ==="

# 设置基础URL
BASE_URL="http://localhost:8008"

echo "1. 测试健康检查接口"
curl -X GET "$BASE_URL/api/health" | jq '.'

echo ""
echo "2. 测试公开布料接口"
curl -X GET "$BASE_URL/api/fabrics/all" | jq '.[0:1]'

echo ""
echo "3. 测试需要认证的接口（应该返回401）"
echo "测试删除布料接口："
curl -X DELETE "$BASE_URL/api/fabrics/50" \
  -H "Content-Type: application/json" | jq '.'

echo ""
echo "测试从订单移除布料接口："
curl -X DELETE "$BASE_URL/api/orders/1/remove-fabric" \
  -H "Content-Type: application/json" \
  -d '{"fabricId": 123}' | jq '.'

echo ""
echo "测试删除文件接口："
curl -X DELETE "$BASE_URL/api/files/test_file_123" \
  -H "Content-Type: application/json" | jq '.'

echo ""
echo "=== 路由测试完成 ===" 