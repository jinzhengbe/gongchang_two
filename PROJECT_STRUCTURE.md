# 项目结构说明

## 根目录
- `README.md`: 项目说明文档
- `docker-compose.yml`: Docker 容器编排配置文件
- `.env.example`: 环境变量示例文件
- `.env`: 实际环境变量配置文件
- `LICENSE`: 项目许可证文件

## backend 目录
后端服务，使用 Go 语言开发

### 主要目录
- `controllers/`: 控制器层，处理 HTTP 请求
  - `user.go`: 用户管理接口
  - `order.go`: 订单管理接口
  - `product.go`: 产品管理接口
  - `file.go`: 文件处理接口
  - `auth.go`: 认证相关接口
- `services/`: 服务层，实现业务逻辑
  - `user.go`: 用户相关业务逻辑
  - `order.go`: 订单相关业务逻辑
  - `product.go`: 产品相关业务逻辑
- `models/`: 数据模型定义
- `routes/`: 路由配置
- `middleware/`: 中间件
- `config/`: 配置文件
- `utils/`: 工具函数
- `database/`: 数据库相关操作
- `api/`: API 接口定义

### 重要文件
- `main.go`: 程序入口文件，包含：
  - 配置加载
  - 数据库连接初始化
  - 数据库表自动迁移
  - 测试数据初始化
  - 路由设置
  - 服务器启动
- `go.mod`: Go 模块依赖管理
- `go.sum`: Go 模块依赖校验
- `Dockerfile`: 后端服务容器构建文件

## 部署相关
- `docker-compose.yml`: 用于部署后端服务
- `Dockerfile`: 用于构建后端服务容器 