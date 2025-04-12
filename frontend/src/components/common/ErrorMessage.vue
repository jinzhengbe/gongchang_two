<template>
  <div class="error-message" v-if="error">
    <div class="icon">⚠️</div>
    <div class="content">
      <div class="title">{{ title }}</div>
      <div class="message">{{ message }}</div>
      <button v-if="retry" class="btn-retry" @click="$emit('retry')">重试</button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  error: boolean
  title?: string
  message: string
  retry?: boolean
}>()

defineEmits<{
  (e: 'retry'): void
}>()
</script>

<style scoped>
.error-message {
  display: flex;
  gap: 1rem;
  padding: 1rem;
  background-color: #fef2f2;
  border: 1px solid #fee2e2;
  border-radius: 0.5rem;
  color: #991b1b;
}

.icon {
  font-size: 1.5rem;
}

.content {
  flex: 1;
}

.title {
  font-weight: 500;
  margin-bottom: 0.5rem;
}

.message {
  margin-bottom: 1rem;
}

.btn-retry {
  padding: 0.5rem 1rem;
  background-color: #ef4444;
  color: white;
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-retry:hover {
  background-color: #dc2626;
}
</style> 