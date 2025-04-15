<template>
  <div class="loading-spinner" v-if="loading">
    <div class="spinner"></div>
    <div class="text">{{ text }}</div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  loading: boolean
  text?: string
}>()
</script>

<style scoped>
.loading-spinner {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f4f6;
  border-top: 4px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.text {
  margin-top: 1rem;
  color: #6b7280;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style> 