#!/bin/sh

# 设置域名
DOMAIN="aneworders.com"
EMAIL="admin@aneworders.com"

# 确保 SSL 目录存在
mkdir -p /app/ssl

# 获取证书
certbot certonly --standalone \
    --non-interactive \
    --agree-tos \
    --email $EMAIL \
    --domains $DOMAIN \
    --test-cert

# 复制证书到应用目录
cp /etc/letsencrypt/live/$DOMAIN/fullchain.pem /app/ssl/cert.pem
cp /etc/letsencrypt/live/$DOMAIN/privkey.pem /app/ssl/key.pem

# 设置正确的权限
chmod 600 /app/ssl/cert.pem /app/ssl/key.pem

# 设置定时任务，每 60 天更新一次证书
echo "0 0 1 */2 * certbot renew --quiet && cp /etc/letsencrypt/live/$DOMAIN/fullchain.pem /app/ssl/cert.pem && cp /etc/letsencrypt/live/$DOMAIN/privkey.pem /app/ssl/key.pem" > /etc/crontabs/root

# 启动 crond
crond -f 