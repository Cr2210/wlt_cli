# 设计：finance CRUD 命令迁移到 cmdutil 工厂

## 1. 行为缺口（最小差异）

**现状**：`cmd/finance/` 下 15 个 CRUD 子命令手写实现，与 `internal/cmdutil/crud.go` 的 `CrudGetCmd`/`CrudCreateCmd`/`CrudUpdateCmd`/`CrudDeleteCmd`/`CrudUpdateStatusCmd` 逻辑逐行重复。

**目标**：同样的命令由工厂生成。

**期望行为**：**完全不变**。这是纯等价重构，不产生任何用户可见差异。

## 2. 行为实际所在位置

命令的可观察行为由三处共同决定：

| 位置 | 决定什么 | 本任务是否触碰 |
|------|---------|---------------|
| `cmd/finance/*.go` 的 `newXxxCmd()` | 子命令名、`Short`、flag 声明、`RunE` 逻辑 | **是**（迁移主体） |
| `internal/cmdutil/crud.go:101-244` | 工厂生成的等价实现 | **否**（只读复用） |
| `internal/client` / `internal/output` | HTTP 与信封 | **否** |

关键：工厂把 `label` 参数插进 `Short` 和 `--id` 帮助文本。**label 取值必须与手写版字符串逐字匹配**，否则 CLI 表面就变了。

## 3. 等价的精确判据

对每条待迁移命令，必须同时满足：

1. `Use` 相同（`get` / `create` / `update` / `delete` / `update-status`）
2. `Short` 相同 — 工厂产出 `fmt.Sprintf("获取%s详情", label)` 等，需 `label` 使结果与手写串**完全一致**
3. flag 集合相同（名字 + 类型 + 默认值 + 帮助文本）
4. 必需的 flag 标记相同（`MarkFlagRequired` 的调用点）
5. HTTP method + 路径 + body/query 形态相同
6. 成功路径输出相同（`OutputJSON` + `ParsePagedJSON` 的选择）
7. 失败路径退出码与错误文案相同

任一条不满足 → 该命令保留手写。

## 4. 护栏设计：命令树快照测试

**为什么**：人眼比对 Short 字符串容易漏。测试是硬证据。

**放在哪**：`cmd/finance/finance_snapshot_test.go`（新增）。

**做法**——不改动生产代码，遍历已注册的命令树：

```go
package finance

import (
    "fmt"
    "sort"
    "strings"
    "testing"

    "github.com/spf13/cobra"
)

// snapshotCmd 把一个子命令渲染成 "name|short|flag:help;flag:help..." 的稳定字符串。
func snapshotCmd(c *cobra.Command) string {
    var flags []string
    c.Flags().VisitAll(func(f *pflag.Flag) {
        flags = append(flags, fmt.Sprintf("%s:%s:%s", f.Name, f.Value.Type(), f.Usage))
    })
    sort.Strings(flags)
    return fmt.Sprintf("%s|%s|%s", c.Name(), c.Short, strings.Join(flags, ";"))
}
```

对每个 finance 子组（`payment-apply`、`invoice`、`prepayment-apply`、`refund`、`write-off`、`receipt-payment`、`invoice-apply`、`account`、`settlement`、`account-settlement`、`freight-apply`）断言其子命令快照。

**基线怎么来**：**迁移前先跑一次，把实际输出写进期望值**（不是手抄）。这样期望值 = 迁移前真实行为，迁移后必须仍相等。

**覆盖范围**：本任务只对 A 类 15 个命令加严格断言；对其余命令加"存在性"断言（子命令在即可），避免把范围外的内容也锁死。

## 5. 修改的文件

| 文件 | 改动 | 为什么必须 |
|------|------|-----------|
| `cmd/finance/finance_snapshot_test.go` | 新增 | 护栏，锁住迁移前行为 |
| `cmd/finance/finance_payment_apply.go` | 4 个命令换工厂 | A 类 |
| `cmd/finance/finance_invoice.go` | 4 个命令换工厂 | A 类 |
| `cmd/finance/finance_prepayment_apply.go` | 4 个命令换工厂 | A 类 |
| `cmd/finance/finance_refund.go` | 3 个命令换工厂 | A 类 |
| `.trellis/spec/cli/command-factories.md` | 更新调用点计数 | spec 与代码同步要求（`cli/index.md:52`） |

其余 finance 文件的 A 类命令（`finance_account.go` 等）在实现时按同一判据核对后决定，**核对结果不是等价就保留**。

## 6. 如何证明没有改变行为

三重验证，缺一不可：

1. **快照测试**：迁移前后同一测试全绿（这是主证据）。
2. **`--help` 实测**：迁移前后分别执行各子组 `--help`，diff 输出。
3. **`go test ./...`**：现有 `finance_receipt_payment_test.go`、`freight_apply_test.go` 等不回归。

## 7. 明确不做

- 不扩展 `DocumentConfig`（`Fixed`/`Summary`/`Export`/`GetBy`/`StatusMode`）。
- 不把 `page` 改名 `list`——`cli/index.md:50-51` 明确禁止当全仓规则。
- 不合并重复的 `[]FlagSpec`（`financeRefundFilters` vs `financeRefundPageCountFilters`）。
- 不改 `internal/` 任何文件。
- 不碰 `context.Background()`、`EnsureClient` 样板收敛（后续任务）。
- 不碰 `summary` / `export` / `page` / `page-count` 的 `RunE` 逻辑（除非实现时证明可换 `CrudPageCountCmdWithFixed` 且逐字等价）。
