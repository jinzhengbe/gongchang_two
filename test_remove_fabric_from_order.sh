#!/bin/bash

# 测试从订单移除布料的API
echo "=== 测试从订单移除布料API ==="

# 设置基础URL
BASE_URL="http://localhost:8008"

# 1. 首先获取一个订单ID（假设订单ID为1）
ORDER_ID=1
echo "使用订单ID: $ORDER_ID"

# 2. 获取该订单的布料列表
echo "获取订单布料列表..."
curl -X GET "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq '.'

# 3. 假设要移除的布料ID
FABRIC_ID=123
echo "移除布料ID: $FABRIC_ID"

# 4. 调用移除布料API
echo "调用移除布料API..."
curl -X DELETE "$BASE_URL/api/orders/$ORDER_ID/remove-fabric" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d "{
    \"fabricId\": $FABRIC_ID
  }" | jq '.'

echo "=== 测试完成 ===" 