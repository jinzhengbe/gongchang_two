<template>
  <el-card class="factory-card" shadow="hover">
    <div class="factory-image">
      <img :data-src="factory.imageUrl" :alt="factory.name">
    </div>
    <div class="factory-info">
      <h3>{{ factory.name }}</h3>
      <p class="factory-type">
        <el-tag :type="getFactoryType(factory.type)">
          {{ $t(`factory.type.${factory.type.toLowerCase()}`) }}
        </el-tag>
      </p>
      <p class="factory-location">
        <el-icon><Location /></el-icon>
        {{ factory.location }}
      </p>
      <p class="factory-capacity">
        {{ $t('factory.monthlyCapacity') }}: {{ formatNumber(factory.monthlyCapacity) }}
      </p>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Location } from '@element-plus/icons-vue'
import { formatNumber } from '@/utils/format'

const { t } = useI18n()

interface Props {
  factory: {
    id: string
    name: string
    imageUrl: string
    type: string
    location: string
    monthlyCapacity: number
  }
}

const props = defineProps<Props>()

const getFactoryType = (type: string) => {
  const types = {
    PREMIUM: 'success',
    STANDARD: 'info'
  }
  return types[type] || 'info'
}
</script>

<style scoped>
.factory-card {
  height: 100%;
  transition: transform 0.3s ease;
}

.factory-card:hover {
  transform: translateY(-5px);
}

.factory-image {
  height: 200px;
  overflow: hidden;
}

.factory-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.factory-info {
  padding: 15px;
}

.factory-info h3 {
  margin: 0 0 10px;
  font-size: 18px;
  color: var(--el-text-color-primary);
  font-weight: bold;
}

.factory-type {
  margin: 10px 0;
}

.factory-location {
  color: var(--el-text-color-secondary);
  font-size: 14px;
  margin: 5px 0;
  display: flex;
  align-items: center;
  gap: 5px;
}

.factory-capacity {
  color: var(--el-color-primary);
  font-weight: bold;
  font-size: 16px;
  margin: 5px 0;
  font-family: var(--el-font-family);
}

/* Responsive styles */
@media (max-width: 768px) {
  .factory-image {
    height: 150px;
  }

  .factory-info {
    padding: 10px;
  }

  .factory-info h3 {
    font-size: 16px;
  }

  .factory-capacity {
    font-size: 14px;
  }
}
</style> 