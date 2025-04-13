<template>
  <div class="language-selector">
    <!-- 下拉菜单形式 -->
    <el-dropdown @command="handleLanguageChange" trigger="click">
      <span class="language-trigger">
        {{ currentLanguage.nativeName }}
        <el-icon class="el-icon--right"><arrow-down /></el-icon>
      </span>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item
            v-for="lang in languages"
            :key="lang.code"
            :command="lang.code"
            :class="{ active: lang.code === current }"
          >
            {{ lang.nativeName }}
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ArrowDown } from '@element-plus/icons-vue'
import { getCurrentLanguage, getSupportedLanguages, setLanguage } from '@/i18n'
import type { SupportedLocale } from '@/i18n'

// 获取支持的语言列表
const languages = getSupportedLanguages()

// 当前语言代码
const current = ref<SupportedLocale>(getCurrentLanguage())

// 当前语言完整信息
const currentLanguage = computed(() => {
  return languages.find(lang => lang.code === current.value) || languages[0]
})

// 处理语言切换
const handleLanguageChange = (lang: SupportedLocale) => {
  current.value = lang
  setLanguage(lang)
}

// 监听语言变更事件
onMounted(() => {
  window.addEventListener('language-changed', ((event: CustomEvent) => {
    current.value = event.detail as SupportedLocale
  }) as EventListener)
})
</script>

<style scoped>
.language-selector {
  display: inline-flex;
  align-items: center;
}

.language-trigger {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 8px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.language-trigger:hover {
  background-color: var(--el-fill-color-light);
}

.active {
  color: var(--el-color-primary);
  font-weight: bold;
}

:deep(.el-dropdown-menu__item) {
  display: flex;
  align-items: center;
  padding: 8px 16px;
}
</style> 