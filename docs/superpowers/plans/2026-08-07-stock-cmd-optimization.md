# 库存命令优化 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 把全仓单据命令（采购入库/销售出库/其他入库/出库/调拨/盘点/退货）归一到一套通用工厂 `DocumentCmds`，消除 stock 手写样板与命名/headers/解析不一致，并重构 query/record。

**Architecture:** 由现有 `SalePurchaseCmds` 演进出通用 `DocumentCmds`（加 `--headers`、`page` 别名、可配时间字段，去掉隐式 base flags）。stock in/out/move/check 与退货改用该工厂；warehouse 改用 `AddCRUDToParent`；query/record 在 stock 包内用共享 `FlagSpec` + 本地辅助函数消除 flag 重复，保留特殊业务逻辑；删除 `crud_legacy.go`。

**Tech Stack:** Go 1.26+ / Cobra / 现有 `internal/cmdutil` CRUD 工厂 / 标准 `testing`（表驱动）

## Global Constraints

- 命令注册沿用各包 `init()` + `Register()` 约定；输出协议不变（stdout 数据 JSON / stderr 错误 / 退出码 0-6）。
- 后端 API 契约不变：`{APIPath}/page|page-count|get|create|update|delete|update-status`；时间范围传数组 `{TimeKey}[0]/{TimeKey}[1]`（已确认 stock in/out/move-check 后端认数组）。
- kebab-case flag → camelCase 参数键由 `CollectStringFlags`/`flagToParamKey` 自动转换（如 `supplier-id`→`supplierId`）。
- stock 单据的 `page` 子命令改名为 `list`，但必须保留 `page` 为 cobra `Aliases` 兼容别名，不破坏旧调用。
- 每个任务结束 `go build ./...` 通过、`make test` 通过、并 commit。
- RunE 依赖全局 `APIClient`（见 `internal/cmdutil/client.go:39`），无法纯单测；工厂层逻辑（时间转换、命令结构）用单元测试守护，迁移类任务靠编译 + `make test` + `--help` 结构核对验证等价性。

## File Structure

- `internal/cmdutil/crud.go`（改）：`SalePurchaseConfig`→`DocumentConfig`、`SalePurchaseCmds`→`DocumentCmds`；加 `Headers`/`PageAlias` 字段；删 `addSalePurchaseBaseFlags`；`collectSalePurchaseTime`→`collectDocumentTimeRange`、`addSalePurchaseTimeFlags`→`addDocumentTimeRangeFlags`、`salePurchaseTimeFlags`→`documentTimeFlags`；`salePurchaseListCmd`→`documentListCmd`、`salePurchasePageCountCmd`→`documentPageCountCmd`。其他工厂（Contract/Plan/Order 等）不动。
- `internal/cmdutil/crud_test.go`（新）：`DocumentCmds` 生成的命令结构 + `collectDocumentTimeRange` 单元测试。
- `cmd/purchase/purchase_in.go`（改）：改用 `DocumentCmds`，`Filters` 补 base 6 字段，`Headers:true`。
- `cmd/sale/sale_out.go`（改）：同上（TimeKey=outTime）。
- `cmd/stock/stock_in.go` / `stock_out.go` / `stock_move.go` / `stock_check.go`（改）：重写为 `DocumentCmds` 调用（各 ~20 行）。
- `cmd/stock/stock_warehouse.go`（改）：重写为 `AddCRUDToParent` + `CrudSimpleListCmd`。
- `cmd/stock/stock_doc.go`（新）：query/record 共享 `FlagSpec`（`queryPageFilters`/`recordPageFilters`/`recordDimensionFilters`）+ 本地辅助 `newStockPagedCmd`/`newStockCountCmd`。
- `cmd/stock/stock_query.go` / `stock_record.go`（改）：用共享 FlagSpec + 辅助函数重构；保留 `count`/`get`/`stock-record-count` 等特殊逻辑。
- `cmd/purchase/purchase_return.go` / `cmd/sale/sale_return.go`（改）：`NewCRUDSubCmd`→`DocumentCmds`。
- `internal/cmdutil/crud_legacy.go`（删）：全仓无调用方后删除。
- `CLAUDE.md`（改）：修正 CRUD 模式描述。
- `skills/references/stock.md` / `sale-purchase.md`（改）：`page`→`list`（别名说明）、时间 flag、`--headers`。

---

### Task 1: 通用单据工厂 DocumentCmds（TDD）

**Files:**
- Modify: `internal/cmdutil/crud.go`（`SalePurchaseConfig` 段 ~866-964 行 + 时间 helper ~776-800 行 + `addSalePurchaseBaseFlags` ~855-864 行）
- Modify: `cmd/purchase/purchase_in.go`
- Modify: `cmd/sale/sale_out.go`
- Create: `internal/cmdutil/crud_test.go`

**Interfaces:**
- Produces: `DocumentConfig{Name,APIPath,Label,Filters,TimeKey,Headers,PageAlias}`、`DocumentCmds(parent *cobra.Command, cfgs ...DocumentConfig)`、`collectDocumentTimeRange(cmd, params, timeKey)`、`addDocumentTimeRangeFlags(cmd)`。

- [ ] **Step 1: 写失败测试 `internal/cmdutil/crud_test.go`**

