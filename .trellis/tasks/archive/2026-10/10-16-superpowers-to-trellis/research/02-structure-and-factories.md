<!-- 本文件是 `conventions-extraction` 的拆分件之一（按研究产物约定：一主题一文件）。
     初衷：单文件 57KB 超过 context_injection.max_file_bytes(32768) 会被截断。
     全量导航见同目录 `00-index.md`。文中形如「§C-1」「§4.5」的交叉引用按该导航定位。 -->

## 2. 命令分包与注册约定

### 2.1 package-per-domain 布局

- `cmd/<domain>/` 一个业务域一个 Go 包，包内暴露 `func Register(parent *cobra.Command)`。25 个包全部具备该函数（`grep -rn "func Register(" cmd/` 命中 25 个文件，含 `cmd/config/config.go:11`、`cmd/contract/contract.go:12`、…、`cmd/weight/weight.go:12`）。
- 典型最小实现（无额外逻辑的域）：`cmd/product/product.go:5-14`

```go
var productCmd = &cobra.Command{Use: "product", Short: "产品管理", Long: "产品模块操作：产品、分类、单位、指标。"}
// Register adds the product command to the parent command.
func Register(parent *cobra.Command) { parent.AddCommand(productCmd) }
```

- 例外形态：`cmd/profit/` 把 `Register` 放在 `register.go:6`（而非 `profit.go`）。

### 2.2 `init()` 注册约定

- 子命令在**各自文件的 `init()`** 内挂到包级父命令变量。例：
  - 工厂式：`cmd/product/product_main.go:14-33`（`init()` 内 `AddCRUDToParent` + 附加命令）
  - 手写式：`cmd/settlement/settlement_main.go:14-16`（挂 `settlementMainCmd`）、`:23-37`（挂 11 个子命令）
  - 新命令组式：`cmd/settlement/logistics.go`、`cmd/sale/accounting.go`、`cmd/finance/freight_apply.go`、`cmd/order/order_prepayment.go` 均在 `init()` 内 `xxxCmd.AddCommand(groupCmd)` 后再挂子命令（见 `sync-plan` Task 3/4/5 Step 3 的「实现全文」，与仓库现状一致）。
- `cmd/stock/stock.go:11-13` 是空 `init()`，注释说明子命令由同包其他文件注册（文档化约定）。
- 根命令统一在 `cmd/root.go:87-123` 的 `init()` 调 25 个 `Register(rootCmd)`；`cmd/api.go:91`、`cmd/version.go:25`、`cmd/completion.go`、`cmd/help.go` 直接在 `init()` 内 `rootCmd.AddCommand(...)`。

### 2.3 新增业务域步骤（可从代码直接推出的最小集）

1. 建 `cmd/<domain>/` 目录与包级父命令变量（`Use`/`Short`/`Long`）。
2. 写 `func Register(parent *cobra.Command) { parent.AddCommand(<domain>Cmd) }`。
3. 业务文件内 `init()` 挂子命令（工厂或手写）。
4. `cmd/root.go` 顶部 import 别名 + `init()` 内加 `<domain>.Register(rootCmd)`（`cmd/root.go:99-122` 两个分区：`// 基础设施` 与 `// 业务域`）。
5. 若复用 `cmd/partner/`：在 `init()` 调 `partner.AddPartnerCommands(...)` 等导出函数（见 `cmd/customer/customer.go:14-20`、`cmd/supplier/supplier.go:22` 的 `Register`）。

### 2.4 域内文件命名约定（实测样本）

| 域 | 文件命名 | 样本 |
|---|---|---|
| 一个域多子模块 | `<域>_<子模块>.go` | `cmd/product/product_main.go`、`product_unit.go`、`product_category.go`、`product_metrics.go`；`cmd/order/order_main.go`、`order_plan.go`、`order_prepayment.go`；`cmd/stock/stock_in.go`、`stock_out.go`、`stock_move.go`、`stock_check.go`、`stock_warehouse.go`、`stock_query.go`、`stock_record.go` |
| 父命令 + Register | `<域>.go` | `cmd/product/product.go`、`cmd/stock/stock.go`、`cmd/settlement/settlement.go` |
| 新命令组（0923 后） | `<子模块>.go`（不带域前缀） | `cmd/settlement/logistics.go`、`cmd/sale/accounting.go`、`cmd/finance/freight_apply.go` |
| 共享内部辅助（包内） | `<域>_doc.go` | `cmd/stock/stock_doc.go`（共享 `FlagSpec` + `newStockPagedCmd`/`newStockCountCmd`） |
| 测试 | `<被测文件>_test.go` | `cmd/stock/stock_doc_test.go`、`cmd/order/order_prepayment_test.go`、`cmd/settlement/settlement_main_test.go`、`cmd/finance/freight_apply_test.go` |

> 注意两种命名风格并存：0923 之前用 `<域>_<子模块>`，0923 新模块用裸 `<子模块>`（`logistics.go`/`accounting.go`/`freight_apply.go`）。代码中不存在强制规则，spec 应据实描述为「两种并存」。

