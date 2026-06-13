# 质检称重 (quality / weight)

## 称重质检 (`wlt quality weight`)

只读报告命令，按订单或运单维度查看称重数据。

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt quality weight order-page` | 订单称重分页 | `--order-id`, `--product-id`, `--page-no`, `--page-size` |
| `wlt quality weight order-summary` | 订单称重汇总 | `--order-id`, `--product-id` |
| `wlt quality weight order-export` | 订单称重导出 | `--order-id`, `--product-id` |
| `wlt quality weight waybill-page` | 运单称重分页 | `--waybill-id`, `--product-id`, `--page-no`, `--page-size` |
| `wlt quality weight waybill-summary` | 运单称重汇总 | `--waybill-id`, `--product-id` |
| `wlt quality weight waybill-export` | 运单称重导出 | `--waybill-id`, `--product-id` |

## 质检单 (`wlt quality inspection`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt quality inspection list` | 分页查询质检单 | `--inspection-no`, `--status`, `--warehouse-id`, `--product-id`, `--business-id`, `--business-type`, `--page-no`, `--page-size` |
| `wlt quality inspection get --id <N>` | 获取质检单详情 | `--id`（必填） |
| `wlt quality inspection create --data '<json>'` | 创建质检单 | `--data`（必填） |
| `wlt quality inspection update --data '<json>'` | 更新质检单 | `--data`（必填） |
| `wlt quality inspection delete --id <N>` | 删除质检单 | `--id`（必填） |
| `wlt quality inspection update-status --data '<json>'` | 更新质检单状态 | `--data`（必填） |
| `wlt quality inspection summary` | 质检汇总 | 无 |
| `wlt quality inspection relate-list` | 关联质检列表 | `--business-id`, `--business-type` |
| `wlt quality inspection order-waybill-inspection` | 订单运单质检 | `--business-id`, `--business-type` |
| `wlt quality inspection export` | 导出 Excel | `--inspection-no`, `--status`, `--warehouse-id`, `--product-id`, `--business-id`, `--business-type` |
| `wlt quality inspection refresh-all-summary` | 刷新汇总 | 无 |
| `wlt quality inspection import --file <path>` | 导入 Excel | `--file`（必填） |

## 称重管理 (`wlt weight waybill`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt weight waybill list` | 分页查询称重 | `--waybill-no`, `--license-plate`, `--weighing-no`, `--page-no`, `--page-size` |
| `wlt weight waybill image --id <N>` | 获取称重图片 | `--id`（必填） |
| `wlt weight waybill create-waybill-source --data '<json>'` | 创建关联运单 | `--data`（必填） |
| `wlt weight waybill batch-create-waybill-source --data '<json>'` | 批量创建关联运单 | `--data`（必填） |
| `wlt weight waybill match-link-waybill-source --data '<json>'` | 匹配关联运单 | `--data`（必填） |
| `wlt weight waybill unlink-waybill-source --data '<json>'` | 取消关联运单 | `--data`（必填） |
| `wlt weight waybill export` | 导出 Excel | `--waybill-no`, `--license-plate`, `--weighing-no` |
