<template>
  <div class="order-card" @click="$emit('click', order)">
    <div class="order-status" :class="getStatusType(order.status)">
      {{ $t(`orders.status.${order.status}`) }}
    </div>
    <div class="order-title">{{ order.title }}</div>
    <div class="order-info">
      <div class="info-item">
        <span class="label">{{ $t('order.quantity') }}:</span>
        <span class="value">{{ order.quantity }} {{ $t('orders.pieces') }}</span>
      </div>
      <div class="info-item">
        <span class="label">{{ $t('order.price') }}:</span>
        <span class="value">¥{{ order.price }}</span>
      </div>
      <div class="info-item">
        <span class="label">{{ $t('orders.deadline') }}:</span>
        <span class="value">{{ order.deadline }}</span>
      </div>
    </div>
    <div class="progress-section" v-if="order.progress !== undefined">
      <div class="progress-bar">
        <div class="progress" :style="{ width: `${order.progress}%` }"></div>
      </div>
      <span class="progress-text">{{ $t('orders.progress', { percentage: order.progress }) }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Goods, Money, Timer } from '@element-plus/icons-vue'

interface Order {
  id: string | number
  title: string
  status: 'pending' | 'processing' | 'completed' | 'cancelled'
  quantity: number
  price: number
  deadline: string
  progress: number
}

defineProps<{
  order: Order
}>()

const getStatusType = (status: string): string => {
  const types: Record<string, string> = {
    'pending': 'warning',
    'processing': 'primary',
    'completed': 'success',
    'cancelled': 'danger'
  }
  return types[status] || 'info'
}
</script>

<style scoped lang="scss">
.order-card {
  height: 100%;
  display: flex;
  flex-direction: column;

  .order-status {
    display: flex;
    justify-content: flex-end;
    margin-bottom: 0.5rem;
  }

  .order-title {
    font-size: 1.1rem;
    font-weight: 500;
    margin: 0.5rem 0;
    color: var(--el-text-color-primary);
  }

  .order-info {
    color: var(--el-text-color-regular);
    font-size: 0.9rem;

    .info-item {
      margin: 0.3rem 0;
      display: flex;
      align-items: center;
      gap: 0.5rem;

      .label {
        color: var(--el-text-color-secondary);
      }

      .value {
        color: var(--el-text-color-primary);
      }
    }
  }

  .progress-section {
    margin-top: 1rem;
    
    .progress-bar {
      height: 0.5rem;
      background-color: var(--el-bg-color-base);
      border-radius: 0.25rem;
      overflow: hidden;
    }

    .progress {
      height: 100%;
      background-color: var(--el-color-primary);
    }

    .progress-text {
      font-size: 0.9rem;
      color: var(--el-text-color-secondary);
      margin-top: 0.3rem;
      display: block;
    }
  }
}

@media (max-width: 768px) {
  .order-card {
    .order-title {
      font-size: 1rem;
    }

    .order-info {
      font-size: 0.85rem;
    }

    .progress-text {
      font-size: 0.85rem;
    }
  }
}
</style> 