#!/bin/bash

# 测试删除布料API
echo "=== 测试删除布料API ==="

# 设置基础URL
BASE_URL="http://localhost:8008"

# 1. 首先获取布料列表
echo "获取布料列表..."
curl -X GET "$BASE_URL/api/fabrics/all" \
  -H "Content-Type: application/json" | jq '.[0:3]'

# 2. 假设要删除的布料ID
FABRIC_ID=1
echo "删除布料ID: $FABRIC_ID"

# 3. 调用删除布料API
echo "调用删除布料API..."
curl -X DELETE "$BASE_URL/api/fabrics/$FABRIC_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq '.'

echo "=== 测试完成 ===" 