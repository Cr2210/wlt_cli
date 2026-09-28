# wlt 后端迭代同步（2026-09-23 基线）设计

日期:2026-09-28 ｜ 状态:已确认,待实现

## 背景

CLI 上次对接后端为 2026-08-21（6318e49）。后端 dev-1.7/ge-dev-1.8 已迭代至 2026-09-23（2e52721b），差异分析结论：**3 个全新模块 + 3 处已有域增量**，其余为后端内部口径调整或响应字段透传（CLI 是薄透传层，无需改动）。

| 类别 | 模块 | 后端前缀 | 创建时间 | 端点数 |
|---|---|---|---|---|
| A1 | 物流结算 | `/erp/logistics-settlement` | 08-28 | 13 |
| A2 | 销售核算 | `/erp/sale-accounting` | 09-15 | 11 |
| A3 | 运费申请（物流付款申请） | `/erp/finance-freight-apply` | 09-17 | 11 |
| B1 | 结算单忽略剩余 | `/erp/settlement/ignore-remainder` + `invoiceFlag` 提交字段 | 09-22/23 | 1 |
| B2 | 预付关联双业务类型 | `bizType`（PURCHASE/LOGISTICS） | 09-15 | flag 增量 |
| B3 | 收付款预付标识 | `isPrepaid` 查询筛选 | 09-14/15 | flag 增量 |

**排除项**（后端内部逻辑或响应透传，CLI 无需改动）：发票核销/开票申请卡口过滤、出入库详情磅单回填、驾驶舱/资金占用/物流预付统计口径、WebSocket 鉴权、通知公告评论回复（CLI 未接 notice）、扫码接单/费用下发。`invoiceFlag`/`isPrepaid` 的提交字段经 `--data` JSON 直传天然支持，仅在 skills 文档补充说明。

## 目标

- 新增 3 个子命令组覆盖 A 类 37 个端点，挂载与命名遵循全仓惯例。
- B 类增量改动 3 处既有命令。
- 核实既有 `unsettle-waybill` 命令路径有效性（结论：有效，无需改动）。
- 同步更新 CLAUDE.md、skills 文档、query-endpoints backlog。

## 非目标

- 不封装组合工作流（如"试算→创建核算单"一键命令），保持单接口单命令。
- 不改动 `DocumentCmds` / `AddCRUDToParent` / `CrudXxxCmd` 通用工厂签名；不足之处用域内自定义命令补齐。
- 不补 notice（通知公告）域——不在本次后端差异主线上，如需要另立任务。

## 挂载决策（已确认）

| 模块 | 挂载 | 文件 |
|---|---|---|
| A1 物流结算 | `wlt settlement logistics <sub>`（settlement 域子命令组，方案 A） | `cmd/settlement/logistics.go` |
| A2 销售核算 | `wlt sale accounting <sub>` | `cmd/sale/accounting.go` |
| A3 运费申请 | `wlt finance freight-apply <sub>` | `cmd/finance/freight_apply.go` |
| B1 ignore-remainder | `wlt settlement main ignore-remainder`（与 settlement main CRUD 平级） | `cmd/settlement/settlement_main.go` |
| B2 biz-type | `wlt order prepayment-relation list/create` 加 `--biz-type` | `cmd/order/` 对应文件 |
| B3 is-prepaid | `wlt finance` 收付款 list 加 `--is-prepaid` | `cmd/finance/finance_receipt_payment.go` |

## 设计

通用约定（全仓既有）：stdout 数据 JSON / stderr 错误 JSON / 退出码 0-6；`EnsureClient()` 鉴权；分页 `--page-no/--page-size`；时间筛选 `--start-time/--end-time`（映射 `{timeKey}[0]/[1]`）；create/update 用 `--data` JSON 直传；无参 GET 直接透传。

### A1 物流结算 `wlt settlement logistics`（13 端点）

| 子命令 | 端点 | flags |
|---|---|---|
| list | GET /page | `--no --name --carrier-name --carrier-enterprise-id --status --start-time/--end-time(settlementDate) --page-no --page-size` |
| page-count | GET /page-count | 同 list 筛选 |
| get | GET /get | `--id` |
| create | POST /create | `--data`（含 invoiceFlag/items） |
| update | PUT /update | `--data` |
| delete | DELETE /delete | `--ids` |
| update-status | PUT /update-status | `--id --status` |
| mismatch-list | GET /mismatch-page | `--keyword --page-no --page-size` |
| available-waybill | GET /available-waybill | `--carrier-id --carrier-name --keyword --settlement-id --page-no --page-size` |
| available-waybill-count | GET /available-waybill/count | 同上筛选 |
| recalculate-all | GET /recalculate-all | 无 |
| ignore-remainder | PUT /ignore-remainder | `--id --type(SETTLE/INVOICE) --reason` |
| export | GET /export-excel | 同 list 筛选 + `--headers`（对齐 settlement main export 模式） |

