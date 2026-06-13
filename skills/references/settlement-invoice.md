# 结算发票 (settlement / invoice)

## 运单结算 (`wlt settlement main`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt settlement main list` | 分页查询结算单 | `--settlement-no`, `--customer-id`, `--supplier-id`, `--status`, `--warehouse-id`, `--page-no`, `--page-size` |
| `wlt settlement main page-count` | 统计结算数量 | `--settlement-no`, `--customer-id`, `--supplier-id`, `--status`, `--warehouse-id` |
| `wlt settlement main get --id <N>` | 获取结算详情 | `--id`（必填） |
| `wlt settlement main create --data '<json>'` | 创建结算单 | `--data`（必填） |
| `wlt settlement main update --data '<json>'` | 更新结算单 | `--data`（必填） |
| `wlt settlement main delete --id <N>` | 删除结算单 | `--id`（必填） |
| `wlt settlement main update-status --data '<json>'` | 更新结算状态 | `--data`（必填） |
| `wlt settlement main unsettle-waybill` | 未结算运单 | `--customer-id`, `--supplier-id`, `--warehouse-id` |
| `wlt settlement main unsettle-waybill-count` | 未结算运单数量 | `--customer-id`, `--supplier-id`, `--warehouse-id` |
| `wlt settlement main export` | 导出 Excel | `--settlement-no`, `--customer-id`, `--supplier-id`, `--status`, `--warehouse-id` |

## 发票管理 (`wlt invoice main`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt invoice main list` | 分页查询发票 | `--invoice-no`, `--customer-id`, `--supplier-id`, `--status`, `--type`, `--page-no`, `--page-size` |
| `wlt invoice main page-count` | 统计发票数量 | `--invoice-no`, `--customer-id`, `--supplier-id`, `--status`, `--type` |
| `wlt invoice main get --id <N>` | 获取发票详情 | `--id`（必填） |
| `wlt invoice main create --data '<json>'` | 创建发票 | `--data`（必填） |
| `wlt invoice main update --data '<json>'` | 更新发票 | `--data`（必填） |
| `wlt invoice main delete --id <N>` | 删除发票 | `--id`（必填） |
| `wlt invoice main update-status --data '<json>'` | 更新发票状态 | `--data`（必填） |
| `wlt invoice main export` | 导出 Excel | `--invoice-no`, `--customer-id`, `--supplier-id`, `--status`, `--type` |

## 关键区分

- **`settlement`（运单结算）**：基于运单的结算管理，处理运输费用结算
- **`finance settlement`（财务结算单据）**：财务模块的结算单据管理
- **`invoice`（业务发票）**：独立的发票管理模块
- **`partner invoice`（客户/供应商发票抬头）**：客户/供应商下的发票抬头信息
