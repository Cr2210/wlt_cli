# 第 2 轮查询端点对接与 export 下载方案（开放待办）

> **状态**：`planning`（尚未实施）。
> 来源：`docs/superpowers/specs/2026-06-27-query-endpoints-backlog.md`（superpowers 交接文档，2026-06-27）。
> 原始文档已从仓库删除，本文档为其内容的 Trellis 化保留（第 2 轮部分逐字保留）。

## Goal

按 `D:\Download\erp.md`（Swagger 导出，510 端点：GET 303 / POST 65 / PUT 92 / DELETE 50）
对接第 2 轮数据接口：约 31 个 JSON 查询端点 + 约 40 个 export-excel 端点（后者阻塞于 `-o` 下载方案）。

**用户已拍板的前置决策**：export-excel 跳过（待 `-o` 下载方案）；路径 bug 顺手修；先做高频模块。

## Background — 第 1 轮已完成（baseline，2026-06-27）

### 基础修复

- **filter 参数名 kebab→camelCase**：`internal/cmdutil/helpers.go` 的 `CollectStringFlag/Flags/IntFlags`
  现把 `--supplier-id` 转成 `supplierId`。**此前所有多词 filter 都不生效**（后端收不到）。
  已加单测 `helpers_test.go`。
  - 注意：单字段 filter 的「字段名」已修，但**部分老命令的 flag 名本身与文档不符**
    （如 `finance refund` 用 `refund-no` 而文档是 `no`、用 `customer-id/supplier-id` 而文档是 `partnerId`）。
    本轮未逐个核对老命令 flag，下次可按需修。

### 路径 bug 修复（4 处，对齐文档）

| 命令 | 旧（错） | 新（文档） | 文件 |
|---|---|---|---|
| `profit-event *` | `/erp/profit-event/*` | `/erp/profit/event/*` | `cmd/profit/profit_event.go` |
| `report stock warehouse/buy-send/finance/produce` | `/warehouse` 等 + 误当 raw JSON | `/stock-warehouse-page` 等 `*-page` 分页端点 | `cmd/report/report_stock.go`（重写） |
| `homepage dashboard2/inventory-backlog/product-ranking` | 全被强制到 dashboard6 | 各自路径 | `internal/cmdutil/stats.go`（删 override）+ `cmd/stats/stats.go`（stock 改自定义 dashboard6 命令） |
| `job-trigger *` | `/erp/finance-job/*` | `/erp/job-trigger/*` | `cmd/job_trigger/job_trigger.go` |

### 新增查询命令（约 36 个）

- **财务新模块**：`finance invoice / payment-apply / prepayment-apply`（各 list·get·summary）
  → `cmd/finance/finance_invoice.go`、`finance_payment_apply.go`、`finance_prepayment_apply.go`。refund 已存在未动。
- **库存**：`stock query page-count / detail-count / stock-record-count`（`stock_query.go`）、
  `stock record page-count / record-page / total-cost`（`stock_record.go`）
- **报表**：`report direct`（新模块 detail·detail-count，`report_direct.go`）、
  `report stock *-count`（5 个）+ `price-change-logs`（`report_stock.go`）、`report purchase/sale *-count`（共 5 个）
- **单据域 page-count**：`internal/cmdutil/crud_legacy.go` 的 `NewCRUDSubCmd` 加了 `page-count`
  → stock in/out/move/check、purchase in/return、sale out/return **8 个域**一次补齐。
  （注：`crud_legacy.go` 已在 2026-08-07 库存命令优化中删除，上句为当时事实。）
- **文档**：`skills/SKILL.md` 实测表移除已修复 ❌，加「🛠 本轮新增/修复（待有效 token 实测）」小节。
- 本轮改动已提交 git（commit `b69f448`）。

### ⚠️ 未完成：实测（本轮遗留）

冒烟时旧 token 已过期（返回 `code:401 账号未登录`）。**链路已验证通**
（flag 解析/校验/请求正常，无 exit 4 / 无 404），但**无有效 token 实测返回数据**。
下次先拿新 token 跑一遍本轮命令。

---

## Requirements

### A. JSON 查询端点（约 31 个，按模块）

