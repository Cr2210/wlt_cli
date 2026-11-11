# export 全景与 `-o` 下载方案设计输入（2026-10-19）

## 后端 export 端点（65 个，ge-dev-1.8 实测 grep `@GetMapping.*export`）

| 模块 | 数量 | 端点 |
|---|---|---|
| finance | 13 | 各子模块 `/export-excel` |
| stock | 7 | `/export-excel` |
| contract | 4 | contract / service-contract / transport-contract / provision-contract 各 1 |
| report | 15 | `export-detail-excel`×3、`export-summer-excel`×2、`export-excel`×1、`export-stock-warehouse-excel`、`export-profit-excel`、`export-produce-profit-excel`、`export-produce-detail-excel`、`export-finance-record`、`export-car-profit-excel`、`export-buy-send-excel` 等 |
| product | 3 | `/export-excel` |
| businesspartner | 3 | `/export-excel` |
| screen | 2 | `/export`（非 -excel 后缀） |
| sale / purchase / produce / order | 各 2 | `/export-excel` |
| weight / waybill / settlement / serviceprovider / saleaccounting | 各 1 | `/export-excel` |

（prd 估的「约 40」偏保守，实际 65；但其中相当部分 CLI 已有对应 export 命令，见下。）

## CLI export 命令现状

**临时形态**（`fmt.Println("导出成功，返回数据：", string(resp.Data))`，非 JSON 信封）：**23 处**，
分布 21 个文件（prd 记录时 21 处，0923 同步新增 sale/accounting、finance/freight_apply 等）：

```
cmd/finance/{finance_account, finance_account_settlement, finance_invoice, finance_invoice_apply,
  finance_receipt_payment, finance_refund, finance_settlement, finance_write_off, freight_apply}.go
cmd/invoice/invoice_main.go
cmd/order/{order_main, order_plan}.go
cmd/produce/{produce_main, produce_plan}.go
cmd/quality/{quality_inspection, quality_weight}.go
cmd/sale/accounting.go
cmd/settlement/{logistics, settlement_main}.go
cmd/system/sys_user.go
cmd/weight/weight_waybill.go
```

**已知较新形态**：`cmd/settlement/logistics.go` 与 `cmd/sale/accounting.go` 的 export 带 `--headers` flag
（commit `398d6ee`）；实现方式仍需在 design 阶段统一核对（哪些直接 `fmt.Println`、哪些有别的处理）。

**关键事实**：后端 export 返回**二进制 Excel**（`HttpServletResponse` 直写流），CLI 侧 `resp.Data` 拿到的是
`json.RawMessage`——若后端真回二进制，现有 23 处「打印字符串」输出本身就不完整/不可用，
这是 `-o` 方案要一并解决的**功能缺陷**而非单纯格式问题（design 阶段需抓一个真实响应确认 Data 的实际内容形态）。

## `-o` 方案设计要点（承接 prd B 节）

1. 全局 flag `--output/-o <文件路径>`（挂 root PersistentFlags；export 命令检测到时写文件）
2. export 命令**不走** `OutputJSON` 信封；无 `-o` 时的行为需定：报错提示用 `-o`，还是打印到 stdout（默认文件名？）
3. 统一改造 23 处 `fmt.Println` 临时形态 → 公共 helper（建议放 `internal/cmdutil/output.go` 或新文件）
4. 写文件成功后 stdout 输出一行**信封 JSON**（如 `{"data":{"file":"xx.xlsx","bytes":12345}}`）以维持
   「stdout 机器可读」协议，还是纯文本行——design 阶段定
5. 涉及 65 后端端点中 CLI **尚无** export 命令的部分（如 report 的新导出、screen 2 个），
   本轮是否补齐需划范围（建议：先统一已有 23 处，缺的按模块逐步补）
