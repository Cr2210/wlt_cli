<!-- 本文件是 `conventions-extraction` 的拆分件之一（按研究产物约定：一主题一文件）。
     初衷：单文件 57KB 超过 context_injection.max_file_bytes(32768) 会被截断。
     全量导航见同目录 `00-index.md`。文中形如「§C-1」「§4.5」的交叉引用按该导航定位。 -->

# D. backlog 开放任务事实（源：`docs/superpowers/specs/2026-06-27-query-endpoints-backlog.md`）

## D.1 第 2 轮待对接 JSON 查询端点（约 31 个，`backlog:44-57`）

| 模块 | 端点 path | 备注（原文） |
|---|---|---|
| `invoice-item` | `/erp/invoice-item/list` | 发票明细；可能并入 invoice |
| `plan-prepayment-relation` | `/available-initials`、`/available-payments`、`/page` | **已接入**（2026-08-21，`wlt order prepayment-relation`，另含 create/update-amount）→ 视为已完成 |
| `position` | `/erp/position/get-history-waybill-follow`、`/get-latest-waybill-follow` | 位置服务/轨迹，新模块 |
| `screen` | `/erp/screen/line`、`/line-price`、`/purchase-enterprise`、`/purchase-in-price`、`/purchase-in-product`、`/settle-amount` | 大屏补 6 个（`cmd/screen/screen.go`） |
| `system/account` | `/list` + `/getCarDetail`、`/getCarrierDetail`、`/getDriverDetail`、`/getEnterpriseInfo`、`/getLatestWaybillFollow`、`/getWaybill`、`/getWaybillFollow`、`/getWaybillSourceInfo` | ERP 账号绑定，新模块，9 个 |
| `system/weight-account` | `/erp/system/weight-account/get` | 称重账号绑定 |
| `waybill`（非 waybill-source） | `/erp/waybill/get`、`/page`、`/page-count`、`/get-events`、`/get-order-info` | 与现有 `/erp/waybill-source/*` 是两套，新模块 |
| `service-contract` | `/page`、`/page-count` | 决策见 D.3 |
| `transport-contract` | `/page`、`/page-count` | 决策见 D.3 |
| `finance-record` | `/erp/finance-record/init` | 名字怪，可能非查询，低优先 |

> 注：`screen` / `waybill` 的实际端点集可能在 0923 同步后已有变化——`sync-design` 未把它们列入 A 类，说明不在 0923 主线上；转任务时应保留原文并标注「0923 后端基线未覆盖」。

## D.2 export-excel 端点（约 40 个）与 `-o` 下载方案（`backlog:59-62`）

- 现状：export 端点返回**二进制 Excel**，与 JSON stdout 协议冲突。
- **设计诉求（待决策）**：加全局 `--output/-o <文件>` flag；export 命令把 body 写文件，**不走** `OutputJSON`。
- 涉及几乎所有模块的 `export-excel` / `export-*-excel`；列表可用 `grep '**接口地址**:`/admin-api/erp/.*export' erp.md` 得到（`backlog:62`）。
- 与 §C-9 呼应：当前 21 处 `fmt.Println` 是临时形态。

## D.3 两处待决策（`backlog:64-68`）

**决策 1 — `service-contract` / `transport-contract` 的 page 路径归属**
- 现状：`contract_service.go` / `contract_transport.go` 走 `/erp/contract/page?type=SERVICE|TRANSPORT`（共享端点 + 类型区分）。
- Swagger 却列出独立 `/erp/service-contract/page`、`/erp/transport-contract/page`。
- 需确认：同一后端两种路由，还是需新增独立命令。
- 原文建议：先 `wlt api GET /erp/service-contract/page --dry-run` + 真 token 试一条，确认是否 200 再决定。
- 代码侧锚点：`ContractTypeConfig.APIPath` 已支持覆盖（`internal/cmdutil/crud.go:370-383`，销售合同已走独立 `/erp/service-contract`）——为「新增独立命令」留好口子。

**决策 2 — 新财务模块的写命令是否补全**
- `finance invoice` / `payment-apply` / `prepayment-apply` 当时只做了查询（list/get/summary）。
- Swagger 里它们也有 create/update/delete/update-status；参考实现 `finance_payment.go` 全套。
- 现状锚点：`cmd/finance/finance_invoice.go`、`finance_payment_apply.go`、`finance_prepayment_apply.go` 存在；**本文未逐命令核对其子命令集（UNTESTED）**，转任务时再核。

## D.4 Round 1 未完成的实测项（`backlog:36-38,88-93`）

- Round 1 新增/修复约 36 个命令，冒烟时 token 已过期（返回 `code:401 账号未登录`）。
- **链路已验证通**（flag 解析/校验/请求正常，无 exit 4 / 无 404），但**无有效 token 实测返回数据**。
- 原文验证步骤（保留）：① 拿有效 token + tenant-id；② 复测 `wlt finance invoice list`、`wlt report stock warehouse ...`、`wlt profit-event list`；③ 测 filter camelCase（`wlt customer list --status 0`）；④ 每加一个模块 `go build ./... && go test ./...`。
- 文件清单与 commit：`backlog:99-117`（含 `cmd/profit/profit_event.go`、`cmd/job_trigger/job_trigger.go`、`cmd/report/report_stock.go`、`cmd/stats/stats.go`、`cmd/stock/stock_query.go`、`stock_record.go`、`cmd/finance/finance_invoice.go` 等）；commit `b69f448`（`backlog:119`）。
- **另一处未核对遗留（`backlog:18`）**：部分老命令 flag 名与文档不符（如 `finance refund` 用 `refund-no` vs 文档 `no`；`customer-id`/`supplier-id` vs 文档 `partnerId`），「本轮未逐个核对，下次可按需修」。

## D.5 转开放任务的建议要点（据 prd R5）

新任务 `prd.md` 需保留：D.1 完整约 31 个 JSON 端点清单 + D.2 约 40 个 export-excel + `-o` 方案诉求 + D.3 两处待决策 + D.4 的 Round 1 实现记录与未实测项作为背景；状态保持 `planning`，本任务不实施。

---

## 附录：锚点速查（写 spec 时按需引用）

- 鉴权链路：`cmd/root.go:45-61`（暂存）→ `internal/cmdutil/client.go:39-66`（校验/组装）→ `internal/client/client.go:192-201`（成头）
- 退出码：定义 `internal/apierr/apierr.go:43-68`；消费 `cmd/root.go:67-85`；构造 `internal/output/output.go:26-44`
- 输出信封：`internal/output/output.go:60-98`；命令层封装 `internal/cmdutil/output.go:12-43`
- 工厂：`internal/cmdutil/crud.go`（全表 §3.1）；时间 `internal/cmdutil/timeflags.go:8-21`；传参 `internal/cmdutil/helpers.go:28-85`
- 三种工厂真实调用点见 §3.2；`DocumentCmds` 配置键 `internal/cmdutil/crud.go:856-864`
- 测试样本：`internal/cmdutil/crud_test.go`（工厂结构）、`cmd/order/order_prepayment_test.go`（flag 装配）、`internal/client/client_test.go`（httptest）
- 构建验证：`Makefile:7-11`；`go build ./... && go vet ./... && go test ./...`
