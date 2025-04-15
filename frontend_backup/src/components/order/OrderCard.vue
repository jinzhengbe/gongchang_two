<template>
  <el-card class="order-card" :body-style="{ padding: '0' }">
    <div class="order-header">
      <span class="order-number">{{ $t('order.number') }}: {{ order.orderNumber }}</span>
      <el-tag :type="getStatusType(order.status)">
        {{ $t(`order.status.${order.status}`) }}
      </el-tag>
    </div>
    
    <div class="order-content">
      <div class="order-info">
        <p class="info-item">
          <span class="label">{{ $t('order.date') }}:</span>
          <span class="value">{{ formatDate(order.orderDate) }}</span>
        </p>
        <p class="info-item">
          <span class="label">{{ $t('order.quantity') }}:</span>
          <span class="value">{{ formatNumber(order.quantity) }}</span>
        </p>
        <p class="info-item">
          <span class="label">{{ $t('order.price') }}:</span>
          <span class="value">{{ formatCurrency(order.price) }}</span>
        </p>
      </div>
      
      <div class="order-progress">
        <span class="progress-label">{{ $t('order.progress') }}</span>
        <el-progress 
          :percentage="order.progress" 
          :status="getProgressStatus(order.progress)"
        />
      </div>
    </div>

    <div class="order-footer">
      <el-button type="text" @click="$emit('view', order)">
        {{ $t('order.detail') }}
      </el-button>
      <el-button 
        v-if="order.status === 'pending'" 
        type="primary" 
        size="small"
        @click="$emit('process', order)"
      >
        {{ $t('common.edit') }}
      </el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { formatDate, formatNumber, formatCurrency } from '@/i18n'

interface OrderProps {
  orderNumber: string
  status: 'pending' | 'processing' | 'completed' | 'cancelled'
  orderDate: string | Date
  quantity: number
  price: number
  progress: number
}

const props = defineProps<{
  order: OrderProps
}>()

const emit = defineEmits<{
  (e: 'view', order: OrderProps): void
  (e: 'process', order: OrderProps): void
}>()

// 获取状态对应的类型
const getStatusType = (status: string): string => {
  const types: Record<string, string> = {
    pending: 'warning',
    processing: 'primary',
    completed: 'success',
    cancelled: 'danger'
  }
  return types[status] || 'info'
}

// 获取进度条状态
const getProgressStatus = (progress: number): 'success' | 'exception' | undefined => {
  if (progress >= 100) return 'success'
  if (progress < 0) return 'exception'
  return undefined
}
</script>

<style scoped>
.order-card {
  margin-bottom: 16px;
  border-radius: 8px;
  overflow: hidden;
}

.order-header {
  padding: 16px;
  background-color: var(--el-fill-color-light);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.order-number {
  font-weight: bold;
  color: var(--el-text-color-primary);
}

.order-content {
  padding: 16px;
}

.order-info {
  margin-bottom: 16px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  color: var(--el-text-color-regular);
}

.info-item:last-child {
  margin-bottom: 0;
}

.label {
  color: var(--el-text-color-secondary);
}

.value {
  font-family: var(--el-font-family-monospace);
  color: var(--el-text-color-primary);
}

.order-progress {
  margin-top: 16px;
}

.progress-label {
  display: block;
  margin-bottom: 8px;
  color: var(--el-text-color-secondary);
}

.order-footer {
  padding: 12px 16px;
  border-top: 1px solid var(--el-border-color-lighter);
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style> 