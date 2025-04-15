<template>
  <div class="i18n-demo">
    <el-card class="demo-section">
      <template #header>
        <div class="card-header">
          <span>{{ $t('demo.title') }}</span>
          <LanguageSelector />
        </div>
      </template>

      <!-- 日期格式化示例 -->
      <div class="demo-block">
        <h3>{{ $t('demo.dateFormat') }}</h3>
        <el-row :gutter="20">
          <el-col :span="12">
            <p>{{ $t('demo.shortDate') }}: {{ formatDate(currentDate) }}</p>
          </el-col>
          <el-col :span="12">
            <p>{{ $t('demo.longDate') }}: {{ formatDate(currentDate, 'PPP') }}</p>
          </el-col>
        </el-row>
      </div>

      <!-- 数字格式化示例 -->
      <div class="demo-block">
        <h3>{{ $t('demo.numberFormat') }}</h3>
        <el-row :gutter="20">
          <el-col :span="8">
            <p>{{ $t('demo.decimal') }}: {{ formatNumber(12345.6789) }}</p>
          </el-col>
          <el-col :span="8">
            <p>{{ $t('demo.currency') }}: {{ formatCurrency(12345.6789) }}</p>
          </el-col>
          <el-col :span="8">
            <p>{{ $t('demo.percent') }}: {{ formatPercent(0.12345) }}</p>
          </el-col>
        </el-row>
      </div>

      <!-- 文本翻译示例 -->
      <div class="demo-block">
        <h3>{{ $t('demo.translation') }}</h3>
        <el-row :gutter="20">
          <el-col :span="12">
            <h4>{{ $t('demo.orderStatus') }}</h4>
            <el-tag v-for="status in orderStatuses" 
                   :key="status" 
                   :type="getStatusType(status)"
                   class="status-tag">
              {{ $t(`order.status.${status}`) }}
            </el-tag>
          </el-col>
          <el-col :span="12">
            <h4>{{ $t('demo.factoryTypes') }}</h4>
            <el-tag v-for="type in factoryTypes" 
                   :key="type" 
                   :type="getFactoryType(type)"
                   class="type-tag">
              {{ $t(`factory.type.${type}`) }}
            </el-tag>
          </el-col>
        </el-row>
      </div>

      <!-- 动态内容示例 -->
      <div class="demo-block">
        <h3>{{ $t('demo.dynamicContent') }}</h3>
        <p>{{ $t('demo.welcome', { name: userName }) }}</p>
        <p>{{ $t('demo.orderCount', { count: orderCount }) }}</p>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import LanguageSelector from '@/components/common/LanguageSelector.vue'
import { formatDate, formatNumber, formatCurrency, formatPercent } from '@/i18n'

// 示例数据
const currentDate = new Date()
const userName = 'John Doe'
const orderCount = 5

// 订单状态
const orderStatuses = ['pending', 'processing', 'completed', 'cancelled']
const factoryTypes = ['premium', 'standard']

// 状态标签类型
const getStatusType = (status: string): string => {
  const types: Record<string, string> = {
    pending: 'warning',
    processing: 'primary',
    completed: 'success',
    cancelled: 'danger'
  }
  return types[status] || 'info'
}

// 工厂类型标签
const getFactoryType = (type: string): string => {
  const types: Record<string, string> = {
    premium: 'success',
    standard: 'info'
  }
  return types[type] || 'info'
}
</script>

<style scoped>
.i18n-demo {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.demo-section {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.demo-block {
  padding: 20px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.demo-block:last-child {
  border-bottom: none;
}

h3 {
  margin-top: 0;
  margin-bottom: 20px;
  color: var(--el-text-color-primary);
}

h4 {
  margin-top: 0;
  margin-bottom: 16px;
  color: var(--el-text-color-regular);
}

.status-tag,
.type-tag {
  margin-right: 8px;
  margin-bottom: 8px;
}

p {
  margin: 8px 0;
  color: var(--el-text-color-regular);
}
</style> 