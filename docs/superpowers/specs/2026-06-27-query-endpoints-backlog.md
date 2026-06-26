# ERP 查询接口对接 — 交接文档

日期:2026-06-27 ｜ 状态:第 1 轮(高频模块 + bug 修复)已完成,第 2 轮待办

## 背景

按 `D:\Download\erp.md`(Swagger 导出,510 端点:GET 303 / POST 65 / PUT 92 / DELETE 50)对接「所有查询接口」。
经差集(文档 GET 端点 − CLI 已对接)后,分两轮推进。本文档记录已完成项 + 第 2 轮待办。

**决策(用户已拍板)**:export-excel 跳过(待 `-o` 下载方案);路径 bug 顺手修;先做高频模块。

---

## ✅ 第 1 轮已完成(本轮,baseline)

### 基础修复
- **filter 参数名 kebab→camelCase**:`internal/cmdutil/helpers.go` 的 `CollectStringFlag/Flags/IntFlags` 现把 `--supplier-id` 转成 `supplierId`。**此前所有多词 filter 都不生效**(后端收不到)。已加单测 `helpers_test.go`。
  - 注意:单字段 filter 的「字段名」已修,但**部分老命令的 flag 名本身与文档不符**(如 `finance refund` 用 `refund-no` 而文档是 `no`、用 `customer-id/supplier-id` 而文档是 `partnerId`)。本轮未逐个核对老命令 flag,下次可按需修。

### 路径 bug 修复(4 处,对齐文档)
| 命令 | 旧(错) | 新(文档) | 文件 |
|---|---|---|---|
| `profit-event *` | `/erp/profit-event/*` | `/erp/profit/event/*` | `cmd/profit/profit_event.go` |
| `report stock warehouse/buy-send/finance/produce` | `/warehouse` 等 + 误当 raw JSON | `/stock-warehouse-page` 等 `*-page` 分页端点 | `cmd/report/report_stock.go`(重写) |
| `homepage dashboard2/inventory-backlog/product-ranking` | 全被强制到 dashboard6 | 各自路径 | `internal/cmdutil/stats.go`(删 override)+ `cmd/stats/stats.go`(stock 改自定义 dashboard6 命令) |
| `job-trigger *` | `/erp/finance-job/*` | `/erp/job-trigger/*` | `cmd/job_trigger/job_trigger.go` |

### 新增查询命令(约 36 个)
- **财务新模块**:`finance invoice / payment-apply / prepayment-apply`(各 list·get·summary)→ `cmd/finance/finance_invoice.go`、`finance_payment_apply.go`、`finance_prepayment_apply.go`。refund 已存在未动。
- **库存**:`stock query page-count / detail-count / stock-record-count`(`stock_query.go`)、`stock record page-count / record-page / total-cost`(`stock_record.go`)
- **报表**:`report direct`(新模块 detail·detail-count,`report_direct.go`)、`report stock *-count`(5 个)+ `price-change-logs`(`report_stock.go`)、`report purchase/sale *-count`(共 5 个)
- **单据域 page-count**:`internal/cmdutil/crud_legacy.go` 的 `NewCRUDSubCmd` 加了 `page-count` → stock in/out/move/check、purchase in/return、sale out/return **8 个域**一次补齐

### 文档
- `skills/SKILL.md` 实测表移除已修复 ❌,加「🛠 本轮新增/修复(待有效 token 实测)」小节

### ⚠️ 未完成:实测
冒烟时旧 token 已过期(返回 `code:401 账号未登录`)。**链路已验证通**(flag 解析/校验/请求正常,无 exit 4 / 无 404),但**无有效 token 实测返回数据**。下次先拿新 token 跑一遍本轮命令。

---

## 📋 第 2 轮待办(低频模块 + export)

### A. JSON 查询端点(约 31 个,按模块)

| 模块 | 待对接端点(path) | 备注 |
|---|---|---|
| `invoice-item` | `/erp/invoice-item/list` | 发票明细;可能并入 invoice |
| `plan-prepayment-relation` | `/available-initials`、`/available-payments`、`/page` | 计划预付关联,新模块 |
| `position` | `/erp/position/get-history-waybill-follow`、`/get-latest-waybill-follow` | 位置服务/轨迹,新模块 |
| `screen` | `/erp/screen/line`、`/line-price`、`/purchase-enterprise`、`/purchase-in-price`、`/purchase-in-product`、`/settle-amount` | 大屏补 6 个(`cmd/screen/screen.go`) |
| `system/account` | `/list` + `/getCarDetail`、`/getCarrierDetail`、`/getDriverDetail`、`/getEnterpriseInfo`、`/getLatestWaybillFollow`、`/getWaybill`、`/getWaybillFollow`、`/getWaybillSourceInfo` | ERP 账号绑定,新模块,9 个 |
| `system/weight-account` | `/erp/system/weight-account/get` | 称重账号绑定 |
| `waybill`(注意:非 waybill-source) | `/erp/waybill/get`、`/page`、`/page-count`、`/get-events`、`/get-order-info` | 文档里 `/erp/waybill/*` 与现有 `/erp/waybill-source/*` 是两套,新模块 |
| `service-contract` | `/page`、`/page-count` | ⚠️ 决策见下 |
| `transport-contract` | `/page`、`/page-count` | ⚠️ 决策见下 |
| `finance-record` | `/erp/finance-record/init` | 名字怪,可能非查询,低优先 |

