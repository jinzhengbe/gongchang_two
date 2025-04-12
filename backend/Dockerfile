# 构建阶段
FROM golang:1.18-alpine as build

WORKDIR /app

# 安装构建工具
RUN apk add --no-cache gcc musl-dev

# 复制 go.mod
COPY go.mod ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 生产阶段
FROM alpine:latest

WORKDIR /app

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 从构建阶段复制二进制文件和配置文件
COPY --from=build /app/main .
COPY --from=build /app/.env .
COPY --from=build /app/config ./config

# 创建上传目录
RUN mkdir -p uploads

# 设置时区
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./main"] 