<template>
  <div class="carousel-container">
    <el-carousel height="400px" :interval="5000" arrow="always">
      <el-carousel-item v-for="(item, index) in carouselItems" :key="index">
        <div class="carousel-content" :style="{ backgroundColor: item.color }">
          <div class="carousel-text">
            <h2>{{ $t(item.title) }}</h2>
            <p>{{ $t(item.description) }}</p>
            <el-button type="primary" @click="handleShowMore(item)">
              {{ $t('common.showMore') }}
            </el-button>
          </div>
        </div>
      </el-carousel-item>
    </el-carousel>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const carouselItems = ref([
  {
    title: 'carousel.hotFabrics',
    description: 'carousel.hotFabricsDesc',
    color: '#409EFF',
    link: '/fabrics'
  },
  {
    title: 'carousel.latestOrders',
    description: 'carousel.latestOrdersDesc',
    color: '#67C23A',
    link: '/orders'
  },
  {
    title: 'carousel.topFactories',
    description: 'carousel.topFactoriesDesc',
    color: '#E6A23C',
    link: '/factories'
  }
])

const handleShowMore = (item: typeof carouselItems.value[0]) => {
  router.push(item.link)
}
</script>

<style scoped>
.carousel-container {
  width: 100%;
  margin: 20px 0;
}

.carousel-content {
  height: 100%;
  display: flex;
  align-items: center;
  position: relative;
}

.carousel-text {
  position: relative;
  color: white;
  padding: 0 50px;
  max-width: 600px;
  z-index: 1;
}

.carousel-text h2 {
  font-size: 2.5rem;
  margin-bottom: 1rem;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
}

.carousel-text p {
  font-size: 1.2rem;
  margin-bottom: 2rem;
  line-height: 1.6;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.3);
}

@media (max-width: 768px) {
  .carousel-text h2 {
    font-size: 1.8rem;
  }

  .carousel-text p {
    font-size: 1rem;
  }

  .carousel-text {
    padding: 0 20px;
  }
}
</style> 