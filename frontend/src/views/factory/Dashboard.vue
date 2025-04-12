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
              <span>{{ $t('dashboard.activeOrders') }}</span>
            </div>
          </template>
          <el-table :data="activeOrders" style="width: 100%">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="title" :label="$t('dashboard.orderTitle')" />
            <el-table-column prop="designer" :label="$t('dashboard.designer')" />
            <el-table-column prop="progress" :label="$t('dashboard.progress')">
              <template #default="{ row }">
                <el-progress :percentage="row.progress" />
              </template>
            </el-table-column>
            <el-table-column prop="dueDate" :label="$t('dashboard.dueDate')">
              <template #default="{ row }">
                {{ formatDate(row.dueDate) }}
              </template>
            </el-table-column>
            <el-table-column :label="$t('dashboard.actions')" width="150">
              <template #default="{ row }">
                <el-button-group>
                  <el-button size="small" @click="viewOrder(row.id)">
                    {{ $t('dashboard.view') }}
                  </el-button>
                  <el-button size="small" type="primary" @click="updateProgress(row.id)">
                    {{ $t('dashboard.update') }}
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
              <span>{{ $t('dashboard.capacity') }}</span>
            </div>
          </template>
          <div class="capacity-chart">
            <div class="chart-header">
              <span>{{ $t('dashboard.currentCapacity') }}</span>
              <el-tag type="success">{{ capacity }}%</el-tag>
            </div>
            <el-progress
              :percentage="capacity"
              :color="getCapacityColor"
              :format="formatCapacity"
            />
            <div class="capacity-info">
              <div class="info-item">
                <span>{{ $t('dashboard.available') }}</span>
                <span>{{ availableCapacity }}</span>
              </div>
              <div class="info-item">
                <span>{{ $t('dashboard.inProgress') }}</span>
                <span>{{ inProgressCapacity }}</span>
              </div>
            </div>
          </div>
        </el-card>

        <el-card class="mt-20">
          <template #header>
            <div class="card-header">
              <span>{{ $t('dashboard.quickActions') }}</span>
            </div>
          </template>
          <div class="quick-actions">
            <el-button type="primary" icon="el-icon-edit" @click="updateCapacity">
              {{ $t('dashboard.updateCapacity') }}
            </el-button>
            <el-button icon="el-icon-setting" @click="manageEquipment">
              {{ $t('dashboard.manageEquipment') }}
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
import { format } from 'date-fns'
import axios from 'axios'

export default defineComponent({
  name: 'FactoryDashboard',
  setup() {
    const router = useRouter()
    const { t } = useI18n()

    const statistics = ref([
      { title: 'totalOrders', value: 0, trend: 0 },
      { title: 'activeOrders', value: 0, trend: 0 },
      { title: 'completedOrders', value: 0, trend: 0 }
    ])

    const activeOrders = ref([])
    const capacity = ref(0)
    const availableCapacity = ref(0)
    const inProgressCapacity = ref(0)

    const formatDate = (date: string) => {
      return format(new Date(date), 'yyyy-MM-dd')
    }

    const getCapacityColor = (percentage: number) => {
      if (percentage < 50) return '#67c23a'
      if (percentage < 80) return '#e6a23c'
      return '#f56c6c'
    }

    const formatCapacity = (percentage: number) => {
      return `${percentage}%`
    }

    const fetchDashboardData = async () => {
      try {
        const response = await axios.get('/api/factory/dashboard')
        const { stats, orders, capacityInfo } = response.data
        statistics.value = stats
        activeOrders.value = orders
        capacity.value = capacityInfo.current
        availableCapacity.value = capacityInfo.available
        inProgressCapacity.value = capacityInfo.inProgress
      } catch (error) {
        console.error('Failed to fetch dashboard data:', error)
      }
    }

    const viewOrder = (id: number) => {
      router.push(`/factory/orders/${id}`)
    }

    const updateProgress = (id: number) => {
      router.push(`/factory/orders/${id}/progress`)
    }

    const updateCapacity = () => {
      router.push('/factory/capacity')
    }

    const manageEquipment = () => {
      router.push('/factory/equipment')
    }

    onMounted(() => {
      fetchDashboardData()
    })

    return {
      statistics,
      activeOrders,
      capacity,
      availableCapacity,
      inProgressCapacity,
      formatDate,
      getCapacityColor,
      formatCapacity,
      viewOrder,
      updateProgress,
      updateCapacity,
      manageEquipment
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

.capacity-chart {
  padding: 20px;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.capacity-info {
  margin-top: 20px;
  display: flex;
  justify-content: space-between;
}

.info-item {
  display: flex;
  flex-direction: column;
  align-items: center;
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