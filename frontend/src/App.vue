<template>
  <div class="app">
    <Navbar />
    <router-view />
  </div>
</template>

<script setup lang="ts">
import Navbar from '@/components/common/Navbar.vue'
</script>

<style>
.app {
  min-height: 100vh;
  background-color: var(--el-bg-color);
}
</style> 