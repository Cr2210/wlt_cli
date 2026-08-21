# 订单预付关联命令实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在 wlt CLI 的 order 域下新增 `wlt order prepayment-relation` 命令组，覆盖后端 `/erp/order-prepayment-relation` 的 5 个接口。

**Architecture:** 新建单个命令文件 `cmd/order/order_prepayment.go`（package order，经 `init()` 注册到已有的 `orderCmd` 变量），5 个子命令全部手写（接口形状各异，不套用 `DocumentCmds`/`AddCRUDToParent` 工厂）；两个"可用来源查询"接口共享一个工厂函数。输出走既有 `cmdutil` 辅助函数，错误走 `output.NewExitError`。

**Tech Stack:** Go 1.26+ / Cobra / 既有内部包 `internal/cmdutil`、`internal/output`。零新依赖。

**设计文档:** `docs/superpowers/specs/2026-08-21-order-prepayment-relation-design.md`

## Global Constraints

- 输出协议：stdout 数据 JSON / stderr 错误 JSON / 退出码 0-6；API 调用失败用 `output.NewExitError(5, ...)`，JSON 解析失败用退出码 4（本计划无 `--data` flag，无解析错误路径）。
- flag 名 kebab-case，后端参数 camelCase —— 一律用 `cmdutil.CollectStringFlag`/`CollectStringFlags`/`CollectIntFlags` 收集（自带转换），必填 ID 直接放入 params（有 `cmd/order/order_main.go` 的 `map[string]any{"orderId": orderId}` 先例，int64 值可直接传）。
- 不新增第三方依赖；不改 `internal/cmdutil`、`internal/client`。
- 提交信息用仓库既有风格（conventional commit + 中文描述，如 `feat(order): ...`）。
- 金额单位为元，flag 类型 `float64`；ID 类 flag 类型 `int64`。
- 平台为 Windows + Git Bash，shell 命令用 POSIX 语法（`go test ./...`，不用 `&&` 之外的 PowerShell 语法）。

---

### Task 1: 命令组骨架 + list / create / update-amount

**Files:**
- Create: `cmd/order/order_prepayment.go`
- Test: `cmd/order/order_prepayment_test.go`

**Interfaces:**
- Consumes: 本包已有的 `orderCmd` 变量（`cmd/order/order.go:5`）；`cmdutil.EnsureClient/GetClient/OutputJSON/CollectStringFlag`、`output.NewExitError`。
- Produces: 包级常量 `orderPrepaymentRelationAPIPath`（string，值 `/erp/order-prepayment-relation`）、包级变量 `orderPrepaymentRelationCmd`（*cobra.Command，Use 为 `prepayment-relation`）、构造函数 `newOrderPrepayRelationListCmd()/newOrderPrepayRelationCreateCmd()/newOrderPrepayRelationUpdateAmountCmd()`（均返回 `*cobra.Command`）。Task 2 依赖该常量、变量与 init() 注册结构。

- [ ] **Step 1: 写失败的 flag 装配测试**

创建 `cmd/order/order_prepayment_test.go`（仿 `cmd/stock/stock_doc_test.go` 风格）：

```go
package order

import "testing"

func TestPrepayRelationListCmdFlags(t *testing.T) {
	c := newOrderPrepayRelationListCmd()
	for _, f := range []string{"order-id", "page-no", "page-size", "relation-type", "headers"} {
		if c.Flags().Lookup(f) == nil {
			t.Errorf("list cmd missing flag %q", f)
		}
	}
}

func TestPrepayRelationCreateCmdFlags(t *testing.T) {
	c := newOrderPrepayRelationCreateCmd()
	for _, f := range []string{"order-id", "relation-type", "relation-id", "relation-amount"} {
		if c.Flags().Lookup(f) == nil {
			t.Errorf("create cmd missing flag %q", f)
		}
	}
}

func TestPrepayRelationUpdateAmountCmdFlags(t *testing.T) {
	c := newOrderPrepayRelationUpdateAmountCmd()
	for _, f := range []string{"id", "relation-amount"} {
		if c.Flags().Lookup(f) == nil {
			t.Errorf("update-amount cmd missing flag %q", f)
		}
	}
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./cmd/order/ -run TestPrepayRelation -v`
Expected: FAIL（编译错误，`newOrderPrepayRelationListCmd` 等 undefined）

