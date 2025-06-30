# 工厂订单进度管理API使用指南

## 概述

工厂订单进度管理模块用于跟踪和管理订单的生产进度。工厂可以为订单创建多个进度记录，记录不同阶段的完成情况。

## 数据库表结构

### order_progress 表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| `id` | bigint unsigned | 主键，自增 |
| `order_id` | bigint unsigned | 订单ID，外键关联orders表 |
| `factory_id` | varchar(191) | 工厂ID |
| `progress_type` | varchar(50) | 进度类型 |
| `percentage` | int | 完成百分比(0-100) |
| `status` | varchar(50) | 进度状态 |
| `description` | text | 进度描述 |
| `estimated_completion_time` | datetime(3) | 预计完成时间 |
| `actual_completion_time` | datetime(3) | 实际完成时间 |
| `creator_id` | varchar(191) | 创建者ID |
| `created_at` | datetime(3) | 创建时间 |
| `updated_at` | datetime(3) | 更新时间 |
| `deleted_at` | datetime(3) | 删除时间（软删除） |

## 进度类型说明

### ProgressType

- `design`: 设计阶段
- `material`: 材料准备
- `production`: 生产阶段
- `quality`: 质检阶段
- `packaging`: 包装阶段
- `shipping`: 发货阶段
- `custom`: 自定义阶段

## 进度状态说明

### ProgressStatus

- `not_started`: 未开始
- `in_progress`: 进行中
- `completed`: 已完成
- `delayed`: 延期
- `on_hold`: 暂停

## API 接口详细说明

### 1. 创建进度记录

**接口地址：** `POST /api/orders/{orderId}/progress`

**请求头：**
```
Content-Type: application/json
Authorization: Bearer {token}
```

**请求体：**
```json
{
  "order_id": 123,
  "factory_id": "factory_user_id",
  "progress_type": "production",
  "percentage": 50,
  "status": "in_progress",
  "description": "生产进度过半，质量良好",
  "estimated_completion_time": "2025-07-15T10:00:00Z",
  "actual_completion_time": null,
  "creator_id": "factory_user_id"
}
```

**响应：**
```json
{
  "id": 1,
  "order_id": 123,
  "factory_id": "factory_user_id",
  "progress_type": "production",
  "percentage": 50,
  "status": "in_progress",
  "description": "生产进度过半，质量良好",
  "estimated_completion_time": "2025-07-15T10:00:00Z",
  "actual_completion_time": null,
  "creator_id": "factory_user_id",
  "created_at": "2025-06-29T10:30:00Z",
  "updated_at": "2025-06-29T10:30:00Z"
}
```

**权限要求：** 只有工厂用户可以创建进度记录

### 2. 获取订单进度列表

**接口地址：** `GET /api/orders/{orderId}/progress`

**请求头：**
```
Authorization: Bearer {token}
```

**响应：**
```json
[
  {
    "id": 1,
    "order_id": 123,
    "factory_id": "factory_user_id",
    "progress_type": "production",
    "percentage": 50,
    "status": "in_progress",
    "description": "生产进度过半，质量良好",
    "estimated_completion_time": "2025-07-15T10:00:00Z",
    "actual_completion_time": null,
    "creator_id": "factory_user_id",
    "created_at": "2025-06-29T10:30:00Z",
    "updated_at": "2025-06-29T10:30:00Z",
    "order": {
      "id": 123,
      "title": "订单标题",
      "description": "订单描述"
    },
    "factory": {
      "user_id": "factory_user_id",
      "company_name": "工厂名称"
    }
  }
]
```

### 3. 更新进度记录

**接口地址：** `PUT /api/orders/{orderId}/progress/{progressId}`

**请求头：**
```
Content-Type: application/json
Authorization: Bearer {token}
```

**请求体：**
```json
{
  "progress_type": "production",
  "percentage": 75,
  "status": "in_progress",
  "description": "生产进度75%，即将完成",
  "estimated_completion_time": "2025-07-10T10:00:00Z",
  "actual_completion_time": null
}
```

**响应：**
```json
{
  "id": 1,
  "order_id": 123,
  "factory_id": "factory_user_id",
  "progress_type": "production",
  "percentage": 75,
  "status": "in_progress",
  "description": "生产进度75%，即将完成",
  "estimated_completion_time": "2025-07-10T10:00:00Z",
  "actual_completion_time": null,
  "creator_id": "factory_user_id",
  "created_at": "2025-06-29T10:30:00Z",
  "updated_at": "2025-06-29T11:00:00Z"
}
```

**权限要求：** 只有工厂用户可以更新自己创建的进度记录

