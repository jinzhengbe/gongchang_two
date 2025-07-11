#!/bin/bash

# 批量尝试常见工厂用户和密码组合，输出第一个成功获取token的账号和token

BASE_URL="http://localhost:8008"
USERNAMES=(gongchang factory1 testfactory admin)
PASSWORDS=(123456 test123 password admin123)

for USER in "${USERNAMES[@]}"; do
  for PASS in "${PASSWORDS[@]}"; do
    echo "尝试登录: $USER / $PASS"
    LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
      -H "Content-Type: application/json" \
      -d "{\"username\":\"$USER\",\"password\":\"$PASS\"}")
    TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    if [ ! -z "$TOKEN" ]; then
      echo "\n✅ 登录成功: $USER / $PASS"
      echo "TOKEN: $TOKEN"
      exit 0
    fi
  done
done

echo "❌ 没有账号密码组合登录成功，请检查数据库或用户信息。"
exit 1 