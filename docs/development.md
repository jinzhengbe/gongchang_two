# 开发指南

## 最近更新

### 2024-04-29
- 更新了数据库配置
  - 将数据库从 MariaDB 改回 MySQL 8.0
  - 更新了 docker-compose.yml 配置
  - 优化了数据库连接配置
- 修复了后端服务启动问题
  - 修复了数据库连接配置
  - 更新了配置文件中的数据库主机名
  - 验证了服务健康检查功能
- 更新了测试用户
  - 用户名: test
  - 邮箱: test@test.com
  - 角色: designer
  - 密码: test123

### 2024-04-28
- 更新了订单接口
  - 添加了分页和状态筛选功能
  - 修复了重复的 GetOrdersByUserID 方法
  - 优化了路由顺序
  - 配置了正确的 JWT 密钥
- 更新了服务器配置
  - 修改为监听 0.0.0.0:8080
  - 优化了错误处理

### 2024-04-27
- 添加了新的测试用户
  - 用户名: testuser1
  - 邮箱: testuser1@example.com
  - 角色: designer
  - 密码: test123
  - 创建时间: 2025-04-27 15:30:40

### 2024-04-24
- 移除了邮箱唯一约束
  - 移除了 User 模型中的邮箱唯一标签
  - 移除了数据库中的邮箱唯一索引
  - 支持使用同一邮箱注册不同角色账号
  - 添加了数据库迁移脚本
- 将数据库从 MySQL 改为 MariaDB
  - 更新了 docker-compose.yml 配置
  - 使用 MariaDB 10.11 版本
  - 优化了数据库连接配置
- 修复了用户注册功能
  - 修改了用户模型，使用 gorm.Model 管理通用字段
  - 修复了重复用户名错误处理，返回 409 状态码
  - 优化了数据库表结构

### 2024-04-23
- 修改了订单模型，支持设计师订单
  - 将 `ProductID` 字段改为可空
  - 更新了相关的验证逻辑
  - 支持创建没有产品 ID 的设计师订单
- 修复了健康检查配置问题
  - 将健康检查从 wget 改为 curl
  - 确保健康检查使用 GET 请求而不是 HEAD 请求
- 修复了设计师登录功能
  - 确认测试用户密码为 "test123"
  - 验证了 JWT token 生成功能
- 更新了项目结构
  - 将后端代码移动到 backend 目录
  - 优化了目录结构

## 开发环境设置

### 1. 使用 Docker Compose 启动开发环境

本项目使用 Docker Compose 管理开发环境，所有服务（包括数据库）都运行在容器中。

#### 启动服务
```bash
docker-compose up -d
```

#### 停止服务
```bash
docker-compose down
```

#### 查看服务状态
```bash
docker-compose ps
```

#### 查看服务日志
```bash
docker-compose logs -f
```

### 2. 数据库配置

#### 数据库服务
- 类型：MySQL 8.0
- 容器名称：gongchang-mysql
- 主机名：mysql
- 端口：3306
- 数据卷：gongchang_mysql_data

#### 数据库连接信息
- 数据库名：gongchang
- 用户名：gongchang
- 密码：gongchang
- 字符集：utf8mb4
- 时区：Asia/Shanghai

#### 数据库初始化
数据库会在容器首次启动时自动创建，并包含以下系统数据库：
- information_schema
- performance_schema

#### 数据持久化
- 数据存储在 Docker 数据卷 `gongchang_mysql_data` 中
- 即使容器重启，数据也会保留
- 要完全重置数据库，需要删除数据卷：
  ```bash
  docker-compose down -v
  ```

#### 数据库管理
- 使用 MySQL 客户端连接：
  ```bash
  docker exec -it gongchang-mysql mysql -ugongchang -pgongchang
  ```
- 查看数据库列表：
  ```bash
  docker exec -it gongchang-mysql mysql -ugongchang -pgongchang -e "SHOW DATABASES;"
  ```

#### 注意事项
1. 数据库配置在 docker-compose.yml 中定义，请勿随意修改
2. 如需修改配置，请先备份数据
3. 确保使用正确的数据卷名称
4. 数据库服务会在容器重启时自动恢复

### 依赖要求
- Docker & Docker Compose
- MySQL 8.0+
- Go 1.20+

### 本地开发步骤

1. 克隆仓库
```bash
git clone <repository-url>
cd gongChang
```

2. 启动服务
```bash
docker-compose up --build
```

3. 测试服务
```bash
# 测试健康检查
curl http://localhost:8080/api/health

# 测试用户注册
curl -X POST -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"test123","email":"test@test.com","role":"designer"}' \
  http://localhost:8080/api/users/register

# 测试设计师登录
curl -X POST -H "Content-Type: application/json" \
  -d '{"username":"designer1","password":"test123"}' \
  http://localhost:8080/api/auth/login
```

## 项目结构

