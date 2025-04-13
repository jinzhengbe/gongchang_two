<template>
  <!-- 移动端导航 -->
  <nav v-if="isMobile" class="navigation-mobile">
    <!-- 顶部栏 -->
    <div class="nav-header">
      <div class="brand">
        <router-link to="/">
          <span class="brand-text">Sewing Mast</span>
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
          设计师登录
        </button>
        <button class="login-btn factory" @click="handleLogin('factory')">
          工厂登录
        </button>
        <button class="login-btn supplier" @click="handleLogin('supplier')">
          供应商登录
        </button>
      </div>

      <!-- Latest Orders 卡片，始终显示 -->
      <div class="latest-orders">
        <h2>Latest Orders</h2>
        <p>View the latest published orders</p>
        <button class="view-more">查看更多</button>
      </div>
    </div>
  </nav>

  <!-- 桌面端导航 -->
  <nav v-else class="navigation-desktop">
    <div class="nav-container">
      <div class="nav-left">
        <div class="nav-logo">
          <router-link to="/">
            <img src="@/assets/logo.svg" alt="SewingMast" />
          </router-link>
        </div>
        <ul class="nav-menu">
          <li class="nav-menu-item" :class="{ active: isActive('/') }">
            <router-link to="/">{{ t('nav.home') }}</router-link>
          </li>
          <li class="nav-menu-item" :class="{ active: isActive('/orders') }">
            <router-link to="/orders">{{ t('nav.orders') }}</router-link>
          </li>
          <li class="nav-menu-item" :class="{ active: isActive('/factories') }">
            <router-link to="/factories">{{ t('nav.factories') }}</router-link>
          </li>
          <li class="nav-menu-item" :class="{ active: isActive('/fabrics') }">
            <router-link to="/fabrics">{{ t('nav.fabrics') }}</router-link>
          </li>
        </ul>
      </div>

      <div class="nav-right">
        <div class="nav-actions">
          <el-button class="designer" @click="handleLogin('designer')">
            {{ t('login.designer') }}
          </el-button>
          <el-button class="factory" @click="handleLogin('factory')">
            {{ t('login.factory') }}
          </el-button>
          <el-button class="supplier" @click="handleLogin('supplier')">
            {{ t('login.supplier') }}
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
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ArrowDown, Grid } from '@element-plus/icons-vue'
import { ElMessageBox } from 'element-plus'

const router = useRouter()
const route = useRoute()
const { t, locale } = useI18n()

const isMobile = ref(false)
const showLoginButtons = ref(false)

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

const isActive = (path: string) => {
  return route.path === path
}

const handleLogin = (role: string) => {
  router.push(`/login/${role}`)
}

const handleLanguageChange = () => {
  locale.value = locale.value === 'zh-CN' ? 'en-US' : 'zh-CN'
}

const toggleMenu = () => {
  showLoginButtons.value = !showLoginButtons.value
}
</script>

<style scoped lang="scss">
/* 移动端样式 */
.navigation-mobile {
  width: 100%;
  background-color: #fff;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.nav-header {
  display: flex;
  align-items: center;
  height: 50px;
  padding: 0;
  border-bottom: 1px solid #eee;
}

.brand {
  flex: 1;
  text-align: left;
  padding-left: 16px;
}

.brand-text {
  font-size: 18px;
  color: #000;
  text-decoration: none;
  border-bottom: 1px solid #000;
}

.current-lang {
  flex: 1;
  text-align: center;
  font-size: 16px;
  color: #333;
  cursor: pointer;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 4px 12px;
  display: inline-block;
  margin: 0 auto;
  background-color: #fff;
  transition: all 0.2s;
}

.current-lang:hover {
  border-color: #c0c4cc;
  background-color: #f5f7fa;
}

.menu-grid {
  flex: 1;
  text-align: right;
  padding-right: 16px;
  cursor: pointer;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 16px;
  gap: 20px;
}

.login-buttons {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.login-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  border: none;
  border-radius: 4px;
  color: #fff;
}

.login-btn.designer {
  background-color: #409EFF;
}

.login-btn.factory {
  background-color: #67C23A;
}

.login-btn.supplier {
  background-color: #E6A23C;
}

.latest-orders {
  background-color: #67C23A;
  color: white;
  padding: 24px;
  border-radius: 8px;
}

.latest-orders h2 {
  font-size: 24px;
  margin: 0 0 8px 0;
}

.latest-orders p {
  font-size: 16px;
  margin: 0 0 20px 0;
  opacity: 0.9;
}

.latest-orders .view-more {
  background-color: #fff;
  color: #67C23A;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
}

/* 桌面端样式 */
.navigation-desktop {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  z-index: 1000;
}

.navigation-desktop .nav-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.navigation-desktop .nav-left {
  display: flex;
  align-items: center;
  gap: 48px;
}

.navigation-desktop .nav-logo {
  display: flex;
  align-items: center;
  
  img {
    height: 32px;
    object-fit: contain;
  }
}

.navigation-desktop .nav-menu {
  display: flex;
  gap: 32px;
  margin: 0;
  padding: 0;
  list-style: none;
}

.navigation-desktop .nav-menu-item {
  font-size: 15px;
  color: #606266;
  cursor: pointer;
  transition: color 0.3s;

  &:hover, &.active {
    color: var(--el-color-primary);
  }
}

.navigation-desktop .nav-right {
  display: flex;
  align-items: center;
  gap: 24px;
}

.navigation-desktop .nav-actions {
  display: flex;
  align-items: center;
  gap: 12px;

  .el-button {
    margin: 0;
    height: 36px;
    min-width: 100px;
    padding: 0 16px;
    font-size: 14px;
    font-weight: 500;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
    
    &.designer {
      background-color: var(--el-color-primary);
      &:hover {
        background-color: var(--el-color-primary-light-3);
      }
    }
    &.factory {
      background-color: var(--el-color-success);
      &:hover {
        background-color: var(--el-color-success-light-3);
      }
    }
    &.supplier {
      background-color: var(--el-color-warning);
      &:hover {
        background-color: var(--el-color-warning-light-3);
      }
    }
  }
}

.navigation-desktop .language-selector {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  border-radius: 4px;
  background-color: #f5f7fa;
  cursor: pointer;
  transition: all 0.3s;

  &:hover {
    background-color: #e6e8eb;
  }
}
</style> 