### 4. 删除进度记录

**接口地址：** `DELETE /api/orders/{orderId}/progress/{progressId}`

**请求头：**
```
Authorization: Bearer {token}
```

**响应：**
```json
{
  "message": "进度记录删除成功"
}
```

**权限要求：** 只有工厂用户可以删除自己创建的进度记录

### 5. 获取工厂进度列表

**接口地址：** `GET /api/factories/{factoryId}/progress?page=1&pageSize=10`

**请求头：**
```
Authorization: Bearer {token}
```

**查询参数：**
- `page`: 页码（默认1）
- `pageSize`: 每页数量（默认10）

**响应：**
```json
{
  "total": 25,
  "page": 1,
  "page_size": 10,
  "progress": [
    {
      "id": 1,
      "order_id": 123,
      "factory_id": "factory_user_id",
      "progress_type": "production",
      "percentage": 75,
      "status": "in_progress",
      "description": "生产进度75%，即将完成",
      "estimated_completion_time": "2025-07-10T10:00:00Z",
      "actual_completion_time": null,
      "creator_id": "factory_user_id",
      "created_at": "2025-06-29T10:30:00Z",
      "updated_at": "2025-06-29T11:00:00Z",
      "order": {
        "id": 123,
        "title": "订单标题"
      },
      "factory": {
        "user_id": "factory_user_id",
        "company_name": "工厂名称"
      }
    }
  ]
}
```

**权限要求：** 只能查看自己工厂的进度记录

### 6. 获取进度统计信息

**接口地址：** `GET /api/factories/{factoryId}/progress-statistics`

**请求头：**
```
Authorization: Bearer {token}
```

**响应：**
```json
{
  "not_started": 5,
  "in_progress": 10,
  "completed": 8,
  "delayed": 2,
  "on_hold": 1
}
```

**权限要求：** 只能查看自己工厂的统计信息

## 业务规则

1. **权限控制**：只有工厂用户可以创建、更新、删除进度记录
2. **数据一致性**：进度记录必须属于指定的订单
3. **身份验证**：工厂只能以自己的身份进行进度管理
4. **百分比范围**：完成百分比必须在0-100之间
5. **时间逻辑**：实际完成时间不能早于创建时间

## 错误码

| 错误码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 权限不足 |
| 404 | 进度记录不存在 |
| 422 | 数据验证失败 |
| 500 | 服务器内部错误 |

## 使用示例

### 工厂进度管理流程

1. **创建进度记录**
```bash
curl -X POST "https://aneworders.com/api/orders/123/progress" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "order_id": 123,
    "factory_id": "factory_user_id",
    "progress_type": "production",
    "percentage": 50,
    "status": "in_progress",
    "description": "生产进度过半，质量良好",
    "estimated_completion_time": "2025-07-15T10:00:00Z",
    "creator_id": "factory_user_id"
  }'
```

2. **查看订单进度**
```bash
curl -X GET "https://aneworders.com/api/orders/123/progress" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

3. **更新进度**
```bash
curl -X PUT "https://aneworders.com/api/orders/123/progress/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "percentage": 75,
    "description": "生产进度75%，即将完成"
  }'
```

4. **查看工厂所有进度**
```bash
curl -X GET "https://aneworders.com/api/factories/factory_user_id/progress?page=1&pageSize=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

5. **查看进度统计**
```bash
curl -X GET "https://aneworders.com/api/factories/factory_user_id/progress-statistics" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

6. **删除进度记录**
```bash
curl -X DELETE "https://aneworders.com/api/orders/123/progress/1" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 前端集成示例

### React组件示例

