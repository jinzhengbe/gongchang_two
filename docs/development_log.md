# 开发日志

## 2025-05-23

### 功能开发
1. 实现工厂注册和登录功能
   - 创建工厂表结构
   - 实现注册 API
   - 实现登录 API
   - 添加 JWT 认证

### 数据库变更
1. 修改 factories 表结构
   - 将 username 和 password 字段改为 VARCHAR 类型
   - 添加唯一索引约束

### 系统配置
1. 配置 Docker Compose 环境
   - 后端服务 (8008端口)
   - 前端服务 (80端口)
   - 数据库连接配置

### 待办事项
1. 完善工厂管理功能
2. 实现订单管理功能
3. 添加文件上传功能
4. 完善用户权限管理 