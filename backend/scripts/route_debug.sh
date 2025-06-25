#!/bin/bash

# 调用布料创建接口
API_URL="http://localhost:8008/api/fabrics"
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTA3MzE3MDAsInJvbGUiOiJmYWN0b3J5IiwidXNlcl9pZCI6IjRjOTE4MzUyLThjYTItNDE3OC04OWQ1LTBhMDBjMWQxODgyNSJ9.J5MdTmcdJkK9L-4-tCTpLOGO4TOYZQUwWBrz8W5N9Ek"

DATA='{"name":"route_debug_test","category":"测试","material":"测试","color":"测试","pattern":"测试","weight":1,"width":1,"price":1,"unit":"米","stock":1,"min_order":1,"description":"route debug test"}'

RESPONSE_FILE="route_debug_response_$(date +%s).log"
LOG_FILE="route_debug_dockerlog_$(date +%s).log"

# 调用接口并记录响应
curl -X POST "$API_URL" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "$DATA" \
  -i > "$RESPONSE_FILE" 2>&1

echo "==== API Response saved to $RESPONSE_FILE ===="

# 记录后端日志
sudo docker logs gongchang_backend --tail 100 > "$LOG_FILE"
echo "==== Docker logs saved to $LOG_FILE ====" 