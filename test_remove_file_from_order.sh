#!/bin/bash

# 测试从订单移除文件的API
echo "=== 测试从订单移除文件API ==="

# 设置基础URL
BASE_URL="https://aneworders.com"

# 1. 使用订单ID 29（数字类型）
ORDER_ID=29
echo "使用订单ID: $ORDER_ID"

# 2. 使用文件ID（UUID格式）
FILE_ID="fbdd3f3e-0a2a-4180-8478-7e334e7d9fe7"
echo "移除文件ID: $FILE_ID"

# 3. 文件类型
FILE_TYPE="image"
echo "文件类型: $FILE_TYPE"

# 4. 调用移除文件API
echo "调用移除文件API..."
curl -X DELETE "$BASE_URL/api/orders/$ORDER_ID/remove-file" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d "{
    \"fileId\": \"$FILE_ID\",
    \"fileType\": \"$FILE_TYPE\"
  }" | jq '.'

echo "=== 测试完成 ===" 