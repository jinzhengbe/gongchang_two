FROM golang:1.20-alpine

WORKDIR /app

# 安装必要的依赖
RUN apk add --no-cache gcc musl-dev wget

# 复制 Go 模块文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制项目文件
COPY . .

# 设置环境变量
ENV GO111MODULE=on
ENV CGO_ENABLED=1

# 构建应用
RUN go build -o main .

# 暴露端口
EXPOSE 8008

# 运行应用
CMD ["./main"] 