```go
package cmdutil

import (
	"testing"

	"github.com/spf13/cobra"
)

// 子命令名断言辅助
func mustHaveSub(t *testing.T, parent *cobra.Command, names ...string) {
	t.Helper()
	got := map[string]bool{}
	for _, c := range parent.Commands() {
		got[c.Name()] = true
	}
	for _, n := range names {
		if !got[n] {
			t.Errorf("missing subcommand %q; got %v", n, got)
		}
	}
}

func TestCollectDocumentTimeRange(t *testing.T) {
	cmd := &cobra.Command{}
	addDocumentTimeRangeFlags(cmd)
	_ = cmd.Flags().Set("start-time", "2026-07-01 00:00:00")
	_ = cmd.Flags().Set("end-time", "2026-07-31 23:59:59")

	params := map[string]any{}
	collectDocumentTimeRange(cmd, params, "inTime")
	if params["inTime[0]"] != "2026-07-01 00:00:00" {
		t.Errorf("inTime[0] = %v, want start-time", params["inTime[0]"])
	}
	if params["inTime[1]"] != "2026-07-31 23:59:59" {
		t.Errorf("inTime[1] = %v, want end-time", params["inTime[1]"])
	}
}

func TestCollectDocumentTimeRangeEmpty(t *testing.T) {
	cmd := &cobra.Command{}
	addDocumentTimeRangeFlags(cmd) // 不 set 任何值
	params := map[string]any{}
	collectDocumentTimeRange(cmd, params, "outTime")
	if _, ok := params["outTime[0]"]; ok {
		t.Errorf("outTime[0] should be absent when flags empty; got %v", params["outTime[0]"])
	}
}

func TestDocumentCmdsStructure(t *testing.T) {
	parent := &cobra.Command{Use: "test"}
	DocumentCmds(parent, DocumentConfig{
		Name:    "in",
		APIPath: "/erp/stock-in",
		Label:   "入库单",
		TimeKey: "inTime",
		Headers: true,
		PageAlias: true,
		Filters: []FlagSpec{
			{Name: "no", Usage: "单号"},
			{Name: "supplier-id", Usage: "供应商"},
		},
	})

	// 子命令齐全
	in := subNamed(parent, "in")
	if in == nil {
		t.Fatal("missing subcommand 'in'")
	}
	mustHaveSub(t, in, "list", "page-count", "get", "create", "update", "delete", "update-status")

	// list 有 page 别名
	list := subNamed(in, "list")
	for _, a := range list.Aliases {
		if a == "page" {
			goto ok
		}
	}
	t.Errorf("list Aliases = %v, want contain 'page'", list.Aliases)
ok:

	// list 注册了分页/时间/headers/filter flag
	for _, f := range []string{"page-no", "page-size", "start-time", "end-time", "headers", "no", "supplier-id"} {
		if list.Flags().Lookup(f) == nil {
			t.Errorf("list missing flag %q", f)
		}
	}
}

func TestDocumentCmdsNoExtraTimeFlags(t *testing.T) {
	// TimeKey 为空时不注册时间 flag
	parent := &cobra.Command{Use: "t"}
	DocumentCmds(parent, DocumentConfig{Name: "x", APIPath: "/erp/x", Label: "X", Filters: nil})
	x := subNamed(parent, "x")
	list := subNamed(x, "list")
	if list.Flags().Lookup("start-time") != nil {
		t.Errorf("start-time should not be registered when TimeKey empty")
	}
}

// subNamed 在 parent 的直接子命令里按名查找。
func subNamed(parent *cobra.Command, name string) *cobra.Command {
	for _, c := range parent.Commands() {
		if c.Name() == name {
			return c
		}
	}
	return nil
}
```

- [ ] **Step 2: 运行测试，确认失败**

Run: `go test ./internal/cmdutil/ -run 'TestCollectDocumentTimeRange|TestDocumentCmds' -v`
Expected: 编译失败（`collectDocumentTimeRange` / `addDocumentTimeRangeFlags` / `DocumentConfig` / `DocumentCmds` 未定义）。

- [ ] **Step 3: 改造 `internal/cmdutil/crud.go`**

把时间 flag 段（约 778-800 行）整体替换为：

```go
// documentTimeFlags 是单据共用的起止时间筛选 flag；
// 查询时转成 {timeKey}[0] / {timeKey}[1] 数组参数
// （采购入库 inTime、销售出库 outTime、调拨 moveTime、盘点 checkTime 等）。
var documentTimeFlags = []FlagSpec{
	{Name: "start-time", Usage: "开始时间（如 2026-07-01 00:00:00）"},
	{Name: "end-time", Usage: "结束时间（如 2026-07-31 23:59:59）"},
}

// collectDocumentTimeRange 把 start-time/end-time 转成 {timeKey}[0]/{timeKey}[1]。
func collectDocumentTimeRange(cmd *cobra.Command, params map[string]any, timeKey string) {
	if v, _ := cmd.Flags().GetString("start-time"); v != "" {
		params[fmt.Sprintf("%s[0]", timeKey)] = v
	}
	if v, _ := cmd.Flags().GetString("end-time"); v != "" {
		params[fmt.Sprintf("%s[1]", timeKey)] = v
	}
}

func addDocumentTimeRangeFlags(c *cobra.Command) {
	for _, f := range documentTimeFlags {
		c.Flags().String(f.Name, "", f.Usage)
	}
}
```

删除 `addSalePurchaseBaseFlags` 函数（约 855-864 行，整段删）。

把 `SalePurchaseConfig` 结构体（约 866-873 行）及其后续所有 `SalePurchase`/`salePurchase` 标识符替换。结构体改为：

