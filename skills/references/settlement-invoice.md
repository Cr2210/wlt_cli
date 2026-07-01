# 结算发票 (settlement / invoice)

## 运单结算 (`wlt settlement main`)

### list / page-count / export 共用筛选字段

`--page-no` / `--page-size` 外，后端 `/erp/settlement/page` 完整筛选字段：

| Flag | 后端参数 | 说明 | 销售/采购结算 |
|------|----------|------|---------------|
| `--no` | `no` | 结算单号 | 两者 |
| `--name` | `name` | 结算单名称（模糊） | 两者 |
| `--customer-id` | `customerId` | 客户 ID | 销售结算 |
| `--supplier-id` | `supplierId` | 供应商 ID | 采购结算 |
| `--user-id` | `userId` | 业务员 ID（后端实际为 `userId` 字段） | 采购结算 |
| `--enterprise-id` | `enterpriseId` | 企业 ID | 两者 |
| `--project-id` | `projectId` | 项目 ID | 采购结算（如有） |
| `--project-name` | `projectName` | 项目名称（模糊） | 采购结算 |
| `--settle-status` | `settleStatus` | 结算状态（如 `PART_SETTLED`；多选用逗号分隔） | 两者 |
| `--invoice-status` | `invoiceStatus` | 发票状态（如 `PART_INVOICED`；多选用逗号分隔） | 两者 |
| `--type` | `type` | **结算类型：`SALE_SETTLEMENT` / `PURCHASE_SETTLEMENT`** | 两者（区分子类型） |
| `--settle-type` | `settleType` | **结算方式：`SALE` / `PURCHASE`** | 两者（注意 ≠ `type`） |
| `--metrics-name` | `metricsName` | 检测指标名 / 规格指标 | 两者 |
| `--start-date` | `settlementDate[0]` | 结算日期起始（如 `2026-07-01 00:00:00`） | 两者 |
| `--end-date` | `settlementDate[1]` | 结算日期结束（如 `2026-07-31 23:59:59` | 两者 |

> ⚠️ 此前 CLI 的 `--settlement-no` 对应后端字段名实际为 `settlementNo`，与后端期望的 `no` 不一致，导致单号筛选从未生效。现修正为 `--no`。
>
> `--type` 与 `--settle-type` 为两个独立字段：`type` 区分销售/采购结算单（`SALE_SETTLEMENT` / `PURCHASE_SETTLEMENT`），`settle-type` 为结算方式（`SALE` / `PURCHASE`）。两者组合使用。

### 子命令一览

| 命令 | 说明 | 必填参数 |
|------|------|----------|
| `wlt settlement main list` | 分页查询结算单 | — |
| `wlt settlement main page-count` | 统计结算单数量 | — |
| `wlt settlement main get --id <N>` | 获取结算单详情 | `--id` |
| `wlt settlement main create` | 创建结算单 | `--data` |
| `wlt settlement main update` | 更新结算单 | `--data` |
| `wlt settlement main delete` | 删除结算单 | `--id` |
| `wlt settlement main update-status` | 更新结算单状态 | `--data`（含 id 和 status） |
| `wlt settlement main unsettle-waybill` | 查看未结算运单 | `--customer-id` / `--supplier-id` / `--warehouse-id` |
| `wlt settlement main unsettle-waybill-count` | 统计未结算运单 | 同上 |
| `wlt settlement main export` | 导出结算单 Excel | 与 list 同筛选字段 |

### 查询示例

```bash
# 销售结算 — 完全还原: /erp/settlement/page?no=JSD20260403000001&enterpriseId=...&settleStatus=PART_SETTLED&invoiceStatus=PART_INVOICED&settlementDate[0]=...&settlementDate[1]=...&metricsName=规格指标&type=SALE_SETTLEMENT&settleType=SALE
wlt settlement main list \
  --no JSD20260403000001 \
  --name 结算单名称 \
  --enterprise-id 2001494968037298178 \
  --settle-status PART_SETTLED \
  --invoice-status PART_INVOICED \
  --start-date "2026-07-08 00:00:00" \
  --end-date "2026-08-05 23:59:59" \
  --metrics-name 规格指标 \
  --type SALE_SETTLEMENT \
  --settle-type SALE

# 采购结算 — 完全还原: /erp/settlement/page?no=JSD20260622000001&userId=144&projectName=项目名称&...&type=PURCHASE_SETTLEMENT
wlt settlement main list \
  --no JSD20260622000001 \
  --name 结算单名称 \
  --enterprise-id 2003369636046413825 \
  --user-id 144 \
  --project-name 项目名称 \
  --settle-status PART_SETTLED \
  --invoice-status PART_INVOICED \
  --start-date "2026-07-09 00:00:00" \
  --end-date "2026-08-13 23:59:59" \
  --metrics-name 指标 \
  --type PURCHASE_SETTLEMENT \
  --settle-type SALE

# 组合筛选：待部分结算的销售结算单
wlt settlement main list \
  --customer-id 2001494968037298178 \
  --settle-status PART_SETTLED \
  --type SALE_SETTLEMENT \
  --start-date "2026-07-01 00:00:00" \
  --end-date "2026-07-31 23:59:59"

# 导出 Excel（参数同 list）
wlt settlement main export --type SALE_SETTLEMENT --settle-status PART_SETTLED
```

---

## 发票管理 (`wlt invoice main`)

| 命令 | 说明 | 关键参数 |
|------|------|----------|
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
