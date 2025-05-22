# 基于官方 Nginx 镜像
FROM nginx:1.25-alpine

# 拷贝 Flutter Web 构建产物到 Nginx 静态目录
COPY . /usr/share/nginx/html

# 拷贝自定义 Nginx 配置
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"] 