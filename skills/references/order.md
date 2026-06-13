# 订单管理 (order)

## 主订单 (`wlt order main`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt order main list` | 分页查询订单 | `--order-no`, `--customer-id`, `--status`, `--warehouse-id`, `--page-no`, `--page-size` |
| `wlt order main page-count` | 统计订单数量 | `--order-no`, `--customer-id`, `--status`, `--warehouse-id` |
| `wlt order main get --id <N>` | 获取订单详情 | `--id`（必填） |
| `wlt order main get-linkorder-by-orderId` | 获取关联订单 | `--order-id`（必填） |
| `wlt order main create --data '<json>'` | 创建订单 | `--data`（必填） |
| `wlt order main update --data '<json>'` | 更新订单 | `--data`（必填） |
| `wlt order main delete --id <N>` | 删除订单 | `--id`（必填） |
| `wlt order main update-status --data '<json>'` | 更新订单状态 | `--data`（必填） |
| `wlt order main cancel --data '<json>'` | 取消订单 | `--data`（必填） |
| `wlt order main reopen --data '<json>'` | 重新打开订单 | `--data`（必填） |
| `wlt order main complete --data '<json>'` | 完成订单 | `--data`（必填） |
| `wlt order main link-waybill --data '<json>'` | 关联运单 | `--data`（必填） |
| `wlt order main unlink-waybill --data '<json>'` | 取消关联运单 | `--data`（必填） |
| `wlt order main export` | 导出 Excel | `--order-no`, `--customer-id`, `--status`, `--warehouse-id` |

## 排产计划 (`wlt order plan`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt order plan list` | 分页查询计划 | `--plan-no`, `--customer-id`, `--status`, `--warehouse-id`, `--page-no`, `--page-size` |
| `wlt order plan page-count` | 统计计划数量 | `--plan-no`, `--customer-id`, `--status`, `--warehouse-id` |
| `wlt order plan get --id <N>` | 获取计划详情 | `--id`（必填） |
| `wlt order plan create --data '<json>'` | 创建计划 | `--data`（必填） |
| `wlt order plan update --data '<json>'` | 更新计划 | `--data`（必填） |
| `wlt order plan delete --id <N>` | 删除计划 | `--id`（必填） |
| `wlt order plan update-status --data '<json>'` | 更新计划状态 | `--data`（必填） |
| `wlt order plan cancel --data '<json>'` | 取消计划 | `--data`（必填） |
| `wlt order plan reopen --data '<json>'` | 重新打开计划 | `--data`（必填） |
| `wlt order plan complete --data '<json>'` | 完成计划 | `--data`（必填） |
| `wlt order plan export` | 导出 Excel | `--plan-no`, `--customer-id`, `--status`, `--warehouse-id` |

## 常见工作流

### 创建订单并关联运单

```bash
wlt order main create --data '{"orderNo":"ORD-001","customerId":1,...}'
wlt waybill source create --data '{"waybillNo":"WB-001",...}'
wlt order main link-waybill --data '{"orderId":1,"waybillId":1}'
```

### 订单生命周期

```bash
wlt order main list --status 0               # 查询待处理
wlt order main update-status --data '...'    # 更新状态
wlt order main complete --data '...'         # 完成订单
```
