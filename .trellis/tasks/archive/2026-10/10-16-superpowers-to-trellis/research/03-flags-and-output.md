<!-- 本文件是 `conventions-extraction` 的拆分件之一（按研究产物约定：一主题一文件）。
     初衷：单文件 57KB 超过 context_injection.max_file_bytes(32768) 会被截断。
     全量导航见同目录 `00-index.md`。文中形如「§C-1」「§4.5」的交叉引用按该导航定位。 -->

## 4. flag 与传参映射

### 4.1 flag 命名：kebab-case

- `FlagSpec.Name` 注释明示 kebab-case（`crud.go:16`）：`"status"`、`"category-id"`、`"carrier-enterprise-id"`。
- 分页 flag：`--page-no`（默认 1）/`--page-size`（默认 20），注册处例 `crud.go:91-92`、`:915-916`、`:455-456`、`:726-727`。
- 必填 flag 用 `c.MarkFlagRequired(...)`（返回 error 显式忽略 `_ =`）：例 `crud.go:118`、`:145`、`:196`、`:215`。

### 4.2 kebab → camelCase 转换位置与算法

- 实现在 **`internal/cmdutil/helpers.go:28-47` `flagToParamKey`**：无连字符原样返回（`status`/`no`/`pageNo`）；有连字符则逐字符扫描，遇 `-` 置 capitalize，下一 ASCII 字符大写化（`toUpperASCII` `:49-54`）。
- 三个收集函数都走它：`CollectStringFlag` `:58-63`（只收非空）、`CollectStringFlags` `:67-74`、`CollectIntFlags` `:78-85`（只收非零，值为 `strconv.FormatInt(v,10)` 字符串）。
- 表驱动单测：`internal/cmdutil/helpers_test.go:5-22`（含 `supplier-id→supplierId`、`approve-status→approveStatus`、`bill-date-start→billDateStart`、**`pageNo→pageNo` 不变**、`""→""`）。
- 后端名称：`{APIPath}/page` + `pageNo`/`pageSize`（`crud.go:74-77`、`:84`）；`GET /page?pageNo=1&pageSize=20` → `{list:[...], total:N}`（`skills/references/api-conventions.md:36-38`）。
- **注意**：转换器不做后端语义纠正；若 flag 名本身与后端字段不符（`backlog:18` 举例 `finance refund` 用 `refund-no` 而后端是 `no`、用 `customer-id/supplier-id` 而后端是 `partnerId`），转换器不救。属**已知未逐个核对的遗留**（见 §D）。

### 4.3 时间范围 → `{timeKey}[0]/[1]` 数组

两类实现**并存**：

**(a) 通用 helper（0923 新增）** `internal/cmdutil/timeflags.go`：
- `AddTimeRangeFlags(c, flagName, paramKey, label)` → flag `--{flagName}-start` / `--{flagName}-end`（`:8-11`）
- `CollectTimeRangeFlags(cmd, params, flagName, paramKey)` → `params[paramKey+"[0]"]` / `"[1]"`（`:14-21`，只收非空）
- 真实调用（`sync-plan` Task 4/5）：`("load-date","realLoadDate","装货日期")`、`("unload-date","realUnloadDate",...)`、`("pay-date","payDate",...)`、`("create-time","createTime",...)`、`("settlement-date","settlementDate",...)`
- 单测：`internal/cmdutil/timeflags_test.go:9-26`

**(b) 域内专用 helper（crud.go 内，命名不统一）**：

| 函数 | 行 | flag | 参数 key |
|---|---|---|---|
| `CollectContractDateRange` / `AddContractDateRangeFlags` | `:341-354` / `:357-361` | `order-start/order-end/end-start/end-end` | `orderDate[0]/[1]`、`endTime[0]/[1]` |
| `CollectPlanDateRange` / `AddPlanDateRangeFlags` | `:520-527` / `:530-534` | `start/end` | `startDate[0]/[1]` |
| `CollectOrderDateRange` / `AddOrderDateRangeFlags` | `:665-672` / `:675-679` | `order-start/order-end` | `orderTime[0]/[1]` |
| `collectDocumentTimeRange` / `addDocumentTimeRangeFlags`（包内私有，由 `DocumentConfig.TimeKey` 驱动） | `:787-794` / `:796-800` | `start-time/end-time` | `fmt.Sprintf("%s[0]", timeKey)` |
| `CollectSettlementDate` / `AddSettlementDateFlags` | `:812-819` / `:822-826` | `start-date/end-date` | `settlementDate[0]/[1]` |
| `CollectOrderTime` / `AddOrderTimeFlags` | `:839-846` / `:849-853` | `start-time/end-time` | `orderTime[0]/[1]` |

