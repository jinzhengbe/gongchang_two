# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /build

# 安装依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码和证书文件
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o app

# 生产阶段
FROM alpine:latest

RUN apk add --no-cache ca-certificates && \
    mkdir /app && \
    mkdir /ssl

WORKDIR /app

COPY --from=builder /build/app .
COPY --from=builder /build/.env .
COPY --from=builder /build/ssl/cert.pem /ssl/cert.pem
COPY --from=builder /build/ssl/key.pem /ssl/key.pem

# 设置权限
RUN chmod +x ./app && \
    chmod 644 /ssl/cert.pem && \
    chmod 600 /ssl/key.pem

# 暴露端口
EXPOSE 8080
EXPOSE 443

# 运行应用
CMD ["./app"] 