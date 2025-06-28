#!/bin/bash

# 测试从订单移除布料的API（修复后）
echo "=== 测试从订单移除布料API（修复后）==="

# 设置基础URL
BASE_URL="https://aneworders.com"

# 1. 使用订单ID 29
ORDER_ID=29
echo "使用订单ID: $ORDER_ID"

# 2. 使用布料ID 44
FABRIC_ID=44
echo "移除布料ID: $FABRIC_ID"

# 3. 调用移除布料API
echo "调用移除布料API..."
curl -X DELETE "$BASE_URL/api/orders/$ORDER_ID/remove-fabric" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d "{
    \"fabricId\": $FABRIC_ID
  }" | jq '.'

echo "=== 测试完成 ===" 