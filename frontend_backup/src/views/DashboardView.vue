<template>
  <div class="dashboard">
    <el-container>
      <el-aside width="200px">
        <el-menu
          :default-active="activeMenu"
          class="el-menu-vertical"
          @select="handleMenuSelect"
        >
          <el-menu-item index="overview">
            <el-icon><DataLine /></el-icon>
            <span>概览</span>
          </el-menu-item>
          <el-menu-item index="orders">
            <el-icon><List /></el-icon>
            <span>订单管理</span>
          </el-menu-item>
          <el-menu-item index="products">
            <el-icon><Goods /></el-icon>
            <span>产品管理</span>
          </el-menu-item>
          <el-menu-item index="profile">
            <el-icon><User /></el-icon>
            <span>个人中心</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <el-container>
        <el-header>
          <div class="header-content">
            <h2>{{ pageTitle }}</h2>
            <div class="user-info">
              <el-dropdown @command="handleCommand">
                <span class="user-dropdown">
                  {{ username }}
                  <el-icon class="el-icon--right"><ArrowDown /></el-icon>
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                    <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>
        </el-header>

        <el-main>
          <router-view></router-view>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import {
  DataLine,
  List,
  Goods,
  User,
  ArrowDown
} from '@element-plus/icons-vue'

export default defineComponent({
  name: 'DashboardView',
  components: {
    DataLine,
    List,
    Goods,
    User,
    ArrowDown
  },
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    const activeMenu = ref('overview')

    const username = computed(() => {
      return authStore.currentUser?.username || '用户'
    })

    const pageTitle = computed(() => {
      switch (activeMenu.value) {
        case 'overview':
          return '概览'
        case 'orders':
          return '订单管理'
        case 'products':
          return '产品管理'
        case 'profile':
          return '个人中心'
        default:
          return '仪表板'
      }
    })

    const handleMenuSelect = (index: string) => {
      activeMenu.value = index
      router.push(`/${index}`)
    }

    const handleCommand = (command: string) => {
      if (command === 'logout') {
        authStore.logout()
        router.push('/login')
      } else if (command === 'profile') {
        router.push('/profile')
      }
    }

    return {
      activeMenu,
      username,
      pageTitle,
      handleMenuSelect,
      handleCommand
    }
  }
})
</script>

<style scoped>
.dashboard {
  height: 100vh;
}

.el-container {
  height: 100%;
}

.el-aside {
  background-color: #304156;
}

.el-menu {
  border-right: none;
  background-color: transparent;
}

.el-menu-item {
  color: #bfcbd9;
}

.el-menu-item:hover {
  color: #fff;
}

.el-menu-item.is-active {
  color: #409eff;
  background-color: #263445;
}

.el-header {
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  padding: 0 20px;
}

.header-content {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.user-info {
  display: flex;
  align-items: center;
}

.user-dropdown {
  cursor: pointer;
  display: flex;
  align-items: center;
  color: #606266;
}

.user-dropdown:hover {
  color: #409eff;
}

.el-main {
  background-color: #f0f2f5;
  padding: 20px;
}
</style> 