### 2.5 `cmd/partner/` 共享包

客户/供应商复用同一包：导出 `AddPartnerCommands(parent, partnerType, label)`（`cmd/partner/partner.go:15`）、`NewPartnerListCmd`（`:30`）、`NewPartnerSimpleListCmd`（`:62`）、`NewPartnerPageCountCmd`（`:81`）、`NewPartnerUpdateAuditStatusCmd`（`:104`）、`NewPartnerDeleteListCmd`（`:129`）、`NewCreditGroup`（`partner_credit.go:20`）、`NewInvoiceGroup`（`partner_invoice.go:15`）、`NewSettlementGroup(partnerType)`（`partner_settlement.go:25`）。调用方 `cmd/customer/customer.go`、`cmd/supplier/supplier.go`。partnerType 常量字符串 `"CUSTOMER"`（`cmd/customer/customer.go:15`）。

---

## 3. 工厂选型

### 3.1 `internal/cmdutil/` 工厂全清单（按文件）

**`crud.go` — 标准 CRUD 工厂族**

| 符号 | 行 | 生成什么 |
|---|---|---|
| `FlagSpec{Name, Usage}` | `:15-18` | 筛选 flag 声明（Name 为 kebab-case） |
| `CRUDConfig{Name, APIPath, Label, ListFilters, SingleDelete, SkipStatus}` | `:21-28` | 标准 CRUD 配置 |
| `NewCRUDGroup(cfg)` | `:31-38` | 新建**具名子命令**并对其调用 `AddCRUDToParent` |
| `AddCRUDToParent(parent, cfg)` | `:41-52` | 直接向 parent 挂 list/get/create/update/delete（+ update-status，除非 `SkipStatus`） |
| `CrudListCmd` | `:57-59` | 只是 `CrudListCmdWithFixed(..., nil)` 的向后兼容包装 |
| `CrudListCmdWithFixed` | `:65-97` | `list` → `GET {apiPath}/page`，注入固定 query（`type` 等）+ `pageNo/pageSize`，输出 `ParsePagedJSON` |
| `CrudGetCmd` | `:101-120` | `get` → `GET {apiPath}/get?id=`，`--id` 必填，输出 `OutputJSON` |
| `CrudCreateCmd` | `:124-147` | `create` → `POST {apiPath}/create`，`--data` 必填 JSON |
| `CrudUpdateCmd` | `:151-174` | `update` → `PUT {apiPath}/update`，`--data` 必填 |
| `CrudDeleteCmd(apiPath,label,single)` | `:178-217` | `delete` → `DELETE {apiPath}/delete?id=`（single）或 `?ids=`（非 single） |
| `CrudUpdateStatusCmd` | `:221-244` | `update-status` → `PUT {apiPath}/update-status`，`--data` JSON body |
| `CrudSimpleListCmd` | `:248-264` | `simple-list` → `GET {apiPath}/simple-list`（无参） |
| `CrudSimpleListCmdWithFlags` | `:268-291` | `simple-list` + 筛选 flag |
| `CrudPageCountCmd` | `:296-298` | 包装 `CrudPageCountCmdWithFixed(..., nil)` |
| `CrudPageCountCmdWithFixed` | `:301-327` | `page-count` → `GET {apiPath}/page-count`，输出 `OutputJSON` |
| `CrudListAllCmd` | `:971-994` | `list` → `GET {apiPath}/list`（**非**分页，`OutputJSON`） |
| `ContractTypeConfig` / `ContractTypeCmds` | `:365-375` / `:395-427` | 合同子类型族：每类型 list/page-count/get/update-status/create/update/delete |
| `PlanListCmdConfig` / `PlanListCmds` | `:538-543` / `:546-630` | 计划子类型族：仅 list + page-count，走 `/erp/order-plan` |
| `OrderTypeConfig` / `OrderTypeCmds` | `:683-688` / `:692-765` | 订单子类型族：仅 list + page-count，走 `/erp/order` |
| `DocumentConfig` / `DocumentCmds` | `:856-864` / `:868-883` | 单据族：list(+可选 page 别名)/page-count/get/create/update/delete/update-status |

**`timeflags.go`** — 通用时间范围：`AddTimeRangeFlags(c, flagName, paramKey, label)` `:8-11`；`CollectTimeRangeFlags(cmd, params, flagName, paramKey)` `:14-21`。
**`stats.go`** — `NewStatsGetCmd(name, apiPath, label, extraFlags)` `:15-40`；`CollectStatsFlags` `:43-62`（默认 `type=month`、`startTime=当月1号`、`sortBy=amount`）；`DefaultMonthRange` `:65-69`；`AddStatsFlags` `:72-76`。
**`helpers.go`** — `ParseJSONData` `:13`、`flagToParamKey` `:28`、`CollectStringFlag` `:58`、`CollectStringFlags` `:67`、`CollectIntFlags` `:78`。
**`output.go`** — 见 §5。 **`client.go`** — 见 §1。

