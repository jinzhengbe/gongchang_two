<template>
  <div class="dashboard-container">
    <el-row :gutter="20">
      <el-col :span="8" v-for="(stat, index) in statistics" :key="index">
        <el-card class="stat-card">
          <template #header>
            <div class="card-header">
              <span>{{ $t(`dashboard.stats.${stat.title}`) }}</span>
            </div>
          </template>
          <div class="stat-content">
            <div class="stat-value">{{ stat.value }}</div>
            <div class="stat-trend" :class="stat.trend > 0 ? 'up' : 'down'">
              {{ stat.trend > 0 ? '+' : '' }}{{ stat.trend }}%
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-20">
      <el-col :span="16">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ $t('dashboard.recentOrders') }}</span>
              <el-button type="primary" @click="createOrder">
                {{ $t('dashboard.createOrder') }}
              </el-button>
            </div>
          </template>
          <el-table :data="recentOrders" style="width: 100%">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="title" :label="$t('dashboard.orderTitle')" />
            <el-table-column prop="status" :label="$t('dashboard.status')">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)">
                  {{ $t(`order.status.${row.status}`) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="progress" :label="$t('dashboard.progress')">
              <template #default="{ row }">
                <el-progress :percentage="row.progress" />
              </template>
            </el-table-column>
            <el-table-column :label="$t('dashboard.actions')" width="150">
              <template #default="{ row }">
                <el-button-group>
                  <el-button size="small" @click="viewOrder(row.id)">
                    {{ $t('dashboard.view') }}
                  </el-button>
                  <el-button size="small" type="primary" @click="editOrder(row.id)">
                    {{ $t('dashboard.edit') }}
                  </el-button>
                </el-button-group>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ $t('dashboard.quickActions') }}</span>
            </div>
          </template>
          <div class="quick-actions">
            <el-button type="primary" icon="el-icon-plus" @click="createOrder">
              {{ $t('dashboard.createOrder') }}
            </el-button>
            <el-button icon="el-icon-search" @click="findFactory">
              {{ $t('dashboard.findFactory') }}
            </el-button>
            <el-button icon="el-icon-shopping-cart" @click="findMaterial">
              {{ $t('dashboard.findMaterial') }}
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import axios from 'axios'

export default defineComponent({
  name: 'DesignerDashboard',
  setup() {
    const router = useRouter()
    const { t } = useI18n()

    const statistics = ref([
      { title: 'totalOrders', value: 0, trend: 0 },
      { title: 'activeOrders', value: 0, trend: 0 },
      { title: 'completedOrders', value: 0, trend: 0 }
    ])

    const recentOrders = ref([])

    const getStatusType = (status: string) => {
      const types: { [key: string]: string } = {
        pending: 'warning',
        processing: 'primary',
        completed: 'success',
        cancelled: 'danger'
      }
      return types[status] || 'info'
    }

    const fetchDashboardData = async () => {
      try {
        const response = await axios.get('/api/designer/dashboard')
        const { stats, orders } = response.data
        statistics.value = stats
        recentOrders.value = orders
      } catch (error) {
        console.error('Failed to fetch dashboard data:', error)
      }
    }

    const createOrder = () => {
      router.push('/designer/orders/create')
    }

    const viewOrder = (id: number) => {
      router.push(`/designer/orders/${id}`)
    }

    const editOrder = (id: number) => {
      router.push(`/designer/orders/${id}/edit`)
    }

    const findFactory = () => {
      router.push('/designer/factories')
    }

    const findMaterial = () => {
      router.push('/designer/materials')
    }

    onMounted(() => {
      fetchDashboardData()
    })

    return {
      statistics,
      recentOrders,
      getStatusType,
      createOrder,
      viewOrder,
      editOrder,
      findFactory,
      findMaterial
    }
  }
})
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
}

.mt-20 {
  margin-top: 20px;
}

.stat-card {
  margin-bottom: 20px;
}

.stat-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
}

.stat-trend {
  font-size: 14px;
}

.stat-trend.up {
  color: #67c23a;
}

.stat-trend.down {
  color: #f56c6c;
}

.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.quick-actions .el-button {
  width: 100%;
}
</style> 