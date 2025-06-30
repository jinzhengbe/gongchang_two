# 工厂职工管理API文档

## 概述

工厂职工管理模块提供完整的职工信息管理功能，包括职工的增删改查、搜索、统计等操作。所有接口都需要工厂角色权限。

## 基础信息

- **基础URL**: `http://localhost:8008/api`
- **认证方式**: Bearer Token
- **权限要求**: 工厂角色 (factory)
- **数据格式**: JSON

## API接口列表

### 1. 创建职工

**接口**: `POST /api/employees`

**描述**: 创建新的职工记录

**请求头**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**:
```json
{
  "name": "张三",
  "position": "生产主管",
  "grade": "高级",
  "work_years": 5,
  "hire_date": "2020-01-15T00:00:00Z",
  "phone": "13800138000",
  "email": "zhangsan@factory.com",
  "department": "生产部",
  "salary": 8000.00,
  "status": "active"
}
```

**响应示例**:
```json
{
  "message": "职工创建成功",
  "employee": {
    "id": 1,
    "name": "张三",
    "position": "生产主管",
    "grade": "高级",
    "work_years": 5,
    "factory_id": "factory_user_id",
    "hire_date": "2020-01-15T00:00:00Z",
    "phone": "13800138000",
    "email": "zhangsan@factory.com",
    "department": "生产部",
    "salary": 8000.00,
    "status": "active",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### 2. 获取职工列表

**接口**: `GET /api/employees`

**描述**: 获取当前工厂的职工列表，支持分页和筛选

**请求头**:
```
Authorization: Bearer <token>
```

**查询参数**:
- `page`: 页码 (默认: 1)
- `page_size`: 每页数量 (默认: 10, 最大: 100)
- `status`: 状态筛选 (active/inactive)
- `department`: 部门筛选

**请求示例**:
```
GET /api/employees?page=1&page_size=10&status=active&department=生产部
```

**响应示例**:
```json
{
  "total": 25,
  "page": 1,
  "page_size": 10,
  "employees": [
    {
      "id": 1,
      "name": "张三",
      "position": "生产主管",
      "grade": "高级",
      "work_years": 5,
      "factory_id": "factory_user_id",
      "hire_date": "2020-01-15T00:00:00Z",
      "phone": "13800138000",
      "email": "zhangsan@factory.com",
      "department": "生产部",
      "salary": 8000.00,
      "status": "active",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z",
      "factory": {
        "user_id": "factory_user_id",
        "factory_name": "示例工厂",
        "contact_person": "李厂长",
        "phone": "13900139000",
        "email": "factory@example.com",
        "address": "广东省深圳市南山区",
        "business_license": "license_number",
        "certification": "certification_info",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    }
  ]
}
```

### 3. 获取单个职工信息

**接口**: `GET /api/employees/{id}`

**描述**: 根据ID获取特定职工的详细信息

**请求头**:
```
Authorization: Bearer <token>
```

**路径参数**:
- `id`: 职工ID

**请求示例**:
```
GET /api/employees/1
```

**响应示例**:
```json
{
  "employee": {
    "id": 1,
    "name": "张三",
    "position": "生产主管",
    "grade": "高级",
    "work_years": 5,
    "factory_id": "factory_user_id",
    "hire_date": "2020-01-15T00:00:00Z",
    "phone": "13800138000",
    "email": "zhangsan@factory.com",
    "department": "生产部",
    "salary": 8000.00,
    "status": "active",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "factory": {
      "user_id": "factory_user_id",
      "factory_name": "示例工厂",
      "contact_person": "李厂长",
      "phone": "13900139000",
      "email": "factory@example.com",
      "address": "广东省深圳市南山区",
      "business_license": "license_number",
      "certification": "certification_info",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

### 4. 更新职工信息

**接口**: `PUT /api/employees/{id}`

**描述**: 更新指定职工的信息

**请求头**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**路径参数**:
- `id`: 职工ID

**请求体** (所有字段都是可选的):
```json
{
  "name": "张三丰",
  "position": "高级生产主管",
  "grade": "专家级",
  "work_years": 6,
  "hire_date": "2020-01-15T00:00:00Z",
  "phone": "13800138001",
  "email": "zhangsanfeng@factory.com",
  "department": "生产部",
  "salary": 9000.00,
  "status": "active"
}
```

**响应示例**:
```json
{
  "message": "职工信息更新成功",
  "employee": {
    "id": 1,
    "name": "张三丰",
    "position": "高级生产主管",
    "grade": "专家级",
    "work_years": 6,
    "factory_id": "factory_user_id",
    "hire_date": "2020-01-15T00:00:00Z",
    "phone": "13800138001",
    "email": "zhangsanfeng@factory.com",
    "department": "生产部",
    "salary": 9000.00,
    "status": "active",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

### 5. 删除职工

**接口**: `DELETE /api/employees/{id}`

**描述**: 删除指定的职工记录

**请求头**:
```
Authorization: Bearer <token>
```

**路径参数**:
- `id`: 职工ID

**请求示例**:
```
DELETE /api/employees/1
```

**响应示例**:
```json
{
  "message": "职工删除成功"
}
```

### 6. 搜索职工

**接口**: `GET /api/employees/search`

**描述**: 根据关键词搜索职工

**请求头**:
```
Authorization: Bearer <token>
```

**查询参数**:
- `q`: 搜索关键词 (必填，支持姓名、职位、部门搜索)
- `page`: 页码 (默认: 1)
- `page_size`: 每页数量 (默认: 10, 最大: 100)

**请求示例**:
```
GET /api/employees/search?q=生产&page=1&page_size=10
```

**响应示例**:
```json
{
  "total": 5,
  "page": 1,
  "page_size": 10,
  "employees": [
    {
      "id": 1,
      "name": "张三",
      "position": "生产主管",
      "grade": "高级",
      "work_years": 5,
      "factory_id": "factory_user_id",
      "hire_date": "2020-01-15T00:00:00Z",
      "phone": "13800138000",
      "email": "zhangsan@factory.com",
      "department": "生产部",
      "salary": 8000.00,
      "status": "active",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

### 7. 获取职工统计

**接口**: `GET /api/employees/statistics`

**描述**: 获取当前工厂的职工统计信息

**请求头**:
```
Authorization: Bearer <token>
```

**请求示例**:
```
GET /api/employees/statistics
```

**响应示例**:
```json
{
  "statistics": {
    "total_employees": 25,
    "active_employees": 23,
    "inactive_employees": 2,
    "average_work_years": 3.5,
    "department_stats": {
      "生产部": 15,
      "质检部": 5,
      "管理部": 3,
      "技术部": 2
    }
  }
}
```

## 数据模型

### FactoryEmployee 职工模型

```json
{
  "id": "uint",
  "name": "string (必填)",
  "position": "string (必填)",
  "grade": "string (可选)",
  "work_years": "int (默认: 0)",
  "factory_id": "string (必填)",
  "hire_date": "date (必填)",
  "phone": "string (可选)",
  "email": "string (可选)",
  "department": "string (可选)",
  "salary": "decimal(10,2) (可选)",
  "status": "string (默认: active)",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### 状态枚举

- `active`: 在职
- `inactive`: 离职

## 错误码说明

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 未授权访问 |
| 403 | 权限不足 (非工厂角色) |
| 404 | 职工不存在 |
| 500 | 服务器内部错误 |

## 使用示例

### cURL示例

```bash
# 1. 登录获取token
curl -X POST "http://localhost:8008/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "factory1",
    "password": "123456"
  }'

# 2. 创建职工
curl -X POST "http://localhost:8008/api/employees" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三",
    "position": "生产主管",
    "grade": "高级",
    "work_years": 5,
    "hire_date": "2020-01-15T00:00:00Z",
    "phone": "13800138000",
    "email": "zhangsan@factory.com",
    "department": "生产部",
    "salary": 8000.00,
    "status": "active"
  }'

# 3. 获取职工列表
curl -X GET "http://localhost:8008/api/employees?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 4. 搜索职工
curl -X GET "http://localhost:8008/api/employees/search?q=生产" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 5. 获取统计信息
curl -X GET "http://localhost:8008/api/employees/statistics" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### JavaScript示例

```javascript
// 创建职工
async function createEmployee(employeeData) {
  const response = await fetch('http://localhost:8008/api/employees', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(employeeData)
  });
  
  return await response.json();
}

// 获取职工列表
async function getEmployees(page = 1, pageSize = 10, status = '', department = '') {
  const params = new URLSearchParams({
    page: page.toString(),
    page_size: pageSize.toString()
  });
  
  if (status) params.append('status', status);
  if (department) params.append('department', department);
  
  const response = await fetch(`http://localhost:8008/api/employees?${params}`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  
  return await response.json();
}

// 搜索职工
async function searchEmployees(keyword, page = 1, pageSize = 10) {
  const params = new URLSearchParams({
    q: keyword,
    page: page.toString(),
    page_size: pageSize.toString()
  });
  
  const response = await fetch(`http://localhost:8008/api/employees/search?${params}`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  
  return await response.json();
}
```

## 注意事项

1. **权限控制**: 所有职工管理接口都需要工厂角色权限
2. **数据隔离**: 每个工厂只能管理自己的职工数据
3. **软删除**: 删除操作采用软删除，数据不会真正从数据库中移除
4. **分页限制**: 每页最大数量限制为100条
5. **搜索功能**: 支持按姓名、职位、部门进行模糊搜索
6. **统计功能**: 提供详细的职工统计信息，包括部门分布等 