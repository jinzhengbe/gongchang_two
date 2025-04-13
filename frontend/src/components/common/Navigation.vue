<template>
  <nav class="navigation">
    <div class="nav-container">
      <div class="nav-logo">
        <router-link to="/">
          <img src="@/assets/logo.svg" alt="SewingMast" />
        </router-link>
      </div>

      <el-menu
        :default-active="activeRoute"
        mode="horizontal"
        class="nav-menu"
        @select="handleSelect"
      >
        <el-menu-item index="/">{{ t('nav.home') }}</el-menu-item>
        <el-menu-item index="/orders">{{ t('nav.orders') }}</el-menu-item>
        <el-menu-item index="/factories">{{ t('nav.factories') }}</el-menu-item>
        <el-menu-item index="/fabrics">{{ t('nav.fabrics') }}</el-menu-item>
      </el-menu>

      <div class="nav-actions">
        <el-dropdown trigger="click" class="language-dropdown">
          <span class="language-selector">
            {{ currentLanguage }}
            <el-icon><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="changeLanguage('zh')">中文</el-dropdown-item>
              <el-dropdown-item @click="changeLanguage('en')">English</el-dropdown-item>
              <el-dropdown-item @click="changeLanguage('ko')">한국어</el-dropdown-item>
              <el-dropdown-item @click="changeLanguage('vi')">Tiếng Việt</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>

        <div class="login-section">
          <el-dropdown v-for="role in roles" :key="role.value" class="login-dropdown">
            <el-button type="primary" :class="role.value">
              {{ t(`login.${role.label}`) }}
              <el-icon><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleLogin(role.value)">
                  {{ t('login.login') }}
                </el-dropdown-item>
                <el-dropdown-item @click="handleRegister(role.value)">
                  {{ t('login.register') }}
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>

      <!-- 移动端菜单按钮 -->
      <div class="mobile-menu-button" @click="toggleMobileMenu">
        <el-icon><Menu /></el-icon>
      </div>
    </div>

    <!-- 移动端菜单 -->
    <transition name="slide-fade">
      <div v-if="showMobileMenu" class="mobile-menu">
        <el-menu mode="vertical" @select="handleSelect">
          <el-menu-item index="/">{{ t('nav.home') }}</el-menu-item>
          <el-menu-item index="/orders">{{ t('nav.orders') }}</el-menu-item>
          <el-menu-item index="/factories">{{ t('nav.factories') }}</el-menu-item>
          <el-menu-item index="/fabrics">{{ t('nav.fabrics') }}</el-menu-item>
        </el-menu>
        <div class="mobile-actions">
          <el-dropdown trigger="click" class="language-dropdown">
            <el-button>{{ currentLanguage }}</el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="changeLanguage('zh')">中文</el-dropdown-item>
                <el-dropdown-item @click="changeLanguage('en')">English</el-dropdown-item>
                <el-dropdown-item @click="changeLanguage('ko')">한국어</el-dropdown-item>
                <el-dropdown-item @click="changeLanguage('vi')">Tiếng Việt</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <div class="mobile-login-buttons">
            <el-button 
              v-for="role in roles" 
              :key="role.value"
              type="primary"
              :class="role.value"
              @click="handleLogin(role.value)"
            >
              {{ t(`login.${role.label}`) }}
            </el-button>
          </div>
        </div>
      </div>
    </transition>
  </nav>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter, useRoute } from 'vue-router'
import { ArrowDown, Menu } from '@element-plus/icons-vue'

const { t, locale } = useI18n()
const router = useRouter()
const route = useRoute()

const showMobileMenu = ref(false)

const roles = [
  { label: 'designer', value: 'designer' },
  { label: 'factory', value: 'factory' },
  { label: 'supplier', value: 'supplier' }
]

const activeRoute = computed(() => route.path)

const currentLanguage = computed(() => {
  const languages = {
    zh: '中文',
    en: 'English',
    ko: '한국어',
    vi: 'Tiếng Việt'
  }
  return languages[locale.value] || languages.en
})

const handleSelect = (path: string) => {
  router.push(path)
  showMobileMenu.value = false
}

const changeLanguage = (lang: string) => {
  locale.value = lang
}

const handleLogin = (role: string) => {
  router.push(`/login/${role}`)
  showMobileMenu.value = false
}

const handleRegister = (role: string) => {
  router.push(`/register/${role}`)
  showMobileMenu.value = false
}

const toggleMobileMenu = () => {
  showMobileMenu.value = !showMobileMenu.value
}
</script>

<style scoped lang="scss">
.navigation {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  z-index: 1000;
}

.nav-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nav-logo {
  img {
    height: 40px;
    object-fit: contain;
  }
}

.nav-menu {
  flex: 1;
  margin: 0 40px;
  border-bottom: none;
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 20px;
}

.language-dropdown {
  .language-selector {
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 4px;
  }
}

.login-section {
  display: flex;
  gap: 10px;

  .login-dropdown {
    .el-button {
      &.designer {
        background-color: #409EFF;
      }
      &.factory {
        background-color: #67C23A;
      }
      &.supplier {
        background-color: #E6A23C;
      }
    }
  }
}

.mobile-menu-button {
  display: none;
  cursor: pointer;
  padding: 8px;
}

.mobile-menu {
  display: none;
  background: white;
  padding: 20px;
  border-top: 1px solid #eee;
}

.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.3s ease;
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateY(-20px);
  opacity: 0;
}

@media (max-width: 1024px) {
  .nav-menu {
    margin: 0 20px;
  }
}

@media (max-width: 768px) {
  .nav-menu,
  .nav-actions {
    display: none;
  }

  .mobile-menu-button {
    display: block;
  }

  .mobile-menu {
    display: block;

    .el-menu {
      border-right: none;
    }

    .mobile-actions {
      margin-top: 20px;
      display: flex;
      flex-direction: column;
      gap: 10px;
    }

    .mobile-login-buttons {
      display: flex;
      flex-direction: column;
      gap: 10px;
    }
  }
}
</style> 