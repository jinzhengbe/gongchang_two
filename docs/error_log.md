# 错误日志

## 2024-03-21 错误记录

### 1. 用户角色类型转换错误
- 文件：`backend/controllers/user.go`
- 错误：`cannot use req.Role (variable of type string) as models.UserRole value in struct literal`
- 位置：第35行
- 原因：`RegisterRequest` 中的 `Role` 字段是 `string` 类型，而 `User` 结构体中的 `Role` 字段是 `models.UserRole` 类型
- 修复：在创建用户时进行类型转换 `models.UserRole(req.Role)`

### 2. 密码哈希函数未定义错误
- 文件：`backend/controllers/user.go`
- 错误：`undefined: services.HashPassword`
- 位置：第32行
- 原因：`HashPassword` 函数未在 `services` 包中定义
- 修复：需要在 `services` 包中实现 `HashPassword` 函数

### 3. 数据库连接问题
- 文件：`backend/config/config.yaml`
- 问题：数据库名称配置错误
- 原配置：`dbname: "gongchang"`
- 修改为：`dbname: "gongChang"`

### 4. 测试数据问题
- 文件：`backend/database/init.go`
- 问题：测试订单中的 `FactoryID` 设置错误
- 原值：`FactoryID: 1`
- 修改为：`FactoryID: 2` (factory1 的 ID)

## 待解决问题
1. 实现 `services.HashPassword` 函数
2. 确保数据库连接正确
3. 验证测试数据是否正确创建
4. 测试公开订单接口是否正常工作 