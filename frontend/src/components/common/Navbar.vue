<template>
  <el-menu
    :default-active="activeIndex"
    class="navbar"
    mode="horizontal"
    router
  >
    <div class="logo-container">
      <img src="@/assets/logo.svg" alt="SewingMast" class="logo" />
    </div>

    <div class="menu-items">
      <el-menu-item index="/">{{ $t('nav.home') }}</el-menu-item>
      <el-menu-item index="/orders">{{ $t('nav.orders') }}</el-menu-item>
      <el-menu-item index="/factories">{{ $t('nav.factories') }}</el-menu-item>
    </div>

    <div class="right-menu">
      <el-dropdown v-if="isLoggedIn" trigger="click">
        <span class="user-profile">
          <el-avatar :size="32" :src="userAvatar" />
          <span class="username">{{ username }}</span>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item @click="goToProfile">
              {{ $t('nav.profile') }}
            </el-dropdown-item>
            <el-dropdown-item @click="handleLogout">
              {{ $t('nav.logout') }}
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>

      <el-button v-else type="primary" @click="goToLogin">
        {{ $t('nav.login') }}
      </el-button>

      <el-dropdown trigger="click" class="language-dropdown">
        <span class="language-selector">
          <el-icon><Monitor /></el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item
              v-for="lang in languages"
              :key="lang.value"
              @click="changeLanguage(lang.value)"
            >
              <el-icon><Monitor /></el-icon>
              {{ lang.label }}
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </el-menu>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Monitor } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const { t, locale } = useI18n()
const authStore = useAuthStore()

const activeIndex = ref('/')
const isLoggedIn = computed(() => authStore.isLoggedIn)
const username = computed(() => authStore.user?.username || '')
const userAvatar = computed(() => authStore.user?.avatar || '')

const languages = [
  { label: '简体中文', value: 'zh' },
  { label: 'English', value: 'en' }
]

const goToLogin = () => {
  router.push('/login')
}

const goToProfile = () => {
  router.push('/profile')
}

const handleLogout = async () => {
  await authStore.logout()
  router.push('/')
}

const changeLanguage = (lang: string) => {
  locale.value = lang
}
</script>

<style scoped>
.navbar {
  display: flex;
  align-items: center;
  padding: 0 20px;
  height: 60px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.logo-container {
  display: flex;
  align-items: center;
  margin-right: 40px;
}

.logo {
  height: 40px;
}

.menu-items {
  flex: 1;
  display: flex;
}

.right-menu {
  display: flex;
  align-items: center;
  gap: 20px;
}

.user-profile {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.username {
  font-size: 14px;
}

.language-selector {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 8px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.language-selector:hover {
  background-color: var(--el-color-primary-light-9);
}

@media (max-width: 768px) {
  .navbar {
    padding: 0 10px;
  }

  .logo-container {
    margin-right: 20px;
  }

  .username {
    display: none;
  }
}
</style> 