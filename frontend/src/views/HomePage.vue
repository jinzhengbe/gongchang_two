<template>
  <div class="home-page" :class="{ 'mobile': mobile }">
    <!-- 轮播展示区 -->
    <section class="carousel-section">
      <el-carousel height="400px">
        <el-carousel-item>
          <div class="carousel-content latest-orders">
            <h2>{{ $t('home.carousel.latestOrders') }}</h2>
            <div class="carousel-items">
              <OrderCard
                v-for="order in latestOrders"
                :key="order.id"
                :order="order"
              />
            </div>
          </div>
        </el-carousel-item>

        <el-carousel-item>
          <div class="carousel-content top-factories">
            <h2>{{ $t('home.carousel.topFactories') }}</h2>
            <div class="carousel-items">
              <FactoryCard
                v-for="factory in topFactories"
                :key="factory.id"
                :factory="factory"
              />
            </div>
          </div>
        </el-carousel-item>

        <el-carousel-item>
          <div class="carousel-content hot-fabrics">
            <h2>{{ $t('home.carousel.hotFabrics') }}</h2>
            <div class="carousel-items">
              <FabricCard
                v-for="fabric in hotFabrics"
                :key="fabric.id"
                :fabric="fabric"
              />
            </div>
          </div>
        </el-carousel-item>
      </el-carousel>
    </section>

    <!-- 内容板块 -->
    <section class="content-section">
      <!-- 设计师订单区 -->
      <div class="content-block designer-orders">
        <h2>{{ $t('home.content.designerOrders') }}</h2>
        <div class="order-list">
          <OrderCard
            v-for="order in designerOrders"
            :key="order.id"
            :order="order"
          />
        </div>
        <el-button type="text" @click="showMore('designer')">
          {{ $t('common.showMore') }}
        </el-button>
      </div>

      <!-- 工厂展示区 -->
      <div class="content-block factory-showcase">
        <h2>{{ $t('home.content.factoryShowcase') }}</h2>
        <div class="factory-list">
          <FactoryCard
            v-for="factory in factories"
            :key="factory.id"
            :factory="factory"
          />
        </div>
        <el-button type="text" @click="showMore('factory')">
          {{ $t('common.showMore') }}
        </el-button>
      </div>

      <!-- 布料展示区 -->
      <div class="content-block fabric-showcase">
        <h2>{{ $t('home.content.fabricShowcase') }}</h2>
        <div class="fabric-list">
          <FabricCard
            v-for="fabric in fabrics"
            :key="fabric.id"
            :fabric="fabric"
          />
        </div>
        <el-button type="text" @click="showMore('fabric')">
          {{ $t('common.showMore') }}
        </el-button>
      </div>
    </section>

    <!-- 返回顶部按钮 -->
    <el-backtop @click="handleScrollToTop" />

    <!-- 错误提示 -->
    <el-dialog
      v-model="error.show"
      :title="t('common.error')"
      width="30%"
    >
      <span>{{ error.message }}</span>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="error.show = false">{{ t('common.close') }}</el-button>
          <el-button type="primary" @click="retryLoading">
            {{ t('common.retry') }}
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { format } from 'date-fns'
import { zhCN } from 'date-fns/locale'
import { ElMessage } from 'element-plus'
import SkeletonLoader from '@/components/common/SkeletonLoader.vue'
import type { CarouselItem, Order, Factory, Fabric } from '@/types'
import { homeApi } from '@/api'
import { handleApiError } from '@/utils/error'
import { formatDate, formatNumber, isMobile, scrollToTop } from '@/utils/common'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import OrderCard from '@/components/designer/OrderCard.vue'
import FactoryCard from '@/components/designer/FactoryCard.vue'
import FabricCard from '@/components/designer/FabricCard.vue'

const { t } = useI18n()
const router = useRouter()
const authStore = useAuthStore()

// 加载状态
const loading = ref({
  carousel: true,
  orders: true,
  factories: true,
  fabrics: true
})

// 错误状态
const error = ref({
  show: false,
  message: '',
  type: ''
})

// 数据状态
const carouselItems = ref<CarouselItem[]>([])
const latestOrders = ref<Order[]>([])
const topFactories = ref<Factory[]>([])
const hotFabrics = ref<Fabric[]>([])
const designerOrders = ref<Order[]>([])
const factories = ref<Factory[]>([])
const fabrics = ref<Fabric[]>([])

// 加载数据
const loadData = async () => {
  try {
    // 加载轮播图数据
    loading.value.carousel = true
    const carouselResponse = await homeApi.getCarouselItems()
    carouselItems.value = carouselResponse.data
    loading.value.carousel = false

    // 加载订单数据
    loading.value.orders = true
    const ordersResponse = await homeApi.getLatestOrders()
    latestOrders.value = ordersResponse.data
    loading.value.orders = false

    // 加载工厂数据
    loading.value.factories = true
    const factoriesResponse = await homeApi.getRecommendedFactories()
    factories.value = factoriesResponse.data
    loading.value.factories = false

    // 加载面料数据
    loading.value.fabrics = true
    const fabricsResponse = await homeApi.getHotFabrics()
    fabrics.value = fabricsResponse.data
    loading.value.fabrics = false
  } catch (err) {
    handleApiError(err)
    error.value = {
      show: true,
      message: err.message || t('common.loadError'),
      type: 'error'
    }
  } finally {
    // 确保所有加载状态都被重置
    loading.value = {
      carousel: false,
      orders: false,
      factories: false,
      fabrics: false
    }
  }
}

