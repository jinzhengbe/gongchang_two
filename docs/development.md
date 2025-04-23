# 开发文档

## 最近更新

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