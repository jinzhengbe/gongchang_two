<template>
  <div class="home-page">
    <!-- 轮播图部分 -->
    <div class="carousel-section">
      <el-carousel :interval="5000" height="500px" class="main-carousel">
        <el-carousel-item v-for="item in carouselItems" :key="item.id">
          <img :src="item.imageUrl" :alt="item.title" class="carousel-image">
          <div class="carousel-content">
            <h2>{{ item.title }}</h2>
            <p>{{ item.description }}</p>
          </div>
        </el-carousel-item>
      </el-carousel>
    </div>

    <div class="container">
      <!-- 新品面料 -->
      <div class="section">
        <div class="section-header">
          <h2>{{ $t('home.newFabrics') }}</h2>
          <p class="section-desc">{{ $t('home.newFabricsDesc') }}</p>
        </div>
        <div class="fabric-grid">
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12" :md="8" :lg="6" v-for="fabric in newFabrics" :key="fabric.id">
              <FabricCard :fabric="fabric" />
            </el-col>
          </el-row>
        </div>
      </div>

      <!-- 最新订单 -->
      <div class="section">
        <LatestOrders 
          :orders="latestOrders" 
          @view-more="handleViewMore"
          @select-order="handleSelectOrder"
        />
      </div>

      <!-- 热门工厂 -->
      <div class="section">
        <div class="section-header">
          <h2>{{ $t('home.topFactories') }}</h2>
          <p class="section-desc">{{ $t('home.topFactoriesDesc') }}</p>
        </div>
        <div class="factory-list">
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12" :md="8" :lg="6" v-for="factory in topFactories" :key="factory.id">
              <FactoryCard :factory="factory" />
            </el-col>
          </el-row>
        </div>
      </div>

      <!-- 热门面料 -->
      <div class="section">
        <div class="section-header">
          <h2>{{ $t('home.hotFabrics') }}</h2>
          <p class="section-desc">{{ $t('home.hotFabricsDesc') }}</p>
        </div>
        <div class="fabric-grid">
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12" :md="8" :lg="6" v-for="fabric in hotFabrics" :key="fabric.id">
              <FabricCard :fabric="fabric" />
            </el-col>
          </el-row>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { homeApi } from '@/api'
import FabricCard from '@/components/fabric/FabricCard.vue'
import FactoryCard from '@/components/factory/FactoryCard.vue'
import LatestOrders from '@/components/order/LatestOrders.vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'

const router = useRouter()
const carouselItems = ref([])
const hotFabrics = ref([])
const newFabrics = ref([])
const topFactories = ref([])
const latestOrders = ref([])

const getStatusType = (status: string) => {
  const types: Record<string, string> = {
    'pending': 'warning',
    'processing': 'primary',
    'completed': 'success',
    'cancelled': 'danger'
  }
  return types[status] || 'info'
}

