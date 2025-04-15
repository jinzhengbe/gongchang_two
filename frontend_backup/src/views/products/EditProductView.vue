<template>
  <div class="edit-product">
    <h2>编辑产品</h2>
    <el-form label-width="120px">
      <el-form-item label="产品名称">
        <el-input v-model="form.name" />
      </el-form-item>
      <el-form-item label="类别">
        <el-input v-model="form.category" />
      </el-form-item>
      <el-form-item label="价格">
        <el-input-number v-model="form.price" :min="0" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary">保存</el-button>
        <el-button>取消</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'

export default defineComponent({
  name: 'EditProductView',
  setup() {
    const form = ref({
      name: '',
      category: '',
      price: 0
    })

    return { form }
  }
})
</script>

<style scoped>
.edit-product {
  padding: 20px;
  max-width: 600px;
  margin: 0 auto;
}
</style> 