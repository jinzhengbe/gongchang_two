<template>
  <div class="home">
    <!-- 顶部导航 -->
    <el-header class="header">
      <div class="logo">
        <img src="@/assets/logo.svg" alt="SewingMast" />
        <span>SewingMast</span>
      </div>
      <div class="nav">
        <el-menu mode="horizontal" :router="true">
          <el-menu-item index="/">{{ $t('nav.home') }}</el-menu-item>
          <el-menu-item index="/orders">{{ $t('nav.orders') }}</el-menu-item>
          <el-menu-item index="/factories">{{ $t('nav.factories') }}</el-menu-item>
          <el-menu-item index="/fabrics">{{ $t('nav.fabrics') }}</el-menu-item>
        </el-menu>
      </div>
      <div class="actions">
        <el-button type="primary" @click="$router.push('/login')">
          {{ $t('login.login') }}
        </el-button>
      </div>
    </el-header>

    <!-- 主要内容 -->
    <el-main class="main">
      <!-- 轮播展示区 -->
      <el-carousel :interval="4000" type="card" height="400px">
        <el-carousel-item v-for="item in carouselItems" :key="item.title">
          <div class="carousel-content">
            <h2>{{ $t(item.title) }}</h2>
            <p>{{ $t(item.description) }}</p>
            <el-button type="primary" @click="handleCarouselClick(item)">
              {{ $t('common.showMore') }}
            </el-button>
          </div>
        </el-carousel-item>
      </el-carousel>

      <!-- 内容板块 -->
      <div class="content-sections">
        <!-- 设计师订单区 -->
        <el-card class="section-card">
          <template #header>
            <div class="card-header">
              <span>{{ $t('home.designerOrders') }}</span>
            </div>
          </template>
          <div class="card-content">
            <el-row :gutter="20">
              <el-col :span="12" v-for="order in designerOrders" :key="order.id">
                <el-card shadow="hover" class="order-card">
                  <h3>{{ order.title }}</h3>
                  <p>{{ $t('order.status.' + order.status) }}</p>
                  <p>{{ order.createTime }}</p>
                </el-card>
              </el-col>
            </el-row>
          </div>
        </el-card>

        <!-- 工厂展示区 -->
        <el-card class="section-card">
          <template #header>
            <div class="card-header">
              <span>{{ $t('home.factories') }}</span>
            </div>
          </template>
          <div class="card-content">
            <el-row :gutter="20">
              <el-col :span="12" v-for="factory in factories" :key="factory.id">
                <el-card shadow="hover" class="factory-card">
                  <h3>{{ factory.name }}</h3>
                  <p>{{ $t('factory.type.' + factory.type) }}</p>
                  <p>{{ factory.capacity }}</p>
                </el-card>
              </el-col>
            </el-row>
          </div>
        </el-card>

        <!-- 布料展示区 -->
        <el-card class="section-card">
          <template #header>
            <div class="card-header">
              <span>{{ $t('home.fabrics') }}</span>
            </div>
          </template>
          <div class="card-content">
            <el-row :gutter="20">
              <el-col :span="12" v-for="fabric in fabrics" :key="fabric.id">
                <el-card shadow="hover" class="fabric-card">
                  <h3>{{ fabric.name }}</h3>
                  <p>{{ $t('fabric.type.' + fabric.type) }}</p>
                  <p>{{ fabric.material }}</p>
                </el-card>
              </el-col>
            </el-row>
          </div>
        </el-card>
      </div>
    </el-main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'

const router = useRouter()
const { t } = useI18n()

// 轮播数据
const carouselItems = [
  {
    title: 'carousel.latestOrders',
    description: 'carousel.latestOrdersDesc',
    path: '/orders'
  },
  {
    title: 'carousel.topFactories',
    description: 'carousel.topFactoriesDesc',
    path: '/factories'
  },
  {
    title: 'carousel.hotFabrics',
    description: 'carousel.hotFabricsDesc',
    path: '/fabrics'
  }
]

// 模拟数据
const designerOrders = [
  {
    id: 1,
    title: '夏季连衣裙订单',
    status: 'PENDING',
    createTime: '2024-04-15'
  },
  {
    id: 2,
    title: '冬季外套订单',
    status: 'PROCESSING',
    createTime: '2024-04-14'
  }
]

const factories = [
  {
    id: 1,
    name: '优质服装厂',
    type: 'premium',
    capacity: '月产能：10000件'
  },
  {
    id: 2,
    name: '标准服装厂',
    type: 'standard',
    capacity: '月产能：5000件'
  }
]

const fabrics = [
  {
    id: 1,
    name: '纯棉面料',
    type: 'cotton',
    material: '100%棉'
  },
  {
    id: 2,
    name: '真丝面料',
    type: 'silk',
    material: '100%桑蚕丝'
  }
]

const handleCarouselClick = (item: any) => {
  router.push(item.path)
}
</script>

<style scoped>
.home {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  background-color: #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo img {
  height: 40px;
}

.main {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}

.carousel-content {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  color: #fff;
  text-align: center;
  padding: 20px;
  background: linear-gradient(rgba(0, 0, 0, 0.5), rgba(0, 0, 0, 0.5));
}

.content-sections {
  margin-top: 40px;
}

.section-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-content {
  padding: 20px 0;
}

.order-card,
.factory-card,
.fabric-card {
  margin-bottom: 20px;
  cursor: pointer;
  transition: transform 0.3s;
}

.order-card:hover,
.factory-card:hover,
.fabric-card:hover {
  transform: translateY(-5px);
}

@media (max-width: 768px) {
  .header {
    flex-direction: column;
    padding: 10px;
  }

  .nav {
    margin: 10px 0;
  }

  .actions {
    margin-top: 10px;
  }

  .el-col {
    width: 100%;
  }
}
</style> 