```jsx
import React, { useState, useEffect } from 'react';

// 进度显示组件
const ProgressDisplay = ({ orderId }) => {
  const [progress, setProgress] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchProgress();
  }, [orderId]);

  const fetchProgress = async () => {
    try {
      const response = await fetch(`/api/orders/${orderId}/progress`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
      const data = await response.json();
      setProgress(data);
    } catch (error) {
      console.error('获取进度失败:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <div>加载中...</div>;

  return (
    <div className="progress-container">
      {progress.map(item => (
        <div key={item.id} className="progress-item">
          <div className="progress-header">
            <span className="progress-type">{item.progress_type}</span>
            <span className="progress-status">{item.status}</span>
          </div>
          <div className="progress-bar">
            <div 
              className="progress-fill" 
              style={{width: `${item.percentage || 0}%`}}
            />
          </div>
          <div className="progress-description">{item.description}</div>
          <div className="progress-time">
            预计完成: {item.estimated_completion_time}
          </div>
        </div>
      ))}
    </div>
  );
};

// 进度创建表单
const ProgressForm = ({ orderId, onSuccess }) => {
  const [formData, setFormData] = useState({
    progress_type: 'production',
    percentage: 0,
    status: 'not_started',
    description: '',
    estimated_completion_time: ''
  });

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch(`/api/orders/${orderId}/progress`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({
          order_id: orderId,
          factory_id: localStorage.getItem('user_id'),
          ...formData,
          creator_id: localStorage.getItem('user_id')
        })
      });
      
      if (response.ok) {
        onSuccess();
        setFormData({
          progress_type: 'production',
          percentage: 0,
          status: 'not_started',
          description: '',
          estimated_completion_time: ''
        });
      }
    } catch (error) {
      console.error('创建进度失败:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="progress-form">
      <div className="form-group">
        <label>进度类型:</label>
        <select 
          value={formData.progress_type}
          onChange={(e) => setFormData({...formData, progress_type: e.target.value})}
        >
          <option value="design">设计阶段</option>
          <option value="material">材料准备</option>
          <option value="production">生产阶段</option>
          <option value="quality">质检阶段</option>
          <option value="packaging">包装阶段</option>
          <option value="shipping">发货阶段</option>
          <option value="custom">自定义阶段</option>
        </select>
      </div>
      
      <div className="form-group">
        <label>完成百分比:</label>
        <input 
          type="number" 
          min="0" 
          max="100"
          value={formData.percentage}
          onChange={(e) => setFormData({...formData, percentage: parseInt(e.target.value)})}
        />
      </div>
      
      <div className="form-group">
        <label>状态:</label>
        <select 
          value={formData.status}
          onChange={(e) => setFormData({...formData, status: e.target.value})}
        >
          <option value="not_started">未开始</option>
          <option value="in_progress">进行中</option>
          <option value="completed">已完成</option>
          <option value="delayed">延期</option>
          <option value="on_hold">暂停</option>
        </select>
      </div>
      
      <div className="form-group">
        <label>描述:</label>
        <textarea 
          value={formData.description}
          onChange={(e) => setFormData({...formData, description: e.target.value})}
        />
      </div>
      
      <div className="form-group">
        <label>预计完成时间:</label>
        <input 
          type="datetime-local"
          value={formData.estimated_completion_time}
          onChange={(e) => setFormData({...formData, estimated_completion_time: e.target.value})}
        />
      </div>
      
      <button type="submit">创建进度</button>
    </form>
  );
};

export { ProgressDisplay, ProgressForm };
```

### Vue.js组件示例