```go
// DocumentConfig 描述一个单据 CRUD 子模块（采购入库/销售出库/其他入库/出库/调拨/盘点/退货）。
type DocumentConfig struct {
	Name      string     // 子命令名：in / out / move / check / return …
	APIPath   string     // 后端根路径：/erp/stock-in
	Label     string     // 中文标签：入库单
	Filters   []FlagSpec // 模块特有筛选 flag（含原 base 字段，如 warehouse-id/product-id/no/status/type 等）
	TimeKey   string     // 时间字段：空=无时间筛选；非空=注册 --start-time/--end-time，传 {TimeKey}[0]/{TimeKey}[1]
	Headers   bool       // 是否注册 --headers 自定义导出表头
	PageAlias bool       // 是否给 list 注册 page 别名（向后兼容）
}

// DocumentCmds 为一组单据子模块生成标准 CRUD 子命令，统一注册到 parent。
// 每个子模块自动拥有：list(+可选 page 别名) / page-count / get / create / update / delete / update-status。
func DocumentCmds(parent *cobra.Command, cfgs ...DocumentConfig) {
	for _, cfg := range cfgs {
		cmd := &cobra.Command{
			Use:   cfg.Name,
			Short: cfg.Label + "管理",
		}
		cmd.AddCommand(documentListCmd(cfg))
		cmd.AddCommand(documentPageCountCmd(cfg))
		cmd.AddCommand(CrudGetCmd(cfg.APIPath, cfg.Label))
		cmd.AddCommand(CrudCreateCmd(cfg.APIPath, cfg.Label))
		cmd.AddCommand(CrudUpdateCmd(cfg.APIPath, cfg.Label))
		cmd.AddCommand(CrudDeleteCmd(cfg.APIPath, cfg.Label, false))
		cmd.AddCommand(CrudUpdateStatusCmd(cfg.APIPath, cfg.Label))
		parent.AddCommand(cmd)
	}
}

func documentListCmd(cfg DocumentConfig) *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:     "list",
		Short:   "分页查询" + cfg.Label,
		Aliases: docPageAlias(cfg.PageAlias),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			for _, f := range cfg.Filters {
				CollectStringFlag(cmd, params, f.Name)
			}
			if cfg.TimeKey != "" {
				collectDocumentTimeRange(cmd, params, cfg.TimeKey)
			}
			if cfg.Headers {
				CollectStringFlag(cmd, params, "headers")
			}
			resp, err := GetClient().Get(context.Background(), cfg.APIPath+"/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询%s失败: %s", cfg.Label, err), "")
			}
			return ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	for _, f := range cfg.Filters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	if cfg.TimeKey != "" {
		addDocumentTimeRangeFlags(c)
	}
	if cfg.Headers {
		c.Flags().String("headers", "", "自定义导出表头")
	}
	return c
}

func documentPageCountCmd(cfg DocumentConfig) *cobra.Command {
	c := &cobra.Command{
		Use:   "page-count",
		Short: "统计" + cfg.Label + "数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			for _, f := range cfg.Filters {
				CollectStringFlag(cmd, params, f.Name)
			}
			if cfg.TimeKey != "" {
				collectDocumentTimeRange(cmd, params, cfg.TimeKey)
			}
			// page-count 不带 headers（与原实现一致）
			resp, err := GetClient().Get(context.Background(), cfg.APIPath+"/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计%s失败: %s", cfg.Label, err), "")
			}
			return OutputJSON(json.RawMessage(resp.Data))
		},
	}
	for _, f := range cfg.Filters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	if cfg.TimeKey != "" {
		addDocumentTimeRangeFlags(c)
	}
	return c
}

// docPageAlias 在启用别名时返回 ["page"]，否则 nil。
func docPageAlias(enabled bool) []string {
	if enabled {
		return []string{"page"}
	}
	return nil
}
```

- [ ] **Step 4: 更新采购/销售调用点**

`cmd/purchase/purchase_in.go` 整体替换为：

```go
package purchase

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	cmdutil.DocumentCmds(purchaseCmd,
		cmdutil.DocumentConfig{
			Name:    "in",
			APIPath: "/erp/purchase-in",
			Label:   "采购入库",
			TimeKey: "inTime",
			Headers: true,
			Filters: []cmdutil.FlagSpec{
				{Name: "warehouse-id", Usage: "仓库 ID"},
				{Name: "product-id", Usage: "产品 ID"},
				{Name: "product-name", Usage: "产品名称（模糊搜索）"},
				{Name: "no", Usage: "单号"},
				{Name: "status", Usage: "状态"},
				{Name: "type", Usage: "类型"},
				{Name: "supplier-id", Usage: "供应商 ID"},
				{Name: "metrics-name", Usage: "检测指标名（如 含水）"},
			},
		},
	)
}
```

`cmd/sale/sale_out.go` 整体替换为：

```go
package sale

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	cmdutil.DocumentCmds(saleCmd,
		cmdutil.DocumentConfig{
			Name:    "out",
			APIPath: "/erp/sale-out",
			Label:   "销售出库",
			TimeKey: "outTime",
			Headers: true,
			Filters: []cmdutil.FlagSpec{
				{Name: "warehouse-id", Usage: "仓库 ID"},
				{Name: "product-id", Usage: "产品 ID"},
				{Name: "product-name", Usage: "产品名称（模糊搜索）"},
				{Name: "no", Usage: "单号"},
				{Name: "status", Usage: "状态"},
				{Name: "type", Usage: "类型"},
				{Name: "customer-id", Usage: "客户 ID"},
				{Name: "batch-no", Usage: "批次号"},
			},
		},
	)
}
```

- [ ] **Step 5: 运行测试 + 构建**

Run: `go test ./internal/cmdutil/ -v`
Expected: 全部 PASS（含新 `TestCollectDocumentTimeRange*`、`TestDocumentCmds*` 与既有 `TestFlagToParamKey`）。

Run: `go build ./...`
Expected: 编译通过（无残留 `SalePurchase` 引用）。若有遗漏的 `SalePurchase`/`salePurchase` 标识符（全仓 grep），一并改名。

- [ ] **Step 6: 回归核对采购/销售 flag 等价**

Run: `go run . purchase in list --help` 与 `go run . sale out list --help`
Expected: list 含 `--warehouse-id --product-id --product-name --no --status --type`（原 base flags）+ 各自特有字段 + `--start-time --end-time --headers --page-no --page-size`；与改造前列表筛选能力等价（多了 `--headers`、时间由 `inTime` 改为 `start-time/end-time` 范围）。

- [ ] **Step 7: Commit**

```bash
git add internal/cmdutil/crud.go internal/cmdutil/crud_test.go cmd/purchase/purchase_in.go cmd/sale/sale_out.go
git commit -m "refactor(cmdutil): SalePurchaseCmds 通用化为 DocumentCmds（+headers/page 别名/可配时间）"
```

---

### Task 2: stock warehouse 迁移到 AddCRUDToParent

**Files:**
- Modify: `cmd/stock/stock_warehouse.go`（整文件重写）

**Interfaces:**
- Consumes: `cmdutil.AddCRUDToParent`（既有）、`cmdutil.CRUDConfig`、`cmdutil.CrudSimpleListCmd`（既有）、`cmdutil.FlagSpec`。
- 仓库现状子命令：list/get/simple-list/create/update/delete/update-status（无 page-count）；list filters：name/status/type；delete 用 `ids`。

- [ ] **Step 1: 重写 `cmd/stock/stock_warehouse.go`**

