#!/bin/bash

# 创建目录
sudo mkdir -p /runData/gongChang/data/nginx/ssl

# 生成自签名证书
sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout /runData/gongChang/data/nginx/ssl/aneworders.com.key \
  -out /runData/gongChang/data/nginx/ssl/aneworders.com.crt \
  -subj "/C=CN/ST=State/L=City/O=Organization/CN=aneworders.com"

# 设置权限
sudo chmod 644 /runData/gongChang/data/nginx/ssl/aneworders.com.crt
sudo chmod 600 /runData/gongChang/data/nginx/ssl/aneworders.com.key 