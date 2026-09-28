# 销售采购 (sale / purchase)

销售和采购模块均使用标准 CRUD 子命令模式（list / page-count / get / create / update / delete / update-status）。

## 共用查询参数（list / page-count）

除通用分页参数 `--page-no` / `--page-size` 外，两个模块共享以下筛选字段：

| Flag | 后端参数 | 说明 |
|------|----------|------|
| `--warehouse-id` | `warehouseId` | 仓库 ID |
| `--product-id` | `productId` | 产品 ID |
| `--product-name` | `productName` | 产品名称（**模糊搜索**） |
| `--no` | `no` | 单号 |
| `--status` | `status` | 状态 |
| `--type` | `type` | 类型 |
| `--start-time` | `{timeKey}[0]` | 时间范围起始（如 `2026-07-01 00:00:00`） |
| `--end-time` | `{timeKey}[1]` | 时间范围结束（如 `2026-07-31 23:59:59`） |

> **时间范围说明**：采购入库的 `timeKey` 为 `inTime`（即后端收到 `inTime[0]`/`inTime[1]`），销售出库为 `outTime`（即 `outTime[0]`/`outTime[1]`）。CLI 统一用 `--start-time`/`--end-time` 暴露，内部按模块自动转换。
>
> **时间格式**：`yyyy-MM-dd HH:mm:ss`（示例中的 URL 使用 `2026-07-08 00:00:00` 格式）。

---

## 销售出库 (`wlt sale out`)

**API 路径**：`/erp/sale-out`

### 模块特有筛选字段

| Flag | 后端参数 | 说明 |
|------|----------|------|
| `--customer-id` | `customerId` | 客户 ID |
| `--batch-no` | `batchNo` | 批次号 |

### 子命令一览

| 命令 | 说明 | 必填参数 |
|------|------|----------|
| `wlt sale out list` | 分页查询销售出库单 | — |
| `wlt sale out page-count` | 统计销售出库单数量 | — |
| `wlt sale out get` | 获取出库单详情 | `--id` |
| `wlt sale out create` | 创建销售出库单 | `--data` |
| `wlt sale out update` | 更新销售出库单 | `--data` |
| `wlt sale out delete` | 删除销售出库单 | `--ids` |
| `wlt sale out update-status` | 更新出库单状态 | `--data`（含 id 和 status） |

### 查询示例

```bash
# 查询指定产品的出库单（模糊搜索产品名）
wlt sale out list --product-name 习惯 --profile sit --token <token> --tenant-id <tid>

# 按时间范围 + 批次号筛选
wlt sale out list --start-time "2026-07-18 00:00:00" --end-time "2026-08-12 23:59:59" \
  --batch-no PC --warehouse-id 6 --profile sit --token <token> --tenant-id <tid>

# 匹配真实 URL：sale-out/page?no=XSCK20260402000001&outTime[0]=...&outTime[1]=...&productName=习惯&batchNo=PC&warehouseId=6
wlt sale out list --no XSCK20260402000001 \
  --start-time "2026-07-18 00:00:00" --end-time "2026-08-12 23:59:59" \
  --product-name 习惯 --batch-no PC --warehouse-id 6 \
  --profile sit --token <token> --tenant-id <tid>
```

---

## 销售退货 (`wlt sale return`)

**API 路径**：`/erp/sale-return`

> ⚠️ `sale return list` 在当前环境返回 403（后端权限配置，非 CLI 缺陷）。

### 子命令一览

| 命令 | 说明 | 必填参数 |
|------|------|----------|
| `wlt sale return list` | 分页查询销售退货单 | — |
| `wlt sale return page-count` | 统计销售退货单数量 | — |
| `wlt sale return get` | 获取退货单详情 | `--id` |
| `wlt sale return create` | 创建销售退货单 | `--data` |
| `wlt sale return update` | 更新销售退货单 | `--data` |
| `wlt sale return delete` | 删除销售退货单 | `--ids` |
| `wlt sale return update-status` | 更新退货单状态 | `--data` |

---

## 销售核算 (`wlt sale accounting`)

**API 路径**：`/erp/sale-accounting`

### list / page-count / export 共用筛选字段

| Flag | 后端参数 | 说明 |
|------|----------|------|
| `--no` | `no` | 核算单号 |
| `--customer-name` | `customerName` | 客户名称 |
| `--waybill-no` | `waybillNo` | 运单号 |
| `--period-start` | `periodRange[0]` | 核算周期起（`yyyy-MM-dd`，如 `2026-07-01`） |
| `--period-end` | `periodRange[1]` | 核算周期止（`yyyy-MM-dd`，如 `2026-07-31`） |

> 注意核算周期为纯日期格式（`yyyy-MM-dd`），与出库单的 `--start-time`/`--end-time`（`yyyy-MM-dd HH:mm:ss`）不同。

### 子命令一览