| 模块 | 待对接端点（path） | 备注 |
|---|---|---|
| `invoice-item` | `/erp/invoice-item/list` | 发票明细；可能并入 invoice |
| `plan-prepayment-relation` | `/available-initials`、`/available-payments`、`/page` | **已接入** `wlt order prepayment-relation`（2026-08-21，另含 create/update-amount）→ 视为已完成 |
| `position` | `/erp/position/get-history-waybill-follow`、`/get-latest-waybill-follow` | 位置服务/轨迹，新模块 |
| `screen` | `/erp/screen/line`、`/line-price`、`/purchase-enterprise`、`/purchase-in-price`、`/purchase-in-product`、`/settle-amount` | 大屏补 6 个（`cmd/screen/screen.go`） |
| `system/account` | `/list` + `/getCarDetail`、`/getCarrierDetail`、`/getDriverDetail`、`/getEnterpriseInfo`、`/getLatestWaybillFollow`、`/getWaybill`、`/getWaybillFollow`、`/getWaybillSourceInfo` | ERP 账号绑定，新模块，9 个 |
| `system/weight-account` | `/erp/system/weight-account/get` | 称重账号绑定 |
| `waybill`（注意：**非** waybill-source） | `/erp/waybill/get`、`/page`、`/page-count`、`/get-events`、`/get-order-info` | 文档里 `/erp/waybill/*` 与现有 `/erp/waybill-source/*` 是两套，新模块 |
| `service-contract` | `/page`、`/page-count` | ⚠️ 待决策（见下） |
| `transport-contract` | `/page`、`/page-count` | ⚠️ 待决策（见下） |
| `finance-record` | `/erp/finance-record/init` | 名字怪，可能非查询，低优先 |

> 注：`screen` / `waybill` 的实际端点集可能在 0923 同步后已有变化——`2026-09-28-backend-sync-0923-design.md`
> 未把它们列入 A 类，说明不在 0923 主线上；本清单保留原文并标注「0923 后端基线未覆盖」。

### B. export-excel 端点（约 40 个）— 需先做 `-o` 下载方案

返回**二进制 Excel**，与现有 JSON stdout 协议冲突。**待决策**：加全局 `--output/-o <文件>` flag，
export 命令把 body 写文件（**不走** `OutputJSON`）。
涉及几乎所有模块的 `export-excel` / `export-*-excel`。
列表可从 erp.md grep `` **接口地址**:`/admin-api/erp/.*export `` 得到。

现状（已知例外）：仓库约 21 处 export 命令用 `fmt.Println("导出成功，返回数据：", string(resp.Data))`
输出文本行，不符合 `{data:...}` 信封——这是待 `-o` 方案统一的临时形态。

涉及设计诉求：
- 全局 `--output/-o <文件>` flag；
- export 命令把 body 写文件，不走 `OutputJSON`；
- 统一现有 21 处 `fmt.Println` 临时形态。

### C. 待决策 1 — `service-contract` / `transport-contract` 的 page 路径归属

现有 `contract_service.go` / `contract_transport.go` 走 `/erp/contract/page?type=SERVICE|TRANSPORT`
（共享端点 + 类型区分）；Swagger 却列出独立的 `/erp/service-contract/page`、`/erp/transport-contract/page`。
需确认：是同一后端的两种路由，还是要新增独立命令。

**原文建议**：先 `wlt api GET /erp/service-contract/page --dry-run` + 真 token 试一条，确认是否 200 再决定。

代码侧口子已留好：`ContractTypeConfig.APIPath` 支持覆盖默认 `/erp/contract`
（`internal/cmdutil/crud.go`，销售合同已走独立 `/erp/service-contract`）。

### D. 待决策 2 — 本轮新财务模块的写命令是否补全

`finance invoice` / `payment-apply` / `prepayment-apply` 本轮只做了查询（list/get/summary）。
Swagger 里它们也有 create/update/delete/update-status；参考实现 `finance_payment.go` 全套。
现状：`cmd/finance/finance_invoice.go`、`finance_payment_apply.go`、`finance_prepayment_apply.go` 存在；
**尚未逐命令核对其子命令集**——转入实施时需先核对。

---

## Acceptance Criteria（实施本任务时用）

- [ ] 第 2 轮 A 类端点按模块逐条对接，每条 `--help` / `--dry-run` 路径核对正确
- [ ] `-o` 全局下载 flag 设计并落地，export 类命令统一走文件输出
- [ ] 待决策 1、2 各有明确结论并记录
- [ ] 第 1 轮遗留的实测项用有效 token 复测通过（`finance invoice list`、`report stock warehouse`、`profit-event list`、filter camelCase）
- [ ] 每加一个模块 `go build ./... && go test ./...` 通过

## Notes

- 本任务为开放待办，当前不实施。原文完整交接文档内容已由本 prd 承载。
- 实现模式参考（原文）：`cmdutil.AddCRUDToParent`、`CrudListCmdWithFixed`、`CrudPageCountCmd(WithFixed)`、
  `CrudGetCmd`、自定义查询参考 `cmd/finance/finance_invoice.go`、filter 用 `cmdutil.FlagSpec{Name, Usage}`。
- 从 Swagger 取参数：每个端点的 query 参数在 erp.md 的 `**请求参数**` 表（请求类型 = query 行），参数名 camelCase。
  grep 例：`grep -nA40 "接口地址**:\`/admin-api/erp/waybill/page\`" D:/Download/erp.md`（`backlog:84`）。
- 后端仓库（只读参考）：`D:\Projects\weiliantong\weiliantong-api`，分支 `ge-dev-1.8`。