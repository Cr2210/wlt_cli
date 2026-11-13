# finance 包 CRUD 子命令零风险迁移到 cmdutil 工厂

## Goal

把 `cmd/finance/` 下 11 个业务文件中**与 `internal/cmdutil` 现成工厂逐字等价**的 CRUD 子命令替换为工厂调用，消除重复，**不改变任何 CLI 对外行为**。

## 背景与已核实事实

`internal/cmdutil/crud.go` 已提供以下可直接复用的工厂（行号为当前代码）：

| 工厂 | 行 | 生成的 `RunE` 行为 |
|------|----|-------------------|
| `CrudGetCmd(apiPath, label)` | `:101` | `GET {apiPath}/get?id=`，`--id` 必填 int64，`OutputJSON`，失败 `NewExitError(5, "获取{label}失败: …")` |
| `CrudCreateCmd(apiPath, label)` | `:124` | `POST {apiPath}/create`，`--data` 必填，解析失败 exit 4 |
| `CrudUpdateCmd(apiPath, label)` | `:151` | `PUT {apiPath}/update`，`--data` 必填，解析失败 exit 4 |
| `CrudDeleteCmd(apiPath, label, single)` | `:178` | `single=true` → `?id=`；`false` → `?ids=` |
| `CrudUpdateStatusCmd(apiPath, label)` | `:221` | `PUT {apiPath}/update-status`，**body `--data`** |

**关键约束（来自 spec，不可违反）：**
- `cli/index.md:45-46` / `command-factories.md:70`：**无隐式 base flags**，筛选字段必须显式配置。
- `cli/index.md:50-51`：**「分页子命令统一叫 `list`」不是全仓规则**；finance 用 `page` 属现状，本任务**不得**改名。
- `command-factories.md:107-116`（判据 4）：形状各异的端点**域内手写 + 复用 `Crud*Cmd` 单件**，正是本任务采用的方式。
- `command-factories.md:145-156`（§8）：禁止引入 legacy 子命令构造器。

## 可替换 / 不可替换的判定（基于逐文件核对）

### A. 可直接替换（行为逐字等价）

| 文件 | 命令 | 替换为 | 依据 |
|------|------|--------|------|
| `finance_payment_apply.go:81` | `get` | `CrudGetCmd("/erp/finance-payment-apply", "付款申请")` | id 必填、OutputJSON、文案 `获取付款申请失败` — 全等 |
| `finance_payment_apply.go:141` | `create` | `CrudCreateCmd(…)` | 全等 |
| `finance_payment_apply.go:168` | `update` | `CrudUpdateCmd(…)` | 全等 |
| `finance_payment_apply.go:222` | `delete` | `CrudDeleteCmd(…, false)` | `--ids` 批量 — 全等 |
| `finance_invoice.go:253` | `delete` | `CrudDeleteCmd(…, false)` | `--ids` 批量 — 全等 |
| `finance_invoice.go:109` | `get` | `CrudGetCmd(…)` | 全等 |
| `finance_invoice.go:172` | `create` | `CrudCreateCmd(…)` | 全等 |
| `finance_invoice.go:199` | `update` | `CrudUpdateCmd(…)` | 全等 |
| `finance_prepayment_apply.go:72` | `get` | `CrudGetCmd(…)` | 全等 |
| `finance_prepayment_apply.go:132` | `create` | `CrudCreateCmd(…)` | 全等 |
| `finance_prepayment_apply.go:159` | `update` | `CrudUpdateCmd(…)` | 全等 |
| `finance_prepayment_apply.go:213` | `delete` | `CrudDeleteCmd(…, false)` | 全等 |
| `finance_refund.go:172` | `create` | `CrudCreateCmd(…)` | 全等 |
| `finance_refund.go:200` | `update` | `CrudUpdateCmd(…)` | 全等 |
| `finance_refund.go:252` | `update-status` | `CrudUpdateStatusCmd(…)` | body `--data` — 全等 |

