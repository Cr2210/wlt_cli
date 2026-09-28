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
| `wlt settlement main ignore-remainder` | 忽略剩余结算/开票金额 | `--id`, `--type`, `--reason` |
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

### 忽略剩余（ignore-remainder）

```bash
# 忽略剩余结算金额（剩余归零，状态翻转为已结算）
wlt settlement main ignore-remainder --id <结算单ID> --type SETTLE --reason "尾差忽略"

# 忽略剩余开票金额（要求结算单 invoiceFlag=无需开票，且仅已审核单可操作）
wlt settlement main ignore-remainder --id <结算单ID> --type INVOICE --reason "无需开票"
```

> ⚠️ `--type` 取值：`SETTLE`（剩余结算）/ `INVOICE`（剩余开票）；`--reason` 必填。忽略开票剩余要求结算单 `invoiceFlag=无需开票` 且已审核，否则后端拒绝。
>
> `create` / `update` 的 `--data` 可携带 `"invoiceFlag": 1`（无需开票标识）；`1→0` 的回退有守卫，已置为无需开票后不能直接改回。

---

## 物流结算 (`wlt settlement logistics`)

**API 路径**：`/erp/logistics-settlement`（承运商维度的物流结算单）

### list / page-count / export 共用筛选字段

| Flag | 后端参数 | 说明 |
|------|----------|------|
| `--no` | `no` | 结算单号 |
| `--name` | `name` | 结算单名称 |
| `--carrier-name` | `carrierName` | 承运商名称 |
| `--carrier-enterprise-id` | `carrierEnterpriseId` | 承运商企业 ID |
| `--status` | `status` | 状态 |
| `--start-date` | `settlementDate[0]` | 结算日期起始（如 `2026-07-01 00:00:00`） |
| `--end-date` | `settlementDate[1]` | 结算日期结束（如 `2026-07-31 23:59:59`） |

### 子命令一览

| 命令 | 说明 | 必填参数 |
|------|------|----------|
| `wlt settlement logistics list` | 分页查询物流结算单 | — |
| `wlt settlement logistics page-count` | 按筛选统计金额合计（Σ含税金额/Σ结算金额/Σ开票金额/Σ忽略额） | — |
| `wlt settlement logistics get --id <N>` | 获取物流结算单详情 | `--id` |
| `wlt settlement logistics create` | 创建物流结算单 | `--data` |
| `wlt settlement logistics update` | 更新物流结算单 | `--data` |
| `wlt settlement logistics delete` | 删除物流结算单 | `--ids`（逗号分隔） |
| `wlt settlement logistics update-status` | 更新物流结算单状态（后端 @RequestParam） | `--id`, `--status` |
| `wlt settlement logistics mismatch-list` | 分页查询金额差异运单（运单金额与结算金额不一致） | —（可选 `--keyword`） |
| `wlt settlement logistics available-waybill` | 分页查询可选运单 | —（可选 `--carrier-id`, `--carrier-name`, `--keyword`, `--settlement-id`） |
| `wlt settlement logistics available-waybill-count` | 统计可选运单合计（金额/装货重量/卸货重量） | 同上 |
| `wlt settlement logistics recalculate-all` | 全量重算所有物流订单聚合结算明细（数据对齐，低频运维操作） | — |
| `wlt settlement logistics ignore-remainder` | 忽略剩余结算/开票金额 | `--id`, `--type`, `--reason` |
| `wlt settlement logistics export` | 导出物流结算单 Excel | 与 list 同筛选字段 + `--headers` 自定义导出表头 |

### 查询示例

```bash
# 按承运商 + 结算日期范围筛选
wlt settlement logistics list \
  --carrier-name 某物流公司 \
  --start-date "2026-07-01 00:00:00" --end-date "2026-07-31 23:59:59"

# 差异运单（结算金额 vs 运单金额不一致）
wlt settlement logistics mismatch-list --keyword 车牌号

# 可选运单（创建结算单时选运单 / 详情回显 / 结算单带出明细）
wlt settlement logistics available-waybill --carrier-id <承运商ID> --settlement-id <结算单ID>

# 忽略剩余（同 settlement main 的规则：SETTLE/INVOICE，INVOICE 要求 invoiceFlag=无需开票）
wlt settlement logistics ignore-remainder --id <物流结算单ID> --type SETTLE --reason "尾差忽略"

# 导出 Excel（参数同 list）
wlt settlement logistics export --carrier-name 某物流公司
```

> ⚠️ `recalculate-all` 为全量重算运维操作，执行前必须获得用户确认。

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

- **`settlement main`（采购/销售结算单）**：客户/供应商维度的结算单管理
- **`settlement logistics`（物流结算单）**：承运商维度的物流结算，处理运输费用结算
- **`finance settlement`（财务结算单据）**：财务模块的结算单据管理
- **`invoice`（业务发票）**：独立的发票管理模块
- **`partner invoice`（客户/供应商发票抬头）**：客户/供应商下的发票抬头信息
