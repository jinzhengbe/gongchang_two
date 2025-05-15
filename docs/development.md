# 开发指南

## 重要说明

### 数据库存储位置
- 数据库数据存储在容器外，具体位置：`/runData/gongChang/mysql_data`
- 数据持久化：即使容器被删除，数据也不会丢失
- 数据备份：可以直接备份主机上的目录
- 数据迁移：可以轻松地将数据迁移到其他主机

## 最近更新

### 2025-05-06
- 优化了文件上传功能
  - 问题：文件上传时出现 SocketException: Connection reset by peer 错误
  - 原因：
    1. Nginx 配置中的文件大小限制不足
    2. 上传路由未正确注册
    3. 文件上传目录权限问题
  - 解决方案：
    1. 增加了 Nginx 配置中的文件大小限制
    2. 修复了文件上传路由注册
    3. 优化了文件上传目录权限设置
    4. 增强了错误处理和日志记录
  - 预防措施：
    1. 定期检查上传目录权限
    2. 监控文件上传错误日志
    3. 实现文件上传进度跟踪
  - 相关文件：
    - `backend/controllers/file.go`
    - `backend/services/file.go`
    - `backend/routes/file.go`
    - `data/nginx/conf.d/default.conf`
- 改进了文件上传功能
  - 变更：将订单ID（orderId）改为可选参数
  - 原因：支持独立文件上传，不强制要求关联订单
  - 实现：
    1. 修改 File 模型，将 OrderID 改为指针类型
    2. 更新文件上传接口，支持无订单ID上传
    3. 优化日志记录，区分有无订单ID的情况
  - 相关文件：
    - `backend/models/file.go`
    - `backend/controllers/file.go`
    - `backend/services/file.go`

### 2025-05-05
- 修复了 MySQL 数据目录权限问题
  - 问题：MySQL 数据目录权限被 dnsmasq 和 systemd-journal 服务修改
  - 原因：k8s worker 节点存储池配置导致权限冲突
  - 解决方案：
    1. 删除不再使用的 k8s worker 节点配置
    2. 修复数据目录权限为 999:999
    3. 重启 MySQL 容器
  - 预防措施：
    1. 避免将 MySQL 数据目录放在系统服务可能访问的位置
    2. 定期检查目录权限
    3. 使用 ACL 权限控制

### 2024-05-01
- 数据库存储配置说明
  - 确认数据库数据存储在容器外
  - 添加了数据库存储位置说明文档
  - 验证了数据持久化功能
- API 路径变更
  - 统一 API 路径格式，从 `/api/v1/xxx` 改为 `/api/xxx`
  - 影响范围：订单管理相关接口
  - 解决方案：
    1. 修改客户端代码，使用新的 API 路径
    2. 添加临时兼容层，支持旧的 API 路径
    3. 更新 API 文档，明确说明路径变更
  - 相关文件：
    - `backend/routes/router.go`
    - `docs/order_api.md`
    - `docs/development.md`

### 2024-04-30
- 数据库配置更新
  - 将数据库从 MariaDB 改回 MySQL 8.0
  - 更新了 docker-compose.yml 配置
  - 优化了数据库连接配置
- 修复了后端服务启动问题
  - 修复了数据库连接问题
  - 优化了错误处理
  - 添加了重试机制

### 2024-04-29
- 更新了订单接口
  - 添加了分页和状态筛选功能
  - 修复了重复的 GetOrdersByUserID 方法
  - 优化了路由顺序
  - 配置了正确的 JWT 密钥
- 更新了服务器配置
  - 修改为监听 0.0.0.0:8080
  - 优化了错误处理

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

### 2025-05-07
- 订单结构与接口字段扩展
  - 新增字段：attachments, models, images, videos，均为 string 数组，可为空
  - 数据库存储类型：JSON（推荐），如不支持则用 TEXT 存储 JSON 字符串
  - 后端接口：支持这四个字段的读写和 JSON 解析，返回前端时保持为数组格式，避免字符串分割
  - 文档同步：已更新 API 文档和数据结构说明，便于团队协作
  - 兼容性：老数据默认空数组，前后端均已兼容
  - 相关文件：
    - backend/models/order.go
    - backend/controllers/order.go
    - docs/order_api.md
    - docs/development.md

