# 两处待决策的事实核查（2026-10-19，后端 ge-dev-1.8）

## 待决策 1：service-contract / transport-contract 路径归属 —— **已被后续重构解决**

prd 记录的前提（2026-06-27）是：`contract_service.go` / `contract_transport.go` 走共享端点
`/erp/contract/page?type=SERVICE|TRANSPORT`，与 Swagger 独立端点冲突，需真机验证。

**现状（两仓库核实）**：

- 后端确有**独立 Controller**：
  - `$C/contract/ErpServiceContractController.java`（`@RequestMapping("/erp/service-contract")`，8 方法全套：create/update/get/page/page-count/delete/update-status/export-excel）
  - `$C/contract/ErpTransportContractController.java`（`@RequestMapping("/erp/transport-contract")`，同样 8 方法）
  - 共享端点 `$C/contract/ErpContractController.java`（`/erp/contract`）仍存在，服务长协等类型
- **CLI 已切独立端点**：`cmd/contract/contract_main.go:24-25`（SALE_CONTRACT → `APIPath:"/erp/service-contract"`）、
  `:39-40`（TRANSPORT → `APIPath:"/erp/transport-contract"`）、`:58`（SERVICE → `/erp/provision-contract`）。
  prd 提到的 `contract_service.go` / `contract_transport.go` 文件已不存在（合同域已重构为 contract.go + contract_main.go）。

**建议结论**：待决策 1 关闭，无需开发；残余动作 = 有效 token 下实测
`wlt contract <销售/运输合同> page` 各跑一条确认 200。

## 待决策 2：财务新模块写命令是否补全 —— 事实表

三个模块后端方法集**完全一致**（8 方法）：create / update / update-status / delete / get / page / summary / export-excel。
锚点：`$C/finance/FinanceInvoiceController.java:33-78`、`ErpFinancePaymentApplyController.java:37-90`、
`ErpFinancePrepaymentApplyController.java:37-90`。

CLI 现状（`cmd/finance/`）：

| 模块 | CLI 子命令 | 调用端点 | 缺口（相对后端 8 方法） |
|---|---|---|---|
| `finance_invoice.go` | page / page-count / get / summary / export | `/erp/finance-invoice/*` | **create / update / update-status / delete** 四个写命令；且 page-count 为漂移端点（后端已无，见 01 号文件） |
| `finance_payment_apply.go` | page / page-count / get / summary | `/erp/finance-payment-apply/*` | **create / update / update-status / delete / export** 五个 |
| `finance_prepayment_apply.go` | list（实际调 `/page`，`finance_prepayment_apply.go:53`）/ get / summary | `/erp/finance-prepayment-apply/*` | **page-count / create / update / update-status / delete / export** 六个；子命令命名 `list` 与另两模块 `page` 不一致（R6b 类不一致又一例） |

**参考实现**：`cmd/finance/finance_receipt_payment.go`（收付款管理，:77-303 全套 8 子命令：
page/page-count/get/create/update/delete/update-status/summary/export），可作模板。

**附带发现**：后端 `FinanceInvoiceController` 无 `page-count`，但 payment-apply / prepayment-apply 两个 Controller
同样没有 `page-count` 方法——而 CLI 的 `finance payment-apply page-count` 已存在（`finance_payment_apply.go:100`）。
即三个财务模块的 `page-count` **全部疑似漂移**（01 号文件只标了 invoice 一处，实施前需对这三个
`page-count` 端点逐一真机确认，404 则删）。

**建议方案（供拍板）**：
- 方案 A（推荐）：本轮连同 A 类一起补全——三模块统一为 page / get / summary / export + 四写命令，
  删除漂移的 page-count；工作量约 3 文件 × 模式化复制，参考 finance_payment.go
- 方案 B：本轮只做查询补齐（export、page-count 清理），写命令另开任务
