# 实施计划：finance CRUD 命令迁移到 cmdutil 工厂

## 前置

```bash
git status                 # 确认工作区干净
go build ./... && go vet ./... && go test ./...   # 基线全绿
```

## 步骤 1 — 建护栏（先做，不可跳过）

1. 新建 `cmd/finance/finance_snapshot_test.go`，实现 `snapshotCmd` 辅助函数。
2. 先写一个**打印模式**的临时测试，遍历 `financeCmd` 全部子命令，把快照输出打到 stdout。
3. 运行它，把真实输出**粘贴**为期望值（不手抄、不臆造）。
4. 改成正式断言测试，`go test ./cmd/finance/` 全绿。
5. 同时记录基线：`wc -l cmd/finance/*.go` 总行数。

**验证**：`go test ./cmd/finance/ -run Snapshot -v` 通过。

## 步骤 2 — 迁移 payment-apply（A 类 4 条）

文件 `cmd/finance/finance_payment_apply.go`：

- 删除 `newFinancePaymentApplyGetCmd`（`:81`）、`newFinancePaymentApplyCreateCmd`（`:141`）、`newFinancePaymentApplyUpdateCmd`（`:168`）、`newFinancePaymentApplyDeleteCmd`（`:222`）
- 在 `init()` 里改为：

```go
cmdutil.CrudGetCmd("/erp/finance-payment-apply", "付款申请"),
cmdutil.CrudCreateCmd("/erp/finance-payment-apply", "付款申请"),
cmdutil.CrudUpdateCmd("/erp/finance-payment-apply", "付款申请"),
cmdutil.CrudDeleteCmd("/erp/finance-payment-apply", "付款申请", false),
```

- **保留** `newFinancePaymentApplyUpdateStatusCmd`（`:195`，`PutQuery` + `approveStatus`，工厂不支持）
- 清理因删除而失效的 import（若 `encoding/json` / `output` 仍被其它函数使用则保留）

**验证**：`go build ./... && go test ./cmd/finance/ -run Snapshot`

## 步骤 3 — 迁移 invoice（A 类 4 条）

文件 `cmd/finance/finance_invoice.go`：删除 `:109` get、`:172` create、`:199` update、`:253` delete，换成工厂；**保留** `:226` update-status（`PutQuery` + `status`）。

**注意**：`addFinanceFilterFlags` / `collectFinanceFilters`（`:275`、`:282`）被其它命令使用，**不要删**。

## 步骤 4 — 迁移 prepayment-apply（A 类 4 条）

文件 `cmd/finance/finance_prepayment_apply.go`：删除 `:72` get、`:132` create、`:159` update、`:213` delete；**保留** `:186` update-status。

## 步骤 5 — 迁移 refund（A 类 3 条）

文件 `cmd/finance/finance_refund.go`：删除 `:172` create、`:200` update、`:252` update-status，换成工厂。

**保留** `:137` get（支持 `--id` 或 `--no`）与 `:228` delete（单删 `--id`；实现时先核对工厂 `single=true` 的 `Short`/flag 帮助文本是否逐字一致，不一致就保留手写）。

## 步骤 6 — 核对剩余文件

逐函数核对 `finance_account.go`、`finance_settlement.go`、`finance_account_settlement.go`、`finance_write_off.go`、`finance_receipt_payment.go`、`finance_invoice_apply.go`、`freight_apply.go`。

对每条命令按 `design.md` §3 的 7 条判据检查；**满足才换，不满足保留并在 commit message 里说明**。不要为了凑数强行替换。

### 实际执行结果（已回填）

- **`finance_account.go` 试替换后回退**：`delete` 换成 `CrudDeleteCmd(…, true)` 后，快照测试报 `--id` 帮助文本由 `"账户 ID"` → `"结算账户 ID"`。根因是工厂把 `label` 同时插进 `Short` 与 flag 帮助文本，而该文件 `Short` 用「结算账户」、flag 用「账户 ID」，二者不可兼得 → 整组回退。
- **label 一致性成为决定性判据**。据此核对：`settlement`、`write-off`、`receipt-payment` 三者 label 一致，完成迁移；`invoice-apply`（开票申请 vs 申请单 ID）、`account-settlement`（账户结算 vs 结算 ID）不一致，保留手写。
- **`freight_apply.go` 无需改动**：`init()` 已是工厂写法。
- 最终迁移 **26 条命令**，`cmd/finance/` **3167 → 2701 行（-466）**。

