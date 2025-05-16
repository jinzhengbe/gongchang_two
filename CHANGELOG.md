# 更新日志

## [2025-05-16] - 订单系统优化与文件关联修复

### 功能优化
- 优化订单详情接口，确保返回完整的订单信息
- 添加文件关联预加载，解决文件显示问题
- 完善订单创建和更新逻辑

### 问题修复
- 修复订单与文件的关联关系
- 解决订单详情中 attachments 和 images 字段为 null 的问题
- 优化 JSON 字段的处理逻辑，确保返回空数组而不是 null

### 代码改进
- 在 Order 模型中添加与 File 的关联关系
- 优化 OrderService 中的 GetOrderByID 方法
- 完善 OrderController 的响应数据

### 已知问题
- 订单的 created_at 和 updated_at 字段为 NULL
- 文件名称显示有乱码，需要处理字符编码问题
- 文件表中的 order_id 字段为 NULL，需要修复关联关系

## [2025-05-03] - MySQL 数据目录权限问题修复

### 问题修复
- 修复 MySQL 数据目录权限问题
- 数据目录所有者从 `dnsmasq` 和 `systemd-journal` 修改为 MySQL 容器内的 `mysql` 用户（UID 999）
- 执行 `sudo chown -R 999:999 mysql_data mysql_config mysql_logs` 修复权限
- 重启服务后数据库表结构成功创建

### 文档更新
- 在 README.md 中添加重要警告，说明 MySQL 数据目录权限问题
- 添加权限问题导致的后果说明
- 添加权限修复方法说明

### 当前状态
- 数据库表已创建：`users`、`designer_profiles`、`factory_profiles`、`supplier_profiles`、`orders`、`order_progresses`、`order_attachments`、`products`、`files`

## [2025-05-02] - 服务测试与重构计划

### 测试结果
- 后端服务运行正常
- 数据库连接正常
- API 接口测试通过
- 设计师功能测试通过

### 重构计划
计划将项目结构重构为模块化设计：
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

### 新增
- 添加了详细的数据库配置文档
- 在 README.md 中添加了醒目的数据库配置说明
- 完善了数据库初始化脚本的说明
- 添加了数据库问题排查指南
- 添加了数据存储和挂载配置说明
- 添加了数据备份和恢复指南

### 修复
- 修复了数据库初始化问题
- 优化了数据库连接配置
- 改进了数据持久化方案
- 完善了数据卷挂载配置

### 配置说明
#### 数据库配置
- 数据库名称：`gongchang`
- 用户名：`gongchang`
- 密码：`gongchang`
- 端口：`3306`
- 数据存储位置：`/runData/gongChang/mysql_data`

#### 数据存储配置
- 主数据目录：`/runData/gongChang/mysql_data`
- 配置文件目录：`/runData/gongChang/mysql_config`
- 日志文件目录：`/runData/gongChang/mysql_logs`
- 初始化脚本位置：`./init.sql`

#### 挂载配置
```yaml
volumes:
  - ./mysql_data:/var/lib/mysql
  - ./mysql_config:/etc/mysql/conf.d
  - ./mysql_logs:/var/log/mysql
  - ./init.sql:/docker-entrypoint-initdb.d/init.sql
```

## [2024-05-02] - 数据库配置优化

### 新增
- 添加数据库初始化脚本 `init.sql`
- 在 docker-compose.yml 中添加初始化脚本挂载
- 更新 README.md 中的数据库配置说明

### 修复
- 解决数据库自动创建失败的问题
- 添加手动创建数据库的备用方案

### 变更
- 优化数据库初始化流程
- 完善数据库配置文档

### 已知问题
- 数据卷持久化可能导致初始化脚本不执行
- 需要手动执行数据库创建命令的情况 

## [2024-06-09] - 前端容器化与文档完善

### 新增
- 前端（Flutter Web）添加 Dockerfile 和 nginx.conf，支持 Nginx 静态部署
- docker-compose.yml 增加 web 服务，支持一键部署前后端
- README.md 增加前端 Docker 部署与 Nginx 配置说明

### 变更
- 项目支持通过 Docker/Nginx 统一发布前后端，便于生产环境部署 