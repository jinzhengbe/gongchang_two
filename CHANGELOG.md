# 更新日志

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