<template>
  <div class="order-detail">
    <LoadingSpinner v-if="loading" />
    <ErrorMessage
      v-else-if="error"
      :error="true"
      title="加载失败"
      :message="errorMessage"
      :retry="true"
      @retry="fetchOrder"
    />
    <template v-else>
      <!-- 订单基本信息 -->
      <el-card class="order-info">
        <template #header>
          <div class="card-header">
            <span>订单信息</span>
            <el-button-group>
              <el-button type="primary" @click="handleEdit">编辑</el-button>
              <el-button type="danger" @click="handleDelete" v-if="canDelete">删除</el-button>
            </el-button-group>
          </div>
        </template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="订单号">{{ order.orderNumber }}</el-descriptions-item>
          <el-descriptions-item label="客户名称">{{ order.customerName }}</el-descriptions-item>
          <el-descriptions-item label="产品名称">{{ order.productName }}</el-descriptions-item>
          <el-descriptions-item label="数量">{{ order.quantity }}</el-descriptions-item>
          <el-descriptions-item label="总价">¥{{ order.totalPrice.toFixed(2) }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(order.status)">
              {{ getStatusLabel(order.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatDate(order.createdAt) }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ formatDate(order.updatedAt) }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 订单状态更新 -->
      <el-card class="status-update">
        <template #header>
          <div class="card-header">
            <span>状态更新</span>
          </div>
        </template>
        <el-form :model="statusForm" label-width="100px">
          <el-form-item label="新状态">
            <el-select v-model="statusForm.status" placeholder="请选择状态">
              <el-option
                v-for="status in availableStatuses"
                :key="status.value"
                :label="status.label"
                :value="status.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="备注">
            <el-input
              v-model="statusForm.notes"
              type="textarea"
              :rows="3"
              placeholder="请输入状态更新备注"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleStatusUpdate" :loading="updating">
              更新状态
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 文件管理 -->
      <el-card class="file-management">
        <template #header>
          <div class="card-header">
            <span>文件管理</span>
            <el-upload
              class="upload-btn"
              action="/api/orders/files"
              :headers="uploadHeaders"
              :data="{ orderId: order.id }"
              :on-success="handleUploadSuccess"
              :on-error="handleUploadError"
              :before-upload="beforeUpload"
            >
              <el-button type="primary">上传文件</el-button>
            </el-upload>
          </div>
        </template>
        <el-table :data="order.files" style="width: 100%">
          <el-table-column prop="name" label="文件名" />
          <el-table-column prop="size" label="大小" width="120">
            <template #default="{ row }">
              {{ formatFileSize(row.size) }}
            </template>
          </el-table-column>
          <el-table-column prop="uploadTime" label="上传时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.uploadTime) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150">
            <template #default="{ row }">
              <el-button-group>
                <el-button type="primary" link @click="downloadFile(row)">下载</el-button>
                <el-button type="danger" link @click="deleteFile(row)">删除</el-button>
              </el-button-group>
            </template>
          </el-table-column>
        </el-table>
        <el-form-item :label="$t('order.detail.files')">
          <file-uploader
            :order-id="order.id"
            :max-size="10 * 1024 * 1024"
            @upload-success="handleFileUploadSuccess"
            @upload-error="handleFileUploadError"
            @remove="handleFileRemove"
          />
        </el-form-item>
      </el-card>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { format } from 'date-fns'
import { useOrderStore } from '@/stores/order'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import ErrorMessage from '@/components/common/ErrorMessage.vue'
import type { Order, OrderStatus } from '@/types/order'
import FileUploader from '@/components/common/FileUploader.vue'

const route = useRoute()
const router = useRouter()
const orderStore = useOrderStore()

// 状态定义
const loading = ref(false)
const error = ref(false)
const errorMessage = ref('')
const updating = ref(false)
const order = ref<Order | null>(null)

// 状态更新表单
const statusForm = ref({
  status: '' as OrderStatus | '',
  notes: ''
})

// 上传配置
const uploadHeaders = {
  Authorization: `Bearer ${localStorage.getItem('token')}`
}

// 获取订单详情
const fetchOrder = async () => {
  loading.value = true
  error.value = false
  try {
    const orderId = route.params.id as string
    order.value = await orderStore.getOrderById(orderId)
  } catch (err) {
    error.value = true
    errorMessage.value = '加载订单详情失败，请稍后重试'
    ElMessage.error('加载订单详情失败')
  } finally {
    loading.value = false
  }
}

// 可用的状态选项
const availableStatuses = computed(() => {
  if (!order.value) return []
  const currentStatus = order.value.status
  const statuses = [
    { label: '待处理', value: 'pending' },
    { label: '处理中', value: 'processing' },
    { label: '已完成', value: 'completed' },
    { label: '已取消', value: 'cancelled' }
  ]
  return statuses.filter(status => status.value !== currentStatus)
})

// 状态类型和标签
const getStatusType = (status: OrderStatus) => {
  const types: Record<OrderStatus, string> = {
    pending: 'warning',
    processing: 'primary',
    completed: 'success',
    cancelled: 'info'
  }
  return types[status]
}

const getStatusLabel = (status: OrderStatus) => {
  const labels: Record<OrderStatus, string> = {
    pending: '待处理',
    processing: '处理中',
    completed: '已完成',
    cancelled: '已取消'
  }
  return labels[status]
}

// 日期格式化
const formatDate = (date: string) => {
  return format(new Date(date), 'yyyy-MM-dd HH:mm:ss')
}

// 文件大小格式化
const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 操作处理
const handleEdit = () => {
  router.push(`/orders/${order.value?.id}/edit`)
}

const handleDelete = async () => {
  try {
    await ElMessageBox.confirm('确定要删除此订单吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await orderStore.deleteOrder(order.value!.id)
    ElMessage.success('订单删除成功')
    router.push('/orders')
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error('删除订单失败')
    }
  }
}

const handleStatusUpdate = async () => {
  if (!statusForm.value.status) {
    ElMessage.warning('请选择新状态')
    return
  }
  updating.value = true
  try {
    await orderStore.updateOrderStatus(order.value!.id, {
      status: statusForm.value.status,
      notes: statusForm.value.notes
    })
    ElMessage.success('状态更新成功')
    fetchOrder()
    statusForm.value = { status: '', notes: '' }
  } catch (err) {
    ElMessage.error('状态更新失败')
  } finally {
    updating.value = false
  }
}

const beforeUpload = (file: File) => {
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过 10MB!')
    return false
  }
  return true
}

const handleUploadSuccess = () => {
  ElMessage.success('文件上传成功')
  fetchOrder()
}

const handleUploadError = () => {
  ElMessage.error('文件上传失败')
}

const downloadFile = async (file: any) => {
  try {
    const response = await orderStore.downloadFile(order.value!.id, file.id)
    const url = window.URL.createObjectURL(new Blob([response]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', file.name)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  } catch (err) {
    ElMessage.error('文件下载失败')
  }
}

const deleteFile = async (file: any) => {
  try {
    await ElMessageBox.confirm('确定要删除此文件吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await orderStore.deleteFile(order.value!.id, file.id)
    ElMessage.success('文件删除成功')
    fetchOrder()
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error('文件删除失败')
    }
  }
}

// 权限检查
const canDelete = computed(() => {
  return order.value?.status === 'pending'
})

const handleFileUploadSuccess = (file: any) => {
  order.value.files.push(file)
}

const handleFileUploadError = (error: any) => {
  console.error('File upload error:', error)
}

const handleFileRemove = (file: any) => {
  const index = order.value.files.findIndex(f => f.id === file.id)
  if (index > -1) {
    order.value.files.splice(index, 1)
  }
}

// 生命周期钩子
onMounted(() => {
  fetchOrder()
})
</script>

<style scoped>
.order-detail {
  padding: 20px;
}

.order-info,
.status-update,
.file-management {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.upload-btn {
  display: inline-block;
}
</style> 