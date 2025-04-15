<template>
  <div class="home-page">
    <!-- 新品面料轮播图 -->
    <div class="carousel-section">
      <h2>新品面料</h2>
      <p>浏览最新面料系列</p>
      <el-button type="primary" class="view-more">去看看</el-button>
    </div>

    <!-- 热门工厂 -->
    <div class="section">
      <div class="section-header">
        <h2>热门工厂</h2>
        <router-link to="/factories" class="view-more">查看更多</router-link>
      </div>
      <div class="factory-list">
        <div class="factory-card">
          <div class="factory-image">
            <el-image src="/images/factory1.jpg" fit="cover">
              <template #error>
                <div class="image-placeholder">优质服装厂</div>
              </template>
            </el-image>
          </div>
          <div class="factory-info">
            <h3>优质服装厂</h3>
            <p>专业女装生产</p>
          </div>
        </div>
        <div class="factory-card">
          <div class="factory-image">
            <el-image src="/images/factory2.jpg" fit="cover">
              <template #error>
                <div class="image-placeholder">精品制衣厂</div>
              </template>
            </el-image>
          </div>
          <div class="factory-info">
            <h3>精品制衣厂</h3>
            <p>高端定制专家</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 最新面料 -->
    <div class="section">
      <div class="section-header">
        <h2>最新面料</h2>
        <router-link to="/fabrics" class="view-more">查看更多</router-link>
      </div>
      <div class="fabric-grid">
        <div class="fabric-card">
          <div class="fabric-image">
            <el-image src="/images/fabric1.jpg" fit="cover">
              <template #error>
                <div class="image-placeholder">优质棉料</div>
              </template>
            </el-image>
          </div>
          <div class="fabric-info">
            <h3>优质棉料</h3>
            <p>¥58/米</p>
          </div>
        </div>
        <div class="fabric-card">
          <div class="fabric-image">
            <el-image src="/images/fabric2.jpg" fit="cover">
              <template #error>
                <div class="image-placeholder">真丝面料</div>
              </template>
            </el-image>
          </div>
          <div class="fabric-info">
            <h3>真丝面料</h3>
            <p>¥128/米</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
</script>

<style scoped lang="scss">
.home-page {
  padding: 60px 16px 16px;
  background-color: #f5f7fa;
  min-height: 100vh;
}

.carousel-section {
  margin: 16px 0;
  padding: 24px;
  background-color: #e6a23c;
  border-radius: 8px;
  color: white;

  h2 {
    margin: 0;
    font-size: 24px;
    font-weight: 600;
  }

  p {
    margin: 8px 0 16px;
    font-size: 16px;
    opacity: 0.9;
  }

  .view-more {
    background-color: white;
    color: #e6a23c;
    border: none;
    font-weight: 500;

    &:hover {
      background-color: #f5f7fa;
    }
  }
}

.section {
  margin: 16px 0;
  padding: 16px;
  background: white;
  border-radius: 8px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;

  h2 {
    margin: 0;
    font-size: 18px;
    color: #333;
  }

  .view-more {
    color: #409eff;
    text-decoration: none;
    font-size: 14px;

    &:hover {
      color: #66b1ff;
    }
  }
}

.factory-list, .fabric-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 16px;
}

.factory-card, .fabric-card {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  transition: transform 0.3s ease;

  &:hover {
    transform: translateY(-4px);
  }
}

.factory-image, .fabric-image {
  height: 120px;
  overflow: hidden;

  .el-image {
    width: 100%;
    height: 100%;
  }

  .image-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #f5f7fa;
    color: #909399;
    font-size: 14px;
  }
}

.factory-info, .fabric-info {
  padding: 12px;

  h3 {
    margin: 0;
    font-size: 14px;
    color: #333;
    font-weight: 500;
  }

  p {
    margin: 4px 0 0;
    font-size: 12px;
    color: #909399;
  }
}

@media (min-width: 768px) {
  .home-page {
    padding: 80px 24px 24px;
  }

  .factory-list, .fabric-grid {
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  }
}
</style> 