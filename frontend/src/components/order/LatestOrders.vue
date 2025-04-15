<template>
  <div class="latest-orders">
    <div class="section-header">
      <h2>{{ $t('home.latestOrders') }}</h2>
      <p class="section-desc">{{ $t('home.latestOrdersDesc') }}</p>
      <button class="view-more" @click="$emit('view-more')">{{ $t('home.viewMore') }}</button>
    </div>
    
    <el-row :gutter="20">
      <el-col :xs="24" :sm="12" :md="8" :lg="6" v-for="order in orders" :key="order.id">
        <el-card class="order-card" shadow="hover" @click="$emit('select-order', order)">
          <div class="order-image">
            <el-image 
              :src="order.imageUrl" 
              :alt="order.title"
              fit="cover"
              :preview-src-list="[order.imageUrl]"
              :initial-index="0"
            />
            <div class="order-status" :class="order.status">
              <span class="status" :class="order.status">{{ $t(`orders.status.${order.status}`) }}</span>
            </div>
          </div>
          <div class="order-content">
            <h3 class="order-title" :title="order.title">{{ order.title }}</h3>
            <p class="order-desc" :title="order.description">{{ order.description }}</p>
            <div class="order-info">
              <div class="info-item">
                <el-icon><ShoppingCart /></el-icon>
                {{ formatNumber(order.quantity) }} {{ $t('orders.pieces') }}
              </div>
              <div class="info-item">
                <el-icon><Money /></el-icon>
                ¥{{ formatNumber(order.price) }}
              </div>
            </div>
            <div class="order-deadline" :class="{ 'urgent': isUrgent(order.deadline) }">
              <el-icon><Timer /></el-icon>
              {{ $t('orders.deadline') }}: {{ formatDate(order.deadline) }}
            </div>
            <el-progress 
              v-if="order.status === 'processing'"
              :percentage="order.progress"
              :format="percentageFormat"
              :status="getProgressStatus(order.progress)"
              class="order-progress"
              striped
              :striped-flow="true"
            />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <div v-if="orders.length === 0" class="no-data">
      <el-empty :description="$t('common.noData')" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ShoppingCart, Money, Timer } from '@element-plus/icons-vue'
import type { Order } from '@/types'

const { t } = useI18n()

defineProps<{
  orders: Order[]
}>()

defineEmits<{
  (e: 'view-more'): void
  (e: 'select-order', order: Order): void
}>()

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString()
}

const formatNumber = (num: number) => {
  return new Intl.NumberFormat().format(num)
}

const percentageFormat = (percentage: number) => {
  return t('orders.progress', { percentage })
}

const isUrgent = (deadline: string) => {
  const deadlineDate = new Date(deadline)
  const now = new Date()
  const diffDays = Math.ceil((deadlineDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
  return diffDays <= 3
}

const getProgressStatus = (progress: number) => {
  if (progress >= 90) return 'success'
  if (progress >= 60) return ''
  return 'warning'
}
</script>

<style scoped lang="scss">
.latest-orders {
  .section-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    margin-bottom: 2rem;
    flex-wrap: wrap;
    gap: 1rem;

    h2 {
      font-size: 1.875rem;
      font-weight: 600;
      color: var(--el-text-color-primary);
      margin: 0;
      flex: 1;
    }

    .section-desc {
      color: var(--el-text-color-secondary);
      font-size: 1rem;
      margin: 0;
      width: 100%;
    }

    .view-more {
      white-space: nowrap;
    }
  }

  .order-card {
    margin-bottom: 1.5rem;
    transition: all 0.3s ease;
    cursor: pointer;

    &:hover {
      transform: translateY(-5px);
      box-shadow: var(--el-box-shadow-light);
    }
  }

  .order-image {
    position: relative;
    height: 200px;
    overflow: hidden;
    border-radius: 4px 4px 0 0;

    .el-image {
      width: 100%;
      height: 100%;
      transition: transform 0.3s ease;

      &:hover {
        transform: scale(1.05);
      }
    }

    .order-status {
      position: absolute;
      top: 1rem;
      right: 1rem;
      padding: 0.25rem 0.75rem;
      border-radius: 1rem;
      color: white;
      font-size: 0.875rem;
      font-weight: 500;
      backdrop-filter: blur(4px);
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

      &.pending {
        background-color: var(--el-color-warning);
      }

      &.processing {
        background-color: var(--el-color-primary);
      }

      &.completed {
        background-color: var(--el-color-success);
      }
    }
  }

  .order-content {
    padding: 1rem;

    .order-title {
      font-size: 1.125rem;
      font-weight: 600;
      margin: 0 0 0.5rem;
      color: var(--el-text-color-primary);
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .order-desc {
      font-size: 0.875rem;
      color: var(--el-text-color-secondary);
      margin: 0 0 1rem;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
      line-height: 1.4;
    }

    .order-info {
      display: flex;
      justify-content: space-between;
      margin-bottom: 0.75rem;

      .info-item {
        display: flex;
        align-items: center;
        color: var(--el-text-color-regular);
        font-size: 0.875rem;
        gap: 0.25rem;

        .el-icon {
          color: var(--el-color-primary);
        }
      }
    }

    .order-deadline {
      display: flex;
      align-items: center;
      color: var(--el-text-color-regular);
      font-size: 0.875rem;
      margin-bottom: 0.75rem;
      gap: 0.25rem;

      .el-icon {
        color: var(--el-color-info);
      }

      &.urgent {
        color: var(--el-color-danger);
        font-weight: 500;

        .el-icon {
          color: var(--el-color-danger);
        }
      }
    }

    .order-progress {
      margin-top: 1rem;

      :deep(.el-progress-bar__inner) {
        transition: all 0.3s ease;
      }
    }
  }

  .no-data {
    padding: 2rem;
    text-align: center;
  }

  @media (max-width: 768px) {
    .section-header {
      h2 {
        font-size: 1.5rem;
      }

      .section-desc {
        font-size: 0.875rem;
      }

      .view-more {
        width: 100%;
      }
    }

    .order-image {
      height: 180px;
    }

    .order-content {
      .order-title {
        font-size: 1rem;
      }

      .order-desc {
        font-size: 0.813rem;
      }

      .order-info,
      .order-deadline {
        font-size: 0.813rem;
      }
    }
  }
}
</style> 