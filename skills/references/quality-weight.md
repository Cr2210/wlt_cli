# 质检称重 (quality / weight)

## 质检称重 (`wlt quality weight`)

后端 `/erp/quality-inspection-weight/{summary-order,summary-waybill,page-order,page-waybill,export-excel-order,export-excel-waybill}` 共用一套 **13** 个筛选字段。

### 共用筛选字段（6 个查询命令共享）

| Flag | 后端参数 | 说明 | 订单端 | 运单端 |
|------|----------|------|:---:|:---:|
| `--type` | `type` | 业务类型 `SALE/PURCHASE` | ✅ (必填) | — |
| `--order-type` | `orderType` | 订单类型 `SALE_OUT/PURCHASE_IN` 等 | ✅ | ✅ |
| `--order-no` | `orderNo` | 订单号 | ✅ | ✅ |
| `--enterprise-id` | `enterpriseId` | 企业 ID | ✅ | ✅ |
| `--product-id` | `productId` | 产品 ID | ✅ | ✅ |
| `--product-name` | `productName` | 产品名称（模糊） | ✅ | ✅ |
| `--metrics-name` | `metricsName` | 规格指标 | ✅ | ✅ |
| `--send-inspection-no` | `sendInspectionNo` | 发货质检单号 | ✅ | ✅ |
| `--send-metrics-name` | `sendMetricsName` | 发货质检指标名 | ✅ | ✅ |
| `--receive-inspection-no` | `receiveInspectionNo` | 收货质检单号 | ✅ | ✅ |
| `--receive-metrics-name` | `receiveMetricsName` | 收货质检指标名 | ✅ | ✅ |
| `--send-address` | `sendAddress` | 发货地 | ✅ | ✅ |
| `--start-time` | `orderTime[0]` | 订单时间起始 `yyyy-MM-dd HH:mm:ss` | ✅ | ✅ |
| `--end-time` | `orderTime[1]` | 订单时间结束 | ✅ | ✅ |

> `--start-time/--end-time` 为统一入口，内部自动折叠为后端数组参数 `orderTime[0]`/`orderTime[1]`。
>
> `--type` 仅订单端（`order-page/order-summary/order-export`）**必填**；运单端（`waybill-*`）无需 `--type`，后端按 waybill 路由处理。

### 命令一览

| 命令 | 说明 | 必填参数 |
|------|------|----------|
| `wlt quality weight order-page` | 订单称重分页 | `--type` |
| `wlt quality weight order-summary` | 订单称重汇总 | `--type` |
| `wlt quality weight order-export` | 订单称重导出 | `--type` |
| `wlt quality weight waybill-page` | 运单称重分页 | — |
| `wlt quality weight waybill-summary` | 运单称重汇总 | — |
| `wlt quality weight waybill-export` | 运单称重导出 | — |

### 查询示例

```bash
# 完全还原: /erp/quality-inspection-weight/summary-order?type=PURCHASE&orderType=PURCHASE_IN&orderNo=CGDD20260402000002&...&orderTime[0]=...&orderTime[1]=...&sendAddress=发货地
wlt quality weight order-summary \
  --type PURCHASE \
  --order-type PURCHASE_IN \
  --order-no CGDD20260402000002 \
  --enterprise-id 2001552070697033730 \
  --product-id 1927612799140261889 \
  --product-name 货物名称 \
  --metrics-name 规格指标 \
  --send-inspection-no 发货质检单 \
  --send-metrics-name 检验结果加权 \
  --receive-inspection-no 收货质检单 \
  --receive-metrics-name 收货质检加权 \
  --send-address 发货地 \
  --start-time "2026-07-09 00:00:00" \
  --end-time "2026-08-13 23:59:59"

# 运单端: /erp/quality-inspection-weight/page-waybill?type=PURCHASE&orderType=PURCHASE_IN&...&pageNo=1&pageSize=10
wlt quality weight waybill-page \
  --order-type PURCHASE_IN \
  --order-no CGDD20260402000002 \
  --enterprise-id 2001552070697033730 \
  --product-id 1927612799140261889 \
  --start-time "2026-07-09 00:00:00" \
  --end-time "2026-08-13 23:59:59"
```

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
