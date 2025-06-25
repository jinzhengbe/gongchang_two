#!/bin/bash

set -e

cd "$(dirname "$0")/.."

# 1. 停止并删除后端容器
sudo docker compose down

# 2. 强制重建镜像
sudo docker compose build --no-cache

# 3. 启动服务
sudo docker compose up -d --force-recreate

# 4. 等待服务启动
sleep 5

# 5. 校验容器内控制器代码是否包含panic
CONTAINER_ID=$(sudo docker ps -qf "name=gongchang_backend")
echo "==== 校验容器内controllers/fabric.go是否包含panic ===="
sudo docker exec "$CONTAINER_ID" grep 'panic("test controller panic")' /app/controllers/fabric.go || echo "[警告] 容器内未找到panic语句！"

# 6. 自动调用布料创建接口并收集响应和日志
cd backend
./scripts/route_debug.sh
cd ..

echo "==== 重建与调试完成，请检查route_debug_response_*.log和route_debug_dockerlog_*.log ====" 