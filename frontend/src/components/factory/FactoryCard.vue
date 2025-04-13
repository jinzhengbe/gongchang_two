<template>
  <el-card class="factory-card">
    <template #header>
      <div class="factory-header">
        <h3 class="factory-name">{{ factory.name }}</h3>
        <el-tag :type="getFactoryType(factory.type)">
          {{ $t(`factory.type.${factory.type}`) }}
        </el-tag>
      </div>
    </template>

    <div class="factory-info">
      <div class="info-row">
        <span class="info-label">{{ $t('factory.capacity') }}</span>
        <span class="info-value">{{ formatNumber(factory.capacity) }}</span>
      </div>

      <div class="info-row">
        <span class="info-label">{{ $t('factory.rating') }}</span>
        <el-rate
          v-model="factory.rating"
          disabled
          show-score
          text-color="#ff9900"
        />
      </div>

      <div class="info-row">
        <span class="info-label">{{ $t('factory.location') }}</span>
        <span class="info-value">{{ factory.location }}</span>
      </div>

      <div class="certifications">
        <h4>{{ $t('factory.certification') }}</h4>
        <div class="cert-tags">
          <el-tag
            v-for="cert in factory.certifications"
            :key="cert"
            size="small"
            class="cert-tag"
          >
            {{ cert }}
          </el-tag>
        </div>
      </div>
    </div>

    <div class="factory-footer">
      <el-button type="text" @click="$emit('view', factory)">
        {{ $t('factory.info') }}
      </el-button>
      <el-button type="primary" size="small" @click="$emit('contact', factory)">
        {{ $t('factory.contact') }}
      </el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { formatNumber } from '@/i18n'

interface FactoryProps {
  name: string
  type: 'premium' | 'standard'
  capacity: number
  rating: number
  location: string
  certifications: string[]
}

const props = defineProps<{
  factory: FactoryProps
}>()

const emit = defineEmits<{
  (e: 'view', factory: FactoryProps): void
  (e: 'contact', factory: FactoryProps): void
}>()

const getFactoryType = (type: string): string => {
  return type === 'premium' ? 'success' : 'info'
}
</script>

<style scoped>
.factory-card {
  margin-bottom: 16px;
  border-radius: 8px;
}

.factory-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.factory-name {
  margin: 0;
  font-size: 18px;
  color: var(--el-text-color-primary);
}

.factory-info {
  padding: 8px 0;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.info-label {
  color: var(--el-text-color-secondary);
  font-size: 14px;
}

.info-value {
  color: var(--el-text-color-primary);
  font-family: var(--el-font-family-monospace);
}

.certifications {
  margin-top: 20px;
}

.certifications h4 {
  margin: 0 0 12px;
  color: var(--el-text-color-primary);
  font-size: 14px;
}

.cert-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.cert-tag {
  margin-right: 8px;
  margin-bottom: 8px;
}

.factory-footer {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--el-border-color-lighter);
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style> 