```
gongChang/
├── backend/           # 后端代码
│   ├── api/          # API 定义
│   ├── config/       # 配置文件
│   ├── controllers/  # 控制器
│   ├── database/     # 数据库相关
│   ├── middleware/   # 中间件
│   ├── models/       # 数据模型
│   ├── routes/       # 路由定义
│   ├── services/     # 业务逻辑
│   └── utils/        # 工具函数
├── docs/             # 文档
├── scripts/          # 脚本文件
└── ssl/              # SSL 证书
```

## 测试用户

系统预置了以下测试用户：

1. 设计师用户
   - 用户名: designer1
   - 密码: test123
   - 邮箱: designer1@test.com

2. 工厂用户
   - 用户名: factory1
   - 密码: test123
   - 邮箱: factory1@test.com

3. 供应商用户
   - 用户名: supplier1
   - 密码: test123
   - 邮箱: supplier1@test.com

## API 文档

### 用户认证

#### 注册
- 请求: POST `/api/users/register`
- 请求体:
```json
{
  "username": "string",
  "password": "string",
  "email": "string",
  "role": "string"
}
```
- 响应:
  - 成功 (201 Created):
  ```json
  {
    "message": "User registered successfully"
  }
  ```
  - 用户名已存在 (409 Conflict):
  ```json
  {
    "error": "Username already exists"
  }
  ```

#### 登录
- 请求: POST `/api/auth/login`
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
  "token": "string",
  "user": {
    "id": "string",
    "username": "string",
    "email": "string",
    "role": "string"
  }
}
```

### 健康检查

- 请求: GET `/api/health`
- 响应:
```json
{
  "status": "ok"
}
```

### 订单管理

#### 获取订单列表
- 请求: GET `/api/orders`
- 查询参数:
  - `status`: 订单状态 (可选)
  - `page`: 页码 (可选，默认 1)
  - `pageSize`: 每页大小 (可选，默认 10)
- 请求头:
  - `Authorization: Bearer <token>`
- 响应:
```json
{
  "total": 0,
  "page": 1,
  "pageSize": 10,
  "orders": [
    {
      "id": 1,
      "title": "string",
      "description": "string",
      "designer_id": "string",
      "customer_id": "string",
      "quantity": 0,
      "unit_price": 0,
      "total_price": 0,
      "status": "string",
      "payment_status": "string",
      "shipping_address": "string",
      "order_date": "string",
      "order_type": "string",
      "fabrics": "string",
      "delivery_date": "string",
      "special_requirements": "string",
      "designer": {
        "id": "string",
        "username": "string",
        "email": "string",
        "role": "string"
      },
      "customer": {
        "id": "string",
        "username": "string",
        "email": "string",
        "role": "string"
      }
    }
  ]
}
```

#### 创建订单
- 请求: POST `/api/orders`
- 请求头:
  - `Authorization: Bearer <token>`
  - `Content-Type: application/json`
- 请求体:
```json
{
  "title": "string",
  "description": "string",
  "quantity": 0,
  "shipping_address": "string",
  "order_type": "string",
  "fabrics": "string",
  "delivery_date": "string",
  "special_requirements": "string"
}
```

#### 获取订单详情
- 请求: GET `/api/orders/:id`
- 请求头:
  - `Authorization: Bearer <token>`
- 响应: 返回单个订单的详细信息

#### 更新订单状态
- 请求: PUT `/api/orders/:id/status`
- 请求头:
  - `Authorization: Bearer <token>`
  - `Content-Type: application/json`
- 请求体:
```json
{
  "status": "string"
}
```

#### 获取最近订单
- 请求: GET `/api/orders/recent`
- 查询参数:
  - `limit`: 返回数量 (可选，默认 5)
- 请求头:
  - `Authorization: Bearer <token>`

#### 获取最新订单
- 请求: GET `/api/orders/latest`
- 请求头:
  - `Authorization: Bearer <token>`

#### 获取热门订单
- 请求: GET `/api/orders/hot`
- 请求头:
  - `Authorization: Bearer <token>`

## 常见问题

### 健康检查失败
如果遇到健康检查失败的问题，请检查：
1. 确保服务正在运行
2. 检查日志中的错误信息
3. 确认健康检查配置正确

### 登录失败
如果登录失败，请检查：
1. 用户名和密码是否正确
2. 用户是否存在于数据库中
3. 数据库连接是否正常

## 开发规范

### 代码提交
- 提交信息格式: `<type>: <description>`
- 类型包括: feat, fix, docs, style, refactor, test, chore
- 描述要简洁明了

### 分支管理
- main: 主分支，用于生产环境
- develop: 开发分支，用于日常开发
- feature/*: 功能分支，用于开发新功能
- hotfix/*: 修复分支，用于紧急修复

## 部署说明

### 生产环境部署
1. 拉取最新代码
2. 构建镜像
3. 启动服务
4. 验证服务状态

### 备份策略
- 数据库每日备份
- 配置文件定期备份
- 日志文件定期归档 