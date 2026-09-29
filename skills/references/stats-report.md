# 统计报表 (stats / report / homepage / screen)

## 数据统计 (`wlt stats`)

所有统计命令共享通用参数：`--type`（month/year，默认 month）, `--start-time`（默认当月1号）, `--sort-by`（默认 amount）

### 总览 (`wlt stats overview`)

| 命令 | 说明 | 附加参数 |
|------|------|---------|
| `wlt stats overview` | 经营总览 | `--type`, `--start-time`, `--sort-by`, `--product-id` |

> **审计备注(2026-07)**:
> - URL 样本 `?productId=...&type=month&startTime=2026-07-01 00:00:00` / `?productId=...&type=year&startTime=2026-01-01 00:00:00` 中,后端实际接收的参数为 `productId` / `type` / `startTime`。
> - CLI 默认 `type=month`、`startTime=当前月第一天 00:00:00`(走 `DefaultMonthRange()`)、`sortBy=amount`,与 URL 样本对齐。
> - 原 `wlt homepage dashboard2` 已删除(避免重复注册;此命令在 `wlt stats overview` 完整覆盖,且后者多一个 `--product-id` flag)。

### 库存统计 (`wlt stats stock`)

| 命令 | 说明 | 附加参数 |
|------|------|---------|
| `wlt stats stock` | 库存统计 | `--type`, `--start-time`, `--sort-by`, `--product-id`, `--warehouse-id` |

> **审计备注(2026-07)**:
> - URL 样本 `?productId=...&type=month&startTime=2026-07-01 00:00:00&sortBy=amount` 确认库存统计 1 个命令的现有 flag 集合(`product-id` + 可选 `warehouse-id`)与后端对齐。**代码无需改动**。
> - `--sort-by` 本次样本为 `amount`(默认)。其他值未确认。
> - 后端路径: `/erp/homepage/dashboard6`(无独立 `/erp/homepage/stock` 端点;`wlt homepage dashboard6` 已删除,避免与 `wlt stats stock` 重复注册)。

### 财务统计 (`wlt stats finance`)

| 命令 | 说明 | 附加参数 |
|------|------|---------|
| `wlt stats finance data-overview` | 财务数据总览 | `--type`, `--start-time`, `--sort-by`, `--product-id` |
| `wlt stats finance receivable-rankings` | 应收排名 | 同上 |
| `wlt stats finance overdue-receivable-rankings` | 逾期应收排名 | 同上 |
| `wlt stats finance payable-rankings` | 应付排名 | 同上 |
| `wlt stats finance overdue-payable-rankings` | 逾期应付排名 | 同上 |

> **审计备注(2026-07)**:
> - URL 样本 `?sortBy=enterprise&productId=...&type=month&startTime=2026-07-01 00:00:00` 确认财务统计 5 个命令均需 `--product-id` flag。
> - `--sort-by` 已知合法值: `amount`(默认)、`enterprise`。其他值未确认。
> - 按"stats 下接口请求参数一致"原则,5 个命令共享同一组 flag(`financeFlags = [{product-id}]`)。

### 销售统计 (`wlt stats sale`)

| 命令 | 说明 | 附加参数 |
|------|------|---------|
| `wlt stats sale data-overview` | 销售数据总览 | `--type`, `--start-time`, `--sort-by`, `--product-id` |
| `wlt stats sale customer-rankings` | 客户排名 | 同上 |
| `wlt stats sale product-rankings` | 产品排名 | 同上 |
| `wlt stats sale employee-rankings` | 员工排名 | 同上 |
| `wlt stats sale region-rankings` | 区域排名 | 同上 |

> **审计备注(2026-07)**:
> - URL 样本 `?productId=...&type=month&startTime=2026-07-01 00:00:00&sortBy=amount` 确认销售统计 5 个命令均需 `--product-id` flag。
> - `--sort-by` 本次样本为 `amount`(默认)。已知其他值:`enterprise`(见 finance 备注)。其他值未确认。
> - 按"stats 下接口请求参数一致"原则,5 个命令共享同一组 flag(`saleFlags = [{product-id}]`)。

### 采购统计 (`wlt stats purchase`)

| 命令 | 说明 | 附加参数 |
|------|------|---------|
| `wlt stats purchase data-overview` | 采购数据总览 | `--type`, `--start-time`, `--sort-by`, `--product-id` |
| `wlt stats purchase supplier-rankings` | 供应商排名 | 同上 |
| `wlt stats purchase product-rankings` | 产品排名 | 同上 |
| `wlt stats purchase employee-rankings` | 员工排名 | 同上 |
| `wlt stats purchase region-rankings` | 区域排名 | 同上 |

> **审计备注(2026-07)**:
> - URL 样本 `?productId=...&type=month&startTime=2026-07-01 00:00:00&sortBy=amount` 确认采购统计 5 个命令均需 `--product-id` flag。
> - `--sort-by` 本次样本为 `amount`(默认)。其他值未确认。
> - 按"stats 下接口请求参数一致"原则,5 个命令共享同一组 flag(`purchaseFlags = [{product-id}]`)。

### 生产统计 (`wlt stats produce`)

| 命令 | 说明 | 附加参数 |
|------|------|---------|
| `wlt stats produce data-overview` | 生产数据总览 | `--type`, `--start-time`, `--sort-by`, `--product-id` |

> **审计备注(2026-07)**:
> - URL 样本 `?productId=...&type=month&startTime=2026-07-01 00:00:00` 确认生产统计 1 个命令需 `--product-id` flag。
> - **特殊**:URL 样本**未带 `sortBy`**,而 finance/sale/purchase 的样本都带 `sortBy=amount` / `sortBy=enterprise`。**CLI 默认仍会发 `sortBy=amount`**(走 `CollectStatsFlags` 通用逻辑),后端可能忽略或报错,待实测确认。
> - 1 个命令共享 `produceFlags = [{product-id}]`。

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
| `wlt homepage inventory-backlog` | 库存积压 | `--type`, `--start-time`, `--sort-by`, `--product-id`, `--warehouse-id` |
| `wlt homepage product-ranking` | 产品排行 | 同上 |

> **审计备注(2026-07)**:
> - `wlt homepage dashboard2`(业务概览)、`wlt homepage dashboard6`(库存分析)分别已迁至 `wlt stats overview` / `wlt stats stock`,避免重复注册。`wlt homepage` 现仅保留 dashboard1 / inventory-backlog / product-ranking 三个原生 homepage 端点。

## 大屏数据 (`wlt screen`)

| 命令 | 说明 | 参数 |
|------|------|------|
| `wlt screen purchase-sale` | 采销大屏 | 无 |
| `wlt screen stock-count` | 库存大屏 | 无 |
| `wlt screen amount-used` | 金额大屏 | 无 |
| `wlt screen project-count` | 项目大屏 | 无 |

> operate-log / data-sync / profit-event / profit-calculation / job-trigger 等辅助模块命令见 [SKILL.md「辅助模块」](../SKILL.md)。
