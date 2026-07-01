package cmdutil

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/output"
)

// FlagSpec defines a filter flag for list commands.
type FlagSpec struct {
	Name  string // kebab-case flag name (e.g. "status", "category-id")
	Usage string // help text
}

// CRUDConfig configures a standard CRUD command group.
type CRUDConfig struct {
	Name         string     // subcommand name (empty = add directly to parent)
	APIPath      string     // base API path (e.g. "/erp/product")
	Label        string     // Chinese label (e.g. "产品")
	ListFilters  []FlagSpec // filter flags for list command
	SingleDelete bool       // true: delete with ?id=; false: delete with ?ids=
	SkipStatus   bool       // true: omit update-status subcommand
}

// NewCRUDGroup creates a named sub-command with standard CRUD operations.
func NewCRUDGroup(cfg CRUDConfig) *cobra.Command {
	parent := &cobra.Command{
		Use:   cfg.Name,
		Short: fmt.Sprintf("%s管理", cfg.Label),
	}
	AddCRUDToParent(parent, cfg)
	return parent
}

// AddCRUDToParent adds CRUD subcommands directly to a parent command.
func AddCRUDToParent(parent *cobra.Command, cfg CRUDConfig) {
	parent.AddCommand(
		CrudListCmd(cfg.APIPath, cfg.Label, cfg.ListFilters),
		CrudGetCmd(cfg.APIPath, cfg.Label),
		CrudCreateCmd(cfg.APIPath, cfg.Label),
		CrudUpdateCmd(cfg.APIPath, cfg.Label),
		CrudDeleteCmd(cfg.APIPath, cfg.Label, cfg.SingleDelete),
	)
	if !cfg.SkipStatus {
		parent.AddCommand(CrudUpdateStatusCmd(cfg.APIPath, cfg.Label))
	}
}

// ---- Paginated list ----

// CrudListCmd builds a paginated list command. Kept for backward compatibility.
func CrudListCmd(apiPath, label string, filters []FlagSpec) *cobra.Command {
	return CrudListCmdWithFixed(apiPath, label, filters, nil)
}

// CrudListCmdWithFixed is like CrudListCmd but merges fixed query params into the
// request. Use it when a shared backend endpoint requires a discriminator the
// caller always knows (e.g. contract type, business-partner type) that should
// never be exposed as a user flag.
func CrudListCmdWithFixed(apiPath, label string, filters []FlagSpec, fixed map[string]any) *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "list",
		Short: fmt.Sprintf("分页查询%s", label),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			for k, v := range fixed {
				params[k] = v
			}
			for _, f := range filters {
				CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := GetClient().Get(context.Background(), apiPath+"/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询%s失败: %s", label, err), "")
			}
			return ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	for _, f := range filters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}

// ---- Get by ID ----

func CrudGetCmd(apiPath, label string) *cobra.Command {
	var id int64
	c := &cobra.Command{
		Use:   "get",
		Short: fmt.Sprintf("获取%s详情", label),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := EnsureClient(); err != nil {
				return err
			}
			resp, err := GetClient().Get(context.Background(), apiPath+"/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取%s失败: %s", label, err), "")
			}
			return OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, fmt.Sprintf("%s ID", label))
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- Create ----

func CrudCreateCmd(apiPath, label string) *cobra.Command {
	var data string
	c := &cobra.Command{
		Use:   "create",
		Short: fmt.Sprintf("创建%s", label),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := EnsureClient(); err != nil {
				return err
			}
			body, err := ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := GetClient().Post(context.Background(), apiPath+"/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建%s失败: %s", label, err), "")
			}
			return OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- Update ----

func CrudUpdateCmd(apiPath, label string) *cobra.Command {
	var data string
	c := &cobra.Command{
		Use:   "update",
		Short: fmt.Sprintf("更新%s", label),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := EnsureClient(); err != nil {
				return err
			}
			body, err := ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := GetClient().Put(context.Background(), apiPath+"/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新%s失败: %s", label, err), "")
			}
			return OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- Delete ----

func CrudDeleteCmd(apiPath, label string, single bool) *cobra.Command {
	if single {
		var id int64
		c := &cobra.Command{
			Use:   "delete",
			Short: fmt.Sprintf("删除%s", label),
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := EnsureClient(); err != nil {
					return err
				}
				resp, err := GetClient().Delete(context.Background(), apiPath+"/delete", map[string]any{"id": id})
				if err != nil {
					return output.NewExitError(5, fmt.Sprintf("删除%s失败: %s", label, err), "")
				}
				return OutputJSON(json.RawMessage(resp.Data))
			},
		}
		c.Flags().Int64Var(&id, "id", 0, fmt.Sprintf("%s ID", label))
		_ = c.MarkFlagRequired("id")
		return c
	}
	var ids string
	c := &cobra.Command{
		Use:   "delete",
		Short: fmt.Sprintf("删除%s", label),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := EnsureClient(); err != nil {
				return err
			}
			resp, err := GetClient().Delete(context.Background(), apiPath+"/delete", map[string]any{"ids": ids})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("删除%s失败: %s", label, err), "")
			}
			return OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&ids, "ids", "", "ID 列表（逗号分隔）")
	_ = c.MarkFlagRequired("ids")
	return c
}