### 2025-05-08
- 订单表字段名称变更
  - 变更内容：
    1. `attachments` 改为 `attachments`
    2. `models` 改为 `models`
    3. `images` 改为 `images`
    4. `videos` 改为 `videos`
  - 变更原因：统一字段命名规范，提高代码可读性
  - 影响范围：
    1. 数据库表结构
    2. Order 模型结构
    3. OrderRequest 结构
    4. OrderController 相关代码
  - 相关文件：
    - backend/models/order.go
    - backend/controllers/order.go
    - docs/order_api.md
    - docs/development.md
  - 注意事项：
    1. 数据库字段已通过 ALTER TABLE 语句更新
    2. 代码中的字段名已同步更新
    3. API 接口保持不变，仅内部字段名变更

### 2025-05-09
- 修复了 GORM datatypes 包依赖问题
  - 问题：运行后端服务时出现 "no required module provides package gorm.io/datatypes" 错误
  - 原因：缺少必要的 GORM datatypes 包依赖
  - 解决方案：
    1. 添加 gorm.io/datatypes 包依赖
    2. 更新 go.mod 文件
  - 执行命令：
    ```bash
    go get gorm.io/datatypes
    ```
  - 相关文件：
    - backend/models/order.go
    - backend/go.mod

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
   - 角色: designer

2. 工厂用户
   - 用户名: factory1
   - 密码: test123
   - 邮箱: factory1@test.com
   - 角色: factory

3. 供应商用户
   - 用户名: supplier1
   - 密码: test123
   - 邮箱: supplier1@test.com
   - 角色: supplier

## API 文档

### 用户认证

#### 登录
- 请求: POST `/api/users/login`
- 请求体:
```json
{
  "username": "string",
  "password": "string",
  "user_type": "string"  // 必须是 "designer", "factory", 或 "supplier"
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

#### 注册
- 请求: POST `/api/users/register`
- 请求体:
```json
{
  "username": "string",
  "password": "string",
  "email": "string",
  "role": "string"  // 必须是 "designer", "factory", 或 "supplier"
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

# 开发文档

## 系统架构

### 后端服务
- 使用 Go + Gin 框架
- MySQL 数据库
- JWT 认证
- Docker 容器化部署

### 目录结构
```
backend/
├── api/          # API 路由和处理
├── config/       # 配置文件
├── controllers/  # 控制器
├── database/     # 数据库相关
├── middleware/   # 中间件
├── models/       # 数据模型
├── routes/       # 路由配置
└── services/     # 业务逻辑
```

## 配置说明

### 数据库配置
- 数据库名：gongchang
- 用户名：gongchang
- 密码：gongchang
- 端口：3306

### JWT 配置
- 密钥：your-secret-key
- 过期时间：24小时

### 测试账号
- 设计师账号：designer1 / test123
- 工厂账号：factory1 / test123
- 供应商账号：supplier1 / test123

## API 接口

### 认证接口
- POST /api/users/login - 用户登录
- POST /api/users/register - 用户注册

### 订单接口
- GET /api/orders - 获取订单列表
- POST /api/orders - 创建订单
- GET /api/orders/:id - 获取订单详情
- PUT /api/orders/:id/status - 更新订单状态

## 注意事项

1. 数据库权限
   - MySQL 数据目录权限问题已解决
   - 确保数据目录权限为 999:999（MySQL 容器用户）
   - 避免与系统服务（如 dnsmasq）的权限冲突

2. 开发环境
   - 使用 Docker Compose 管理服务
   - 开发时注意检查容器日志
   - 确保配置文件正确加载

3. 部署注意事项
   - 确保环境变量正确设置
   - 检查数据库连接配置
   - 验证 JWT 密钥配置

## 常见问题解决

1. 权限问题
   - 如果遇到权限问题，检查目录权限
   - 使用 `chown -R 999:999` 修复权限
   - 重启相关容器

2. 数据库连接
   - 检查数据库配置
   - 确保数据库容器正常运行
   - 验证网络连接

3. API 访问
   - 确认接口路径正确
   - 检查认证 token
   - 验证请求参数格式 