```go
package stock

import (
	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
)

var warehouseCmd = &cobra.Command{
	Use:  "warehouse",
	Short: "仓库管理",
	Long: "管理 ERP 仓库：查询、创建、更新、删除、状态管理。",
}

func init() {
	stockCmd.AddCommand(warehouseCmd)
	cmdutil.AddCRUDToParent(warehouseCmd, cmdutil.CRUDConfig{
		APIPath: "/erp/warehouse",
		Label:   "仓库",
		ListFilters: []cmdutil.FlagSpec{
			{Name: "name", Usage: "仓库名称（模糊）"},
			{Name: "status", Usage: "状态"},
			{Name: "type", Usage: "类型"},
		},
		SingleDelete: false, // DELETE ?ids=
	})
	warehouseCmd.AddCommand(cmdutil.CrudSimpleListCmd("/erp/warehouse", "仓库"))
}
```

- [ ] **Step 2: 构建并核对结构**

Run: `go build ./...`
Expected: 编译通过。

Run: `go run . stock warehouse --help`
Expected: 子命令含 `list get simple-list create update delete update-status`（不再有手写 page-count，仓库本就没有）。

Run: `go run . stock warehouse list --help`
Expected: flag 含 `--page-no --page-size --name --status --type`（无 `--headers`，仓库原本就没有）。

- [ ] **Step 3: 跑全量测试并提交**

Run: `make test`
Expected: PASS。

```bash
git add cmd/stock/stock_warehouse.go
git commit -m "refactor(stock/warehouse): 迁移到 AddCRUDToParent，分页解析统一为 ParsePagedJSON"
```

---

### Task 3: stock in 迁移到 DocumentCmds

**Files:**
- Modify: `cmd/stock/stock_in.go`（整文件重写）

**Interfaces:**
- Consumes: `DocumentCmds`、`DocumentConfig`（Task 1 产出）。
- 现状：21 个业务筛选 flag（不含 `in-time`/`headers`）+ `in-time` 时间单值；子命令 page/page-count/get/create/update/delete/update-status；delete 用 `ids`。

> 迁移后：`in-time` 单值 → `--start-time/--end-time`（转 `inTime[0]/inTime[1]` 数组，已确认后端认数组）；`page` → `list`（保留 `page` 别名）；新增 `--headers`；分页解析改用 `ParsePagedJSON`。

- [ ] **Step 1: 重写 `cmd/stock/stock_in.go`**

```go
package stock

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	cmdutil.DocumentCmds(stockCmd,
		cmdutil.DocumentConfig{
			Name:      "in",
			APIPath:   "/erp/stock-in",
			Label:     "入库单",
			TimeKey:   "inTime",
			Headers:   true,
			PageAlias: true,
			Filters: []cmdutil.FlagSpec{
				{Name: "no", Usage: "入库单号"},
				{Name: "supplier-id", Usage: "供应商编号"},
				{Name: "supplier-name", Usage: "供应商"},
				{Name: "status", Usage: "状态"},
				{Name: "remark", Usage: "备注"},
				{Name: "creator", Usage: "创建者"},
				{Name: "product-id", Usage: "产品编号"},
				{Name: "product-name", Usage: "产品名"},
				{Name: "warehouse-id", Usage: "仓库编号"},
				{Name: "warehouse-name", Usage: "仓库名"},
				{Name: "metrics-name", Usage: "产品指标名称"},
				{Name: "creator-name", Usage: "创建人名称"},
				{Name: "user-id", Usage: "业务人员ID"},
				{Name: "receive-address", Usage: "收货地址"},
				{Name: "send-address", Usage: "发货地址"},
				{Name: "batch-no", Usage: "批次号"},
				{Name: "create-time", Usage: "创建时间"},
				{Name: "updater-name", Usage: "更新人名称"},
				{Name: "update-time", Usage: "更新时间"},
				{Name: "custom-order", Usage: "前端自定义排序规则"},
				{Name: "keyword", Usage: "关键字"},
			},
		},
	)
}
```

- [ ] **Step 2: 构建并核对结构与别名**

Run: `go build ./...`
Expected: 编译通过。

Run: `go run . stock in list --help`
Expected: flag 含上述 21 个 + `--start-time --end-time --headers --page-no --page-size`；不再有 `--in-time`。

Run: `go run . stock in page --help`
Expected: 与 `list --help` 相同（`page` 是 `list` 的别名，仍可调用）。

Run: `go run . stock in --help`
Expected: 子命令含 `list page-count get create update delete update-status`。

- [ ] **Step 3: 跑全量测试并提交**

Run: `make test`
Expected: PASS。

```bash
git add cmd/stock/stock_in.go
git commit -m "refactor(stock/in): 迁移到 DocumentCmds，时间改 start/end 数组，page→list(别名)"
```

---

### Task 4: stock out 迁移到 DocumentCmds

**Files:**
- Modify: `cmd/stock/stock_out.go`（整文件重写）

**Interfaces:**
- Consumes: `DocumentCmds`、`DocumentConfig`（Task 1）。
- 现状：与 `stock in` 对称，supplier→customer；`out-time` 单值；21 个业务筛选 flag（不含 `out-time`/`headers`）。

- [ ] **Step 1: 重写 `cmd/stock/stock_out.go`**

```go
package stock

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	cmdutil.DocumentCmds(stockCmd,
		cmdutil.DocumentConfig{
			Name:      "out",
			APIPath:   "/erp/stock-out",
			Label:     "出库单",
			TimeKey:   "outTime",
			Headers:   true,
			PageAlias: true,
			Filters: []cmdutil.FlagSpec{
				{Name: "no", Usage: "出库单号"},
				{Name: "customer-id", Usage: "客户编号"},
				{Name: "customer-name", Usage: "客户"},
				{Name: "status", Usage: "状态"},
				{Name: "remark", Usage: "备注"},
				{Name: "creator", Usage: "创建者"},
				{Name: "product-id", Usage: "产品编号"},
				{Name: "product-name", Usage: "产品名"},
				{Name: "warehouse-id", Usage: "仓库编号"},
				{Name: "warehouse-name", Usage: "仓库名"},
				{Name: "metrics-name", Usage: "产品指标名称"},
				{Name: "creator-name", Usage: "创建人名称"},
				{Name: "user-id", Usage: "业务人员ID"},
				{Name: "receive-address", Usage: "收货地址"},
				{Name: "send-address", Usage: "发货地址"},
				{Name: "batch-no", Usage: "批次号"},
				{Name: "create-time", Usage: "创建时间"},
				{Name: "updater-name", Usage: "更新人名称"},
				{Name: "update-time", Usage: "更新时间"},
				{Name: "custom-order", Usage: "前端自定义排序规则"},
				{Name: "keyword", Usage: "关键字"},
			},
		},
	)
}
```