// ---- Update status ----

func CrudUpdateStatusCmd(apiPath, label string) *cobra.Command {
	var data string
	c := &cobra.Command{
		Use:   "update-status",
		Short: fmt.Sprintf("更新%s状态", label),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := EnsureClient(); err != nil {
				return err
			}
			body, err := ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "")
			}
			resp, err := GetClient().Put(context.Background(), apiPath+"/update-status", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新%s状态失败: %s", label, err), "")
			}
			return OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据（含 id 和 status）")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- Simple list (no params) ----

func CrudSimpleListCmd(apiPath, label string) *cobra.Command {
	c := &cobra.Command{
		Use:   "simple-list",
		Short: fmt.Sprintf("获取%s精简列表", label),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := EnsureClient(); err != nil {
				return err
			}
			resp, err := GetClient().Get(context.Background(), apiPath+"/simple-list", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取%s列表失败: %s", label, err), "")
			}
			return OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

// ---- Simple list with query params ----

func CrudSimpleListCmdWithFlags(apiPath, label string, flags []FlagSpec) *cobra.Command {
	c := &cobra.Command{
		Use:   "simple-list",
		Short: fmt.Sprintf("获取%s精简列表", label),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			for _, f := range flags {
				CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := GetClient().Get(context.Background(), apiPath+"/simple-list", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取%s列表失败: %s", label, err), "")
			}
			return OutputJSON(json.RawMessage(resp.Data))
		},
	}
	for _, f := range flags {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}

// ---- Page count ----

// CrudPageCountCmd builds a page-count command. Backward compatible.
func CrudPageCountCmd(apiPath, label string, filters []FlagSpec) *cobra.Command {
	return CrudPageCountCmdWithFixed(apiPath, label, filters, nil)
}

// CrudPageCountCmdWithFixed merges fixed query params, mirroring CrudListCmdWithFixed.
func CrudPageCountCmdWithFixed(apiPath, label string, filters []FlagSpec, fixed map[string]any) *cobra.Command {
	c := &cobra.Command{
		Use:   "page-count",
		Short: fmt.Sprintf("统计%s数量", label),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			for k, v := range fixed {
				params[k] = v
			}
			for _, f := range filters {
				CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := GetClient().Get(context.Background(), apiPath+"/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计%s失败: %s", label, err), "")
			}
			return OutputJSON(json.RawMessage(resp.Data))
		},
	}
	for _, f := range filters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}

// ---- Contract type variants (长协、采购长协 等共享 /erp/contract 端点的合同子类型) ----

// contractDateRangeFlags 是合同列表共用的时间范围筛选 flag。
var contractDateRangeFlags = []FlagSpec{
	{Name: "order-start", Usage: "下单日期起始（如 2026-07-06 00:00:00）"},
	{Name: "order-end", Usage: "下单日期结束（如 2026-08-03 23:59:59）"},
	{Name: "end-start", Usage: "合同到期起始（如 2026-07-13 00:00:00）"},
	{Name: "end-end", Usage: "合同到期结束（如 2026-08-03 23:59:59）"},
}

// CollectContractDateRange 把 order-start/order-end/end-start/end-end 转成
// orderDate[0]/orderDate[1]/endTime[0]/endTime[1] 数组参数。
func CollectContractDateRange(cmd *cobra.Command, params map[string]any) {
	if v, _ := cmd.Flags().GetString("order-start"); v != "" {
		params["orderDate[0]"] = v
	}
	if v, _ := cmd.Flags().GetString("order-end"); v != "" {
		params["orderDate[1]"] = v
	}
	if v, _ := cmd.Flags().GetString("end-start"); v != "" {
		params["endTime[0]"] = v
	}
	if v, _ := cmd.Flags().GetString("end-end"); v != "" {
		params["endTime[1]"] = v
	}
}

// AddContractDateRangeFlags 把时间范围 flag 注册到命令上。
func AddContractDateRangeFlags(c *cobra.Command) {
	for _, f := range contractDateRangeFlags {
		c.Flags().String(f.Name, "", f.Usage)
	}
}

// ContractTypeConfig 描述一种合同子类型（长协、采购长协、销售、业务、运输……），
// 后端统一通过「API 路径 + type」区分。
type ContractTypeConfig struct {
	Type    string // 后端 type 枚举值，如 "LONG" / "PURCHASE_LONG_COOPERATE" / "SALE_CONTRACT"
	Label   string // 中文标签，如 "长协合同" / "采购长协" / "销售合同"
	Use     string // 自定义 cobra 子命令名（默认由 Type 自动转写：SERVICE→service）
	Filters []FlagSpec
	// APIPath 是合同子类型对应的 HTTP 端点根路径（默认 /erp/contract）。
	// 部分子类型（如销售合同 / SALE_CONTRACT）使用独立端点 /erp/service-contract。
	APIPath string
	// HasDateRange 为 true 时追加 order-start/order-end/end-start/end-end 四个 flag。
	HasDateRange bool
}

// resolveAPIPath 返回配置中的端点路径，未设置时回退到默认 /erp/contract。
func (c ContractTypeConfig) resolveAPIPath() string {
	if c.APIPath != "" {
		return c.APIPath
	}
	return "/erp/contract"
}

// resolveUse 返回 cobra 子命令名；Use 未设置时由 Type 自动转写（SERVICE→service）。
func (c ContractTypeConfig) resolveUse() string {
	if c.Use != "" {
		return c.Use
	}
	return commandSlug(c.Type)
}

// ContractTypeCmds 为一组合同子类型生成标准 CRUD+page-count 子命令，统一注册到 parent。
// 每个子类型自动拥有：list / page-count / get / update-status / create / update / delete。
func ContractTypeCmds(parent *cobra.Command, types ...ContractTypeConfig) {
	for _, tc := range types {
		apiPath := tc.resolveAPIPath()
		filters := tc.Filters
		cmd := &cobra.Command{
			Use:   tc.resolveUse(),
			Short: tc.Label + "管理",
		}

		// list
		if tc.HasDateRange {
			cmd.AddCommand(contractListCmd(tc, apiPath, filters))
		} else {
			cmd.AddCommand(CrudListCmdWithFixed(apiPath, tc.Label, filters, map[string]any{"type": tc.Type}))
		}

		// page-count
		if tc.HasDateRange {
			cmd.AddCommand(contractPageCountCmd(tc, apiPath, filters))
		} else {
			cmd.AddCommand(CrudPageCountCmdWithFixed(apiPath, tc.Label, filters, map[string]any{"type": tc.Type}))
		}

		// get / update-status / create / update / delete
		cmd.AddCommand(CrudGetCmd(apiPath, tc.Label))
		cmd.AddCommand(CrudUpdateStatusCmd(apiPath, tc.Label))
		cmd.AddCommand(CrudCreateCmd(apiPath, tc.Label))
		cmd.AddCommand(CrudUpdateCmd(apiPath, tc.Label))
		cmd.AddCommand(CrudDeleteCmd(apiPath, tc.Label, false))

		parent.AddCommand(cmd)
	}
}

// contractListCmd 构建带日期范围数组参数的 list 命令。
func contractListCmd(tc ContractTypeConfig, apiPath string, filters []FlagSpec) *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询" + tc.Label,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"type":     tc.Type,
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			for _, f := range filters {
				CollectStringFlag(cmd, params, f.Name)
			}
			CollectContractDateRange(cmd, params)
			resp, err := GetClient().Get(context.Background(), apiPath+"/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询%s失败: %s", tc.Label, err), "")
			}
			return ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	for _, f := range filters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	AddContractDateRangeFlags(c)
	return c
}

