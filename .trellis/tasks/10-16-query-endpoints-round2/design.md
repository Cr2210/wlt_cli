# 技术设计 — Round 2 查询端点对接与 export `-o` 方案

> 事实源：后端 `D:\Projects\weiliantong\weiliantong-api` 分支 `ge-dev-1.8`（2026-10-19，HEAD `7011df50`）。
> 盘点依据：`research/01-endpoint-inventory.md`；决策依据：`research/02-decisions-facts.md`。
> `$C` = `wlclouds-module-erp/src/main/java/com/wlclouds/module/erp/controller/admin/`。

## 1. 命令树变更总览

| 业务域 | 新增命令 | 端点 | 后端锚点（参数来源） |
|---|---|---|---|
| `wlt position`（**新包** `cmd/position/`） | `latest-follow`、`history-follow` | `/erp/position/get-latest-waybill-follow`、`/get-history-waybill-follow` | `$C/position/ErpPositionController.java:30 :59` |
| `wlt screen` | `line`、`line-price`、`purchase-enterprise`、`purchase-in-price`、`purchase-in-product`、`settle-amount` | `/erp/screen/*` | `$C/screen/ErpScreenController.java:96 :90 :103 :70 :76 :83` |
| `wlt system account`（**新包** `cmd/system/` 下 `sys_account.go`） | `list` + `car-detail`、`carrier-detail`、`driver-detail`、`enterprise-info`、`latest-waybill-follow`、`waybill`、`waybill-follow`、`waybill-source-info` | `/erp/system/account/list` + 8 个详情 | `$C/system/ErpSystemAccountController.java:47 :92 :98 :86 :34 :73 :80 :67 :61` |
| `wlt system weight-account` | `get` | `/erp/system/weight-account/get` | `$C/system/ErpSystemWeightAccountController.java:24` |
| `wlt waybill` | `events`、`order-info` | `/erp/waybill/get-events`、`/get-order-info` | `$C/waybill/WaybillSourceController.java:98 :104` |
| `wlt invoice item`（`cmd/invoice/invoice_item.go`） | `list` | `/erp/invoice-item/list` | `$C/invoice/InvoiceItemController.java:31` |
| `wlt finance invoice` | `create`、`update`、`update-status`、`delete`；**删** `page-count` | `/erp/finance-invoice/*` | `$C/finance/FinanceInvoiceController.java:33 :39 :45 :51` |
| `wlt finance payment-apply` | `create`、`update`、`update-status`、`delete`、`export`；**删** `page-count` | `/erp/finance-payment-apply/*` | `$C/finance/ErpFinancePaymentApplyController.java:37 :44 :52 :59 :89` |
| `wlt finance prepayment-apply` | `page-count` 之外补齐 `create`、`update`、`update-status`、`delete`、`export`；`list` 子命令**改名** `page`（对齐另两模块，见 §4） | `/erp/finance-prepayment-apply/*` | `$C/finance/ErpFinancePrepaymentApplyController.java:37 :44 :52 :59 :89` |

## 2. 实现模式选型（对照 `.trellis/spec/cli/command-factories.md`）

- **screen 6 个**：无参/少参 GET → 照抄 `cmd/screen/screen.go:32-50` 现有模式
  （`EnsureClient` → `Get` → `OutputJSON`）。若后端方法签名带 query 参数（实施时核对），
  按参数表加离散 flag。
- **position / waybill events·order-info / system account 8 详情 / weight-account get / invoice item list**：
  端点形态各异（单值参、分页），不适合 CRUD 工厂 → 自定义查询命令，
  模板 `cmd/finance/finance_invoice.go`（分页用 `ParsePagedJSON`，单查询用 `OutputJSON`；
  裸数组响应**不可**用 `ParsePagedJSON`，见 spec output-protocol）。