- 值格式：flag Usage 示例统一 `2026-07-01 00:00:00` / `2026-07-31 23:59:59`（`crud.go:783-784`）；核算周期例外为 `yyyy-MM-dd`（`sync-plan` Task 4 `period-start/period-end` → `periodRange[0]/[1]`）。
- **反例（不转数组）**：`stock query/record` 的 `in-time`/`create-time` 保持**单值**语义（`stock-design:108`：「走 `/erp/stock-record` 端点，与 stock-in 的数组范围是不同契约」），仅作普通 `FlagSpec` 收集（`cmd/stock/stock_doc.go:36-37`）。同域内两种契约并存。

### 4.4 `--headers` 导出 flag

- 定义面（`DocumentConfig.Headers: true` → `c.Flags().String("headers","","自定义导出表头")`，`crud.go:923-925`），list 收集 `CollectStringFlag(cmd, params, "headers")`（`crud.go:905-907`）；**page-count 不收集 headers**（注释「page-count 不带 headers（与原实现一致）」，`crud.go:944`）。
- 手写 export 命令自行注册：`cmd/settlement/logistics.go:297`（收集）+ `:307`（注册）；`cmd/finance/finance_receipt_payment.go:51`、`:75`（两个 filter 切片各含 `headers`）。
- **不对称事实**：`wlt settlement main export` **没有** `--headers`（`cmd/settlement/settlement_main.go:331-347` 与其 `registerSettlementListFlags` 均未注册 `headers`）；logistics / 销售核算 / 运费申请的 export 有。`skills/references/settlement-invoice.md:47` 对 main export 写「与 list 同筛选字段」，与代码一致。

### 4.5 分页 flag 与后端名

- flag `--page-no`/`--page-size`（kebab）→ 后端 `pageNo`/`pageSize`（camel）。二者**不是**由 `flagToParamKey` 生成，而是工厂里硬编码字符串键（`crud.go:74-77`）。
- 实测用 `Use: "page"` 而非 `list` 的有 10 处（`grep 'Use:   "page"'`）：`cmd/finance/finance_{invoice,invoice_apply,payment_apply,receipt_payment,refund,write_off}.go`、`cmd/produce/produce_main.go:41`、`cmd/produce/produce_plan.go:39`、`cmd/quality/quality_inspection.go:44`、`cmd/waybill/waybill_page.go:101`。⇒「分页子命令统一叫 `list`」**不成立**，只有 stock 单据 + DocumentCmds/CRUD 工厂族是 `list`。

### 4.6 `page` 别名-on-`list` 约定

- 机制：`DocumentConfig.PageAlias` → `docPageAlias(enabled)` 返回 `[]string{"page"}`（`crud.go:962-967`），挂在 `documentListCmd` 的 `Aliases`（`crud.go:890`）。
- 仅 4 个 stock 单据启用（`cmd/stock/stock_in.go:13`、`stock_out.go:13`、`stock_move.go:13`、`stock_check.go:13`）；采购/销售/退货 `PageAlias` 缺省 false（`stock-design:83`）。
- 单测断言 `list.Aliases` 含 `"page"`：`internal/cmdutil/crud_test.go:71-79`。

### 4.7 `--data` JSON 直传 与 小固定请求体离散 flag

- `--data`：`ParseJSONData(s)`（`helpers.go:13-22`，`s==""` → error；非对象 JSON → error）→ 解析失败 → `NewExitError(4, "解析 data 失败: ...", "data 应为 JSON 对象")`（`crud.go:133-136`）。用于 create/update/update-status（`crud.go:144`、`:171`、`:241`）。
- **离散 flag 构造 `map[string]any` 请求体**（小固定形状）：`order-design:51` 明确「create / update-amount 用离散 flag 而非 `--data` JSON：字段少且固定（4 个 / 2 个），`relationType` 枚举写入 flag 帮助文本，Agent 免拼 JSON；有 `data_sync resend --id`(map 请求体)先例」。实现：`cmd/order/order_prepayment.go`（`create` 构造 `{"orderId","relationType","relationId","relationAmount"}`；`update-amount` 构造 `{"id","relationAmount"}`）。
- 同类：`settlement main ignore-remainder` / `logistics ignore-remainder` 用 `body := map[string]any{"id":id,"type":ignoreType,"reason":reason}` + `MarkFlagRequired("id"/"type"/"reason")`（`sync-plan` Task 1/3；`cmd/settlement/settlement_main.go:32` 已注册，测试 `settlement_main_test.go:5-12`）。

### 4.8 PUT 的两种形态（`Put` vs `PutQuery`）

- `@RequestBody` 型 → `GetClient().Put(ctx, path, body)`（`internal/client/client.go:64`）。
- `@RequestParam` 型 → `GetClient().PutQuery(ctx, path, params)`（`internal/client/client.go:75-77`，无 body，走 query）。
- 判据来源 `sync-plan:19`：「`@RequestParam` 型 PUT（update-status）必须用 `PutQuery`，不能用 Put（那是 JSON body）」。实例：物流结算 / 运费申请的 `update-status` 用 `PutQuery`；`CrudUpdateStatusCmd` 与 ignore-remainder 用 `Put`（JSON body）。
- 另有 `DeleteWithBody`（`internal/client/client.go:81-83`）用于 DELETE + `@RequestBody`。

