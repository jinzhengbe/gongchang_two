<template>
  <el-card class="order-card" shadow="hover">
    <div class="order-image">
      <img :data-src="order.imageUrl" :alt="order.title">
    </div>
    <div class="order-info">
      <h3>{{ order.title }}</h3>
      <p class="order-date">{{ formatDate(order.date) }}</p>
      <p class="order-status">
        <el-tag :type="getStatusType(order.status)">
          {{ $t(`order.status.${order.status.toLowerCase()}`) }}
        </el-tag>
      </p>
      <p class="order-quantity">{{ formatNumber(order.quantity) }} {{ $t('order.pieces') }}</p>
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
    imageUrl: string
    date: string
    status: string
    quantity: number
  }
}

const props = defineProps<Props>()

const getStatusType = (status: string) => {
  const types = {
    PENDING: 'warning',
    PROCESSING: 'primary',
    COMPLETED: 'success'
  }
  return types[status] || 'info'
}
</script>

<style scoped>
.order-card {
  height: 100%;
  transition: transform 0.3s ease;
}

.order-card:hover {
  transform: translateY(-5px);
}

.order-image {
  height: 200px;
  overflow: hidden;
}

.order-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.order-info {
  padding: 15px;
}

.order-info h3 {
  margin: 0 0 10px;
  font-size: 18px;
  color: var(--el-text-color-primary);
  font-weight: bold;
}

.order-date {
  color: var(--el-text-color-secondary);
  font-size: 14px;
  margin: 5px 0;
}

.order-status {
  margin: 10px 0;
}

.order-quantity {
  color: var(--el-color-primary);
  font-weight: bold;
  font-size: 16px;
  margin: 5px 0;
  font-family: var(--el-font-family);
}

/* Responsive styles */
@media (max-width: 768px) {
  .order-image {
    height: 150px;
  }

  .order-info {
    padding: 10px;
  }

  .order-info h3 {
    font-size: 16px;
  }

  .order-quantity {
    font-size: 14px;
  }
}
</style> 