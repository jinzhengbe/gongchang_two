<template>
  <div class="edit-order">
    <h2>编辑订单</h2>
    <el-form label-width="120px">
      <el-form-item label="订单编号">
        <el-input v-model="form.orderNumber" disabled />
      </el-form-item>
      <el-form-item label="产品名称">
        <el-input v-model="form.productName" />
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="form.status">
          <el-option label="待处理" value="待处理" />
          <el-option label="进行中" value="进行中" />
          <el-option label="已完成" value="已完成" />
          <el-option label="已取消" value="已取消" />
        </el-select>
      </el-form-item>
      <el-form-item label="备注">
        <el-input type="textarea" v-model="form.notes" />
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
  name: 'EditOrderView',
  setup() {
    const form = ref({
      orderNumber: '',
      productName: '',
      status: '待处理',
      notes: ''
    })

    return { form }
  }
})
</script>

<style scoped>
.edit-order {
  padding: 20px;
  max-width: 600px;
  margin: 0 auto;
}
</style> 