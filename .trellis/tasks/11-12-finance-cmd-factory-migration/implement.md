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

## 步骤 7 — 收尾验证

```bash
gofmt -l cmd/finance/          # 期望无输出
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
