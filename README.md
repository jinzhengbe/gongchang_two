# 工厂订单管理系统

## 项目说明
这是一个基于 Go 语言开发的工厂订单管理系统，提供订单管理、文件上传、用户认证等功能。

## 环境要求
- Go 1.16+
- MySQL 8.0+
- Docker & Docker Compose

## 重要说明：数据库配置
数据库配置已在 docker-compose.yml 中完成自动化设置：
- 数据库名称：gongchang（自动创建）
- 数据存储位置：/runData/gongChang/mysql_data
- 用户名：gongchang
- 密码：gongchang
- 端口：3306

初始化脚本：
- 位置：./init.sql
- 内容：
  ```sql
  CREATE DATABASE IF NOT EXISTS gongchang;
  GRANT ALL PRIVILEGES ON gongchang.* TO 'gongchang'@'%';
  FLUSH PRIVILEGES;
  ```

注意：如果遇到数据库未创建的问题，可以手动执行以下命令：
```bash
docker-compose exec mysql mysql -uroot -proot -e "CREATE DATABASE IF NOT EXISTS gongchang; GRANT ALL PRIVILEGES ON gongchang.* TO 'gongchang'@'%'; FLUSH PRIVILEGES;"
```

无需手动创建数据库，首次启动时会自动完成配置。

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
├── scripts/          # 脚本目录
│   ├── backup/      # 备份脚本
│   └── deploy/      # 部署脚本
└── README.md         # 项目文档
```

## 认证与授权

### JWT Token 认证
系统使用 JWT (JSON Web Token) 进行身份认证。每个请求都需要在 Authorization 头中携带 token。

#### Token 格式
```
Authorization: Bearer <token>
```

#### Token 结构
```json
{
  "user_id": "string",
  "role": "string",
  "exp": "number",  // 过期时间
  "iat": "number"   // 签发时间
}
```

#### 错误处理
- 401 Unauthorized: token 无效或过期
- 403 Forbidden: 权限不足
- 400 Bad Request: 请求格式错误

#### 客户端处理
1. 登录获取 token
```javascript
// 登录请求
POST /api/auth/login
{
  "username": "string",
  "password": "string"
}

// 响应
{
  "token": "string",
  "user": {
    "id": "string",
    "username": "string",
    "role": "string"
  }
}
```

2. 存储 token
```javascript
// 保存 token
localStorage.setItem('token', token);

// 获取 token
const token = localStorage.getItem('token');
```

3. 请求拦截器
```javascript
// 添加请求头
axios.interceptors.request.use(config => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// 处理 401 错误
axios.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      // 清除 token 并重定向到登录页
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);
```

### 角色权限
- designer: 设计师，可以创建订单
- factory: 工厂，可以处理订单
- supplier: 供应商，可以查看相关订单

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
