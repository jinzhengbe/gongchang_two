<template>
  <el-card class="fabric-card" shadow="hover">
    <div class="fabric-image">
      <img :src="fabric.imageUrl" :alt="fabric.name">
    </div>
    <div class="fabric-info">
      <h3 class="fabric-name">{{ fabric.name }}</h3>
      <p class="fabric-desc">{{ fabric.description }}</p>
      <div class="fabric-price">
        <span class="price">¥{{ fabric.price }}</span>
        <span class="unit">/{{ fabric.unit }}</span>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
interface Fabric {
  id: number | string
  name: string
  description: string
  price: number
  unit: string
  imageUrl: string
}

defineProps<{
  fabric: Fabric
}>()
</script>

<style scoped lang="scss">
.fabric-card {
  height: 100%;
  display: flex;
  flex-direction: column;

  .fabric-image {
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

  .fabric-info {
    padding: 1rem;

    .fabric-name {
      font-size: 1.1rem;
      font-weight: 500;
      margin: 0 0 0.5rem;
      color: var(--el-text-color-primary);
    }

    .fabric-desc {
      font-size: 0.9rem;
      color: var(--el-text-color-secondary);
      margin: 0 0 1rem;
      line-height: 1.4;
    }

    .fabric-price {
      display: flex;
      align-items: baseline;
      gap: 0.25rem;

      .price {
        font-size: 1.2rem;
        font-weight: 500;
        color: var(--el-color-danger);
      }

      .unit {
        font-size: 0.9rem;
        color: var(--el-text-color-secondary);
      }
    }
  }
}

@media (max-width: 768px) {
  .fabric-card {
    .fabric-image {
      height: 160px;
    }

    .fabric-info {
      padding: 0.75rem;

      .fabric-name {
        font-size: 1rem;
      }

      .fabric-desc {
        font-size: 0.85rem;
      }

      .fabric-price {
        .price {
          font-size: 1.1rem;
        }

        .unit {
          font-size: 0.85rem;
        }
      }
    }
  }
}
</style> 