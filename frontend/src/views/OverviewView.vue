<template>
  <div class="overview">
    <!-- 统计卡片 -->
    <el-row :gutter="20" v-loading="loading.stats">
      <el-col :span="8">
        <el-card class="stat-card">
          <template #header>
            <div class="card-header">
              <el-icon><Document /></el-icon>
              <span>总订单数</span>
            </div>
          </template>
          <div class="card-content">
            <span class="number">{{ stats.totalOrders }}</span>
            <span class="label">单</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="stat-card">
          <template #header>
            <div class="card-header">
              <el-icon><Timer /></el-icon>
              <span>进行中订单</span>
            </div>
          </template>
          <div class="card-content">
            <span class="number">{{ stats.activeOrders }}</span>
            <span class="label">单</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="stat-card">
          <template #header>
            <div class="card-header">
              <el-icon><CircleCheck /></el-icon>
              <span>已完成订单</span>
            </div>
          </template>
          <div class="card-content">
            <span class="number">{{ stats.completedOrders }}</span>
            <span class="label">单</span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 快捷操作 -->
    <div class="quick-actions">
      <h3>快捷操作</h3>
      <el-row :gutter="20" class="mt-4">
        <el-col :span="8">
          <el-card class="action-card" @click="$router.push('/orders/create')">
            <el-icon><Plus /></el-icon>
            <h4>新建订单</h4>
            <p>创建新的生产订单</p>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card class="action-card" @click="$router.push('/products/create')">
            <el-icon><ShoppingCart /></el-icon>
            <h4>添加产品</h4>
            <p>添加新的产品信息</p>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card class="action-card" @click="$router.push('/products')">
            <el-icon><List /></el-icon>
            <h4>产品管理</h4>
            <p>查看和管理产品</p>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 最近订单 -->
    <div class="recent-orders">
      <h3>最近订单</h3>
      <el-table
        :data="recentOrders"
        style="width: 100%"
        class="mt-4"
        v-loading="loading.orders"
      >
        <el-table-column
          prop="orderNumber"
          label="订单号"
          width="180"
        />
        <el-table-column
          prop="productName"
          label="产品名称"
          width="180"
        />
        <el-table-column
          prop="status"
          label="状态"
          width="120"
        >
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="createdAt"
          label="创建时间"
          width="180"
        >
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column
          label="操作"
          width="120"
        >
          <template #default="scope">
            <el-button
              link
              type="primary"
              @click="$router.push(`/orders/${scope.row.id}`)"
            >
              查看详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from 'vue'
import { Document, Timer, CircleCheck, Plus, ShoppingCart, List } from '@element-plus/icons-vue'
import { format } from 'date-fns'
import { orderService } from '@/services/orderService'
import { ElMessage } from 'element-plus'

export default defineComponent({
  name: 'OverviewView',
  components: {
    Document,
    Timer,
    CircleCheck,
    Plus,
    ShoppingCart,
    List
  },
  setup() {
    const stats = ref({
      totalOrders: 0,
      activeOrders: 0,
      completedOrders: 0
    })

    const recentOrders = ref([])
    const loading = ref({
      stats: false,
      orders: false
    })

    const getStatusType = (status: string) => {
      const types = {
        '待处理': 'warning',
        '进行中': 'primary',
        '已完成': 'success',
        '已取消': 'danger'
      }
      return types[status] || 'info'
    }

    const formatDate = (date: string) => {
      return format(new Date(date), 'yyyy-MM-dd HH:mm')
    }

    const fetchOverviewData = async () => {
      try {
        loading.value.stats = true
        const statistics = await orderService.getOrderStatistics()
        stats.value = {
          totalOrders: statistics.totalOrders,
          activeOrders: statistics.statusCounts['processing'] || 0,
          completedOrders: statistics.statusCounts['completed'] || 0
        }
      } catch (error) {
        console.error('获取统计数据失败:', error)
        ElMessage.error('获取统计数据失败')
      } finally {
        loading.value.stats = false
      }
    }

    const fetchRecentOrders = async () => {
      try {
        loading.value.orders = true
        recentOrders.value = await orderService.getRecentOrders(5)
      } catch (error) {
        console.error('获取最近订单失败:', error)
        ElMessage.error('获取最近订单失败')
      } finally {
        loading.value.orders = false
      }
    }

    onMounted(() => {
      fetchOverviewData()
      fetchRecentOrders()
    })

    return {
      stats,
      recentOrders,
      loading,
      getStatusType,
      formatDate
    }
  }
})
</script>

<style scoped>
.overview {
  padding: 20px;
}

.stat-card {
  height: 100%;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-content {
  text-align: center;
  padding: 20px 0;
}

.number {
  font-size: 36px;
  font-weight: bold;
  color: #409eff;
}

.label {
  font-size: 14px;
  color: #909399;
  margin-left: 4px;
}

.quick-actions {
  margin-top: 24px;
}

.action-card {
  cursor: pointer;
  text-align: center;
  transition: all 0.3s;
}

.action-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
}

.action-card .el-icon {
  font-size: 32px;
  color: #409eff;
  margin-bottom: 12px;
}

.action-card h4 {
  margin: 12px 0 8px;
  font-size: 16px;
  color: #303133;
}

.action-card p {
  margin: 0;
  font-size: 14px;
  color: #909399;
}

.recent-orders {
  margin-top: 24px;
}

h3 {
  font-size: 18px;
  color: #303133;
  margin-bottom: 16px;
}

.mt-4 {
  margin-top: 16px;
}
</style> 