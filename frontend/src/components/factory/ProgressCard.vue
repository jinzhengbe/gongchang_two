<template>
  <el-card class="progress-card" shadow="hover">
    <div class="progress-header">
      <h3>{{ $t('progress.title') }}</h3>
      <el-select v-model="selectedPeriod" size="small">
        <el-option
          v-for="period in periods"
          :key="period.value"
          :label="$t(`progress.period.${period.value}`)"
          :value="period.value"
        />
      </el-select>
    </div>

    <div class="progress-content">
      <div class="progress-stats">
        <div class="stat-item">
          <span class="label">{{ $t('progress.totalOrders') }}</span>
          <span class="value">{{ formatNumber(stats.totalOrders) }}</span>
        </div>
        <div class="stat-item">
          <span class="label">{{ $t('progress.completedOrders') }}</span>
          <span class="value">{{ formatNumber(stats.completedOrders) }}</span>
        </div>
        <div class="stat-item">
          <span class="label">{{ $t('progress.onTimeRate') }}</span>
          <span class="value">{{ formatPercentage(stats.onTimeRate) }}</span>
        </div>
      </div>

      <div class="progress-chart">
        <el-progress
          type="circle"
          :percentage="stats.completionRate"
          :status="getProgressStatus(stats.completionRate)"
        >
          <template #default>
            <div class="progress-info">
              <span class="rate">{{ formatPercentage(stats.completionRate) }}</span>
              <span class="label">{{ $t('progress.completionRate') }}</span>
            </div>
          </template>
        </el-progress>
      </div>

      <div class="progress-timeline">
        <el-timeline>
          <el-timeline-item
            v-for="milestone in stats.recentMilestones"
            :key="milestone.id"
            :timestamp="formatDate(milestone.date)"
            :type="milestone.type"
          >
            {{ milestone.description }}
          </el-timeline-item>
        </el-timeline>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { formatDate, formatNumber, formatPercentage } from '@/utils/format'

const { t } = useI18n()

interface Props {
  stats: {
    totalOrders: number
    completedOrders: number
    onTimeRate: number
    completionRate: number
    recentMilestones: Array<{
      id: string
      date: string
      description: string
      type: 'primary' | 'success' | 'warning' | 'danger'
    }>
  }
}

const props = defineProps<Props>()
const emit = defineEmits(['period-change'])

const periods = [
  { value: 'day', label: 'Daily' },
  { value: 'week', label: 'Weekly' },
  { value: 'month', label: 'Monthly' }
]

const selectedPeriod = ref('week')

const getProgressStatus = (rate: number) => {
  if (rate >= 90) return 'success'
  if (rate >= 70) return 'warning'
  return 'exception'
}

watch(selectedPeriod, (newValue) => {
  emit('period-change', newValue)
})
</script>

<style scoped>
.progress-card {
  height: 100%;
}

.progress-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.progress-header h3 {
  margin: 0;
  font-size: 18px;
  color: var(--el-text-color-primary);
}

.progress-content {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.progress-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 1rem;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.stat-item .label {
  color: var(--el-text-color-secondary);
  font-size: 14px;
  margin-bottom: 0.5rem;
}

.stat-item .value {
  color: var(--el-text-color-primary);
  font-size: 24px;
  font-weight: 600;
}

.progress-chart {
  display: flex;
  justify-content: center;
  padding: 1rem 0;
}

.progress-info {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.progress-info .rate {
  font-size: 20px;
  font-weight: 600;
  line-height: 1;
}

.progress-info .label {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-top: 0.5rem;
}

.progress-timeline {
  margin-top: 1rem;
}

/* Responsive styles */
@media (max-width: 768px) {
  .progress-header {
    flex-direction: column;
    align-items: stretch;
    gap: 1rem;
  }

  .progress-stats {
    grid-template-columns: 1fr;
  }

  .stat-item {
    padding: 1rem;
    background-color: var(--el-bg-color-page);
    border-radius: var(--el-border-radius-base);
  }
}
</style> 