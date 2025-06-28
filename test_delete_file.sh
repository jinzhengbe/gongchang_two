#!/bin/bash

# 测试删除文件API
echo "=== 测试删除文件API ==="

# 设置基础URL
BASE_URL="http://localhost:8008"

# 1. 首先获取文件列表（通过订单文件）
echo "获取订单文件列表..."
curl -X GET "$BASE_URL/api/files/order/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq '.[0:3]'

# 2. 假设要删除的文件ID
FILE_ID="test_file_123"
echo "删除文件ID: $FILE_ID"

# 3. 调用删除文件API
echo "调用删除文件API..."
curl -X DELETE "$BASE_URL/api/files/$FILE_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq '.'

echo "=== 测试完成 ===" 