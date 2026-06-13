# 统计报表 (stats / report / homepage / screen)

## 数据统计 (`wlt stats`)

所有统计命令共享通用参数：`--type`, `--start-time`, `--end-time`, `--sort-by`

### 总览 (`wlt stats overview`)

| 命令 | 说明 | 附加参数 |
|------|------|---------|
| `wlt stats overview` | 总览仪表盘 | `--type`, `--start-time`, `--end-time`, `--sort-by`, `--product-id` |

### 库存统计 (`wlt stats stock`)

| 命令 | 说明 | 附加参数 |
|------|------|---------|
| `wlt stats stock` | 库存统计 | `--type`, `--start-time`, `--end-time`, `--sort-by`, `--product-id`, `--warehouse-id` |

### 财务统计 (`wlt stats finance`)

| 命令 | 说明 | 附加参数 |
|------|------|---------|
| `wlt stats finance data-overview` | 财务数据总览 | `--type`, `--start-time`, `--end-time`, `--sort-by` |
| `wlt stats finance receivable-rankings` | 应收排名 | 同上 |
| `wlt stats finance overdue-receivable-rankings` | 逾期应收排名 | 同上 |
| `wlt stats finance payable-rankings` | 应付排名 | 同上 |
| `wlt stats finance overdue-payable-rankings` | 逾期应付排名 | 同上 |

### 销售统计 (`wlt stats sale`)

| 命令 | 说明 | 附加参数 |
|------|------|---------|
| `wlt stats sale data-overview` | 销售数据总览 | `--type`, `--start-time`, `--end-time`, `--sort-by` |
| `wlt stats sale customer-rankings` | 客户排名 | 同上 |
| `wlt stats sale product-rankings` | 产品排名 | 同上 |
| `wlt stats sale employee-rankings` | 员工排名 | 同上 |
| `wlt stats sale region-rankings` | 区域排名 | 同上 |

### 采购统计 (`wlt stats purchase`)

| 命令 | 说明 | 附加参数 |
|------|------|---------|
| `wlt stats purchase data-overview` | 采购数据总览 | `--type`, `--start-time`, `--end-time`, `--sort-by` |
| `wlt stats purchase supplier-rankings` | 供应商排名 | 同上 |
| `wlt stats purchase product-rankings` | 产品排名 | 同上 |
| `wlt stats purchase employee-rankings` | 员工排名 | 同上 |
| `wlt stats purchase region-rankings` | 区域排名 | 同上 |

### 生产统计 (`wlt stats produce`)

| 命令 | 说明 | 附加参数 |
|------|------|---------|
| `wlt stats produce data-overview` | 生产数据总览 | `--type`, `--start-time`, `--end-time`, `--sort-by` |

## 报表 (`wlt report`)

### 库存报表 (`wlt report stock`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt report stock detail` | 库存明细报表 | `--product-id`, `--category-id`, `--warehouse-id`, `--start-time`, `--end-time`, `--page-no`, `--page-size` |
| `wlt report stock warehouse` | 仓库库存报表 | `--product-id`, `--category-id`, `--warehouse-id`, `--start-time`, `--end-time` |
| `wlt report stock buy-send` | 采发报表 | 同上 |
| `wlt report stock finance` | 财务库存报表 | 同上 |
| `wlt report stock produce` | 生产库存报表 | 同上 |

### 采购报表 (`wlt report purchase`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt report purchase detail` | 采购明细报表 | `--supplier-id`, `--product-id`, `--category-id`, `--warehouse-id`, `--status`, `--start-time`, `--end-time`, `--page-no`, `--page-size` |
| `wlt report purchase summer` | 采购汇总报表 | 同上 |

### 销售报表 (`wlt report sale`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt report sale detail` | 销售明细报表 | `--customer-id`, `--product-id`, `--category-id`, `--warehouse-id`, `--status`, `--start-time`, `--end-time`, `--page-no`, `--page-size` |
| `wlt report sale summer` | 销售汇总报表 | 同上 |
| `wlt report sale profit` | 利润报表 | 同上 |

## 首页仪表盘 (`wlt homepage`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt homepage dashboard1` | 仪表盘1（无需参数） | 无 |
| `wlt homepage dashboard2` | 仪表盘2 | `--type`, `--start-time`, `--end-time`, `--sort-by` |
| `wlt homepage dashboard6` | 仪表盘6 | `--type`, `--start-time`, `--end-time`, `--sort-by`, `--product-id`, `--warehouse-id` |
| `wlt homepage inventory-backlog` | 库存积压 | 同 dashboard6 |
| `wlt homepage product-ranking` | 产品排行 | 同 dashboard6 |

## 大屏数据 (`wlt screen`)

| 命令 | 说明 | 参数 |
|------|------|------|
| `wlt screen purchase-sale` | 采销大屏 | 无 |
| `wlt screen stock-count` | 库存大屏 | 无 |
| `wlt screen amount-used` | 金额大屏 | 无 |
| `wlt screen project-count` | 项目大屏 | 无 |

## 其他操作命令

### 操作日志 (`wlt operate-log`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt operate-log list` | 操作日志列表 | `--module`, `--type`, `--user-name`, `--page-no`, `--page-size` |

### 数据同步 (`wlt data-sync`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt data-sync list` | 同步消息列表 | `--status`, `--type`, `--page-no`, `--page-size` |
| `wlt data-sync get --id <N>` | 获取同步详情 | `--id`（必填） |
| `wlt data-sync resend --id <N>` | 重新发送 | `--id`（必填） |

### 定时任务 (`wlt job-trigger`)

| 命令 | 说明 | 参数 |
|------|------|------|
| `wlt job-trigger execute-product-cost` | 执行产品成本计算 | 无 |
| `wlt job-trigger execute-receivable-balance` | 执行应收余额计算 | 无 |

### 利润事件 (`wlt profit-event`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt profit-event list` | 利润事件列表 | `--event-type`, `--status`, `--page-no`, `--page-size` |
| `wlt profit-event statistics` | 利润统计 | 无 |
| `wlt profit-event types` | 事件类型 | 无 |
| `wlt profit-event retry --event-id <N>` | 重试事件 | `--event-id`（必填） |
| `wlt profit-event clean-expired` | 清理过期事件 | 无 |
| `wlt profit-event health` | 健康检查 | 无 |

### 利润计算 (`wlt profit-calculation`)

| 命令 | 说明 | 参数 |
|------|------|------|
| `wlt profit-calculation batch-recalculate-all` | 重新计算全部利润 | 无 |