- **system account list**：分页列表 → `page-no`/`page-size` flag + `ParsePagedJSON`。
- **财务写命令**：模板 `cmd/finance/finance_receipt_payment.go:174-278`
  （create/update 用 `--data` JSON 直传 + Post/Put；delete 用 id 查询参；
  update-status 同 create 模式）。
- **参数取法**（每批次第一步）：读 Controller 方法签名的 `@RequestParam` / `@Valid` VO 字段
  （VO 类在方法签名内，跳转定义取字段），flag 名 kebab-case，参数名 camelCase（spec flags-and-params）。

## 3. `-o` 导出方案设计

### 现状与核心事实

- 后端 export 端点经 `HttpServletResponse` 直写**二进制 Excel 流**（非 CommonResult 信封）；
  现有 23 处用 `Get()`（走 CommonResult 解析，`internal/client/client.go:150-158`）——
  链路本身对二进制就不正确，改造是**功能修复**。
- **`Client.DoRaw` 已存在**（`internal/client/client.go:85-97`）：返回原始 body、
  不做 CommonResult 解析——方案落地的关键能力已具备，无需改 client。

### 设计

1. **全局 flag**：`root.go:88` 附近加
   `rootCmd.PersistentFlags().StringP("output", "o", "", "导出文件写入路径(export 命令生效;缺省用自动文件名)")`
2. **公共 helper**：新文件 `internal/cmdutil/export.go`：
   ```go
   // RunExport 调用 export 端点并将二进制响应写入文件，stdout 输出信封 JSON。
   func RunExport(ctx context.Context, path string, params map[string]any) error
   ```
   内部：`DoRaw` → 读 `--output` flag（空则自动名 `wlt-export-20060102-150405.xlsx`）→
   `os.WriteFile` → stdout 输出 `{"data":{"file":"<绝对路径>","bytes":N}}`（维持 stdout 机器可读协议）。
3. **23 处改造**：`Get(...)+fmt.Println` → `cmdutil.RunExport(ctx, apiPath+"/export-excel", params)`，
   参数装配代码保留；`--headers` flag（logistics/sale-accounting/freight-apply）继续收集进 params。
4. **不做**（留后续）：Content-Disposition 文件名解析、后端有而 CLI 无 export 命令的补齐（research/03 §5）。

### 兼容性

- 行为变更：export 命令从「stdout 打印字符串」变为「写文件 + 信封 JSON」——
  原 23 处输出本就不可用（二进制被 CommonResult 解析破坏），无真实调用方依赖，直接切换。
- `-o` 是新 flag，无破坏面。

## 4. prepayment-apply `list` 改名 `page` 的理由

`finance_prepayment_apply.go:42` 子命令名 `list` 实调 `/page` 端点，与同域另两模块（`page`）不一致
（research/02 事实表）。趁本轮补写命令一并统一为 `page`，并在 `skills/SKILL.md` 对应表格同步改名
（属本轮改动文件的文档行级更新，不算"修改 skills 内容"的范围外事项——若实施中发现 SKILL.md
无该行则跳过）。

## 5. 验证与回滚

- 每批次：`go build ./... && go test ./...`；命令树核对 `wlt <domain> --help`。
- helper 单测：`internal/cmdutil/export_test.go`（临时目录写入 + 信封输出断言 + `-o` flag 读取）。
- 真机实测（需用户提供有效 token，最后一批）：prd 验收清单第 4 条的 Round 1 遗留项 +
  本轮新命令各抽 1 条 + 3 个已删 page-count 确认 404。
- 回滚：纯新增命令零风险；`-o` 改造 23 处集中在 helper 一个函数，回滚 = revert 提交。

## 6. 风险

| 风险 | 缓解 |
|---|---|
| screen 6 端点 0923 基线未覆盖，后端可能不稳定 | 已核实 ge-dev-1.8 存在；实测批统一验证 |
| 后端 VO 字段与 erp.md 时代差异 | 参数一律取自 Controller/VO，不参考旧文档 |
| 无 token 阻塞实测 | 开发与单测不依赖 token；实测批放最后，可与交付解耦 |
