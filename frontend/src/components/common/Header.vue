<template>
  <header class="header">
    <nav class="nav-container">
      <div class="logo">
        <img src="@/assets/logo.png" alt="SewingMast" />
      </div>

      <div class="nav-menu">
        <router-link to="/" class="nav-item">{{ $t('nav.home') }}</router-link>
        <router-link to="/orders" class="nav-item">{{ $t('nav.orders') }}</router-link>
        <router-link to="/factories" class="nav-item">{{ $t('nav.factories') }}</router-link>
      </div>

      <div class="nav-right">
        <div class="language-switcher">
          <el-dropdown @command="handleLanguageChange">
            <span class="el-dropdown-link">
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

        <div class="login-entries">
          <template v-if="!isLoggedIn">
            <el-dropdown>
              <span class="el-dropdown-link">
                {{ $t('auth.designer') }}
                <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="navigateToLogin('designer')">
                    {{ $t('auth.login') }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>

            <el-dropdown>
              <span class="el-dropdown-link">
                {{ $t('auth.factory') }}
                <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="navigateToLogin('factory')">
                    {{ $t('auth.login') }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>

            <el-dropdown>
              <span class="el-dropdown-link">
                {{ $t('auth.fabric') }}
                <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="navigateToLogin('fabric')">
                    {{ $t('auth.login') }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>

          <template v-else>
            <el-dropdown>
              <span class="el-dropdown-link">
                {{ userInfo.name }}
                <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="navigateToDashboard">
                    {{ $t('auth.dashboard') }}
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleLogout">
                    {{ $t('auth.logout') }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </div>
      </div>
    </nav>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { ArrowDown } from '@element-plus/icons-vue'

const router = useRouter()
const { t, locale } = useI18n()
const authStore = useAuthStore()

const isLoggedIn = computed(() => authStore.isLoggedIn)
const userInfo = computed(() => authStore.userInfo)
const currentLanguage = computed(() => {
  const langMap = {
    'zh-CN': '中文',
    'en-US': 'English',
    'ko-KR': '한국어',
    'vi-VN': 'Tiếng Việt'
  }
  return langMap[locale.value] || 'Language'
})

const handleLanguageChange = (lang: string) => {
  locale.value = lang
}

const navigateToLogin = (role: string) => {
  router.push(`/login/${role}`)
}

const navigateToDashboard = () => {
  const role = authStore.userInfo?.role
  router.push(`/${role}/dashboard`)
}

const handleLogout = async () => {
  await authStore.logout()
  router.push('/')
}
</script>

<style scoped>
.header {
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 1rem;
}

.nav-container {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo img {
  height: 40px;
}

.nav-menu {
  display: flex;
  gap: 2rem;
}

.nav-item {
  text-decoration: none;
  color: #333;
  font-weight: 500;
}

.nav-right {
  display: flex;
  align-items: center;
  gap: 2rem;
}

.login-entries {
  display: flex;
  gap: 1rem;
}

/* Responsive styles */
@media (max-width: 768px) {
  .nav-menu {
    display: none;
  }

  .nav-right {
    gap: 1rem;
  }

  .login-entries {
    flex-direction: column;
    gap: 0.5rem;
  }
}
</style> 