#!/bin/bash

# 综合测试所有删除API
echo "=== 综合测试所有删除API ==="

# 设置基础URL
BASE_URL="http://localhost:8008"

echo "1. 测试删除布料（从订单移除）API"
echo "----------------------------------------"
ORDER_ID=1
FABRIC_ID=123

echo "移除订单 $ORDER_ID 中的布料 $FABRIC_ID"
curl -X DELETE "$BASE_URL/api/orders/$ORDER_ID/remove-fabric" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d "{
    \"fabricId\": $FABRIC_ID
  }" | jq '.'

echo ""
echo "2. 测试删除布料（物理删除）API"
echo "----------------------------------------"
FABRIC_ID_TO_DELETE=1

echo "删除布料ID: $FABRIC_ID_TO_DELETE"
curl -X DELETE "$BASE_URL/api/fabrics/$FABRIC_ID_TO_DELETE" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq '.'

echo ""
echo "3. 测试删除文件API"
echo "----------------------------------------"
FILE_ID="test_file_123"

echo "删除文件ID: $FILE_ID"
curl -X DELETE "$BASE_URL/api/files/$FILE_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq '.'

echo ""
echo "=== 所有测试完成 ===" 