<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-header">
        <img src="@/assets/logo.svg" alt="SewingMast" class="logo" />
        <h1>{{ $t('login.title') }}</h1>
      </div>

      <el-form
        ref="loginForm"
        :model="loginForm"
        :rules="loginRules"
        label-position="top"
        @submit.prevent="handleLogin"
      >
        <el-form-item :label="$t('login.username')" prop="username">
          <el-input
            v-model="loginForm.username"
            :placeholder="$t('login.usernamePlaceholder')"
          />
        </el-form-item>

        <el-form-item :label="$t('login.password')" prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            :placeholder="$t('login.passwordPlaceholder')"
            show-password
          />
        </el-form-item>

        <el-button
          type="primary"
          native-type="submit"
          :loading="loading"
          class="login-button"
        >
          {{ $t('login.login') }}
        </el-button>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const { t } = useI18n()
const authStore = useAuthStore()

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules = reactive<FormRules>({
  username: [
    { required: true, message: t('login.usernameRequired'), trigger: 'blur' },
    { min: 3, max: 20, message: t('login.usernameLength'), trigger: 'blur' }
  ],
  password: [
    { required: true, message: t('login.passwordRequired'), trigger: 'blur' },
    { min: 6, max: 20, message: t('login.passwordLength'), trigger: 'blur' }
  ]
})

const loading = ref(false)
const loginFormRef = ref<FormInstance>()

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    await loginFormRef.value.validate()
    loading.value = true
    
    const success = await authStore.login(loginForm.username, loginForm.password)
    
    if (success) {
      ElMessage.success(t('login.success'))
      router.push('/dashboard')
    } else {
      ElMessage.error(t('login.error'))
    }
  } catch (error) {
    console.error('Login failed:', error)
    ElMessage.error(t('login.error'))
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
}

.login-container {
  width: 100%;
  max-width: 400px;
  padding: 40px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo {
  height: 60px;
  margin-bottom: 20px;
}

.login-button {
  width: 100%;
  margin-top: 20px;
}

@media (max-width: 768px) {
  .login-container {
    padding: 20px;
    margin: 20px;
  }
}
</style> 