- [ ] **Step 3: 实现命令文件**

创建 `cmd/order/order_prepayment.go`，完整内容：

```go
package order

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

// orderPrepaymentRelationAPIPath 是订单预付关联的后端端点根路径。
const orderPrepaymentRelationAPIPath = "/erp/order-prepayment-relation"

var orderPrepaymentRelationCmd = &cobra.Command{
	Use:   "prepayment-relation",
	Short: "订单预付关联管理",
}

func init() {
	orderCmd.AddCommand(orderPrepaymentRelationCmd)
	orderPrepaymentRelationCmd.AddCommand(
		newOrderPrepayRelationListCmd(),
		newOrderPrepayRelationCreateCmd(),
		newOrderPrepayRelationUpdateAmountCmd(),
	)
}

// ---- 分页查询关联记录 ----

func newOrderPrepayRelationListCmd() *cobra.Command {
	var orderId int64
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询预付关联",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"orderId":  orderId,
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlag(cmd, params, "relation-type")
			cmdutil.CollectStringFlag(cmd, params, "headers")
			resp, err := cmdutil.GetClient().Get(context.Background(), orderPrepaymentRelationAPIPath+"/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询预付关联失败: %s", err), "")
			}
			// 该接口响应 data 为裸数组（无 list/total 包裹），直接输出。
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&orderId, "order-id", 0, "订单 ID")
	_ = c.MarkFlagRequired("order-id")
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().String("relation-type", "", "关联类型：PAYMENT-付款单，SUPPLIER-供应商期初")
	c.Flags().String("headers", "", "自定义导出表头")
	return c
}

// ---- 创建关联 ----

func newOrderPrepayRelationCreateCmd() *cobra.Command {
	var orderId, relationId int64
	var relationType string
	var relationAmount float64
	c := &cobra.Command{
		Use:   "create",
		Short: "创建预付关联",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body := map[string]any{
				"orderId":        orderId,
				"relationType":   relationType,
				"relationId":     relationId,
				"relationAmount": relationAmount,
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), orderPrepaymentRelationAPIPath+"/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建预付关联失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&orderId, "order-id", 0, "订单 ID")
	c.Flags().StringVar(&relationType, "relation-type", "", "关联类型：PAYMENT-付款单，SUPPLIER-供应商期初")
	c.Flags().Int64Var(&relationId, "relation-id", 0, "关联 ID（付款单 ID 或供应商期初 ID）")
	c.Flags().Float64Var(&relationAmount, "relation-amount", 0, "关联预付金额（单位：元）")
	_ = c.MarkFlagRequired("order-id")
	_ = c.MarkFlagRequired("relation-type")
	_ = c.MarkFlagRequired("relation-id")
	_ = c.MarkFlagRequired("relation-amount")
	return c
}

// ---- 更新关联金额 ----

func newOrderPrepayRelationUpdateAmountCmd() *cobra.Command {
	var id int64
	var relationAmount float64
	c := &cobra.Command{
		Use:   "update-amount",
		Short: "更新关联金额",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body := map[string]any{
				"id":             id,
				"relationAmount": relationAmount,
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), orderPrepaymentRelationAPIPath+"/update-amount", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新关联金额失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "关联记录 ID")
	c.Flags().Float64Var(&relationAmount, "relation-amount", 0, "关联预付金额（单位：元）")
	_ = c.MarkFlagRequired("id")
	_ = c.MarkFlagRequired("relation-amount")
	return c
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `go test ./cmd/order/ -run TestPrepayRelation -v`
Expected: PASS（3 个测试全绿）

- [ ] **Step 5: 构建验证后提交**

Run: `go build ./...`
Expected: 无输出（成功）

```bash
git add cmd/order/order_prepayment.go cmd/order/order_prepayment_test.go
git commit -m "feat(order): 订单预付关联命令组（list/create/update-amount）"
```

---

### Task 2: available-payments / available-initials 共享工厂

**Files:**
- Modify: `cmd/order/order_prepayment.go`（扩展 init() 注册 + 追加工厂函数）
- Test: `cmd/order/order_prepayment_test.go`（追加测试）

**Interfaces:**
- Consumes: Task 1 的 `orderPrepaymentRelationAPIPath`、`orderPrepaymentRelationCmd` 及其 init()；`cmdutil.CollectIntFlags/CollectStringFlags/ParsePagedJSON`。
- Produces: `newPrepayAvailableCmd(name, label string) *cobra.Command`，产出两个子命令 `available-payments`、`available-initials`，注册进 `orderPrepaymentRelationCmd`。

- [ ] **Step 1: 追加失败的测试**

在 `cmd/order/order_prepayment_test.go` 末尾追加：

```go
func TestPrepayAvailableCmdFlags(t *testing.T) {
	for _, name := range []string{"available-payments", "available-initials"} {
		c := newPrepayAvailableCmd(name, "查询")
		for _, f := range []string{"supplier-id", "order-id", "biz-date-start", "biz-date-end", "no", "headers", "page-no", "page-size"} {
			if c.Flags().Lookup(f) == nil {
				t.Errorf("%s cmd missing flag %q", name, f)
			}
		}
		if c.Name() != name {
			t.Errorf("name = %q, want %q", c.Name(), name)
		}
	}
}