### 3.2 分工判据（能从代码推断的规则）

**判据 1 — 是否「单据类」（进销存单据 + 退货）→ `DocumentCmds`。**
8 个真实调用点：`cmd/stock/stock_in.go:6`、`stock_out.go:6`、`stock_move.go:6`、`stock_check.go:6`、`cmd/purchase/purchase_in.go:6`、`purchase_return.go:6`、`cmd/sale/sale_out.go:6`、`sale_return.go:6`。
触发能力（`DocumentConfig` 4 个可配项，`crud.go:856-864`）：`TimeKey`（非空→注册 `--start-time/--end-time` 并转 `{TimeKey}[0]/[1]`）、`Headers`（→ 注册 `--headers`）、`PageAlias`（→ list 加 `page` cobra 别名）。`Filters` 必须**显式列全** —— `stock-design:57` 明确「**无隐式 base flags**」。
`PageAlias: true` 仅见于 stock 四个单据：`cmd/stock/stock_{in,out,move,check}.go:13`。

**判据 2 — 是「标准单实体 CRUD」（一实体一端点）→ `AddCRUDToParent` / `NewCRUDGroup`。**
- `AddCRUDToParent` 真实调用点仅 2 处：`cmd/product/product_main.go:16`（产品，`SingleDelete: true`）、`cmd/stock/stock_warehouse.go:17`（仓库）。
- `NewCRUDGroup`（建具名子命令）调用点 4 处：`cmd/product/product_unit.go:6`、`product_metrics.go:15`、`cmd/partner/partner_credit.go:21`、`partner_invoice.go:16`。
- 边界：一个 parent 下要挂**多个**具名 CRUD 子组 → `NewCRUDGroup`；parent 本身就是该实体 → `AddCRUDToParent`。

**判据 3 — 需要「共享端点 + 类型区分」→ `ContractTypeCmds` / `PlanListCmds` / `OrderTypeCmds`。**
共同前提：多个子类型共用同一后端端点、以 `type=` 区分（`crud.go:363-364` 注释、`:408`/`:415` 注入 `fixed{"type": tc.Type}`）。`ContractTypeConfig.APIPath` 可覆盖默认 `/erp/contract`（`crud.go:370-372`、`resolveAPIPath` `:378-383`，销售合同走 `/erp/service-contract`）。`HasDateRange=true` 时换用带日期范围的 `contractListCmd`（`crud.go:405-409`、`:430-462`）。

**判据 4 — 形状各异 / 非 CRUD → 域内手写 + 复用 `Crud*Cmd` 单件。**
- 0923 三个新组即此模式：`sync-plan:7`「CRUD 骨架复用 `cmdutil.CrudGetCmd/CrudCreateCmd/CrudUpdateCmd/CrudDeleteCmd` 工厂，特殊端点在域文件内自定义」。实际：`cmd/settlement/logistics.go`、`cmd/sale/accounting.go`、`cmd/finance/freight_apply.go` 的 `init()` 直接 `AddCommand(cmdutil.CrudGetCmd(apiPath, label), ...)`。
- `order-design:28` 明确「本模块 5 个接口形状各异，不适合套用标准 CRUD 工厂」→ `cmd/order/order_prepayment.go` 全手写。
- `stock-design:122` 决策 4 明确 query/record **不套工厂**（「count/ledger/record-page/total-cost 等非 CRUD 端点有真实业务逻辑，套工厂会过度抽象」）→ 改为包内共享 `FlagSpec` + `newStockPagedCmd`/`newStockCountCmd`（`cmd/stock/stock_doc.go:69` / `:100`）。

**判据 5 — 写命令用 `--data` JSON 直传，除非「字段少且固定」。**
`sync-design:46`、`sync-plan:16`：create/update 一律 `--data`（`ParseJSONData`）。反例（离散 flag）见 §4.7。

> **代码不支持「一条硬规则」的部分（spec 应显式写明）**：`DocumentCmds` vs `AddCRUDToParent` 的分界是「单据（分页查询/单据语义/时间范围/导出 headers）」vs「主数据实体」，无强制约束；两者并存且域间不一致（`stock_warehouse`/`product` 用 AddCRUDToParent，但 `product_unit`/`product_metrics` 用 NewCRUDGroup）。spec 应写「按语义选择」+ 上述调用点清单作为现实样本。

### 3.3 已删除的 legacy 工厂

`internal/cmdutil/crud_legacy.go` **已不存在**（`ls internal/cmdutil/` 仅 client.go / crud.go / crud_test.go / helpers.go / helpers_test.go / output.go / stats.go / timeflags.go / timeflags_test.go）。`stock-plan` Task 8 要求删除；`stock-design:113` 记录「退货迁离后删除 crud_legacy.go（224 行，全仓再无调用方）」。全仓 `NewCRUDSubCmd` **UNVERIFIED：无任何命中**。

---

