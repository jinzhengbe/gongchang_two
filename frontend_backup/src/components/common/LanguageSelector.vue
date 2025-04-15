<template>
  <el-dropdown trigger="click" class="language-selector">
    <el-button>
      {{ currentLanguage }}
      <el-icon class="el-icon--right"><arrow-down /></el-icon>
    </el-button>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item @click="changeLanguage('zh-CN')">中文</el-dropdown-item>
        <el-dropdown-item @click="changeLanguage('en-US')">English</el-dropdown-item>
        <el-dropdown-item @click="changeLanguage('ko-KR')">한국어</el-dropdown-item>
        <el-dropdown-item @click="changeLanguage('vi-VN')">Tiếng Việt</el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ArrowDown } from '@element-plus/icons-vue'

const { locale } = useI18n()

const currentLanguage = computed(() => {
  const languages = {
    'zh-CN': '中文',
    'en-US': 'English',
    'ko-KR': '한국어',
    'vi-VN': 'Tiếng Việt'
  }
  return languages[locale.value] || languages['en-US']
})

const changeLanguage = (lang: string) => {
  locale.value = lang
}
</script>

<style scoped lang="scss">
.language-selector {
  margin-left: 16px;
  
  .el-button {
    display: flex;
    align-items: center;
    gap: 4px;
  }
}
</style> 