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
CREATE DATABASE gongchang;
GRANT ALL PRIVILEGES ON gongchang.* TO 'root'@'%' IDENTIFIED BY '123456';
FLUSH PRIVILEGES;
```

4. 配置环境变量
```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=123456
export DB_NAME=gongchang
```

5. 运行服务
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
- 路径: `http://aneworders.com:8080/api/users/login` 或 `https://aneworders.com/api/users/login`
- 方法: `POST`
- 请求头: `Content-Type: application/json`
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
    "role": "string",
    "email": "string"
  }
}
```

### 测试账号

系统初始化时会自动创建以下测试账号：

1. 设计师账号
   - 用户名: `designer1`
   - 密码: `test123`
   - 角色: `designer`

2. 工厂账号
   - 用户名: `factory1`
   - 密码: `test123`
   - 角色: `factory`

3. 供应商账号
   - 用户名: `supplier1`
   - 密码: `test123`
   - 角色: `supplier`

### 用户接口

#### 注册用户
- 路径: `http://aneworders.com:8080/api/users/register` 或 `https://aneworders.com/api/users/register`
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
- 路径: `http://aneworders.com:8080/api/users/:id` 或 `https://aneworders.com/api/users/:id`
- 方法: `GET`
- 需要认证: 是

## 开发指南

### 开发环境设置

1. 安装开发工具
```bash
# 安装 Go
brew install go  # macOS
sudo apt install golang-go  # Ubuntu

# 安装 MySQL
brew install mysql  # macOS
sudo apt install mysql-server  # Ubuntu

# 安装 Docker
brew install docker docker-compose  # macOS
sudo apt install docker.io docker-compose  # Ubuntu
```

2. 配置开发环境
```bash
# 设置 Go 环境变量
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# 安装开发依赖
go install github.com/cosmtrek/air@latest  # 热重载工具
go install github.com/swaggo/swag/cmd/swag@latest  # API 文档生成
```

3. 启动开发服务器
```bash
# 使用 air 进行热重载开发
air

# 或直接运行
go run main.go
```

### 代码规范

1. 命名规范
   - 包名：小写字母，不使用下划线
   - 文件名：小写字母，使用下划线分隔
   - 结构体：大驼峰命名
   - 变量和函数：小驼峰命名

2. 注释规范
   - 包注释：说明包的功能
   - 函数注释：说明功能、参数和返回值
   - 关键代码注释：说明实现逻辑

3. 错误处理
   - 使用 errors.New 创建简单错误
   - 使用 fmt.Errorf 包装错误
   - 记录关键错误日志

### 测试规范

1. 单元测试
```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./controllers

# 显示测试覆盖率
go test -cover ./...
```

2. 集成测试
```bash
# 使用 Docker Compose 运行测试环境
docker-compose -f docker-compose.test.yml up -d

# 运行集成测试
go test -tags=integration ./...
```

### 部署流程

1. 开发环境
```bash
# 构建开发镜像
docker build -t gongchang-backend:dev .

# 运行开发环境
docker-compose -f docker-compose.dev.yml up -d
```

2. 测试环境
```bash
# 构建测试镜像
docker build -t gongchang-backend:test .

# 运行测试环境
docker-compose -f docker-compose.test.yml up -d
```

3. 生产环境
```bash
# 构建生产镜像
docker build -t gongchang-backend:prod .

# 运行生产环境
docker-compose -f docker-compose.prod.yml up -d
```

### 监控和日志

1. 日志配置
```yaml
# config/logging.yaml
level: info
format: json
output: file
file:
  path: /var/log/gongchang
  maxSize: 100
  maxBackups: 10
  maxAge: 30
```

2. 监控指标
   - API 请求数
   - 响应时间
   - 错误率
   - 数据库连接数
   - 内存使用率

### 安全最佳实践

1. 认证和授权
   - 使用 JWT 进行身份验证
   - 实现基于角色的访问控制
   - 定期轮换密钥

2. 数据安全
   - 使用 HTTPS
   - 加密敏感数据
   - 实现数据备份

3. 代码安全
   - 定期更新依赖
   - 进行安全扫描
   - 实现输入验证

## 故障排除

### 常见问题

1. 数据库连接失败
   - 检查数据库服务是否运行
   - 验证连接参数
   - 检查网络连接

2. API 请求失败
   - 检查服务是否运行
   - 验证请求参数
   - 查看错误日志

3. 性能问题
   - 检查数据库索引
   - 优化查询语句
   - 增加缓存

### 日志分析

1. 错误日志
```bash
# 查看最近错误
tail -f /var/log/gongchang/error.log

# 搜索特定错误
grep "ERROR" /var/log/gongchang/error.log
```

2. 访问日志
```bash
# 查看访问统计
awk '{print $1}' /var/log/gongchang/access.log | sort | uniq -c | sort -nr
```

## 贡献指南

1. 提交 Issue
   - 描述问题
   - 提供复现步骤
   - 添加相关日志

2. 提交 Pull Request
   - 遵循代码规范
   - 添加测试用例
   - 更新文档

3. 代码审查
   - 检查代码质量
   - 验证功能正确性
   - 确保文档更新