- [ ] **Step 2: 构建 + 别名核对**

Run: `go build ./... && go run . stock out list --help`
Expected: flag 含 21 个 + `--start-time --end-time --headers`；无 `--out-time`。

Run: `go run . stock out page --help`
Expected: 同 `list --help`（别名生效）。

- [ ] **Step 3: 测试并提交**

Run: `make test` → PASS。

```bash
git add cmd/stock/stock_out.go
git commit -m "refactor(stock/out): 迁移到 DocumentCmds，时间改 outTime start/end 数组，page→list(别名)"
```

---

### Task 5: stock move 迁移到 DocumentCmds

**Files:**
- Modify: `cmd/stock/stock_move.go`（整文件重写）

**Interfaces:**
- Consumes: `DocumentCmds`、`DocumentConfig`（Task 1）。
- 现状：`move-time` 单值；19 个业务筛选 flag（含 `from-warehouse-id`/`to-warehouse-id`/`updater`，不含 `move-time`/`headers`）。

- [ ] **Step 1: 重写 `cmd/stock/stock_move.go`**

```go
package stock

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	cmdutil.DocumentCmds(stockCmd,
		cmdutil.DocumentConfig{
			Name:      "move",
			APIPath:   "/erp/stock-move",
			Label:     "调拨单",
			TimeKey:   "moveTime",
			Headers:   true,
			PageAlias: true,
			Filters: []cmdutil.FlagSpec{
				{Name: "no", Usage: "调拨单号"},
				{Name: "create-time", Usage: "创建时间"},
				{Name: "update-time", Usage: "更新时间"},
				{Name: "status", Usage: "状态"},
				{Name: "remark", Usage: "备注"},
				{Name: "creator", Usage: "创建者编号"},
				{Name: "creator-name", Usage: "创建者姓名"},
				{Name: "updater", Usage: "更新者编号"},
				{Name: "updater-name", Usage: "更新者姓名"},
				{Name: "product-id", Usage: "产品编号"},
				{Name: "from-warehouse-id", Usage: "调出仓库编号"},
				{Name: "to-warehouse-id", Usage: "调入仓库编号"},
				{Name: "product-name", Usage: "产品名称"},
				{Name: "metrics-name", Usage: "产品指标"},
				{Name: "batch-no", Usage: "批次号"},
				{Name: "user-id", Usage: "业务负责人ID"},
				{Name: "custom-order", Usage: "前端自定义排序规则"},
				{Name: "keyword", Usage: "关键字"},
			},
		},
	)
}
```

- [ ] **Step 2: 构建 + 别名核对**

Run: `go build ./... && go run . stock move list --help`
Expected: flag 含 18 个 + `--start-time --end-time --headers`；无 `--move-time`。

Run: `go run . stock move page --help`
Expected: 同 `list --help`。

- [ ] **Step 3: 测试并提交**

Run: `make test` → PASS。

```bash
git add cmd/stock/stock_move.go
git commit -m "refactor(stock/move): 迁移到 DocumentCmds，时间改 moveTime start/end 数组，page→list(别名)"
```

---

### Task 6: stock check 迁移到 DocumentCmds

**Files:**
- Modify: `cmd/stock/stock_check.go`（整文件重写）

**Interfaces:**
- Consumes: `DocumentCmds`、`DocumentConfig`（Task 1）。
- 现状：`check-time` 单值；18 个业务筛选 flag（含 `warehouse-id`，不含 `check-time`/`headers`）。

- [ ] **Step 1: 重写 `cmd/stock/stock_check.go`**

```go
package stock

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	cmdutil.DocumentCmds(stockCmd,
		cmdutil.DocumentConfig{
			Name:      "check",
			APIPath:   "/erp/stock-check",
			Label:     "盘点单",
			TimeKey:   "checkTime",
			Headers:   true,
			PageAlias: true,
			Filters: []cmdutil.FlagSpec{
				{Name: "no", Usage: "盘点单号"},
				{Name: "warehouse-id", Usage: "仓库编号"},
				{Name: "status", Usage: "状态"},
				{Name: "remark", Usage: "备注"},
				{Name: "creator", Usage: "创建者"},
				{Name: "creator-name", Usage: "创建者姓名"},
				{Name: "create-time", Usage: "创建时间"},
				{Name: "update-time", Usage: "更新时间"},
				{Name: "updater", Usage: "更新者编号"},
				{Name: "updater-name", Usage: "更新者姓名"},
				{Name: "product-id", Usage: "产品编号"},
				{Name: "product-name", Usage: "产品名称"},
				{Name: "metrics-name", Usage: "产品指标"},
				{Name: "batch-no", Usage: "批次号"},
				{Name: "user-id", Usage: "业务负责人ID"},
				{Name: "custom-order", Usage: "前端自定义排序规则"},
				{Name: "keyword", Usage: "关键字"},
			},
		},
	)
}
```

- [ ] **Step 2: 构建 + 别名核对**

Run: `go build ./... && go run . stock check list --help`
Expected: flag 含 17 个 + `--start-time --end-time --headers`；无 `--check-time`。

Run: `go run . stock check page --help`
Expected: 同 `list --help`。

- [ ] **Step 3: 测试并提交**

Run: `make test` → PASS。

```bash
git add cmd/stock/stock_check.go
git commit -m "refactor(stock/check): 迁移到 DocumentCmds，时间改 checkTime start/end 数组，page→list(别名)"
```

---

### Task 7: stock query/record 重构（共享 FlagSpec + 本地辅助）

**Files:**
- Create: `cmd/stock/stock_doc.go`（共享 FlagSpec + `newStockPagedCmd`/`newStockCountCmd` + flag 注册/收集辅助）
- Create: `cmd/stock/stock_doc_test.go`（辅助函数结构单测）
- Modify: `cmd/stock/stock_query.go`（用辅助函数 + 共享 flag 重构，保留 get/count/ledger 特殊逻辑）
- Modify: `cmd/stock/stock_record.go`（同上）

