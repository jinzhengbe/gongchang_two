<template>
  <el-card class="order-card" shadow="hover">
    <div class="order-header">
      <h3>{{ order.title }}</h3>
      <el-tag :type="getStatusType(order.status)">
        {{ $t(`order.status.${order.status.toLowerCase()}`) }}
      </el-tag>
    </div>

    <div class="order-content">
      <div class="order-details">
        <p class="detail-item">
          <span class="label">{{ $t('order.designer') }}:</span>
          <span class="value">{{ order.designer }}</span>
        </p>
        <p class="detail-item">
          <span class="label">{{ $t('order.quantity') }}:</span>
          <span class="value">{{ formatNumber(order.quantity) }}</span>
        </p>
        <p class="detail-item">
          <span class="label">{{ $t('order.deadline') }}:</span>
          <span class="value">{{ formatDate(order.deadline) }}</span>
        </p>
      </div>

      <div class="order-progress">
        <el-progress 
          :percentage="order.progress" 
          :status="getProgressStatus(order.progress)"
        />
      </div>

      <div class="order-actions">
        <el-button type="primary" size="small" @click="handleViewDetails">
          {{ $t('order.viewDetails') }}
        </el-button>
        <el-button 
          v-if="order.status === 'PENDING'" 
          type="success" 
          size="small" 
          @click="handleAcceptOrder"
        >
          {{ $t('order.accept') }}
        </el-button>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { formatDate, formatNumber } from '@/utils/format'

const { t } = useI18n()

interface Props {
  order: {
    id: string
    title: string
    designer: string
    quantity: number
    deadline: string
    status: string
    progress: number
  }
}

const props = defineProps<Props>()
const emit = defineEmits(['view-details', 'accept-order'])

const getStatusType = (status: string) => {
  const types = {
    PENDING: 'warning',
    PROCESSING: 'primary',
    COMPLETED: 'success'
  }
  return types[status] || 'info'
}

const getProgressStatus = (progress: number) => {
  if (progress >= 100) return 'success'
  if (progress >= 80) return 'warning'
  return ''
}

const handleViewDetails = () => {
  emit('view-details', props.order.id)
}

const handleAcceptOrder = () => {
  emit('accept-order', props.order.id)
}
</script>

<style scoped>
.order-card {
  margin-bottom: 1rem;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.order-header h3 {
  margin: 0;
  font-size: 18px;
  color: var(--el-text-color-primary);
}

.order-content {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.order-details {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 0.5rem;
}

.detail-item {
  margin: 0;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.label {
  color: var(--el-text-color-secondary);
  font-size: 14px;
}

.value {
  color: var(--el-text-color-primary);
  font-weight: 500;
  font-family: var(--el-font-family);
}

.order-progress {
  padding: 0.5rem 0;
}

.order-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}

/* Responsive styles */
@media (max-width: 768px) {
  .order-header h3 {
    font-size: 16px;
  }

  .order-details {
    grid-template-columns: 1fr;
  }

  .order-actions {
    flex-direction: column;
  }

  .order-actions .el-button {
    width: 100%;
  }
}
</style> 