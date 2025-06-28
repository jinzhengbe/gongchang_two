# 修复移除布料 API 类型错误

## 问题描述

前端调用 `DELETE /api/orders/{id}/remove-fabric` API 时出现 400 错误：

```
json: cannot unmarshal number into Go struct field RemoveFabricFromOrderRequest.fabricId of type string
```

前端发送的数据：
```json
{
  "fabricId": 44
}
```

但后端期望的是字符串类型。

## 根本原因

在 `backend/models/order_fabric.go` 中，`RemoveFabricFromOrderRequest` 结构体的 `FabricID` 字段被定义为 `string` 类型，但前端发送的是数字类型。

## 修复方案

### 1. 修改请求结构体类型

**文件：** `backend/models/order_fabric.go`

将 `RemoveFabricFromOrderRequest` 结构体中的 `FabricID` 字段从 `string` 改为 `uint`：

```go
// 修复前
type RemoveFabricFromOrderRequest struct {
    FabricID string `json:"fabricId" binding:"required"`
}

// 修复后
type RemoveFabricFromOrderRequest struct {
    FabricID uint `json:"fabricId" binding:"required"`
}
```

### 2. 更新服务层逻辑

**文件：** `backend/services/order.go`

修改 `RemoveFabricFromOrder` 方法，适配 `uint` 类型的 `FabricID`：

- 简化数据库查询逻辑
- 使用 `FabricIDList` 工具类来处理布料ID列表
- 移除不必要的 `OrderFabric` 关联表操作

### 3. 修复控制器验证逻辑

**文件：** `backend/controllers/order.go`

将验证逻辑从字符串空值检查改为 `uint` 零值检查：

```go
// 修复前
if req.FabricID == "" {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "布料ID不能为空"})
    return
}

// 修复后
if req.FabricID == 0 {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "布料ID不能为空"})
    return
}
```

### 4. 更新测试文件

更新了以下测试文件中的 `fabricId` 格式：

- `test_remove_fabric_from_order.sh`
- `test_api_routes.sh`
- `test_all_delete_apis.sh`
- `test_remove_fabric_fixed.sh` (新建)

## 修复后的效果

现在前端可以正常发送数字类型的 `fabricId`：

```json
{
  "fabricId": 44
}
```

API 将正确处理请求并返回成功响应。

## 测试

可以使用以下命令测试修复后的 API：

```bash
curl -X DELETE "https://aneworders.com/api/orders/29/remove-fabric" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"fabricId": 44}'
```

## 注意事项

- 确保前端发送的 `fabricId` 是有效的数字
- 布料ID 0 被视为无效值
- 修复后的逻辑更加简洁，直接使用逗号分隔的字符串来管理订单的布料列表 