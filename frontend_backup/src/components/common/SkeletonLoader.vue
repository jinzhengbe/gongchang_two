<template>
  <div class="skeleton-loader" :class="{ 'is-loading': loading }">
    <template v-if="loading">
      <div class="skeleton-container" :class="type">
        <!-- 卡片骨架屏 -->
        <template v-if="type === 'card'">
          <div class="skeleton-card" v-for="i in count" :key="i">
            <div class="skeleton-header">
              <div class="skeleton-title"></div>
              <div class="skeleton-tag"></div>
            </div>
            <div class="skeleton-content">
              <div class="skeleton-text"></div>
              <div class="skeleton-text"></div>
            </div>
          </div>
        </template>

        <!-- 轮播图骨架屏 -->
        <template v-else-if="type === 'carousel'">
          <div class="skeleton-carousel">
            <div class="skeleton-image"></div>
            <div class="skeleton-caption">
              <div class="skeleton-title"></div>
              <div class="skeleton-text"></div>
            </div>
          </div>
        </template>

        <!-- 列表骨架屏 -->
        <template v-else-if="type === 'list'">
          <div class="skeleton-list" v-for="i in count" :key="i">
            <div class="skeleton-avatar"></div>
            <div class="skeleton-content">
              <div class="skeleton-title"></div>
              <div class="skeleton-text"></div>
            </div>
          </div>
        </template>
      </div>
    </template>
    <slot v-else></slot>
  </div>
</template>

<script setup lang="ts">
defineProps({
  loading: {
    type: Boolean,
    default: true
  },
  type: {
    type: String,
    default: 'card',
    validator: (value: string) => ['card', 'carousel', 'list'].includes(value)
  },
  count: {
    type: Number,
    default: 1
  }
})
</script>

<style scoped lang="scss">
.skeleton-loader {
  width: 100%;
}

@keyframes shimmer {
  0% {
    background-position: -200% 0;
  }
  100% {
    background-position: 200% 0;
  }
}

.skeleton-container {
  width: 100%;
  
  .skeleton-card,
  .skeleton-carousel,
  .skeleton-list {
    background: #fff;
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 16px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  }
}

// 通用骨架屏样式
.skeleton-title,
.skeleton-text,
.skeleton-tag,
.skeleton-image,
.skeleton-avatar {
  background: linear-gradient(90deg, #f2f2f2 25%, #e6e6e6 50%, #f2f2f2 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
}

// 卡片骨架屏
.skeleton-card {
  .skeleton-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;

    .skeleton-title {
      width: 60%;
      height: 24px;
    }

    .skeleton-tag {
      width: 80px;
      height: 24px;
    }
  }

  .skeleton-content {
    .skeleton-text {
      height: 16px;
      margin-bottom: 8px;

      &:last-child {
        width: 80%;
      }
    }
  }
}

// 轮播图骨架屏
.skeleton-carousel {
  .skeleton-image {
    width: 100%;
    height: 200px;
    margin-bottom: 16px;
  }

  .skeleton-caption {
    .skeleton-title {
      width: 70%;
      height: 28px;
      margin-bottom: 12px;
    }

    .skeleton-text {
      width: 90%;
      height: 16px;
    }
  }
}

// 列表骨架屏
.skeleton-list {
  display: flex;
  align-items: center;
  gap: 16px;

  .skeleton-avatar {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    flex-shrink: 0;
  }

  .skeleton-content {
    flex: 1;

    .skeleton-title {
      width: 50%;
      height: 20px;
      margin-bottom: 8px;
    }

    .skeleton-text {
      width: 80%;
      height: 16px;
    }
  }
}

// 响应式调整
@media (max-width: 768px) {
  .skeleton-card {
    .skeleton-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 8px;

      .skeleton-title {
        width: 80%;
      }
    }
  }

  .skeleton-carousel {
    .skeleton-image {
      height: 150px;
    }
  }
}
</style> 