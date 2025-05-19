# 工厂订单管理系统

## 项目说明
这是一个基于 Go 语言开发的工厂订单管理系统，提供订单管理、文件上传、用户认证等功能。

## 项目结构
```
backend/
├── cmd/                    # 命令行入口
├── internal/              # 内部包
│   ├── auth/             # 认证相关
│   ├── designer/         # 设计师模块
│   ├── factory/          # 工厂模块
│   ├── supplier/         # 供应商模块
│   ├── order/            # 订单模块
│   ├── file/             # 文件模块
│   └── common/           # 公共模块
├── pkg/                  # 可导出的包
├── scripts/             # 脚本文件
├── uploads/             # 上传文件目录
└── docs/                # 文档
```

## 模块说明
### 设计师模块 (internal/designer)
- 设计师档案管理
- 订单创建和管理
- 文件上传和下载
- 订单状态跟踪

### 工厂模块 (internal/factory)
- 工厂档案管理
- 订单接收和处理
- 生产进度更新
- 质量检查记录

### 供应商模块 (internal/supplier)
- 供应商档案管理
- 材料供应管理
- 订单材料跟踪
- 库存管理

### 订单模块 (internal/order)
- 订单创建和编辑
- 订单状态管理
- 订单查询和统计
- 订单进度跟踪

### 文件模块 (internal/file)
- 文件上传和下载
- 文件类型验证
- 文件存储管理
- 文件访问控制

### 认证模块 (internal/auth)
- 用户认证
- 权限管理
- JWT token 管理
- 会话管理

## ⚠️ 重要：数据库配置说明
数据库配置是系统正常运行的关键，请确保以下配置正确：

### ⚠️ 重要警告：MySQL 数据目录权限
MySQL 数据目录的所有者权限必须正确设置，否则会导致数据库无法正常启动或表结构无法创建：
- 数据目录所有者必须是 MySQL 容器内的 `mysql` 用户（UID 999）
- 如果发现数据目录所有者是 `dnsmasq` 或 `systemd-journal`，需要执行以下命令修复：
  ```bash
  sudo chown -R 999:999 mysql_data mysql_config mysql_logs
  ```
- 权限问题会导致：
  - 数据库无法正常启动
  - 表结构无法创建
  - 数据无法持久化
  - 服务无法正常运行

### 数据库基本信息
- 数据库名称：`gongchang`
- 用户名：`gongchang`
- 密码：`gongchang`
- 端口：`3306`
- 主机：`mysql` (Docker 网络中的服务名)

### 数据存储配置
#### 存储位置
- 主数据目录：`/runData/gongChang/mysql_data`
- 配置文件目录：`/runData/gongChang/mysql_config`
- 日志文件目录：`/runData/gongChang/mysql_logs`

#### 挂载配置
在 `docker-compose.yml` 中的挂载配置：
```yaml
volumes:
  - ./mysql_data:/var/lib/mysql
  - ./mysql_config:/etc/mysql/conf.d
  - ./mysql_logs:/var/log/mysql
  - ./init.sql:/docker-entrypoint-initdb.d/init.sql
```

#### 数据持久化说明
1. 数据卷类型：本地目录挂载
2. 数据备份：建议定期备份 `/runData/gongChang/mysql_data` 目录
3. 权限设置：确保挂载目录有正确的读写权限
4. 数据迁移：可以通过复制整个 `mysql_data` 目录进行数据迁移

### 初始化脚本
数据库初始化脚本位于 `./init.sql`，包含以下内容：
```sql
CREATE DATABASE IF NOT EXISTS gongchang;
GRANT ALL PRIVILEGES ON gongchang.* TO 'gongchang'@'%';
FLUSH PRIVILEGES;
```

### 常见问题解决
#### 数据库未创建问题
1. 停止服务：`docker-compose down`
2. 删除数据卷：`docker volume rm gongchang_mysql_data`
3. 重新启动：`docker-compose up -d`

