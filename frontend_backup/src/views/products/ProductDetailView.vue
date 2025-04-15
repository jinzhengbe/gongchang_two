<template>
  <div class="product-detail">
    <h2>产品详情</h2>
    <el-descriptions :column="1" border>
      <el-descriptions-item label="产品名称">{{ product.name }}</el-descriptions-item>
      <el-descriptions-item label="类别">{{ product.category }}</el-descriptions-item>
      <el-descriptions-item label="价格">{{ product.price }}</el-descriptions-item>
      <el-descriptions-item label="创建时间">{{ product.createdAt }}</el-descriptions-item>
    </el-descriptions>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'

export default defineComponent({
  name: 'ProductDetailView',
  setup() {
    const product = ref({
      name: '',
      category: '',
      price: 0,
      createdAt: ''
    })

    return { product }
  }
})
</script>

<style scoped>
.product-detail {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}
</style> 