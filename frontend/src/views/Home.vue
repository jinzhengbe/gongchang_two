<template>
  <div class="home">
    <Carousel />
    
    <div class="content-section">
      <el-row :gutter="20">
        <el-col :span="24">
          <h2 class="section-title">{{ $t('home.designerOrders') }}</h2>
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12" v-for="order in designerOrders" :key="order.id">
              <el-card class="order-card">
                <h3>{{ order.title }}</h3>
                <div class="order-info">
                  <el-tag>{{ $t(`order.status.${order.status}`) }}</el-tag>
                  <span class="order-date">{{ order.date }}</span>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </el-col>
      </el-row>

      <el-row :gutter="20" class="mt-4">
        <el-col :span="24">
          <h2 class="section-title">{{ $t('home.factories') }}</h2>
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12" v-for="factory in factories" :key="factory.id">
              <el-card class="factory-card">
                <h3>{{ factory.name }}</h3>
                <div class="factory-info">
                  <el-tag>{{ $t(`factory.type.${factory.type}`) }}</el-tag>
                  <span>{{ $t('home.capacity') }}: {{ factory.capacity }}</span>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </el-col>
      </el-row>

      <el-row :gutter="20" class="mt-4">
        <el-col :span="24">
          <h2 class="section-title">{{ $t('home.fabrics') }}</h2>
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12" v-for="fabric in fabrics" :key="fabric.id">
              <el-card class="fabric-card">
                <h3>{{ fabric.name }}</h3>
                <div class="fabric-info">
                  <el-tag>{{ $t(`fabric.type.${fabric.type}`) }}</el-tag>
                  <span>{{ fabric.material }}</span>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Carousel from '@/components/common/Carousel.vue'

const designerOrders = ref([
  {
    id: 1,
    title: '夏季连衣裙订单',
    status: 'pending',
    date: '2024-04-15'
  },
  {
    id: 2,
    title: '冬季外套订单',
    status: 'processing',
    date: '2024-04-14'
  }
])

const factories = ref([
  {
    id: 1,
    name: '优质服装厂',
    type: 'premium',
    capacity: '10000件/月'
  },
  {
    id: 2,
    name: '标准服装厂',
    type: 'standard',
    capacity: '5000件/月'
  }
])

const fabrics = ref([
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
])
</script>

<style scoped>
.home {
  padding: 20px;
}

.content-section {
  max-width: 1200px;
  margin: 40px auto;
}

.section-title {
  margin-bottom: 20px;
  font-size: 24px;
  color: var(--el-text-color-primary);
}

.order-card,
.factory-card,
.fabric-card {
  margin-bottom: 20px;
  transition: transform 0.3s;
}

.order-card:hover,
.factory-card:hover,
.fabric-card:hover {
  transform: translateY(-5px);
}

.order-info,
.factory-info,
.fabric-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 10px;
}

.order-date {
  color: var(--el-text-color-secondary);
  font-size: 14px;
}

.mt-4 {
  margin-top: 40px;
}

@media (max-width: 768px) {
  .home {
    padding: 10px;
  }

  .section-title {
    font-size: 20px;
  }
}
</style> 