#### 手动创建数据库
```bash
docker-compose exec mysql mysql -uroot -proot -e "CREATE DATABASE IF NOT EXISTS gongchang; GRANT ALL PRIVILEGES ON gongchang.* TO 'gongchang'@'%'; FLUSH PRIVILEGES;"
```

#### 数据备份与恢复
1. 备份数据：
```bash
# 备份整个数据目录
tar -czvf mysql_backup.tar.gz /runData/gongChang/mysql_data
```

2. 恢复数据：
```bash
# 停止服务
docker-compose down

# 恢复数据
tar -xzvf mysql_backup.tar.gz -C /runData/gongChang/

# 重新启动服务
docker-compose up -d
```

### 数据库监控
系统内置了数据库性能监控功能，可以实时监控数据库连接状态和查询性能：

#### 连接池监控
- 最大连接数
- 当前打开连接数
- 使用中连接数
- 空闲连接数
- 等待连接数
- 等待时间
- 连接关闭统计

#### 慢查询监控
- 可配置慢查询阈值
- 自动记录超过阈值的查询
- 记录查询执行时间
- 记录完整 SQL 语句

使用示例：
```go
// 启动数据库监控，每5秒记录一次统计信息
MonitorDatabase(db, 5*time.Second)

// 设置慢查询监控，记录执行时间超过1秒的查询
MonitorSlowQueries(db, 1*time.Second)
```

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
- 前端 Web: http://localhost:80

## 前端（Flutter Web）Docker 部署说明

1. 进入 web 目录，构建 Docker 镜像：
```bash
cd web
# 构建镜像（可自定义 tag）
docker build -t gongchang-web .
```

2. 启动服务（推荐使用 docker-compose）：
```bash
cd ..
docker-compose up -d web
```

3. 访问前端页面：
- http://localhost:80
- 若部署到服务器，配置域名（如 aneworder.com）指向服务器公网 IP

4. Nginx 配置说明：
- 已内置于 web/nginx.conf，支持 Flutter Web SPA 路由
- 静态资源缓存 30 天

5. 常见问题
- 若端口冲突，请修改 docker-compose.yml 中 web 服务的端口映射
- 若需 HTTPS，请在服务器上配置反向代理或修改 Nginx 配置

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

## 开发指南

### 后端开发

1. 环境配置
```bash
cd backend
go mod download
```

2. 运行测试
```bash
go test ./...
```

3. 启动服务
```bash
go run main.go
```

### 前端开发

1. 环境配置
```bash
cd frontend
flutter pub get
```

2. 运行测试
```bash
flutter test
```

3. 启动开发服务器
```bash
flutter run
```

## 部署指南

1. 构建镜像
```bash
docker-compose build
```

2. 启动服务
```bash
docker-compose up -d
```

3. 查看日志
```bash
docker-compose logs -f
```

## 维护说明

1. 数据库备份
   - 自动备份：每天凌晨 2 点执行
   - 手动备份：执行 `./scripts/backup/backup.sh`
   - 主从复制：实时同步数据

2. 服务监控
   - 监控服务日志
   - 检查服务状态
   - 监控主从同步状态

3. 依赖更新
   - 定期更新 Go 依赖
   - 更新 Flutter 依赖

4. 安全维护
   - 定期检查安全漏洞
   - 更新依赖包修复漏洞
   - 定期更换数据库密码
   - 定期更新 JWT secret

## 故障排除

### Token 验证失败
1. 检查 token 格式是否正确
2. 确认 token 未过期
3. 验证 JWT secret 是否正确
4. 检查服务器日志获取详细错误信息

### 数据库连接问题
1. 检查数据库服务是否运行
2. 验证连接信息是否正确
3. 检查网络连接

### 文件上传问题
1. 检查文件大小限制
2. 验证文件类型是否允许
3. 检查存储空间是否充足

## 贡献指南
1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证
MIT License