**Interfaces:**
- Produces（包内私有）：`queryPageFilters`/`recordPageFilters`/`recordDimensionFilters`、`newStockPagedCmd(use, apiPath, endpoint, short, filters, withHeaders)`、`newStockCountCmd(use, apiPath, endpoint, short, filters)`、`addFilterFlags`/`collectFilterFlags`。
- 行为保持：query/record 子命令名不变（`page`/`page-count`/`ledger`/`ledger-count`/`record-page`/`total-cost`/`count`/`get`/`stock-record-count`），仅消除 flag 重复 + 统一分页解析为 `ParsePagedJSON`。`query ledger` 的"解析失败回退原样输出"保留（手写）。

> query/record 的时间 flag（`in-time`/`create-time`）保持单值语义不动（走 `/erp/stock-record`，与 stock-in 数组契约不同），仅作为普通 `FlagSpec` 收集。

- [ ] **Step 1: 写失败测试 `cmd/stock/stock_doc_test.go`**

```go
package stock

import (
	"testing"

	"github.com/weiliantong/cli/internal/cmdutil"
)

func TestNewStockPagedCmdFlags(t *testing.T) {
	c := newStockPagedCmd("page", "/erp/stock", "page", "分页查询产品库存",
		[]cmdutil.FlagSpec{{Name: "product-id"}, {Name: "keyword"}}, true)
	for _, f := range []string{"page-no", "page-size", "product-id", "keyword", "headers"} {
		if c.Flags().Lookup(f) == nil {
			t.Errorf("paged cmd missing flag %q", f)
		}
	}
	if c.Name() != "page" {
		t.Errorf("name = %q, want page", c.Name())
	}
}

func TestNewStockCountCmdFlags(t *testing.T) {
	c := newStockCountCmd("page-count", "/erp/stock", "page-count", "统计",
		[]cmdutil.FlagSpec{{Name: "warehouse-id"}})
	if c.Flags().Lookup("warehouse-id") == nil {
		t.Errorf("count cmd missing flag warehouse-id")
	}
	for _, f := range []string{"page-no", "page-size", "headers"} {
		if c.Flags().Lookup(f) != nil {
			t.Errorf("count cmd should not have flag %q", f)
		}
	}
}
```

- [ ] **Step 2: 运行测试，确认失败**

Run: `go test ./cmd/stock/ -run 'TestNewStock' -v`
Expected: 编译失败（`newStockPagedCmd`/`newStockCountCmd` 未定义）。

- [ ] **Step 3: 创建 `cmd/stock/stock_doc.go`**

```go
package stock

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

// queryPageFilters 是 stock query 的 page/page-count/ledger/ledger-count 共享筛选 flag。
var queryPageFilters = []cmdutil.FlagSpec{
	{Name: "product-id", Usage: "产品编号"},
	{Name: "warehouse-id", Usage: "仓库编号"},
	{Name: "metrics-name", Usage: "指标名称"},
	{Name: "batch-no", Usage: "批次号"},
	{Name: "send-address", Usage: "发货地"},
	{Name: "plan-no", Usage: "方案号"},
	{Name: "supplier-id", Usage: "供应商编号"},
	{Name: "supplier-name", Usage: "供应商名称"},
	{Name: "is-detail", Usage: "是否明细"},
	{Name: "count-positive", Usage: "库存为正数"},
	{Name: "keyword", Usage: "关键字"},
}

// recordPageFilters 是 stock record 的 page/page-count 共享筛选 flag。
var recordPageFilters = []cmdutil.FlagSpec{
	{Name: "product-id", Usage: "产品编号"},
	{Name: "category-id", Usage: "产品分类编号"},
	{Name: "warehouse-id", Usage: "仓库编号"},
	{Name: "biz-type", Usage: "业务类型"},
	{Name: "biz-no", Usage: "业务单号"},
	{Name: "create-time", Usage: "操作时间"},
	{Name: "in-time", Usage: "出入库时间"},
	{Name: "metrics-name", Usage: "指标名称"},
	{Name: "product-name", Usage: "产品名称"},
	{Name: "batch-no", Usage: "批次号"},
	{Name: "keyword", Usage: "关键字"},
}

// recordDimensionFilters 是 stock record 的 record-page/total-cost 共享维度 flag。
var recordDimensionFilters = []cmdutil.FlagSpec{
	{Name: "product-id", Usage: "产品 ID"},
	{Name: "warehouse-id", Usage: "仓库 ID"},
	{Name: "metrics-name", Usage: "指标名称"},
	{Name: "batch-no", Usage: "批次号"},
	{Name: "supplier-id", Usage: "供应商 ID"},
	{Name: "source-supplier-id", Usage: "关联供应商 ID"},
	{Name: "out-count", Usage: "出库数量"},
}

func addFilterFlags(c *cobra.Command, filters []cmdutil.FlagSpec) {
	for _, f := range filters {
		c.Flags().String(f.Name, "", f.Usage)
	}
}

func collectFilterFlags(cmd *cobra.Command, params map[string]any, filters []cmdutil.FlagSpec) {
	for _, f := range filters {
		cmdutil.CollectStringFlag(cmd, params, f.Name)
	}
}

// newStockPagedCmd 构建分页查询命令：GET {apiPath}/{endpoint}，输出 ParsePagedJSON。
// use=命令名；endpoint=page/batch-detail/record-page；withHeaders=是否注册 --headers。
func newStockPagedCmd(use, apiPath, endpoint, short string, filters []cmdutil.FlagSpec, withHeaders bool) *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{"pageNo": pageNo, "pageSize": pageSize}
			collectFilterFlags(cmd, params, filters)
			if withHeaders {
				cmdutil.CollectStringFlag(cmd, params, "headers")
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), apiPath+"/"+endpoint, params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("%s失败: %s", short, err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	addFilterFlags(c, filters)
	if withHeaders {
		c.Flags().String("headers", "", "自定义导出表头")
	}
	return c
}

// newStockCountCmd 构建统计命令：GET {apiPath}/{endpoint}，输出 OutputJSON。
func newStockCountCmd(use, apiPath, endpoint, short string, filters []cmdutil.FlagSpec) *cobra.Command {
	c := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			collectFilterFlags(cmd, params, filters)
			resp, err := cmdutil.GetClient().Get(context.Background(), apiPath+"/"+endpoint, params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("%s失败: %s", short, err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	addFilterFlags(c, filters)
	return c
}
```