| 命令 | 说明 | 必填参数 |
|------|------|----------|
| `wlt sale accounting list` | 分页查询销售核算单 | — |
| `wlt sale accounting page-count` | 按筛选统计（Σ总金额、Σ运单数量） | — |
| `wlt sale accounting get --id <N>` | 获取核算单详情 | `--id` |
| `wlt sale accounting create` | 创建销售核算单 | `--data` |
| `wlt sale accounting update` | 更新销售核算单 | `--data` |
| `wlt sale accounting delete` | 删除销售核算单 | `--ids`（逗号分隔） |
| `wlt sale accounting waybill-list` | 分页查询核算候选运单（游离+已绑定，`settlementId` 非空标记 locked） | — |
| `wlt sale accounting waybill-count` | 统计核算候选运单总数 | — |
| `wlt sale accounting calculate` | 试算核算价（纯内存计算，不写表） | `--data` |
| `wlt sale accounting default-config` | 获取默认结算税率（infra 参数 `erp.sale-accounting.tax-rate`，缺省 13.000） | — |
| `wlt sale accounting export` | 导出销售核算单 Excel | 与 list 同筛选字段（另支持 `--headers`） |

### 核算候选运单筛选字段（waybill-list / waybill-count）

| Flag | 后端参数 | 说明 |
|------|----------|------|
| `--customer-id` | `customerId` | 客户 ID |
| `--exclude-ids` | `excludeIds` | 排除的运单 ID（逗号分隔） |
| `--waybill-no` | `waybillNo` | 运单号 |
| `--car-number` | `carNumber` | 车牌号 |
| `--load-date-start` / `--load-date-end` | `realLoadDate[0]` / `realLoadDate[1]` | 实际装货时间范围 |
| `--unload-date-start` / `--unload-date-end` | `realUnloadDate[0]` / `realUnloadDate[1]` | 实际卸货时间范围 |

### 查询示例

```bash
# 按核算周期 + 客户名称筛选
wlt sale accounting list --customer-name 某客户 \
  --period-start 2026-07-01 --period-end 2026-07-31

# 核算候选运单（按客户 + 实际装卸时间）
wlt sale accounting waybill-list --customer-id <客户ID> \
  --load-date-start "2026-07-01 00:00:00" --load-date-end "2026-07-31 23:59:59"

# 试算核算价（纯内存计算，不写表；按规则质检类型取运单最新已审核质检单指标）
wlt sale accounting calculate --data '{"rule":{...},"items":[{...}],"settleQtyType":"...","priceType":"...","taxRate":13.000}'

# 默认结算税率
wlt sale accounting default-config

# 导出 Excel（参数同 list）
wlt sale accounting export --period-start 2026-07-01 --period-end 2026-07-31
```

---

## 采购入库 (`wlt purchase in`)

**API 路径**：`/erp/purchase-in`

### 模块特有筛选字段

| Flag | 后端参数 | 说明 |
|------|----------|------|
| `--supplier-id` | `supplierId` | 供应商 ID |
| `--metrics-name` | `metricsName` | 检测指标名（如 `含水`） |

### 子命令一览

| 命令 | 说明 | 必填参数 |
|------|------|----------|
| `wlt purchase in list` | 分页查询采购入库单 | — |
| `wlt purchase in page-count` | 统计采购入库单数量 | — |
| `wlt purchase in get` | 获取入库单详情 | `--id` |
| `wlt purchase in create` | 创建采购入库单 | `--data` |
| `wlt purchase in update` | 更新采购入库单 | `--data` |
| `wlt purchase in delete` | 删除采购入库单 | `--ids` |
| `wlt purchase in update-status` | 更新入库单状态 | `--data`（含 id 和 status） |

### 查询示例

```bash
# 按供应商 + 产品名模糊搜索
wlt purchase in list --supplier-id 2001489305039032322 --product-name 西瓜 \
  --profile sit --token <token> --tenant-id <tid>

# 按时间范围 + 检测指标筛选
wlt purchase in list --start-time "2026-07-08 00:00:00" --end-time "2026-08-05 23:59:59" \
  --metrics-name 含水 --warehouse-id 7 --profile sit --token <token> --tenant-id <tid>

# 匹配真实 URL：purchase-in/page?no=CGRK20260508000009&supplierId=...&productName=西瓜&warehouseId=7&inTime[0]=...&inTime[1]=...&metricsName=含水
wlt purchase in list --no CGRK20260508000009 \
  --supplier-id 2001489305039032322 --product-name 西瓜 --warehouse-id 7 \
  --start-time "2026-07-08 00:00:00" --end-time "2026-08-05 23:59:59" \
  --metrics-name 含水 --profile sit --token <token> --tenant-id <tid>
```

---

## 采购退货 (`wlt purchase return`)

**API 路径**：`/erp/purchase-return`

> ⚠️ `purchase return list` 在当前环境返回 403（后端权限配置，非 CLI 缺陷）。

### 子命令一览

| 命令 | 说明 | 必填参数 |
|------|------|----------|
| `wlt purchase return list` | 分页查询采购退货单 | — |
| `wlt purchase return page-count` | 统计采购退货单数量 | — |
| `wlt purchase return get` | 获取退货单详情 | `--id` |
| `wlt purchase return create` | 创建采购退货单 | `--data` |
| `wlt purchase return update` | 更新采购退货单 | `--data` |
| `wlt purchase return delete` | 删除采购退货单 | `--ids` |
| `wlt purchase return update-status` | 更新退货单状态 | `--data` |