func TestPrepayRelationCmdSubcommands(t *testing.T) {
	want := map[string]bool{"list": false, "create": false, "update-amount": false, "available-payments": false, "available-initials": false}
	for _, sub := range orderPrepaymentRelationCmd.Commands() {
		if _, ok := want[sub.Name()]; ok {
			want[sub.Name()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("prepayment-relation missing subcommand %q", name)
		}
	}
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./cmd/order/ -run TestPrepay -v`
Expected: FAIL（编译错误，`newPrepayAvailableCmd` undefined）

- [ ] **Step 3: 实现工厂函数并注册**

3a. 在 `cmd/order/order_prepayment.go` 的 init() 中，把：

```go
	orderPrepaymentRelationCmd.AddCommand(
		newOrderPrepayRelationListCmd(),
		newOrderPrepayRelationCreateCmd(),
		newOrderPrepayRelationUpdateAmountCmd(),
	)
```

替换为：

```go
	orderPrepaymentRelationCmd.AddCommand(
		newOrderPrepayRelationListCmd(),
		newOrderPrepayRelationCreateCmd(),
		newOrderPrepayRelationUpdateAmountCmd(),
		newPrepayAvailableCmd("available-payments", "查询可用的付款单列表"),
		newPrepayAvailableCmd("available-initials", "查询可用的供应商期初列表"),
	)
```

3b. 在文件末尾追加：

```go
// ---- 可用付款单 / 供应商期初查询（共享工厂，flag 与响应结构一致） ----

func newPrepayAvailableCmd(name, label string) *cobra.Command {
	var supplierId int64
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   name,
		Short: label,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"supplierId": supplierId,
				"pageNo":     pageNo,
				"pageSize":   pageSize,
			}
			cmdutil.CollectIntFlags(cmd, params, "order-id")
			cmdutil.CollectStringFlags(cmd, params, "biz-date-start", "biz-date-end", "no", "headers")
			resp, err := cmdutil.GetClient().Get(context.Background(), orderPrepaymentRelationAPIPath+"/"+name, params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("%s失败: %s", label, err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().Int64Var(&supplierId, "supplier-id", 0, "供应商 ID")
	_ = c.MarkFlagRequired("supplier-id")
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().Int64("order-id", 0, "当前订单 ID（排除已关联的）")
	c.Flags().String("biz-date-start", "", "单据日期起（如 2026-07-01 00:00:00）")
	c.Flags().String("biz-date-end", "", "单据日期止（如 2026-07-31 23:59:59）")
	c.Flags().String("no", "", "付款单号/期初编号模糊查询")
	c.Flags().String("headers", "", "自定义导出表头")
	return c
}
```

注意：`CollectStringFlags` 自动把 `biz-date-start` → `bizDateStart`；`CollectIntFlags` 把非零 `--order-id` → `orderId`（字符串形式）。

- [ ] **Step 4: 运行测试确认通过**

Run: `go test ./cmd/order/ -run TestPrepay -v`
Expected: PASS（5 个测试全绿）

- [ ] **Step 5: 构建验证后提交**

Run: `go build ./...`
Expected: 无输出（成功）

```bash
git add cmd/order/order_prepayment.go cmd/order/order_prepayment_test.go
git commit -m "feat(order): 预付关联可用付款单/期初查询命令"
```

---

### Task 3: 文档同步 + 全量验证

**Files:**
- Modify: `skills/references/order.md`
- Modify: `skills/SKILL.md`
- Modify: `CLAUDE.md`

**Interfaces:**
- Consumes: Task 1/2 落地的命令树与 flag 集（文档内容须与之一致）。
- Produces: 无代码接口。

- [ ] **Step 1: 更新 `skills/references/order.md`**

1a. 总览表（第 9-10 行附近）追加一行：

```markdown
| `wlt order prepayment-relation` | 订单预付关联：关联记录查询 / 创建 / 改金额 / 可用付款单与期初查询 |
```

1b. 在「## 常见工作流」章节之前插入新章节：

````markdown
## 订单预付关联 (`wlt order prepayment-relation`)

把订单与预付款来源（付款单 / 供应商期初）建立关联并维护关联金额。后端端点 `/erp/order-prepayment-relation`，`relationType` 枚举：`PAYMENT`（付款单）/ `SUPPLIER`（供应商期初）。

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `list` | 分页查询订单的关联记录 | `--order-id`（必填）、`--relation-type`、`--headers`、`--page-no`、`--page-size` |
| `create` | 创建关联，输出新关联 ID | `--order-id`、`--relation-type`、`--relation-id`、`--relation-amount`（均必填） |
| `update-amount` | 更新关联金额 | `--id`、`--relation-amount`（均必填） |
| `available-payments` | 查询供应商可用的付款单（含可用余额） | `--supplier-id`（必填）、`--order-id`、`--biz-date-start`、`--biz-date-end`、`--no`、`--headers`、`--page-no`、`--page-size` |
| `available-initials` | 查询供应商可用的期初 | 同 `available-payments` |

`list` 响应为裸数组（无 total）；`available-*` 为标准分页结果。金额单位均为元。

```bash
# 查某订单已关联的预付记录
wlt order prepayment-relation list --order-id 123

# 查供应商可用付款单余额 → 挑选后创建关联
wlt order prepayment-relation available-payments --supplier-id 456 --order-id 123
wlt order prepayment-relation create --order-id 123 --relation-type PAYMENT --relation-id 789 --relation-amount 1000

# 调整关联金额
wlt order prepayment-relation update-amount --id 1 --relation-amount 800
```
````

- [ ] **Step 2: 更新 `skills/SKILL.md`**

2a. 模块总览表中 order 一行：

旧：
```markdown
| `order` | 订单管理：主订单 / 计划（采购/销售 + CRUD） | [order.md](./references/order.md) |
```
新：
```markdown
| `order` | 订单管理：主订单 / 计划（采购/销售 + CRUD）/ 预付关联 | [order.md](./references/order.md) |
```

2b. 意图判断决策树中 order 一行：

旧：
```markdown
用户提到"订单/主订单/排产/关联运单/取消订单/完成订单/采购计划/销售计划/运输计划" → `order`
```
新：
```markdown
用户提到"订单/主订单/排产/关联运单/取消订单/完成订单/采购计划/销售计划/运输计划/预付关联/订单预付" → `order`
```

- [ ] **Step 3: 更新 `CLAUDE.md`**

项目结构中 order 目录注释：

旧：
```markdown
├── order/               # 订单管理（main/plan）
```
新：
```markdown
├── order/               # 订单管理（main/plan/prepayment-relation）
```

- [ ] **Step 4: 全量验证**

Run: `go build ./... && go test ./...`
Expected: 构建无输出、全部测试 PASS

Run: `go run . order prepayment-relation --help`
Expected: 帮助文本列出 `available-initials`、`available-payments`、`create`、`list`、`update-amount` 五个子命令

Run: `go run . order prepayment-relation create --help`
Expected: 帮助文本显示 4 个必填 flag（order-id/relation-type/relation-id/relation-amount，均标注 required）

- [ ] **Step 5: 提交**

```bash
git add skills/references/order.md skills/SKILL.md CLAUDE.md
git commit -m "docs(order): 补充预付关联 skills 与 CLAUDE.md 文档"
```