- [ ] **Step 4: 运行测试，确认通过**

Run: `go test ./cmd/stock/ -run 'TestNewStock' -v`
Expected: PASS。

- [ ] **Step 5: 重写 `cmd/stock/stock_query.go`**

仅 `init()`、`newStockQueryLedgerCmd` 改动；`newStockQueryGetCmd`、`newStockQueryCountCmd`、`newStockQueryRecordCountCmd` 三个函数体**逐字保留当前实现，本次不改**（get 的三选一、count 的合并 get-count/get-warehouse-count、stock-record-count 的维度键逻辑均不变）。把文件改为：

```go
package stock

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var queryCmd = &cobra.Command{
	Use:  "query",
	Short: "库存查询",
	Long: "查询产品库存：获取、分页、统计、批次明细。",
}

func init() {
	stockCmd.AddCommand(queryCmd)
	queryCmd.AddCommand(newStockQueryGetCmd())
	queryCmd.AddCommand(newStockPagedCmd("page", "/erp/stock", "page", "分页查询产品库存", queryPageFilters, true))
	queryCmd.AddCommand(newStockCountCmd("page-count", "/erp/stock", "page-count", "按筛选统计产品库存数量与总成本", queryPageFilters))
	queryCmd.AddCommand(newStockQueryCountCmd())
	queryCmd.AddCommand(newStockQueryLedgerCmd())
	queryCmd.AddCommand(newStockCountCmd("ledger-count", "/erp/stock", "detail-count", "按筛选统计库存台账数量与总成本", queryPageFilters))
	queryCmd.AddCommand(newStockQueryRecordCountCmd())
}

// newStockQueryGetCmd、newStockQueryCountCmd、newStockQueryRecordCountCmd：
// 保持当前文件中这三个函数的实现不变（逐字保留，本次不修改）。

// newStockQueryLedgerCmd 查询库存台账（/batch-detail）；解析失败时回退原样输出。
func newStockQueryLedgerCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "ledger",
		Short: "查询库存台账（库存行的历史进出明细）",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{"pageNo": pageNo, "pageSize": pageSize}
			collectFilterFlags(cmd, params, queryPageFilters)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock/batch-detail", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询库存台账失败: %s", err), "")
			}
			var paged struct {
				List  json.RawMessage `json:"list"`
				Total int64           `json:"total"`
			}
			if err := json.Unmarshal(resp.Data, &paged); err != nil {
				return cmdutil.OutputJSON(json.RawMessage(resp.Data))
			}
			var list any
			if err := json.Unmarshal(paged.List, &list); err != nil {
				list = []any{}
			}
			return cmdutil.OutputPagedJSON(list, paged.Total, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	addFilterFlags(c, queryPageFilters)
	return c
}
```

- [ ] **Step 6: 重写 `cmd/stock/stock_record.go`**

仅 `init()` 改动；`newRecordGetCmd`、`newRecordCountCmd` 两个函数体**逐字保留当前实现，本次不改**。把文件改为：

```go
package stock

import (
	"github.com/spf13/cobra"
)

var recordCmd = &cobra.Command{
	Use:  "record",
	Short: "库存明细",
	Long: "查询库存明细记录：分页、详情、统计。",
}

func init() {
	stockCmd.AddCommand(recordCmd)
	recordCmd.AddCommand(newStockPagedCmd("page", "/erp/stock-record", "page", "分页查询出入库明细", recordPageFilters, true))
	recordCmd.AddCommand(newStockCountCmd("page-count", "/erp/stock-record", "page-count", "按筛选统计出入库明细数量", recordPageFilters))
	recordCmd.AddCommand(newRecordGetCmd())
	recordCmd.AddCommand(newRecordCountCmd())
	recordCmd.AddCommand(newStockPagedCmd("record-page", "/erp/stock-record", "record-page", "按维度分页查询库存明细", recordDimensionFilters, false))
	recordCmd.AddCommand(newStockCountCmd("total-cost", "/erp/stock-record", "total-cost", "按维度统计库存明细总成本", recordDimensionFilters))
}

// newRecordGetCmd、newRecordCountCmd：保持当前文件中这两个函数的实现不变（逐字保留）。
```

> 注意：`stock_record.go` 重写后不再直接使用 `context`/`encoding/json`/`fmt`/`cmdutil`/`output`，故 import 只留 `cobra`（get/count 函数体保留时若仍用到这些包，则按实际保留对应 import——这两个函数内部调用了 `cmdutil`/`output`/`context`/`fmt`/`json`，所以 import 块应保持原来的五项：`context`、`encoding/json`、`fmt`、`cobra`、`cmdutil`、`output`）。编译时若报 unused import 再按 `go build` 提示增删。

- [ ] **Step 7: 构建 + 结构核对 + 全量测试**

Run: `go build ./...`
Expected: 编译通过（必要时按提示修正 import）。

Run: `go run . stock query page --help` 与 `go run . stock record page --help`
Expected: 各自 flag 与改造前列表一致（query page 含 `--headers`；record page 含 `--headers`；record-page/total-cost 不含）。

Run: `make test`
Expected: PASS（含 `TestNewStock*`）。

```bash
git add cmd/stock/stock_doc.go cmd/stock/stock_doc_test.go cmd/stock/stock_query.go cmd/stock/stock_record.go
git commit -m "refactor(stock/query,record): 共享 FlagSpec + 本地辅助函数消除 flag 重复，分页解析统一"
```

---

### Task 8: 退货迁移到 DocumentCmds + 删除 crud_legacy.go

**Files:**
- Modify: `cmd/purchase/purchase_return.go`
- Modify: `cmd/sale/sale_return.go`
- Delete: `internal/cmdutil/crud_legacy.go`