### B. export-excel 端点(约 40 个)— 需先做 `-o` 下载方案

返回二进制 Excel,与现有 JSON stdout 协议冲突。**待决策**:加全局 `--output/-o <文件>` flag,export 命令把 body 写文件(不走 OutputJSON)。
涉及几乎所有模块的 `export-excel` / `export-*-excel`。列表可从 erp.md grep `` **接口地址**:`/admin-api/erp/.*export `` 得到。

### ⚠️ 待决策:service-contract / transport-contract 的 page
现有 `contract_service.go` / `contract_transport.go` 走 `/erp/contract/page?type=SERVICE|TRANSPORT`(共享端点 + 类型区分);文档却列出独立的 `/erp/service-contract/page`、`/erp/transport-contract/page`。需确认:是同一后端的两种路由,还是要新增独立命令。**建议**:先 `wlt api GET /erp/service-contract/page --dry-run` + 真 token 试一条,确认是否 200 再决定。

### 可选:本轮新财务模块的写命令
`finance invoice/payment-apply/prepayment-apply` 本轮只做了查询(list/get/summary)。文档里它们也有 create/update/delete/update-status,按需补(参考 `finance_payment.go` 全套)。

---

## 🔧 实现模式(下次照抄)

- **标准 CRUD**:`cmdutil.AddCRUDToParent(parent, CRUDConfig{Name, APIPath, Label, ListFilters, SingleDelete, SkipStatus})` — 自动生成 list/get/create/update/delete(+update-status)
- **分页 list**:`CrudListCmdWithFixed(apiPath, label, filters, fixed)`(fixed 用于强制 type 等)
- **page-count**:`CrudPageCountCmd` / `CrudPageCountCmdWithFixed`;legacy 域已自带(本轮加的)
- **get**:`CrudGetCmd`(--id 必填)
- **自定义查询**:参考 `cmd/finance/finance_invoice.go`(list/get/summary + FlagSpec filter 切片 + `collectFinanceFilters`)
- **filter flag**:`cmdutil.FlagSpec{Name, Usage}`,kebab-case 名会自动转 camelCase 参数
- **新模块注册**:在 `cmd/<domain>/` 建 `.go`,`func init() { <domain>Cmd.AddCommand(...) }`,并在 `cmd/root.go` 的 `init()` 加 `<domain>.Register(rootCmd)`(若新包)

### 从文档取参数
每个端点的 query 参数在 erp.md 的 `**请求参数**` 表(请求类型=query 行),参数名 camelCase。grep 例:
`grep -nA40 "接口地址**:\`/admin-api/erp/waybill/page\`" D:/Download/erp.md`

---

## 🔍 验证步骤(下次先做)

1. 拿一个有效 token + tenant-id
2. 复测本轮:`wlt finance invoice list --token <t> --tenant-id <n>`、`wlt report stock warehouse ...`、`wlt profit-event list ...`(应返回数据,非 401)
3. 测 filter:`wlt customer list --token <t> --tenant-id <n --status 0` 确认 camelCase 参数生效
4. 第 2 轮每加一个模块:`go build ./... && go test ./...`,再用 `--dry-run` 或真请求验证路径

---

## 📁 本轮改动文件清单(参考)

```
internal/cmdutil/helpers.go        kebab→camel 转换
internal/cmdutil/helpers_test.go   新增单测
internal/cmdutil/crud_legacy.go    NewCRUDSubCmd 加 page-count
internal/cmdutil/stats.go          删 homepage→dashboard6 override
cmd/stats/stats.go                 stock 改自定义 dashboard6 命令
cmd/profit/profit_event.go         路径修复
cmd/job_trigger/job_trigger.go     路径修复
cmd/report/report_stock.go         重写(修路径+加 count)
cmd/report/report_purchase.go      加 detail-count/summer-count
cmd/report/report_sale.go          加 detail/summer/profit-count
cmd/report/report_direct.go        新模块
cmd/stock/stock_query.go           加 page-count/detail-count/stock-record-count
cmd/stock/stock_record.go          加 page-count/record-page/total-cost
cmd/finance/finance_invoice.go         新模块
cmd/finance/finance_payment_apply.go   新模块
cmd/finance/finance_prepayment_apply.go 新模块
skills/SKILL.md                    实测表更新
```

本轮改动已提交 git(commit `b69f448`)。
