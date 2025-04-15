# SewingMast Backend API

## 项目简介
SewingMast 是一个基于 Go 语言的后端 API 服务，提供用户管理、产品管理、订单管理等功能。系统采用现代化的技术栈，提供高性能、可扩展的 API 服务。

## 技术栈
- 后端: Go (Gin 框架)
- 数据库: MySQL
- 认证: JWT
- 部署: Docker
- 日志: Zap
- 配置: Viper
- HTTPS: Let's Encrypt

## 环境要求
- Go 1.16+
- MySQL 5.7+
- Docker (可选)
- SSL 证书 (生产环境)

## HTTPS 配置说明

### 1. 证书获取
#### 使用 Let's Encrypt
```bash
# 安装 certbot
sudo apt-get update
sudo apt-get install certbot

# 获取证书
sudo certbot certonly --standalone -d your-domain.com
```

#### 证书文件位置
```bash
/etc/letsencrypt/live/your-domain.com/
├── cert.pem      # 证书文件
├── chain.pem     # 中间证书
├── fullchain.pem # 完整证书链
└── privkey.pem   # 私钥文件
```

### 2. 后端 HTTPS 配置
```go
// main.go
package main

import (
    "github.com/gin-gonic/gin"
    "log"
    "net/http"
)

func main() {
    router := gin.Default()
    
    // 配置 HTTPS
    certFile := "/etc/letsencrypt/live/your-domain.com/fullchain.pem"
    keyFile := "/etc/letsencrypt/live/your-domain.com/privkey.pem"
    
    // 启动 HTTPS 服务
    log.Fatal(http.ListenAndServeTLS(":443", certFile, keyFile, router))
}
```

### 3. Nginx 反向代理配置
```nginx
# /etc/nginx/sites-available/your-domain.com
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl;
    server_name your-domain.com;

    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;
    
    # SSL 配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_timeout 1d;
    ssl_session_cache shared:SSL:50m;
    ssl_session_tickets off;
    ssl_stapling on;
    ssl_stapling_verify on;
    
    # HSTS
    add_header Strict-Transport-Security "max-age=63072000" always;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 4. Docker 配置
```yaml
# docker-compose.yml
version: '3'
services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - SSL_CERT_FILE=/etc/letsencrypt/live/your-domain.com/fullchain.pem
      - SSL_KEY_FILE=/etc/letsencrypt/live/your-domain.com/privkey.pem
    volumes:
      - /etc/letsencrypt:/etc/letsencrypt:ro
```

### 5. 证书自动续期
```bash
# 创建续期脚本
sudo nano /etc/cron.d/certbot-renew

# 添加以下内容
0 0 * * * root certbot renew --quiet --deploy-hook "systemctl reload nginx"
```

### 6. 安全配置
```go
// 强制 HTTPS 重定向
func forceSSL() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Request.TLS == nil {
            url := "https://" + c.Request.Host + c.Request.URL.Path
            c.Redirect(http.StatusMovedPermanently, url)
            c.Abort()
        }
    }
}

// 使用中间件
router.Use(forceSSL())
```

### 7. 环境变量配置
```env
# .env
SSL_ENABLED=true
SSL_CERT_FILE=/etc/letsencrypt/live/your-domain.com/fullchain.pem
SSL_KEY_FILE=/etc/letsencrypt/live/your-domain.com/privkey.pem
```

### 8. 测试 HTTPS
```bash
# 使用 curl 测试
curl -v https://your-domain.com/api/health

# 使用 openssl 测试
openssl s_client -connect your-domain.com:443 -servername your-domain.com
```

## 前端调用说明

### 基础配置
```typescript
// api/config.ts
const API_BASE_URL = process.env.VITE_API_URL || 'https://your-domain.com'

// 请求配置
const requestConfig = {
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
}
```

### 认证配置
```typescript
// 添加认证头
const addAuthHeader = (config: AxiosRequestConfig) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
}

// 请求拦截器
axios.interceptors.request.use(addAuthHeader)
```

### API 调用示例

#### 1. 用户认证
```typescript
// 登录
const login = async (credentials: LoginCredentials) => {
  const response = await axios.post('/api/auth/login', credentials)
  return response.data
}

// 获取用户信息
const getUserProfile = async () => {
  const response = await axios.get('/api/users/profile')
  return response.data
}
```

#### 2. 产品管理
```typescript
// 获取产品列表
const getProducts = async (params: ProductQueryParams) => {
  const response = await axios.get('/api/products', { params })
  return response.data
}

// 创建产品
const createProduct = async (product: ProductCreateDto) => {
  const response = await axios.post('/api/products', product)
  return response.data
}
```

#### 3. 订单管理
```typescript
// 获取订单列表
const getOrders = async (params: OrderQueryParams) => {
  const response = await axios.get('/api/orders', { params })
  return response.data
}

// 创建订单
const createOrder = async (order: OrderCreateDto) => {
  const response = await axios.post('/api/orders', order)
  return response.data
}
```

### 错误处理
```typescript
// 响应拦截器
axios.interceptors.response.use(
  response => response,
  error => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          // 处理未授权
          router.push('/login')
          break
        case 403:
          // 处理权限不足
          showError('权限不足')
          break
        case 404:
          // 处理资源不存在
          showError('资源不存在')
          break
        case 500:
          // 处理服务器错误
          showError('服务器错误')
          break
        default:
          showError(error.response.data.message || '请求失败')
      }
    }
    return Promise.reject(error)
  }
)
```

### 环境变量配置
```env
# .env.development
VITE_API_URL=https://your-domain.com