---

## 5. 输出协议

### 5.1 stdout 数据 / stderr 错误 分离

- **stdout**：仅数据 JSON，写者在 `internal/output/output.go`：
  - `WriteJSON(out, data)` `:60-64` → `{"data": <data>}`
  - `WritePagedJSON(out, list, total, pageNo, pageSize)` `:67-78` → `{"data": <list>, "meta": {"page_no": N, "page_size": N, "total": N}}`
  - `WriteRaw(out, data)` `:96-98` → 原样字节（供 `wlt api`）
- **stderr**：错误 JSON，`WriteError(errOut, exitErr)` `:81-93` → `{"error": {"type","code","message","hint?","detail?"}}`。
- 调用方封装（`internal/cmdutil/output.go`）：`OutputJSON(data)` `:12-14`、`OutputPagedJSON(list,total,pageNo,pageSize)` `:17-19`、`OutputRaw(data)` `:22-24`；导出 `fmt.Println("导出成功，返回数据：", string(resp.Data))`（**非 JSON**，`cmd/` 下 21 处，如 `cmd/finance/finance_account.go:365`）——仓库级既有模式。
- 退出码与错误写入总入口：`cmd/root.go:67-85`（见 §1.2/§1.6）。
- 单测：`internal/output/output_test.go`（`TestWriteJSON` `:10-27` 断言 `data` 字段；`TestWritePagedJSON` `:29-47` 断言 `meta.total`；`TestWriteError` `:49-62` 断言 `"type":"authentication"` 与 `"code":3`；`TestWriteRaw` `:64-70`）。

### 5.2 精确函数签名与位置

```go
// internal/cmdutil/output.go
func OutputJSON(data any) error                                    // :12
func OutputPagedJSON(list any, total int64, pageNo, pageSize int) error  // :17
func OutputRaw(data []byte)                                        // :22
func ParsePagedJSON(data json.RawMessage, pageNo, pageSize int) error    // :29

// internal/output/output.go
func WriteJSON(out io.Writer, data any) error                                       // :60
func WritePagedJSON(out io.Writer, list any, total int64, pageNo, pageSize int) error // :67
func WriteError(errOut io.Writer, exitErr *ExitError) error                          // :81
func WriteRaw(out io.Writer, data []byte)                                            // :96
func NewExitError(code int, msg string, hint string) *ExitError                       // :26
```

### 5.3 `ParsePagedJSON` vs `OutputJSON` 的选择

- **`ParsePagedJSON`**（`internal/cmdutil/output.go:29-43`）用于**分页端点**（`/page`、`*-page`、`available-*` 等返回 `{list,total}` 者）。实现要点：
  - `Total json.Number` —— 注释明示「后端把 total 返回为 JSON 字符串（如 `/erp/weight/page`）也能解析」（`:26-28`）；`total, _ := paged.Total.Int64()`（`:41`）。
  - 反序列化 `data` 失败 → `NewExitError(5, "解析响应失败: ...", "")`（`:34-36`）。
  - `list` 反序列化失败 → 回退 `list = []any{}`（`:38-40`）。
  - 输出走 `OutputPagedJSON`（`:42`）→ 统一信封 `{data, meta:{page_no,page_size,total}}`。
- **`OutputJSON(json.RawMessage(resp.Data))`** 用于**非分页端点**（`get`、`page-count`、`simple-list`、`list`(非分页)、`/summary`、`calculate`、`default-config`、`recalculate-all`、`update-status` 等）。
- **裸数组响应危险（必须记入 spec）**：`order-design:57` 明确「`page` 接口响应 `data` 为**裸数组**（无 `{list,total}` 包裹）→ 直接 `cmdutil.OutputJSON(resp.Data)`；**`ParsePagedJSON` 会因数组解进 struct 而报错，不可用**」。实例：`wlt order prepayment-relation list`（`cmd/order/order_prepayment.go` 的 `newOrderPrepayRelationListCmd` 用 `OutputJSON`）；同组的 `available-payments`/`available-initials` 返回标准 PageResult → 用 `ParsePagedJSON`（`order-design:58`）。
- **判据（可直接执行）**：先看接口文档/实测响应是 `{list,total}` 还是裸数组/标量；前者 `ParsePagedJSON`，后者 `OutputJSON`。`order-design:57` 也给出兜底：「若实测响应为包裹结构，切换为 `ParsePagedJSON`（一行改动）」。
- 兼容旧解析：`stock query ledger` 手写解析失败时回退 `OutputJSON`（`stock-plan` Task 7 Step 5 的实现；`stock-design:17` 记录旧代码「遇字符串型 total 报错」是本类 bug 的由来）。

---

