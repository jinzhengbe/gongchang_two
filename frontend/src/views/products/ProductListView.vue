<template>
  <div class="product-list">
    <h2>产品列表</h2>
    <el-table :data="[]" style="width: 100%">
      <el-table-column prop="name" label="产品名称" />
      <el-table-column prop="category" label="类别" />
      <el-table-column prop="price" label="价格" />
      <el-table-column label="操作">
        <template #default="scope">
          <el-button link type="primary" size="small">查看</el-button>
          <el-button link type="primary" size="small">编辑</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'ProductListView'
})
</script>

<style scoped>
.product-list {
  padding: 20px;
}
</style> 