**Interfaces:**
- Consumes: `DocumentCmds`、`DocumentConfig`（Task 1）。
- 退货现状（legacy `NewCRUDSubCmd`）：list/get/page-count/create/update/delete/update-status，list 带 7 个通用 flag（warehouse-id/product-id/no/status/start-time/end-time/type，单值）。迁移后行为保持：`TimeKey=""`（不转数组，保留单值），`PageAlias=false`（退货原本就是 `list`）。

> 退货的真实业务筛选字段本次不臆改，沿用 legacy 的 7 个通用字段以保持行为不变（设计已标注：后续可按后端/前端核对再补）。

- [ ] **Step 1: 重写 `cmd/purchase/purchase_return.go`**

```go
package purchase

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	cmdutil.DocumentCmds(purchaseCmd,
		cmdutil.DocumentConfig{
			Name:    "return",
			APIPath: "/erp/purchase-return",
			Label:   "采购退货",
			Filters: []cmdutil.FlagSpec{
				{Name: "warehouse-id", Usage: "仓库 ID"},
				{Name: "product-id", Usage: "产品 ID"},
				{Name: "no", Usage: "单号"},
				{Name: "status", Usage: "状态"},
				{Name: "start-time", Usage: "开始时间"},
				{Name: "end-time", Usage: "结束时间"},
				{Name: "type", Usage: "类型"},
			},
		},
	)
}
```

- [ ] **Step 2: 重写 `cmd/sale/sale_return.go`**

```go
package sale

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	cmdutil.DocumentCmds(saleCmd,
		cmdutil.DocumentConfig{
			Name:    "return",
			APIPath: "/erp/sale-return",
			Label:   "销售退货",
			Filters: []cmdutil.FlagSpec{
				{Name: "warehouse-id", Usage: "仓库 ID"},
				{Name: "product-id", Usage: "产品 ID"},
				{Name: "no", Usage: "单号"},
				{Name: "status", Usage: "状态"},
				{Name: "start-time", Usage: "开始时间"},
				{Name: "end-time", Usage: "结束时间"},
				{Name: "type", Usage: "类型"},
			},
		},
	)
}
```

- [ ] **Step 3: 确认无残留 legacy 调用**

Run: `grep -rn "NewCRUDSubCmd\|AddLegacyOptionalFlags\|newLegacy" --include=*.go .`
Expected: 仅命中 `internal/cmdutil/crud_legacy.go` 自身（无其他调用方）。

- [ ] **Step 4: 删除 `internal/cmdutil/crud_legacy.go`**

```bash
git rm internal/cmdutil/crud_legacy.go
```

- [ ] **Step 5: 构建并核对退货结构**

Run: `go build ./...`
Expected: 编译通过。

Run: `go run . purchase return list --help` 与 `go run . sale return list --help`
Expected: flag 含 7 个通用筛选 + `--page-no --page-size`；子命令含 list/page-count/get/create/update/delete/update-status。

- [ ] **Step 6: 全量测试并提交**

Run: `make test` → PASS。

```bash
git add cmd/purchase/purchase_return.go cmd/sale/sale_return.go
git commit -m "refactor(purchase,sale/return): 迁移到 DocumentCmds，删除 crud_legacy.go"
```

---

### Task 9: 修正 CLAUDE.md 与 skills 文档

**Files:**
- Modify: `CLAUDE.md`
- Modify: `skills/references/stock.md`
- Modify: `skills/references/sale-purchase.md`（如涉及）

- [ ] **Step 1: 修正 `CLAUDE.md` 的 CRUD 模式段**

找到「架构约定 → CRUD 模式」与「项目结构 → cmdutil」两处。

CRUD 模式段，把：
```
- **CRUD 模式**：
  - 现代：`cmdutil.AddCRUDToParent(parent, CRUDConfig{...})` — 大多数域使用
  - 遗留：`cmdutil.NewCRUDSubCmd(name, apiPath, label)` — 库存/采购/销售子域使用
```
改为：
```
- **CRUD 模式**：
  - 通用单据：`cmdutil.DocumentCmds(parent, DocumentConfig{...})` — 单据类（采购入库/销售出库/其他入库/出库/调拨/盘点/退货）
  - 标准 CRUD：`cmdutil.AddCRUDToParent(parent, CRUDConfig{...})` — 大多数域使用
```

cmdutil 结构段，删除这一行（legacy 已删）：
```
│   ├── crud_legacy.go  #   NewCRUDSubCmd（库存/采购/销售子域通用 CRUD）
```

- [ ] **Step 2: 更新 `skills/references/stock.md` 的分页子命令**

把文档里 `wlt stock in/out/move/check page` 四行表格的命令名改为 `list`（并在表头或脚注注明 `page` 仍可用作别名），时间 flag 由 `--in-time`/`--out-time`/`--move-time`/`--check-time` 改为 `--start-time`、`--end-time`，保留 `--headers`。示例（in 行）：

把
```
| `wlt stock in page` | 分页查询入库单 | `--page-no`, `--page-size`, `--no`, ..., `--in-time`, ..., `--headers` |
```
改为
```
| `wlt stock in list`（`page` 别名） | 分页查询入库单 | `--page-no`, `--page-size`, `--no`, ..., `--start-time`, `--end-time`, ..., `--headers` |
```
out/move/check 三行同理（out 用 `--start-time`/`--end-time`，move/check 同理）。

> `stock query page` / `stock record page` 等行**不改**（query/record 子命令名未变、时间 flag 未变）。

文末示例行：
```
wlt stock in page --status 0                       # 查询待审核
```
改为：
```
wlt in=list: wlt stock in list --status 0          # 查询待审核（page 仍可用作别名）
```

- [ ] **Step 3: 更新 `skills/references/sale-purchase.md`（仅 `--headers` 补充）**

`purchase in` / `sale out` 现支持 `--headers`。在该文件对应的采购入库/销售出库分页命令说明处，注明可加 `--headers`（若该文件未逐条列 flag，则无需改动）。

- [ ] **Step 4: 提交**

```bash
git add CLAUDE.md skills/references/stock.md skills/references/sale-purchase.md
git commit -m "docs(stock): 修正 CRUD 模式描述，stock 单据 page→list(别名)、时间 flag、headers"
```

---
