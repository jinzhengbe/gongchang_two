<template>
  <div class="create-order-container">
    <div class="create-order-header">
      <h2>创建新订单</h2>
      <button class="btn-back" @click="goBack">返回列表</button>
    </div>

    <div class="create-order-form">
      <form @submit.prevent="handleSubmit">
        <div class="form-section">
          <h3>基本信息</h3>
          <div class="form-grid">
            <div class="form-group">
              <label for="customerName">客户名称</label>
              <input
                id="customerName"
                v-model="form.customerName"
                type="text"
                required
                placeholder="请输入客户名称"
              />
            </div>

            <div class="form-group">
              <label for="productName">产品名称</label>
              <input
                id="productName"
                v-model="form.productName"
                type="text"
                required
                placeholder="请输入产品名称"
              />
            </div>

            <div class="form-group">
              <label for="quantity">数量</label>
              <input
                id="quantity"
                v-model.number="form.quantity"
                type="number"
                required
                min="1"
                placeholder="请输入数量"
              />
            </div>

            <div class="form-group">
              <label for="totalPrice">总价</label>
              <input
                id="totalPrice"
                v-model.number="form.totalPrice"
                type="number"
                required
                min="0"
                step="0.01"
                placeholder="请输入总价"
              />
            </div>
          </div>
        </div>

        <div class="form-section">
          <h3>生产要求</h3>
          <div class="form-group">
            <label for="requirements">具体要求</label>
            <textarea
              id="requirements"
              v-model="form.requirements"
              rows="4"
              placeholder="请输入生产要求"
            ></textarea>
          </div>
        </div>

        <div class="form-section">
          <h3>文件上传</h3>
          <div class="form-group">
            <label for="files">相关文件</label>
            <input
              id="files"
              type="file"
              multiple
              @change="handleFileChange"
            />
            <div class="file-list" v-if="form.files.length > 0">
              <div v-for="(file, index) in form.files" :key="index" class="file-item">
                {{ file.name }}
                <button @click="removeFile(index)" class="btn-remove">删除</button>
              </div>
            </div>
          </div>
        </div>

        <div class="form-actions">
          <button type="submit" class="btn-submit">创建订单</button>
          <button type="button" class="btn-cancel" @click="goBack">取消</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { orderApi } from '@/services/api'

const router = useRouter()

const form = ref({
  customerName: '',
  productName: '',
  quantity: 1,
  totalPrice: 0,
  requirements: '',
  files: [] as File[]
})

const handleFileChange = (event: Event) => {
  const input = event.target as HTMLInputElement
  if (input.files) {
    form.value.files = Array.from(input.files)
  }
}

const removeFile = (index: number) => {
  form.value.files.splice(index, 1)
}

const goBack = () => {
  router.push('/orders')
}

const handleSubmit = async () => {
  try {
    const formData = new FormData()
    Object.entries(form.value).forEach(([key, value]) => {
      if (key !== 'files') {
        formData.append(key, value as string)
      }
    })
    form.value.files.forEach(file => {
      formData.append('files', file)
    })

    await orderApi.createOrder(formData)
    router.push('/orders')
  } catch (error) {
    console.error('创建订单失败:', error)
  }
}
</script>

<style scoped>
.create-order-container {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.create-order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.btn-back {
  padding: 0.5rem 1rem;
  border-radius: 0.375rem;
  background-color: #f3f4f6;
  color: #374151;
  border: 1px solid #d1d5db;
  cursor: pointer;
  transition: all 0.2s;
}

.create-order-form {
  background-color: white;
  padding: 2rem;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
}

.form-section {
  margin-bottom: 2rem;
}

.form-section h3 {
  margin-bottom: 1rem;
  color: #111827;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #374151;
  font-weight: 500;
}

.form-group input[type="text"],
.form-group input[type="number"],
.form-group textarea {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  transition: all 0.2s;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.file-list {
  margin-top: 0.5rem;
}

.file-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem;
  background-color: #f3f4f6;
  border-radius: 0.375rem;
  margin-bottom: 0.5rem;
}

.btn-remove {
  padding: 0.25rem 0.5rem;
  background-color: #ef4444;
  color: white;
  border: none;
  border-radius: 0.25rem;
  cursor: pointer;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
}

.btn-submit {
  padding: 0.5rem 1rem;
  background-color: #3b82f6;
  color: white;
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-cancel {
  padding: 0.5rem 1rem;
  background-color: #f3f4f6;
  color: #374151;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-submit:hover,
.btn-cancel:hover {
  opacity: 0.9;
}
</style> 