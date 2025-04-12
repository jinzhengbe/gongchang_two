<template>
  <div class="order-create">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>创建新订单</span>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        class="order-form"
      >
        <!-- 基本信息 -->
        <el-form-item label="客户名称" prop="customerName">
          <el-input v-model="form.customerName" placeholder="请输入客户名称" />
        </el-form-item>

        <el-form-item label="产品名称" prop="productName">
          <el-input v-model="form.productName" placeholder="请输入产品名称" />
        </el-form-item>

        <el-form-item label="数量" prop="quantity">
          <el-input-number
            v-model="form.quantity"
            :min="1"
            :max="9999"
            controls-position="right"
          />
        </el-form-item>

        <el-form-item label="单价" prop="unitPrice">
          <el-input-number
            v-model="form.unitPrice"
            :min="0"
            :precision="2"
            :step="0.1"
            controls-position="right"
          />
        </el-form-item>

        <el-form-item label="总价">
          <span class="total-price">¥{{ totalPrice.toFixed(2) }}</span>
        </el-form-item>

        <!-- 文件上传 -->
        <el-form-item label="设计文件" prop="files">
          <el-upload
            v-model:file-list="form.files"
            action="/api/orders/files"
            :headers="uploadHeaders"
            :before-upload="beforeUpload"
            :on-success="handleUploadSuccess"
            :on-error="handleUploadError"
            :on-remove="handleRemove"
            multiple
            :limit="5"
          >
            <el-button type="primary">上传文件</el-button>
            <template #tip>
              <div class="el-upload__tip">
                支持上传设计图纸、工艺说明等文件，单个文件不超过10MB
              </div>
            </template>
          </el-upload>
        </el-form-item>

        <!-- 备注信息 -->
        <el-form-item label="备注" prop="notes">
          <el-input
            v-model="form.notes"
            type="textarea"
            :rows="4"
            placeholder="请输入订单备注信息"
          />
        </el-form-item>

        <!-- 提交按钮 -->
        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            创建订单
          </el-button>
          <el-button @click="handleCancel">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { useOrderStore } from '@/stores/order'
import type { CreateOrderRequest } from '@/types/order'

const router = useRouter()
const orderStore = useOrderStore()
const formRef = ref<FormInstance>()

// 表单数据
const form = ref<CreateOrderRequest>({
  customerName: '',
  productName: '',
  quantity: 1,
  unitPrice: 0,
  files: [],
  notes: ''
})

// 上传配置
const uploadHeaders = {
  Authorization: `Bearer ${localStorage.getItem('token')}`
}

// 提交状态
const submitting = ref(false)

// 表单验证规则
const rules = ref<FormRules>({
  customerName: [
    { required: true, message: '请输入客户名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  productName: [
    { required: true, message: '请输入产品名称', trigger: 'blur' },
    { min: 2, max: 100, message: '长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  quantity: [
    { required: true, message: '请输入数量', trigger: 'blur' },
    { type: 'number', min: 1, message: '数量必须大于0', trigger: 'blur' }
  ],
  unitPrice: [
    { required: true, message: '请输入单价', trigger: 'blur' },
    { type: 'number', min: 0, message: '单价不能为负数', trigger: 'blur' }
  ]
})

// 计算总价
const totalPrice = computed(() => {
  return form.value.quantity * form.value.unitPrice
})

// 文件上传处理
const beforeUpload = (file: File) => {
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过 10MB!')
    return false
  }
  return true
}

const handleUploadSuccess = (response: any, file: any) => {
  ElMessage.success('文件上传成功')
}

const handleUploadError = () => {
  ElMessage.error('文件上传失败')
}

const handleRemove = (file: any) => {
  console.log('文件已移除:', file)
}

// 表单提交
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    submitting.value = true

    await orderStore.createOrder(form.value)
    ElMessage.success('订单创建成功')
    router.push('/orders')
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error('订单创建失败')
    }
  } finally {
    submitting.value = false
  }
}

// 取消创建
const handleCancel = () => {
  ElMessageBox.confirm('确定要取消创建订单吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(() => {
      router.push('/orders')
    })
    .catch(() => {})
}
</script>

<style scoped>
.order-create {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.order-form {
  max-width: 800px;
  margin: 0 auto;
}

.total-price {
  font-size: 18px;
  font-weight: bold;
  color: #f56c6c;
}

.el-upload__tip {
  color: #909399;
  font-size: 12px;
  margin-top: 8px;
}
</style> 