### B. 必须保留手写（工厂无法表达）

| 文件:行 | 原因 |
|---------|------|
| `finance_payment_apply.go:195` | `update-status` 用 `PutQuery` + `{id, approveStatus}`，非 body `--data` |
| `finance_prepayment_apply.go:186` | 同上（`approveStatus`） |
| `finance_invoice.go:226` | `update-status` 用 `PutQuery` + `{id, status}` |
| `finance_refund.go:137` | `get` 支持 `--id` **或** `--no` 二选一，工厂只有 `--id` |
| `finance_refund.go:228` | `delete` 用 `--id` 单删，但工厂 `single=true` 的 flag 名/文案需逐一核对（**见下方待确认项**） |
| 各文件 `page` | 分页子命令名为 `page` 而非 `list`，工厂固定生成 `list`（且 spec 禁止改名） |
| 各文件 `page-count` | 是否可换 `CrudPageCountCmdWithFixed` 需按文件核对筛选口径（部分文件 `page-count` 用独立 Filters 集合） |
| 各文件 `summary` / `export` | 工厂不生成 |

### C. 待实现时确认

- `finance_refund.go:228` 的 `delete`：`CrudDeleteCmd(…, true)` 的 `Use`/`Short`/flag 帮助文本是否与现有逐字一致；若 `Short` 文案不同则**保留手写**（本任务零行为变更优先）。
- `finance_account.go` / `finance_settlement.go` / `finance_account_settlement.go` / `finance_write_off.go` / `finance_receipt_payment.go` / `finance_invoice_apply.go` / `freight_apply.go`：实现前逐函数核对，只迁移落在 A 类的命令；**不在 A 类的一律保留**。

## Requirements

1. 只替换证明为逐字等价的命令；无法证明等价的**保留手写**，不发明新工厂、不扩展 `cmdutil` 公共契约。
2. 不改变任何 CLI 表面（子命令名、flag 名、flag 帮助文本、`Short` 文案、输出信封、退出码）。
3. 不碰 `internal/` 下任何文件（本任务**不修改** `internal/cmdutil`、`internal/client`、`internal/output`）。
4. 不碰 `cmd/root.go`，不新增命令组。
5. 迁移后删除因此变为未使用的私有函数（如 `newFinancePaymentApplyGetCmd`），以及确已无引用的 import。

## Acceptance Criteria

- [ ] `go build ./... && go vet ./... && go test ./...` 全部通过。
- [ ] `gofmt -l cmd/finance/` 无输出。
- [ ] 迁移前后 `wlt finance payment-apply --help`、`wlt finance invoice --help`、`wlt finance prepayment-apply --help`、`wlt finance refund --help` 的子命令列表与每个子命令的 flag 集合**完全一致**。
- [ ] 未迁移的 B 类手写命令源码保持原样（字符级不变）。
- [ ] 新增命令树装配测试，断言迁移后各 finance 子组仍包含预期子命令名（防回归）。
- [ ] `cmd/finance/` 总行数相比迁移前减少。
- [ ] 无新增 legacy 构造器写法（`command-factories.md` §8）。

## 非目标

- 不扩展 `DocumentConfig`（不加 `Fixed`/`Summary`/`Export`/`GetBy`/`StatusMode`）。
- 不统一 `page` → `list` 命名（spec 明确禁止当作全仓规则）。
- 不合并 `financeRefundFilters` / `financeRefundPageCountFilters` 两份筛选声明。
- 不处理 `context.Background()` 替换、`EnsureClient` 样板收敛（属后续任务）。
- 不修改 `skills/` 文档。
- 不修 `finance_account.go` / `finance_settlement.go` 等未被 A 类覆盖的手写命令。

## Notes

- 本任务范围经用户确认选择「只做零风险替换」，不做 `DocumentCmds` 扩展。
- 迁移前先加命令树快照测试作为护栏，再逐组替换。
