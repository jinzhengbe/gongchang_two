#!/bin/bash

# 启动服务
echo "启动服务..."
docker-compose down
docker-compose up -d

# 等待 MySQL 就绪
echo "等待 MySQL 就绪..."
until docker-compose exec mysql mysqladmin ping -h localhost --silent; do
    echo "等待 MySQL 启动..."
    sleep 2
done

echo "服务启动完成！" 