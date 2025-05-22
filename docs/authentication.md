# 用户认证文档

## 登录流程

### 请求格式
```json
POST /api/users/login
Content-Type: application/json

{
    "username": "string",
    "password": "string",
    "user_type": "string"  // 必须是 "designer", "factory", 或 "supplier"
}
```

### 响应格式
```json
// 成功响应 (200 OK)
{
    "token": "JWT_TOKEN_STRING",
    "user": {
        "id": "USER_ID",
        "username": "USERNAME",
        "email": "USER_EMAIL",
        "role": "USER_ROLE"
    }
}

// 错误响应 (401 Unauthorized)
{
    "error": "错误信息"
}
```

### 错误类型
- `user not found`: 用户不存在
- `invalid password`: 密码错误
- `invalid user type`: 用户类型无效
- `user type mismatch`: 用户类型与注册类型不匹配

### 注意事项
1. 密码要求：
   - 最小长度：6个字符
   - 最大长度：20个字符

2. 用户类型：
   - designer: 设计师
   - factory: 工厂
   - supplier: 供应商

3. 字段命名要求：
   - 登录接口参数字段名必须为 user_type（下划线），不能为 userType（驼峰）。否则会导致 400 错误。

4. 安全建议：
   - 使用 HTTPS 进行传输
   - 不要在客户端保存密码
   - Token 有效期为 24 小时

## 开发者说明

### 密码存储
- 使用 bcrypt 进行密码加密
- 加密强度：默认 cost=10
- 密码哈希格式：`$2a$10$...`

## 工厂用户注册与登录

### 工厂用户注册
- 接口：`/api/users/register`
- 请求格式：
  ```json
  {
    "username": "your_factory_username",
    "password": "your_password",
    "email": "your_email@example.com",
    "role": "factory"
  }
  ```
- 说明：工厂用户注册到 users 表，role 字段必须为 "factory"。

### 工厂用户登录
- 接口：`/api/users/login`
- 请求格式：
  ```json
  {
    "username": "your_factory_username",
    "password": "your_password",
    "user_type": "factory"
  }
  ```
- 说明：登录时 user_type 必须为 "factory"，否则会返回错误。

## 开发日志
- 2025-05-23: 更新工厂用户注册和登录接口，统一使用 /api/users 路径，并更新文档。 