### check 阶段的修正（回退 4 条 update-status）

首版把 `refund` / `settlement` / `write-off` / `receipt-payment` 的 `update-status` 也迁到
`CrudUpdateStatusCmd`。`trellis-check` 用真实二进制对比发现：手写版 `--data` 解析失败时
`hint` 为 `"data 应为 JSON 对象"`，工厂（`internal/cmdutil/crud.go:232`）该参数是 `""`，
**迁移后 hint 字段从 stderr JSON 中消失**——`--help` diff 与命令树快照都覆盖不到。

已回退这 4 条（恢复手写实现），并补两条护栏：

- `TestMigratedRequiredFlags`：断言迁移命令的 `MarkFlagRequired`（cobra annotation，同样不进 `--help`）。
- `TestRetainedCommandsPresent`：断言保留手写的命令仍然存在，防后续被误删/误改形状。

> 注：`CrudCreateCmd`（`crud.go:135`）与 `CrudUpdateCmd`（`:162`）的 hint 是完好的，
> 只有 `CrudUpdateStatusCmd`（`:232`）缺 hint。这也是不能为 finance 单独改工厂的原因——
> 该工厂另有 4 处调用方（`order_main`、`partner`、`product_category`、`DocumentCmds`），改它会影响它们。

### 等价性证据（回退后重跑）

- `cmd/finance/finance_snapshot_test.go` 的 `TestMigratedCommandSurfaces` 锁定 26 条命令的子命令名 / `Short` / flag 集合（名字、类型、帮助文本）。
- 重构前后命令树全量快照 `diff` **逐字节一致**。
- **11 个 finance 子组** `wlt finance <group> --help` 重构前后 `diff` **逐字节一致**（含未改动的 account / invoice-apply / account-settlement / freight-apply）。
- 4 条回退的 `update-status` 实测 stderr：`"hint":"data 应为 JSON 对象"` 已恢复。
- `go build ./... && go vet ./... && go test ./... -count=1` 全绿。
- `gofmt -l` 对本包全部文件（含未改动的）都报 CRLF——`core.autocrlf=true` 的 Windows 既有现象，非本次引入；本次未对其做格式化以免产生整文件伪 diff。

## 步骤 7 — 收尾验证

```bash
gofmt -l cmd/finance/          # 本包因 CRLF 全部报出，属既有现象，见上
go build ./... && go vet ./... && go test ./...
```

`--help` 对比（迁移前已保存基线输出的话直接 diff）：

```bash
for g in payment-apply invoice prepayment-apply refund; do
  go run . finance $g --help > /tmp/after-$g.txt
done
```

行数对比：`wc -l cmd/finance/*.go` 应比步骤 1 记录的基线小。

## 步骤 8 — 同步 spec

更新 `.trellis/spec/cli/command-factories.md` §5（判据 4）的调用点说明：把 finance 下新采用 `Crud*Cmd` 单件的文件列入。若 `crud.go` 行号未变则无需改 §1 表格。

## 回滚

纯等价重构，若快照测试失败：

```bash
git checkout -- cmd/finance/
```
然后定位是哪条不等价，把它移出 A 类，重新核对。

## 风险

| 风险 | 缓解 |
|------|------|
| label 字符串不匹配导致 `Short` 变化 | 快照测试直接断言 `Short`，会在第一时间失败 |
| 误删仍被引用的私有函数 | 步骤 2-5 每条都明确列出删除与保留清单；`go build` 会报未定义 |
| `finance_refund.go` delete 的 `Short` 不一致 | 步骤 5 明确要求先核对，不一致则保留 |
| 把 `page` 误改成 `list` | PRD/design 明确禁止；快照测试断言 `Name()` 会立即捕获 |