# .env.production
VITE_API_URL=https://your-production-domain.com
```

### 跨域配置
```typescript
// vite.config.ts
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'https://your-domain.com',
        changeOrigin: true,
        secure: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      }
    }
  }
})
```

## 端口使用说明

### 开发环境
```bash
# 默认端口配置
PORT=8080        # API 服务端口
DB_PORT=3306     # MySQL 数据库端口
```

### 生产环境
```bash
# 生产环境端口配置
PORT=8082        # API 服务端口（建议使用非默认端口）
DB_PORT=3306     # MySQL 数据库端口
```

### 端口映射说明
1. API 服务端口
   - 开发环境：8080
   - 生产环境：8082
   - 可以通过环境变量 PORT 修改
   - 确保防火墙开放相应端口

2. 数据库端口
   - 默认端口：3306
   - 可以通过环境变量 DB_PORT 修改
   - 生产环境建议修改默认端口

3. Docker 部署端口映射
```yaml
# docker-compose.yml 示例
services:
  api:
    ports:
      - "8080:8080"  # 开发环境
      # - "8082:8082"  # 生产环境
  mysql:
    ports:
      - "3306:3306"  # 数据库端口
```

### 端口冲突解决方案
1. 修改 API 服务端口
```bash
# 方法1：通过环境变量
export PORT=8082

# 方法2：修改 .env 文件
PORT=8082
```

2. 修改数据库端口
```bash
# 方法1：通过环境变量
export DB_PORT=3307

# 方法2：修改 .env 文件
DB_PORT=3307
```

3. 检查端口占用
```bash
# Linux/Mac
lsof -i :8080

# Windows
netstat -ano | findstr :8080
```

### 安全建议
1. 生产环境使用非默认端口
2. 配置防火墙规则
3. 使用 HTTPS
4. 定期检查端口使用情况

## 快速开始

### 1. 环境配置
```bash
# 复制环境配置文件
cp backend/.env.example backend/.env

# 修改数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=sewingmast
DB_PASSWORD=your_password
DB_NAME=sewingmast
```

### 2. 数据库设置
```sql
# 创建数据库和用户
CREATE DATABASE sewingmast;
CREATE USER 'sewingmast'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON sewingmast.* TO 'sewingmast'@'localhost';
FLUSH PRIVILEGES;
```

### 3. 安装依赖
```bash
cd backend
go mod download
```

### 4. 运行服务
```bash
# 开发模式
go run main.go

# 生产模式
export GIN_MODE=release
go run main.go
```

### 5. 测试 API
```bash
# 测试健康检查接口
curl http://localhost:8080/api/health | jq

# 预期响应
{
  "message": "API is running",
  "status": "ok"
}
```

## 故障排除

### 1. 数据库连接错误
错误信息：
```
Error 1045 (28000): Access denied for user 'sewingmast'@'localhost'
```
解决方案：
1. 检查数据库用户是否存在
2. 验证密码是否正确
3. 确保用户有正确的权限
4. 检查 MySQL 服务是否运行

### 2. 端口占用
错误信息：
```
Error starting server: listen tcp :8080: bind: address already in use
```
解决方案：
1. 查找占用端口的进程：
```bash
lsof -i :8080
```
2. 终止占用进程：
```bash
kill -9 <PID>
```
3. 或修改服务端口：
```bash
export PORT=8081
```

### 3. 代理警告
警告信息：
```
[WARNING] You trusted all proxies, this is NOT safe.
```
解决方案：
在生产环境中设置可信代理：
```go
router.SetTrustedProxies([]string{"127.0.0.1"})
```

## API 文档

### 基础 URL
- 本地开发: `http://localhost:8080`
- 生产环境: `https://<服务器IP>:8080`

### 健康检查
- 端点: `GET /api/health`
- 响应: 
```json
{
  "status": "ok",
  "message": "API is running"
}
```

## 日志说明

### 开发模式日志
```
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
[GIN-debug] GET    /api/health               --> backend/routes.SetupRouter.func1 (3 handlers)
[GIN] 2025/04/15 - 19:18:42 | 200 | 21.952µs | ::1 | GET "/api/health"
```

### 生产模式日志
```bash
# 设置生产模式
export GIN_MODE=release
```

## 数据库配置
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=sewingmast
DB_PASSWORD=your_password
DB_NAME=sewingmast
```

## 部署说明

### 1. 环境变量配置
```env
PORT=8080
DB_HOST=localhost
DB_PORT=3306
DB_USER=sewingmast
DB_PASSWORD=your_password
DB_NAME=sewingmast
JWT_SECRET=your-secret-key
```

### 2. Docker 部署
```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

## 系统架构

### 后端架构
- RESTful API 设计
- 分层架构（Controller/Service/Repository）
- 中间件机制
- 统一错误处理
- 数据验证
- 日志系统

## 安全说明
1. 生产环境请修改默认的 JWT_SECRET
2. 确保数据库密码强度足够
3. 使用 HTTPS
4. 定期备份数据库
5. 设置可信代理
6. 使用生产模式运行

## 更新日志

### v1.0.1 (2024-04-15)
- 添加故障排除指南
- 更新日志说明
- 添加安全配置说明
- 优化数据库配置

### v1.0.0 (2024-04-15)
- 项目初始化
- 基础功能实现
- 健康检查接口
- 基础架构搭建

## 许可证
MIT License

## 联系方式
- 邮箱：support@sewingmast.com
- 电话：+86-XXX-XXXX-XXXX