const retryLoading = () => {
  error.value.show = false
  loadData()
}

const getStatusType = (status: string) => {
  const types = {
    PENDING: 'warning',
    PROCESSING: 'primary'
  }
  return types[status] || 'info'
}

const formatCapacity = (capacity: number): string => {
  return t('factory.monthlyCapacity', { capacity: formatNumber(capacity) })
}

const handleShowMore = (item: CarouselItem) => {
  const route = item.type === 'order' ? '/orders' : 
                item.type === 'factory' ? '/factories' : '/fabrics'
  router.push(route)
}

// 检查是否为移动设备
const mobile = ref(isMobile())

// 监听窗口大小变化
onMounted(() => {
  loadData()
  
  // 监听窗口大小变化
  window.addEventListener('resize', () => {
    mobile.value = isMobile()
  })

  // 实现图片懒加载
  const lazyImages = document.querySelectorAll('img[data-src]')
  const imageObserver = new IntersectionObserver((entries, observer) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        const img = entry.target as HTMLImageElement
        img.src = img.dataset.src || ''
        observer.unobserve(img)
      }
    })
  })

  lazyImages.forEach(img => imageObserver.observe(img))
})

// 在组件卸载时移除事件监听
onUnmounted(() => {
  window.removeEventListener('resize', () => {
    mobile.value = isMobile()
  })
})

// 处理返回顶部
const handleScrollToTop = () => {
  scrollToTop()
}

// 图片懒加载
const observeImages = () => {
  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        const img = entry.target as HTMLImageElement
        if (img.dataset.src) {
          img.src = img.dataset.src
          observer.unobserve(img)
        }
      }
    })
  })

  document.querySelectorAll('img[data-src]').forEach(img => {
    observer.observe(img)
  })
}

// 响应式处理
const handleResize = () => {
  mobile.value = isMobile()
}

onMounted(() => {
  observeImages()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})

// 显示更多
const showMore = (type: string) => {
  if (!authStore.isLoggedIn) {
    router.push(`/login/${type}`)
  } else {
    router.push(`/${type}/dashboard`)
  }
}
</script>

<style scoped>
.home-page {
  padding: 20px;
}

.home-page.mobile {
  padding: 10px;
}

.section {
  margin-bottom: 40px;
}

.section-title {
  font-size: 24px;
  margin-bottom: 20px;
  color: var(--el-text-color-primary);
}

.mobile .section-title {
  font-size: 20px;
  margin-bottom: 15px;
}

.carousel {
  margin-bottom: 40px;
  border-radius: 8px;
  overflow: hidden;
}

.carousel-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.order-card,
.factory-card,
.fabric-card {
  height: 100%;
  transition: transform 0.3s ease;
}

.order-card:hover,
.factory-card:hover,
.fabric-card:hover {
  transform: translateY(-5px);
}

.order-image,
.factory-image,
.fabric-image {
  height: 200px;
  overflow: hidden;
}

.mobile .order-image,
.mobile .factory-image,
.mobile .fabric-image {
  height: 150px;
}

.order-image img,
.factory-image img,
.fabric-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.order-info,
.factory-info,
.fabric-info {
  padding: 15px;
}

.mobile .order-info,
.mobile .factory-info,
.mobile .fabric-info {
  padding: 10px;
}

h3 {
  margin: 0 0 10px;
  font-size: 18px;
  color: var(--el-text-color-primary);
}

.mobile h3 {
  font-size: 16px;
}

p {
  margin: 5px 0;
  color: var(--el-text-color-regular);
}

.order-date,
.factory-location,
.fabric-type {
  color: var(--el-text-color-secondary);
  font-size: 14px;
}

.order-quantity,
.factory-capacity,
.fabric-price {
  color: var(--el-color-primary);
  font-weight: bold;
}

.carousel-section {
  margin-bottom: 3rem;
}

.carousel-content {
  height: 100%;
  padding: 2rem;
  background: #f8f9fa;
}

.carousel-items {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
  margin-top: 1rem;
}

.content-section {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 2rem;
}

.content-block {
  background: #fff;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.content-block h2 {
  margin-bottom: 1.5rem;
  color: #333;
}

/* Responsive styles */
@media (max-width: 1024px) {
  .content-section {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .content-section {
    grid-template-columns: 1fr;
  }

  .carousel-items {
    grid-template-columns: 1fr;
  }

  .home-page {
    padding: 1rem;
  }
}
</style> 