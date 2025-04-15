<template>
  <!-- 移动端导航 -->
  <nav v-if="isMobile" class="navigation-mobile">
    <!-- 顶部栏 -->
    <div class="nav-header">
      <div class="brand">
        <router-link to="/">
          <span class="brand-text">{{ t('nav.brand') }}</span>
        </router-link>
      </div>
      
      <div class="current-lang" @click="handleLanguageChange">
        {{ currentLanguage }}
      </div>

      <div class="menu-grid" @click="toggleMenu">
        <el-icon><grid /></el-icon>
      </div>
    </div>

    <!-- 主体内容 -->
    <div class="main-content">
      <!-- 登录按钮，仅在点击九宫格后显示 -->
      <div v-if="showLoginButtons" class="login-buttons">
        <button class="login-btn designer" @click="handleLogin('designer')">
          {{ t('nav.login.designer') }}
        </button>
        <button class="login-btn factory" @click="handleLogin('factory')">
          {{ t('nav.login.factory') }}
        </button>
        <button class="login-btn supplier" @click="handleLogin('supplier')">
          {{ t('nav.login.supplier') }}
        </button>
      </div>
    </div>
  </nav>

  <!-- 桌面端导航 -->
  <nav v-else class="navigation-desktop">
    <div class="nav-container">
      <div class="nav-left">
        <div class="nav-logo">
          <router-link to="/">
            <span class="brand-text">{{ t('nav.brand') }}</span>
          </router-link>
        </div>
      </div>

      <div class="nav-right">
        <div class="nav-actions">
          <el-button class="designer" @click="handleLogin('designer')">
            {{ t('nav.login.designer') }}
          </el-button>
          <el-button class="factory" @click="handleLogin('factory')">
            {{ t('nav.login.factory') }}
          </el-button>
          <el-button class="supplier" @click="handleLogin('supplier')">
            {{ t('nav.login.supplier') }}
          </el-button>
        </div>
        
        <div class="language-selector">
          <el-dropdown trigger="click" @command="handleLanguageChange">
            <span class="language-text">
              {{ currentLanguage }}
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="zh-CN">中文</el-dropdown-item>
                <el-dropdown-item command="en-US">English</el-dropdown-item>
                <el-dropdown-item command="ko-KR">한국어</el-dropdown-item>
                <el-dropdown-item command="vi-VN">Tiếng Việt</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ArrowDown, Grid } from '@element-plus/icons-vue'
import { getSupportedLanguages, setLanguage, getCurrentLanguage } from '@/i18n'

const router = useRouter()
const { t, locale } = useI18n()

const isMobile = ref(false)
const showLoginButtons = ref(false)
const languages = getSupportedLanguages()

const handleResize = () => {
  isMobile.value = window.innerWidth < 768
}

onMounted(() => {
  handleResize()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})

const currentLanguage = computed(() => {
  const languages = {
    'zh-CN': '中文',
    'en-US': 'English',
    'ko-KR': '한국어',
    'vi-VN': 'Tiếng Việt'
  }
  return languages[locale.value] || '中文'
})

const currentLanguageLabel = computed(() => {
  const currentLang = languages.find(lang => lang.code === getCurrentLanguage())
  return currentLang?.nativeName || 'Language'
})

const handleLogin = (role: string) => {
  router.push(`/login/${role}`)
}

const handleLanguageChange = (lang?: string) => {
  if (lang) {
    locale.value = lang
  } else {
    locale.value = locale.value === 'zh-CN' ? 'en-US' : 'zh-CN'
  }
}

const toggleMenu = () => {
  showLoginButtons.value = !showLoginButtons.value
}
</script>

<style scoped lang="scss">
.navigation-mobile {
  width: 100%;
  background-color: #fff;
  position: fixed;
  top: 0;
  left: 0;
  z-index: 1000;
}

.nav-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 50px;
  padding: 0 16px;
  border-bottom: 1px solid #eee;
  background-color: #fff;
}

.brand {
  flex: 0 0 auto;
  text-align: left;

  a {
    text-decoration: none;
  }

  .brand-text {
    font-size: 18px;
    font-weight: 500;

    &::before {
      content: "Sewing";
      color: #409EFF;
    }

    &::after {
      content: " Mast";
      color: #333;
    }
  }
}

.current-lang {
  flex: 0 0 auto;
  margin: 0 10px;
  padding: 4px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
  color: #333;
  cursor: pointer;
}

.menu-grid {
  flex: 0 0 auto;
  cursor: pointer;
  padding: 8px;
}

.main-content {
  background-color: #fff;
}

.login-buttons {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 16px;
}

.login-btn {
  width: 100%;
  padding: 12px;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;

  &.designer {
    background-color: #409eff;
    color: white;
    &:hover {
      background-color: #66b1ff;
    }
  }

  &.factory {
    background-color: #67c23a;
    color: white;
    &:hover {
      background-color: #85ce61;
    }
  }

  &.supplier {
    background-color: #e6a23c;
    color: white;
    &:hover {
      background-color: #ebb563;
    }
  }
}

.navigation-desktop {
  height: 64px;
  background-color: #fff;
  border-bottom: 1px solid #eee;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
}

.nav-container {
  max-width: 1200px;
  margin: 0 auto;
  height: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
}

.nav-logo {
  img {
    height: 40px;
  }
}

.nav-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.nav-actions {
  display: flex;
  gap: 10px;
}

.language-selector {
  .language-text {
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 4px;
  }
}
</style> 