```vue
<template>
  <div class="progress-management">
    <!-- 进度显示 -->
    <div class="progress-list">
      <h3>订单进度</h3>
      <div v-for="item in progressList" :key="item.id" class="progress-item">
        <div class="progress-header">
          <span class="progress-type">{{ getProgressTypeName(item.progress_type) }}</span>
          <span class="progress-status">{{ getStatusName(item.status) }}</span>
        </div>
        <div class="progress-bar">
          <div 
            class="progress-fill" 
            :style="{width: (item.percentage || 0) + '%'}"
          ></div>
        </div>
        <div class="progress-description">{{ item.description }}</div>
        <div class="progress-actions">
          <button @click="editProgress(item)">编辑</button>
          <button @click="deleteProgress(item.id)">删除</button>
        </div>
      </div>
    </div>

    <!-- 进度表单 -->
    <div class="progress-form">
      <h3>{{ isEditing ? '编辑进度' : '创建进度' }}</h3>
      <form @submit.prevent="submitProgress">
        <div class="form-group">
          <label>进度类型:</label>
          <select v-model="form.progress_type">
            <option value="design">设计阶段</option>
            <option value="material">材料准备</option>
            <option value="production">生产阶段</option>
            <option value="quality">质检阶段</option>
            <option value="packaging">包装阶段</option>
            <option value="shipping">发货阶段</option>
            <option value="custom">自定义阶段</option>
          </select>
        </div>
        
        <div class="form-group">
          <label>完成百分比:</label>
          <input 
            type="number" 
            min="0" 
            max="100"
            v-model="form.percentage"
          />
        </div>
        
        <div class="form-group">
          <label>状态:</label>
          <select v-model="form.status">
            <option value="not_started">未开始</option>
            <option value="in_progress">进行中</option>
            <option value="completed">已完成</option>
            <option value="delayed">延期</option>
            <option value="on_hold">暂停</option>
          </select>
        </div>
        
        <div class="form-group">
          <label>描述:</label>
          <textarea v-model="form.description"></textarea>
        </div>
        
        <div class="form-group">
          <label>预计完成时间:</label>
          <input 
            type="datetime-local"
            v-model="form.estimated_completion_time"
          />
        </div>
        
        <button type="submit">{{ isEditing ? '更新' : '创建' }}</button>
        <button type="button" @click="resetForm">取消</button>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ProgressManagement',
  props: {
    orderId: {
      type: Number,
      required: true
    }
  },
  data() {
    return {
      progressList: [],
      isEditing: false,
      editingId: null,
      form: {
        progress_type: 'production',
        percentage: 0,
        status: 'not_started',
        description: '',
        estimated_completion_time: ''
      }
    }
  },
  mounted() {
    this.fetchProgress();
  },
  methods: {
    async fetchProgress() {
      try {
        const response = await fetch(`/api/orders/${this.orderId}/progress`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });
        this.progressList = await response.json();
      } catch (error) {
        console.error('获取进度失败:', error);
      }
    },
    
    async submitProgress() {
      try {
        const url = this.isEditing 
          ? `/api/orders/${this.orderId}/progress/${this.editingId}`
          : `/api/orders/${this.orderId}/progress`;
        
        const method = this.isEditing ? 'PUT' : 'POST';
        const body = this.isEditing ? this.form : {
          order_id: this.orderId,
          factory_id: localStorage.getItem('user_id'),
          ...this.form,
          creator_id: localStorage.getItem('user_id')
        };
        
        const response = await fetch(url, {
          method,
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          },
          body: JSON.stringify(body)
        });
        
        if (response.ok) {
          this.fetchProgress();
          this.resetForm();
        }
      } catch (error) {
        console.error('提交进度失败:', error);
      }
    },
    
    editProgress(item) {
      this.isEditing = true;
      this.editingId = item.id;
      this.form = {
        progress_type: item.progress_type,
        percentage: item.percentage || 0,
        status: item.status,
        description: item.description,
        estimated_completion_time: item.estimated_completion_time || ''
      };
    },
    
    async deleteProgress(id) {
      if (confirm('确定要删除这个进度记录吗？')) {
        try {
          const response = await fetch(`/api/orders/${this.orderId}/progress/${id}`, {
            method: 'DELETE',
            headers: {
              'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
          });
          
          if (response.ok) {
            this.fetchProgress();
          }
        } catch (error) {
          console.error('删除进度失败:', error);
        }
      }
    },
    
    resetForm() {
      this.isEditing = false;
      this.editingId = null;
      this.form = {
        progress_type: 'production',
        percentage: 0,
        status: 'not_started',
        description: '',
        estimated_completion_time: ''
      };
    },
    
    getProgressTypeName(type) {
      const types = {
        design: '设计阶段',
        material: '材料准备',
        production: '生产阶段',
        quality: '质检阶段',
        packaging: '包装阶段',
        shipping: '发货阶段',
        custom: '自定义阶段'
      };
      return types[type] || type;
    },
    
    getStatusName(status) {
      const statuses = {
        not_started: '未开始',
        in_progress: '进行中',
        completed: '已完成',
        delayed: '延期',
        on_hold: '暂停'
      };
      return statuses[status] || status;
    }
  }
}
</script>

<style scoped>
.progress-management {
  padding: 20px;
}

.progress-item {
  border: 1px solid #ddd;
  padding: 15px;
  margin-bottom: 15px;
  border-radius: 5px;
}

.progress-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}

.progress-bar {
  width: 100%;
  height: 20px;
  background-color: #f0f0f0;
  border-radius: 10px;
  overflow: hidden;
  margin-bottom: 10px;
}

.progress-fill {
  height: 100%;
  background-color: #4CAF50;
  transition: width 0.3s ease;
}

.progress-form {
  margin-top: 30px;
  border-top: 1px solid #ddd;
  padding-top: 20px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.form-group textarea {
  height: 80px;
  resize: vertical;
}

button {
  padding: 10px 20px;
  margin-right: 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button[type="submit"] {
  background-color: #4CAF50;
  color: white;
}

button[type="button"] {
  background-color: #f44336;
  color: white;
}
</style>
```

## 总结

工厂订单进度管理API提供了完整的进度跟踪功能，包括：

1. **创建进度记录** - 工厂可以为订单创建多个进度记录
2. **查看进度列表** - 查看订单的所有进度记录
3. **更新进度** - 更新进度信息，如完成百分比、状态等
4. **删除进度** - 删除不需要的进度记录
5. **进度统计** - 查看工厂的整体进度统计信息

这些API已经完整实现并测试通过，可以直接用于前端开发。记得在使用时进行适当的权限控制和错误处理。 