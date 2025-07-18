# 更新日志

## [v1.1.0] - 2025-01-XX - 工厂信息编辑API发布

### 🎉 新增功能
- **工厂信息编辑API**: 完整的工厂详细信息编辑功能
- **工厂资料管理**: 支持工厂照片、员工数量、视频等详细信息管理
- **工厂资料查询**: 提供工厂详细资料的查询接口

### 新增API接口

#### 工厂信息管理
- `GET /api/factories/profile` - 获取工厂详细信息
- `PUT /api/factories/profile` - 更新工厂详细信息

### 数据模型更新
- 为 `FactoryProfile` 模型添加新字段：
  - `photos` (JSON) - 工厂照片数组
  - `employee_count` (INT) - 员工数量
  - `videos` (JSON) - 工厂视频数组
- 更新数据库迁移脚本，支持新字段

### 技术实现
- 创建 `FactoryProfileController` 控制器
- 实现 `GetFactoryProfile` 和 `UpdateFactoryProfile` 方法
- 添加数据验证和错误处理
- 支持JSON字段的序列化和反序列化

### 测试验证
- 创建完整的API测试脚本
- 验证登录认证和权限控制
- 测试数据验证和错误处理
- 提供cURL和Flutter使用示例

### 文档更新
- 创建详细的Flutter API使用文档
- 包含数据模型、API调用方法和使用示例
- 提供完整的错误处理说明

### 部署说明
- 支持Docker Compose一键部署
- 包含数据库迁移脚本
- 提供详细的开发环境配置说明

## [v1.0.0] - 2025-07-04 - 完整的工厂和设计师搜索评分系统

### 🎉 重大功能发布
- **工厂搜索系统**: 完整的工厂搜索、筛选和评分功能
- **设计师搜索系统**: 完整的设计师搜索、筛选和评分功能
- **智能搜索建议**: 提供智能搜索提示，提升用户体验
- **专业领域管理**: 支持为工厂和设计师添加专业领域标签
- **评分系统**: 完整的评分和评价功能，支持统计分析

### 新增功能

#### 工厂搜索模块
- `GET /api/factories/search` - 工厂搜索（支持关键词、地区、专业领域、评分筛选）
- `GET /api/factories/search/suggestions` - 搜索建议
- `POST /api/factories/{factory_id}/specialties` - 创建专业领域
- `POST /api/factories/{factory_id}/ratings` - 创建评分
- `GET /api/factories/{factory_id}/ratings` - 获取评分列表
- `GET /api/factories/{factory_id}/ratings/stats` - 获取评分统计

#### 设计师搜索模块
- `GET /api/designers/search` - 设计师搜索（支持关键词、地区、专业领域、评分筛选）
- `GET /api/designers/search/suggestions` - 搜索建议
- `POST /api/designers/{designer_id}/specialties` - 创建专业领域
- `POST /api/designers/{designer_id}/ratings` - 创建评分
- `GET /api/designers/{designer_id}/ratings` - 获取评分列表
- `GET /api/designers/{designer_id}/ratings/stats` - 获取评分统计

### 技术实现
- 创建 `FactorySearchService` 和 `DesignerSearchService` 服务层
- 创建 `FactorySearchController` 和 `DesignerSearchController` 控制器
- 新增 `FactorySpecialty`、`FactoryRating`、`DesignerSpecialty`、`DesignerRating` 数据模型
- 为 `FactoryProfile` 和 `DesignerProfile` 模型添加评分和状态字段
- 创建数据库索引优化搜索性能
- 实现智能搜索建议功能

### 性能优化
- 创建15个数据库索引优化搜索性能
- 支持分页查询，默认每页20条记录
- 使用子查询优化评分筛选和排序
- 实现智能搜索建议，提升用户体验

### 测试验证
- 创建完整的API测试脚本
- 支持基础搜索、高级搜索、搜索建议测试
- 包含错误处理和性能测试
- 提供cURL和JavaScript使用示例
- 验证登录认证、权限控制、数据验证

### 文档更新
- 更新开发文档，添加完整的API文档
- 创建详细的使用示例和测试脚本
- 完善Docker Compose使用说明
- 添加数据库索引和性能优化说明

### 部署说明
- 支持Docker Compose一键部署
- 包含完整的数据库初始化脚本
- 提供详细的开发环境配置说明
- 包含常见问题排查指南

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
- 主机：`192.168.0.10` (外部数据库主机)

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

## [2025-05-16] - 文档补充

### 文档更新
- 在 docs/authentication.md 中补充说明：登录接口参数字段名必须为 user_type，不能为 userType，否则会导致 400 错误，便于前端开发者排查问题。 

## [2025-05-17] - 数据库问题修复与订单数据更新

### 问题修复
- 彻底解决 MySQL 数据目录权限和初始化问题
  - 删除数据卷 `gongchang_mysql_data`，确保 `init.sql` 重新执行
  - 修正数据目录权限为 999:999
  - 重启服务后数据库和表结构成功创建

### 数据更新
- 订单表最新数据状态：
  - 最新订单（id=4）：attachments 字段有值 `["bf742079-9a93-46a4-93c5-7fbd67bc7ee1"]`，其他字段（models、images、videos）为 NULL
  - 其他订单（id=1,2,3）：attachments、models、images、videos 均为 NULL 