settlement 域 Long 描述同步更新为"结算单CRUD、物流结算、取消结算运单"。

### A2 销售核算 `wlt sale accounting`（11 端点）

| 子命令 | 端点 | flags |
|---|---|---|
| list | GET /page | `--no --customer-name --period-start/--period-end(periodRange，yyyy-MM-dd) --waybill-no --page-no --page-size` |
| page-count | GET /page-count | 同 list 筛选 |
| get | GET /get | `--id` |
| create | POST /create | `--data`（含 periodRange 数组、ruleSnapshot、items） |
| update | PUT /update | `--data` |
| delete | DELETE /delete | `--ids` |
| waybill-list | GET /waybill-page | `--customer-id --exclude-ids(可多次) --waybill-no --car-number --start-time/--end-time(realLoadDate) --real-unload-start/--real-unload-end --page-no --page-size` |
| waybill-count | GET /waybill-page/count | 同 waybill-list 筛选 |
| calculate | POST /calculate | `--data`（rule + items + settleQtyType/priceType/taxRate） |
| default-config | GET /default-config | 无（返回默认税率） |
| export | GET /export-excel | 同 list 筛选 + `--headers` |

### A3 运费申请 `wlt finance freight-apply`（11 端点）

| 子命令 | 端点 | flags |
|---|---|---|
| list | GET /page | `--no --carrier-enterprise-id --carrier-name --payment-account-id --receipt-account-id --service-user-id --approve-status --start-time/--end-time(payDate) --page-no --page-size` |
| get | GET /get | `--id` |
| create | POST /create | `--data`（含 links 运单明细） |
| update | PUT /update | `--data` |
| delete | DELETE /delete | `--ids` |
| update-status | PUT /update-status | `--id --approve-status` |
| summary | GET /summary | 同 list 筛选（含 createTime 时间） |
| available-waybill | GET /available-waybill | `--waybill-no --carrier-id --settlement-id --apply-id --start-time/--end-time(realLoadDate) --real-unload 时间对 --page-no --page-size` |
| available-waybill-count | GET /available-waybill/count | 同上筛选 |
| available-settlement | GET /available-settlement | `--no --carrier-enterprise-id --carrier-name --start-time/--end-time(settlementDate) --apply-id --page-no --page-size` |
| export | GET /export-excel | 同 list 筛选 + `--headers` |

### B 类增量

1. **B1** `wlt settlement main ignore-remainder`：PUT `/erp/settlement/ignore-remainder`，`--id --type(SETTLE/INVOICE) --reason`，reason 必填（后端要求填原因）。
2. **B2** prepayment-relation 的 list 与 create 加 `--biz-type`（默认 PURCHASE，可选 LOGISTICS）。注意：LOGISTICS 时 supplier 期初不可用（后端 ORDER_PREPAYMENT_BIZ_TYPE_NOT_SUPPORTED），usage 文案注明。
3. **B3** 收付款 list 加 `--is-prepaid`（Boolean 三态：未传=不过滤）。

## 核实结论（原疑点已排除）

- CLI 现有 `wlt settlement main unsettle-waybill(-count)` 指向 `/erp/settlement/unsettle/waybill`，经核实该端点确属 `ErpSettlementController`，命令有效，无需改动。物流结算模块自身无 unsettle 端点（首次差异分析两 Controller 输出合并导致误判，已纠正为 13 端点）。

## 文档与验证

- **代码验证**：`go build ./...`；`wlt settlement logistics --help`、`wlt sale accounting --help`、`wlt finance freight-apply --help` 命令树核对；`make test`（如有相关测试）。SIT 环境（erpsit profile）可用则对 list/get 类只读命令冒烟。
- **文档**：CLAUDE.md 项目结构（settlement/sale/finance 行）；skills 相关文档补充三个新命令组与 invoiceFlag/isPrepaid/--data 说明；`docs/superpowers/specs/2026-06-27-query-endpoints-backlog.md` 增量更新。
- **提交**：按子项目分 4+1 次 commit（B 增量、A1、A2、A3、文档），遵循 conventional commits（历史惯例中文描述）。

## 实施顺序

① B 类增量（小，先行） → ② A1 物流结算 → ③ A2 销售核算 → ④ A3 运费申请 → ⑤ 文档收尾与 backlog 更新。
