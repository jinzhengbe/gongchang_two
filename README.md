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
├── backend/           # 后端代码
│   ├── controllers/  # 控制器
│   ├── models/       # 数据模型
│   ├── routes/       # 路由配置
│   ├── services/     # 业务逻辑
│   └── main.go       # 入口文件
├── docs/             # 文档
│   └── order_api.md  # 订单 API 文档
└── docker-compose.yml # Docker 配置
```

## API 文档
- [订单管理 API](./docs/order_api.md)

## 开发指南

### 数据库迁移
```bash
# 进入 MySQL 容器
docker-compose exec mysql mysql -uroot -p123456 gongchang

# 执行 SQL 文件
source /path/to/schema.sql
```

### 测试数据
系统启动时会自动创建测试用户：
- 设计师: designer1
- 工厂: factory1
- 供应商: supplier1

### 开发流程
1. 创建新分支
2. 实现功能
3. 编写测试
4. 提交代码
5. 创建 Pull Request

## 部署说明
1. 确保服务器已安装 Docker 和 Docker Compose
2. 配置 SSL 证书（如果需要 HTTPS）
3. 设置环境变量
4. 使用 Docker Compose 部署

## 维护说明
- 定期备份数据库
- 监控服务日志
- 更新依赖包
- 定期检查安全漏洞

## 贡献指南
1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证
MIT License
