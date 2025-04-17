# 服装制造管理系统 - 后端服务

## 项目简介

这是一个基于 Go 语言开发的服装制造管理系统后端服务，提供 RESTful API 接口。

## 技术栈

- Go 1.20+
- Gin Web Framework
- GORM
- MySQL 8.0
- JWT 认证

## 目录结构

```
backend/
├── config/         # 配置文件
├── controllers/    # 控制器
├── database/       # 数据库初始化
├── middleware/     # 中间件
├── models/        # 数据模型
├── routes/        # 路由配置
├── services/      # 业务逻辑
├── utils/         # 工具函数
├── main.go        # 主程序
└── Dockerfile     # Docker 配置文件
```

## 快速开始

### 环境要求

- Go 1.20+
- MySQL 8.0+
- Docker & Docker Compose (可选)

### 本地开发

1. 克隆项目
```bash
git clone <repository-url>
cd backend
```

2. 安装依赖
```bash
go mod download
```

3. 配置数据库
```bash
# 创建数据库和用户
mysql -u root -p
CREATE DATABASE sewingmast;
CREATE USER 'sewingmast'@'localhost' IDENTIFIED BY 'sewingmast123';
GRANT ALL PRIVILEGES ON sewingmast.* TO 'sewingmast'@'localhost';
FLUSH PRIVILEGES;
```

4. 运行服务
```bash
go run main.go
```

### Docker 部署

1. 构建并启动服务
```bash
docker-compose up -d
```

2. 查看日志
```bash
docker-compose logs -f
```

## API 接口

### 认证接口

#### 登录
- 路径: `/api/auth/login`
- 方法: `POST`
- 请求体:
```json
{
  "username": "string",
  "password": "string"
}
```
- 响应:
```json
{
  "success": true,
  "data": {
    "token": "string",
    "user": {
      "id": "string",
      "username": "string",
      "role": "string"
    }
  }
}
```

### 用户接口

#### 注册用户
- 路径: `/api/auth/register`
- 方法: `POST`
- 请求体:
```json
{
  "username": "string",
  "password": "string",
  "email": "string",
  "role": "string"
}
```

#### 获取用户信息
- 路径: `/api/users/:id`
- 方法: `GET`
- 需要认证: 是

## 配置说明

### 环境变量

- `DB_HOST`: 数据库主机地址
- `DB_PORT`: 数据库端口
- `DB_USER`: 数据库用户名
- `DB_PASSWORD`: 数据库密码
- `DB_NAME`: 数据库名称
- `JWT_SECRET`: JWT 密钥

## 开发说明

### 添加新的 API 接口

1. 在 `models` 目录下定义数据模型
2. 在 `services` 目录下实现业务逻辑
3. 在 `controllers` 目录下创建控制器
4. 在 `routes` 目录下注册路由

### 数据库迁移

项目使用 GORM 自动迁移功能，在启动时会自动创建/更新数据库表结构。

## 测试

运行单元测试：
```bash
go test ./...
```

## 部署

### 使用 Docker

1. 构建镜像：
```bash
docker build -t sewingmast-backend .
```

2. 运行容器：
```bash
docker run -p 8080:8080 sewingmast-backend
```

### 使用 Docker Compose

```bash
docker-compose up -d
```

## 端口使用说明

- 后端服务端口：8080
- MySQL 端口：3306

### 端口配置

1. 开发环境
   - 后端服务默认使用 8080 端口
   - MySQL 默认使用 3306 端口

2. Docker 环境
   - 可以通过修改 docker-compose.yml 中的端口映射来更改端口

### 端口冲突解决

如果遇到端口冲突，可以：

1. 修改 main.go 中的服务端口
2. 修改 docker-compose.yml 中的端口映射
3. 确保没有其他服务占用相同端口

### 安全建议

1. 在生产环境中使用非默认端口
2. 配置防火墙只允许必要的端口访问
3. 使用反向代理（如 Nginx）来管理端口转发
