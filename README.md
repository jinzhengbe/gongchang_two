# 工厂订单管理系统

## 项目说明
这是一个基于 Go 语言开发的工厂订单管理系统，提供订单管理、文件上传、用户认证等功能。

## 环境要求
- Go 1.16+
- MySQL 8.0+
- Docker & Docker Compose

## 安装步骤

1. 克隆项目
```bash
git clone <repository-url>
cd gongChang
```

2. 配置环境变量
```bash
cp backend/.env.example backend/.env
# 编辑 .env 文件，设置数据库连接信息等
```

3. 使用 Docker Compose 启动服务
```bash
docker-compose up -d
```

4. 访问服务
- 后端 API: https://localhost:443 或 http://localhost:8080
- MySQL: localhost:3307

## 项目结构
```
.
├── backend/           # 后端服务
│   ├── controllers/  # 控制器
│   ├── models/       # 数据模型
│   ├── routes/       # 路由
│   ├── services/     # 业务逻辑
│   ├── utils/        # 工具函数
│   ├── main.go       # 主程序
│   └── Dockerfile    # 后端 Docker 配置
├── frontend/         # 前端应用
│   ├── lib/         # 库文件
│   ├── test/        # 测试文件
│   └── pubspec.yaml # Flutter 配置
├── docker-compose.yml # Docker 编排配置
└── README.md         # 项目文档
```

## 开发环境配置

### 后端服务

1. 安装依赖：
   ```bash
   cd backend
   go mod download
   ```

2. 配置环境变量：
   ```bash
   cp .env.example .env
   # 编辑 .env 文件，设置数据库连接信息
   ```

3. 启动服务：
   ```bash
   go run main.go
   ```

### 前端应用

1. 安装 Flutter 开发环境
2. 安装依赖：
   ```bash
   cd frontend
   flutter pub get
   ```

3. 启动开发服务器：
   ```bash
   flutter run
   ```

## Docker 部署

1. 构建镜像：
   ```bash
   docker-compose build
   ```

2. 启动服务：
   ```bash
   docker-compose up -d
   ```

3. 查看日志：
   ```bash
   docker-compose logs -f
   ```

## API 文档

### 用户认证

- POST `/api/auth/login`
  - 请求体：
    ```json
    {
      "username": "string",
      "password": "string"
    }
    ```
  - 响应：
    ```json
    {
      "token": "string",
      "user": {
        "id": "string",
        "username": "string",
        "role": "string"
      }
    }
    ```

### 用户管理

- GET `/api/users`
  - 需要认证
  - 响应：用户列表

- POST `/api/users`
  - 需要认证
  - 请求体：
    ```json
    {
      "username": "string",
      "password": "string",
      "role": "string"
    }
    ```

### 订单管理

- GET `/api/orders`
  - 需要认证
  - 响应：订单列表

- POST `/api/orders`
  - 需要认证
  - 请求体：
    ```json
    {
      "productId": "string",
      "quantity": "number",
      "customer": "string"
    }
    ```

## 测试

### 后端测试

```bash
cd backend
go test ./...
```

### 前端测试

```bash
cd frontend
flutter test
```

## 维护说明

1. 数据库备份
   - 定期备份数据库文件
   - 使用 `mysqldump` 导出数据

2. 服务监控
   - 监控服务日志
   - 检查服务状态

3. 依赖更新
   - 定期更新 Go 依赖
   - 更新 Flutter 依赖

4. 安全维护
   - 定期检查安全漏洞
   - 更新依赖包修复漏洞

## 贡献指南
1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证
MIT License