// contractPageCountCmd 构建带日期范围数组参数的 page-count 命令。
func contractPageCountCmd(tc ContractTypeConfig, apiPath string, filters []FlagSpec) *cobra.Command {
	c := &cobra.Command{
		Use:   "page-count",
		Short: "统计" + tc.Label + "数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{"type": tc.Type}
			for _, f := range filters {
				CollectStringFlag(cmd, params, f.Name)
			}
			CollectContractDateRange(cmd, params)
			resp, err := GetClient().Get(context.Background(), apiPath+"/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计%s失败: %s", tc.Label, err), "")
			}
			return OutputJSON(json.RawMessage(resp.Data))
		},
	}
	for _, f := range filters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	AddContractDateRangeFlags(c)
	return c
}

// commandSlug 把类型枚举值（如 LONG / PURCHASE_LONG_COOPERATE）转成
// cobra 子命令名（如 long / purchase-long-cooperate）。
func commandSlug(t string) string {
	s := strings.ToLower(t)
	s = strings.ReplaceAll(s, "_", "-")
	return s
}

// ---- Non-paginated list (uses /list endpoint) ----

func CrudListAllCmd(apiPath, label string, filters []FlagSpec) *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: fmt.Sprintf("查询所有%s", label),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			for _, f := range filters {
				CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := GetClient().Get(context.Background(), apiPath+"/list", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询%s失败: %s", label, err), "")
			}
			return OutputJSON(json.RawMessage(resp.Data))
		},
	}
	for _, f := range filters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}
