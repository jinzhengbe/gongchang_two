<template>
  <el-card class="factory-card" shadow="hover">
    <div class="factory-image">
      <img :src="factory.imageUrl" :alt="factory.name">
    </div>
    <div class="factory-info">
      <h3 class="factory-name">{{ factory.name }}</h3>
      <p class="factory-desc">{{ factory.description }}</p>
      <div class="factory-meta">
        <div class="rating">
          <el-rate
            v-model="factory.rating"
            disabled
            show-score
            text-color="#ff9900"
            score-template="{value}"
          />
        </div>
        <div class="location">
          <i class="el-icon-location"></i>
          <span>{{ factory.location }}</span>
        </div>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
interface Factory {
  id: number | string
  name: string
  description: string
  rating: number
  location: string
  imageUrl: string
}

defineProps<{
  factory: Factory
}>()
</script>

<style scoped lang="scss">
.factory-card {
  height: 100%;
  display: flex;
  flex-direction: column;

  .factory-image {
    width: 100%;
    height: 200px;
    overflow: hidden;
    
    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      transition: transform 0.3s ease;

      &:hover {
        transform: scale(1.05);
      }
    }
  }

  .factory-info {
    padding: 1rem;

    .factory-name {
      font-size: 1.1rem;
      font-weight: 500;
      margin: 0 0 0.5rem;
      color: var(--el-text-color-primary);
    }

    .factory-desc {
      font-size: 0.9rem;
      color: var(--el-text-color-secondary);
      margin: 0 0 1rem;
      line-height: 1.4;
    }

    .factory-meta {
      .rating {
        margin-bottom: 0.5rem;
      }

      .location {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: 0.9rem;
        color: var(--el-text-color-secondary);

        i {
          color: var(--el-color-primary);
        }
      }
    }
  }
}

@media (max-width: 768px) {
  .factory-card {
    .factory-image {
      height: 160px;
    }

    .factory-info {
      padding: 0.75rem;

      .factory-name {
        font-size: 1rem;
      }

      .factory-desc {
        font-size: 0.85rem;
      }

      .factory-meta {
        .location {
          font-size: 0.85rem;
        }
      }
    }
  }
}
</style> 