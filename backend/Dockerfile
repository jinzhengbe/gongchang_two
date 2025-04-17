FROM golang:1.20-alpine

WORKDIR /app

# 安装必要的依赖
RUN apk add --no-cache gcc musl-dev

# 复制 Go 模块文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制项目文件
COPY . .

# 创建 SSL 目录并复制证书
RUN mkdir -p /app/ssl && \
    chmod 600 /app/ssl/cert.pem /app/ssl/key.pem

# 构建应用
RUN go build -o app .

# 暴露端口
EXPOSE 443

# 运行应用
CMD ["./app"] 