const fetchHomeData = async () => {
  try {
    const [
      carouselResponse, 
      newFabricsResponse,
      hotFabricsResponse, 
      factoriesResponse,
      latestOrdersResponse
    ] = await Promise.all([
      homeApi.getCarouselItems(),
      homeApi.getNewFabrics(),
      homeApi.getHotFabrics(),
      homeApi.getRecommendedFactories(),
      homeApi.getLatestOrders()
    ])

    carouselItems.value = carouselResponse.data || []
    newFabrics.value = newFabricsResponse.data || []
    hotFabrics.value = hotFabricsResponse.data || []
    topFactories.value = factoriesResponse.data || []
    latestOrders.value = latestOrdersResponse.data || []

    // 测试数据
    if (!carouselItems.value || carouselItems.value.length === 0) {
      carouselItems.value = [
        {
          id: 1,
          title: '2024春季新品发布',
          description: '探索最新的时尚趋势和创新设计',
          imageUrl: 'https://img.freepik.com/free-photo/fashion-clothes-hanging-rack_23-2147669975.jpg'
        },
        {
          id: 2,
          title: '优质工厂直供',
          description: '连接优质工厂，打造高品质服装',
          imageUrl: 'https://img.freepik.com/free-photo/clothing-factory_23-2147669986.jpg'
        }
      ]
    }

    if (!hotFabrics.value || hotFabrics.value.length === 0) {
      hotFabrics.value = [
        {
          id: 1,
          name: '优质棉料',
          description: '100%纯棉面料，柔软舒适',
          price: 58,
          unit: '米',
          imageUrl: 'https://img.freepik.com/free-photo/cotton-fabric_1203-5624.jpg'
        },
        {
          id: 2,
          name: '真丝面料',
          description: '高档真丝面料，光泽优雅',
          price: 128,
          unit: '米',
          imageUrl: 'https://img.freepik.com/free-photo/silk-fabric_1203-5625.jpg'
        },
        {
          id: 3,
          name: '亚麻面料',
          description: '天然亚麻，清爽透气',
          price: 88,
          unit: '米',
          imageUrl: 'https://img.freepik.com/free-photo/linen-fabric_1203-5626.jpg'
        }
      ]
    }

    if (!newFabrics.value || newFabrics.value.length === 0) {
      newFabrics.value = [
        {
          id: 4,
          name: '新款提花面料',
          description: '精美提花工艺，立体质感',
          price: 98,
          unit: '米',
          imageUrl: 'https://img.freepik.com/free-photo/jacquard-fabric_1203-5627.jpg'
        },
        {
          id: 5,
          name: '弹力针织面料',
          description: '高弹力，适合运动服装',
          price: 68,
          unit: '米',
          imageUrl: 'https://img.freepik.com/free-photo/knit-fabric_1203-5628.jpg'
        },
        {
          id: 6,
          name: '印花雪纺面料',
          description: '时尚印花，飘逸轻盈',
          price: 78,
          unit: '米',
          imageUrl: 'https://img.freepik.com/free-photo/chiffon-fabric_1203-5629.jpg'
        }
      ]
    }

    if (!topFactories.value || topFactories.value.length === 0) {
      topFactories.value = [
        {
          id: 1,
          name: '优质服装厂',
          description: '专业女装生产，十年经验',
          rating: 4.8,
          location: '广东省深圳市',
          imageUrl: 'https://img.freepik.com/free-photo/clothing-factory-1_23-2147669987.jpg'
        },
        {
          id: 2,
          name: '精品制衣厂',
          description: '高端定制专家，品质保证',
          rating: 4.9,
          location: '浙江省杭州市',
          imageUrl: 'https://img.freepik.com/free-photo/clothing-factory-2_23-2147669988.jpg'
        },
        {
          id: 3,
          name: '时尚服饰厂',
          description: '快时尚生产线，交期准时',
          rating: 4.7,
          location: '江苏省苏州市',
          imageUrl: 'https://img.freepik.com/free-photo/clothing-factory-3_23-2147669989.jpg'
        }
      ]
    }

    if (!latestOrders.value || latestOrders.value.length === 0) {
      latestOrders.value = [
        {
          id: '1',
          title: '春季女装系列订单',
          description: '2024春季新款连衣裙系列生产订单',
          status: 'pending',
          quantity: 2000,
          price: 99.99,
          deadline: '2024-05-15',
          progress: 0,
          imageUrl: 'https://img.freepik.com/free-photo/spring-collection-dresses_23-2148769754.jpg'
        },
        {
          id: '2',
          title: '儿童校服套装',
          description: '中小学生春秋季校服套装定制',
          status: 'processing',
          quantity: 1500,
          price: 59.99,
          deadline: '2024-05-01',
          progress: 45,
          imageUrl: 'https://img.freepik.com/free-photo/school-uniforms_23-2148769755.jpg'
        },
        {
          id: '3',
          title: '男士运动装订单',
          description: '专业运动服装系列生产',
          status: 'completed',
          quantity: 1000,
          price: 79.99,
          deadline: '2024-04-10',
          progress: 100,
          imageUrl: 'https://img.freepik.com/free-photo/sportswear-collection_23-2148769756.jpg'
        },
        {
          id: '4',
          title: '工作制服定制',
          description: '企业员工制服批量定制订单',
          status: 'pending',
          quantity: 800,
          price: 129.99,
          deadline: '2024-05-20',
          progress: 0,
          imageUrl: 'https://img.freepik.com/free-photo/corporate-uniforms_23-2148769757.jpg'
        }
      ]
    }
  } catch (error) {
    console.error('Failed to fetch home data:', error)
    ElMessage.error('获取数据失败，请稍后重试')
  }
}

const handleViewMore = () => {
  router.push('/orders')
}

const handleSelectOrder = (order: any) => {
  router.push(`/orders/${order.id}`)
}

onMounted(() => {
  fetchHomeData()
})
</script>

<style scoped lang="scss">
.home-page {
  min-height: 100vh;
  background-color: #f5f7fa;
  padding-bottom: 2rem;
}

.main-carousel {
  margin-bottom: 2rem;
}

.carousel-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.carousel-content {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  text-align: center;
  color: white;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
  width: 80%;
  max-width: 800px;

  h2 {
    font-size: 2.5rem;
    margin-bottom: 1rem;
  }

  p {
    font-size: 1.25rem;
  }
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
}

.section {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  margin-bottom: 2rem;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
}

.section-header {
  text-align: center;
  margin-bottom: 2rem;

  h2 {
    font-size: 1.875rem;
    font-weight: 600;
    color: var(--el-text-color-primary);
    margin-bottom: 0.5rem;
  }

  .section-desc {
    color: var(--el-text-color-secondary);
    font-size: 1rem;
  }
}

.order-card {
  margin-bottom: 1rem;
  
  .order-status {
    display: flex;
    justify-content: flex-end;
    margin-bottom: 0.5rem;
  }

  .order-title {
    font-size: 1.1rem;
    margin: 0.5rem 0;
    color: var(--el-text-color-primary);
  }

  .order-info {
    color: var(--el-text-color-regular);
    font-size: 0.9rem;

    p {
      margin: 0.3rem 0;
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }

    i {
      color: var(--el-color-primary);
    }
  }

  .order-progress {
    margin-top: 1rem;
    
    span {
      font-size: 0.9rem;
      color: var(--el-text-color-secondary);
      margin-bottom: 0.3rem;
      display: block;
    }
  }
}

@media (max-width: 768px) {
  .carousel-content {
    h2 {
      font-size: 1.75rem;
    }

    p {
      font-size: 1rem;
    }
  }

  .container {
    padding: 0 0.5rem;
  }

  .section {
    padding: 1.5rem 1rem;
  }

  .section-header {
    h2 {
      font-size: 1.5rem;
    }

    .section-desc {
      font-size: 0.875rem;
    }
  }
}
</style> 