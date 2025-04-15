<template>
  <div class="hot-orders">
    <h2 class="section-title">热门订单</h2>
    <div class="orders-grid">
      <div v-for="order in orders" :key="order.id" class="order-card">
        <el-card :body-style="{ padding: '0px' }">
          <img :src="order.imageUrl" class="order-image" />
          <div class="order-content">
            <h3 class="order-title">{{ order.title }}</h3>
            <el-tag :type="getStatusType(order.status)" size="small">
              {{ getStatusText(order.status) }}
            </el-tag>
            <div class="order-info">
              <p>数量: {{ order.quantity }} 件</p>
              <p>单价: ¥{{ order.price }}</p>
              <p>截止日期: {{ formatDate(order.deadline) }}</p>
            </div>
            <div class="order-progress" v-if="order.status === 'processing'">
              <el-progress :percentage="order.progress" />
            </div>
            <div class="order-description">{{ order.description }}</div>
            <div class="order-actions">
              <el-button type="primary" @click="viewOrderDetails(order.id)">查看详情</el-button>
            </div>
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElCard, ElTag, ElButton, ElProgress } from 'element-plus'

const props = defineProps<{
  orders: Array<{
    id: string
    title: string
    status: 'pending' | 'processing' | 'completed'
    quantity: number
    price: number
    imageUrl: string
    deadline: string
    progress: number
    description: string
  }>
}>()

const router = useRouter()

const getStatusType = (status: string) => {
  const types = {
    pending: 'warning',
    processing: 'primary',
    completed: 'success'
  }
  return types[status] || 'info'
}

const getStatusText = (status: string) => {
  const texts = {
    pending: '待处理',
    processing: '进行中',
    completed: '已完成'
  }
  return texts[status] || '未知状态'
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('zh-CN')
}

const viewOrderDetails = (orderId: string) => {
  router.push(`/orders/${orderId}`)
}
</script>

<style scoped lang="scss">
.hot-orders {
  padding: 20px;

  .section-title {
    font-size: 24px;
    margin-bottom: 20px;
    color: #333;
  }

  .orders-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 20px;
  }

  .order-card {
    transition: transform 0.3s ease;

    &:hover {
      transform: translateY(-5px);
    }

    .order-image {
      width: 100%;
      height: 200px;
      object-fit: cover;
    }

    .order-content {
      padding: 15px;

      .order-title {
        margin: 0 0 10px;
        font-size: 18px;
        color: #333;
      }

      .order-info {
        margin: 10px 0;
        font-size: 14px;
        color: #666;

        p {
          margin: 5px 0;
        }
      }

      .order-progress {
        margin: 10px 0;
      }

      .order-description {
        margin: 10px 0;
        font-size: 14px;
        color: #666;
        line-height: 1.4;
      }

      .order-actions {
        margin-top: 15px;
        text-align: right;
      }
    }
  }
}

@media (max-width: 768px) {
  .hot-orders {
    padding: 15px;

    .orders-grid {
      grid-template-columns: 1fr;
    }

    .order-card {
      .order-image {
        height: 180px;
      }
    }
  }
}
</style> 