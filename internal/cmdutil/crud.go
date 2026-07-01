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
			CollectPlanDateRange(cmd, params)
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
	AddPlanDateRangeFlags(c)
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
			CollectPlanDateRange(cmd, params)
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
	AddPlanDateRangeFlags(c)
	return c
}

// commandSlug 把类型枚举值（如 LONG / PURCHASE_LONG_COOPERATE）转成
// cobra 子命令名（如 long / purchase-long-cooperate）。
func commandSlug(t string) string {
	s := strings.ToLower(t)
	s = strings.ReplaceAll(s, "_", "-")
	return s
}

// ---- 计划类型变体（采购运输计划 / 销售运输计划 等共享 /erp/order-plan 端点的子类型） ----

// planPartnerFilters 是计划列表共用的业务伙伴筛选字段（采购用 supplier-id，销售用 customer-id）。
var planPartnerFilters = []FlagSpec{
	{Name: "supplier-id", Usage: "供应商 ID"},
	{Name: "customer-id", Usage: "客户 ID"},
}

// CollectPlanPartnerFilters 收集采购/销售的 supplierId / customerId 到 params。
func CollectPlanPartnerFilters(cmd *cobra.Command, params map[string]any) {
	CollectStringFlags(cmd, params, "supplier-id", "customer-id")
}

// planDateRangeFlags 是订单计划共用的起止日期筛选 flag（转 startDate[0] / startDate[1]）。
var planDateRangeFlags = []FlagSpec{
	{Name: "start", Usage: "计划开始日期起始（如 2026-07-21 00:00:00）"},
	{Name: "end", Usage: "计划开始日期结束（如 2026-08-19 23:59:59）"},
}

// CollectPlanDateRange 把 start/end 转成 startDate[0] / startDate[1] 数组参数。
func CollectPlanDateRange(cmd *cobra.Command, params map[string]any) {
	if v, _ := cmd.Flags().GetString("start"); v != "" {
		params["startDate[0]"] = v
	}
	if v, _ := cmd.Flags().GetString("end"); v != "" {
		params["startDate[1]"] = v
	}
}

// AddPlanDateRangeFlags 把起止日期 flag 注册到命令上。
func AddPlanDateRangeFlags(c *cobra.Command) {
	for _, f := range planDateRangeFlags {
		c.Flags().String(f.Name, "", f.Usage)
	}
}

// PlanListCmdConfig 描述一个计划子类型（采购运输计划 / 销售运输计划 ……），
// 后端统一通过 /erp/order-plan 端点 + type 区分。
type PlanListCmdConfig struct {
	Type    string // 后端 type 枚举值，如 "PURCHASE_TRANSPORT_PLAN"
	Label   string // 中文标签，如 "采购计划"
	Filters []FlagSpec
	Use     string // 自定义 cobra 子命令名（默认由 Type 自动转写）
}

// PlanListCmds 为一组计划子类型生成 list + page-count 子命令，统一注册到 parent。
func PlanListCmds(parent *cobra.Command, types ...PlanListCmdConfig) {
	for _, tc := range types {
		filters := tc.Filters
		cmd := &cobra.Command{
			Use:   planCommandSlug(tc.Use, tc.Type),
			Short: tc.Label + "管理",
		}

		// list
		{
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
					CollectStringFlags(cmd, params, "product-id")
					CollectPlanPartnerFilters(cmd, params)
					CollectPlanDateRange(cmd, params)
					resp, err := GetClient().Get(context.Background(), "/erp/order-plan/page", params)
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
			c.Flags().String("product-id", "", "产品 ID")
			c.Flags().String("supplier-id", "", "供应商 ID")
			c.Flags().String("customer-id", "", "客户 ID")
			AddPlanDateRangeFlags(c)
			cmd.AddCommand(c)
		}

		// page-count：复用 list 的 filters（含 product-id / supplier-id / 日期范围）。
		// page-count 不需要 customer-id（后端校验严格，销售计划带 customerId 到 page-count 更稳妥，这里都带上）。
		{
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
					CollectStringFlags(cmd, params, "product-id")
					CollectPlanPartnerFilters(cmd, params)
					CollectPlanDateRange(cmd, params)
					resp, err := GetClient().Get(context.Background(), "/erp/order-plan/page-count", params)
					if err != nil {
						return output.NewExitError(5, fmt.Sprintf("统计%s失败: %s", tc.Label, err), "")
					}
					return OutputJSON(json.RawMessage(resp.Data))
				},
			}
			for _, f := range filters {
				c.Flags().String(f.Name, "", f.Usage)
			}
			c.Flags().String("product-id", "", "产品 ID")
			c.Flags().String("supplier-id", "", "供应商 ID")
			c.Flags().String("customer-id", "", "客户 ID")
			AddPlanDateRangeFlags(c)
			cmd.AddCommand(c)
		}

		parent.AddCommand(cmd)
	}
}

// planCommandSlug 把计划类型枚举值（如 PURCHASE_TRANSPORT_PLAN）转成 cobra 子命令名（如 purchase）。
// 规则：小写后去掉末尾的 _plan / _transport，取第一个 _ 之前的段。
// 例：PURCHASE_TRANSPORT_PLAN → purchase，SALE_TRANSPORT_PLAN → sale。
// Use 显式指定时由调用方控制。
func planCommandSlug(use, t string) string {
	if use != "" {
		return use
	}
	s := strings.ToLower(t)
	s = strings.ReplaceAll(s, "_transport", "")
	s = strings.TrimSuffix(s, "_plan")
	if idx := strings.Index(s, "_"); idx >= 0 {
		s = s[:idx]
	}
	return s
}

// ---- 订单类型变体（采购订单 / 销售订单 共享 /erp/order 端点的子类型） ----

// orderCommonFilters 是订单列表共用的筛选字段。
var orderCommonFilters = []FlagSpec{
	{Name: "no", Usage: "订单号"},
	{Name: "enterprise-id", Usage: "企业 ID"},
	{Name: "product-id", Usage: "产品 ID"},
}

// orderDateRangeFlags 是订单列表共用的下单时间范围筛选 flag（转 orderTime[0] / orderTime[1]）。
var orderDateRangeFlags = []FlagSpec{
	{Name: "order-start", Usage: "下单时间起始（如 2026-07-14 00:00:00）"},
	{Name: "order-end", Usage: "下单时间结束（如 2026-08-11 23:59:59）"},
}

// CollectOrderDateRange 把 order-start/order-end 转成 orderTime[0] / orderTime[1] 数组参数。
func CollectOrderDateRange(cmd *cobra.Command, params map[string]any) {
	if v, _ := cmd.Flags().GetString("order-start"); v != "" {
		params["orderTime[0]"] = v
	}
	if v, _ := cmd.Flags().GetString("order-end"); v != "" {
		params["orderTime[1]"] = v
	}
}

// AddOrderDateRangeFlags 把时间范围 flag 注册到命令上。
func AddOrderDateRangeFlags(c *cobra.Command) {
	for _, f := range orderDateRangeFlags {
		c.Flags().String(f.Name, "", f.Usage)
	}
}

// OrderTypeConfig 描述一个订单子类型（采购订单 / 销售订单 ……），后端统一通过
// /erp/order 端点 + type 区分。
type OrderTypeConfig struct {
	Type    string // 后端 type 枚举值，如 "PURCHASE" / "SALE"
	Label   string // 中文标签，如 "采购订单"
	Use     string // 自定义 cobra 子命令名（默认由 Type 自动转写：PURCHASE → purchase）
	Filters []FlagSpec
}

// OrderTypeCmds 为一组订单子类型生成 list + page-count 子命令，统一注册到 parent。
// 每个子类型自动拥有：list / page-count。
func OrderTypeCmds(parent *cobra.Command, types ...OrderTypeConfig) {
	for _, tc := range types {
		filters := tc.Filters
		cmd := &cobra.Command{
			Use:   orderCommandSlug(tc.Use, tc.Type),
			Short: tc.Label + "管理",
		}

		// list
		{
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
					CollectOrderDateRange(cmd, params)
					resp, err := GetClient().Get(context.Background(), "/erp/order/page", params)
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
			AddOrderDateRangeFlags(c)
			cmd.AddCommand(c)
		}

		// page-count
		{
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
					CollectOrderDateRange(cmd, params)
					resp, err := GetClient().Get(context.Background(), "/erp/order/page-count", params)
					if err != nil {
						return output.NewExitError(5, fmt.Sprintf("统计%s失败: %s", tc.Label, err), "")
					}
					return OutputJSON(json.RawMessage(resp.Data))
				},
			}
			for _, f := range filters {
				c.Flags().String(f.Name, "", f.Usage)
			}
			AddOrderDateRangeFlags(c)
			cmd.AddCommand(c)
		}

		parent.AddCommand(cmd)
	}
}

// orderCommandSlug 把订单类型枚举值（如 PURCHASE / SALE）转成 cobra 子命令名（小写形式）。
// Use 显式指定时由调用方控制。
func orderCommandSlug(use, t string) string {
	if use != "" {
		return use
	}
	return strings.ToLower(t)
}

// ---- 采购入库 / 销售出库 变体（含时间范围数组参数 + 模块特有筛选字段） ----

// salePurchaseTimeFlags 是采购/销售共用的起止时间筛选 flag；
// 实际查询时会被转成 {timeKey}[0] / {timeKey}[1] 数组参数
// （采购入库为 inTime[0]/inTime[1]，销售出库为 outTime[0]/outTime[1]）。
var salePurchaseTimeFlags = []FlagSpec{
	{Name: "start-time", Usage: "开始时间（如 2026-07-01 00:00:00）"},
	{Name: "end-time", Usage: "结束时间（如 2026-07-31 23:59:59）"},
}

// collectSalePurchaseTime 把 start-time/end-time 转成 {timeKey}[0]/{timeKey}[1]。
func collectSalePurchaseTime(cmd *cobra.Command, params map[string]any, timeKey string) {
	if v, _ := cmd.Flags().GetString("start-time"); v != "" {
		params[fmt.Sprintf("%s[0]", timeKey)] = v
	}
	if v, _ := cmd.Flags().GetString("end-time"); v != "" {
		params[fmt.Sprintf("%s[1]", timeKey)] = v
	}
}

func addSalePurchaseTimeFlags(c *cobra.Command) {
	for _, f := range salePurchaseTimeFlags {
		c.Flags().String(f.Name, "", f.Usage)
	}
}

// ---- 结算模块时间范围（settlementDate[0] / settlementDate[1]） ----

// settlementDateFlags 是结算模块共用的结算日期范围筛选 flag；
// 实际查询时会被转成 settlementDate[0] / settlementDate[1] 数组参数。
var settlementDateFlags = []FlagSpec{
	{Name: "start-date", Usage: "结算日期起始（如 2026-07-01 00:00:00）"},
	{Name: "end-date", Usage: "结算日期结束（如 2026-07-31 23:59:59）"},
}

// CollectSettlementDate 把 start-date/end-date 转成 settlementDate[0]/settlementDate[1]。
func CollectSettlementDate(cmd *cobra.Command, params map[string]any) {
	if v, _ := cmd.Flags().GetString("start-date"); v != "" {
		params["settlementDate[0]"] = v
	}
	if v, _ := cmd.Flags().GetString("end-date"); v != "" {
		params["settlementDate[1]"] = v
	}
}

// AddSettlementDateFlags 把结算日期范围 flag 注册到命令上。
func AddSettlementDateFlags(c *cobra.Command) {
	for _, f := range settlementDateFlags {
		c.Flags().String(f.Name, "", f.Usage)
	}
}

// ---- 质检称重模块时间范围（orderTime[0] / orderTime[1]） ----

// orderTimeFlags 是质检称重共用的订单时间范围筛选 flag；
// 实际查询时会被转成 orderTime[0] / orderTime[1] 数组参数。
// 对 SALE / PURCHASE 的订单称重与运单称重通用。
var orderTimeFlags = []FlagSpec{
	{Name: "start-time", Usage: "订单时间起始（如 2026-07-01 00:00:00）"},
	{Name: "end-time", Usage: "订单时间结束（如 2026-07-31 23:59:59）"},
}

// CollectOrderTime 把 start-time/end-time 转成 orderTime[0]/orderTime[1] 数组参数。
func CollectOrderTime(cmd *cobra.Command, params map[string]any) {
	if v, _ := cmd.Flags().GetString("start-time"); v != "" {
		params["orderTime[0]"] = v
	}
	if v, _ := cmd.Flags().GetString("end-time"); v != "" {
		params["orderTime[1]"] = v
	}
}

// AddOrderTimeFlags 把订单时间范围 flag 注册到命令上。
func AddOrderTimeFlags(c *cobra.Command) {
	for _, f := range orderTimeFlags {
		c.Flags().String(f.Name, "", f.Usage)
	}
}

// addSalePurchaseBaseFlags 注册采购/销售 list/page-count 共用的筛选 flag。
// 保留 legacy 已有字段（含 start-time/end-time），并新增 product-name（产品名称模糊搜索）。
func addSalePurchaseBaseFlags(c *cobra.Command) {
	c.Flags().String("warehouse-id", "", "仓库 ID")
	c.Flags().String("product-id", "", "产品 ID")
	c.Flags().String("product-name", "", "产品名称（模糊搜索）")
	c.Flags().String("no", "", "单号")
	c.Flags().String("status", "", "状态")
	c.Flags().String("type", "", "类型")
}

// SalePurchaseConfig 描述一个采购入库或销售出库子模块。
type SalePurchaseConfig struct {
	Name    string     // 子命令名（"in" / "out"）
	APIPath string     // 后端根路径（"/erp/purchase-in" / "/erp/sale-out"）
	Label   string     // 中文标签（"采购入库" / "销售出库"）
	TimeKey string     // 时间范围数组参数前缀（"inTime" / "outTime"）
	Filters []FlagSpec // 模块特有筛选字段（supplier-id / metrics-name / batch-no …）
}

// SalePurchaseCmds 为采购/销售子模块生成标准 CRUD 子命令：
// list / page-count / get / create / update / delete / update-status。
// list/page-count 共享 base flags + 时间范围数组 + cfg.Filters。
func SalePurchaseCmds(parent *cobra.Command, cfgs ...SalePurchaseConfig) {
	for _, cfg := range cfgs {
		cmd := &cobra.Command{
			Use:   cfg.Name,
			Short: cfg.Label + "管理",
		}
		cmd.AddCommand(salePurchaseListCmd(cfg))
		cmd.AddCommand(salePurchasePageCountCmd(cfg))
		cmd.AddCommand(CrudGetCmd(cfg.APIPath, cfg.Label))
		cmd.AddCommand(CrudCreateCmd(cfg.APIPath, cfg.Label))
		cmd.AddCommand(CrudUpdateCmd(cfg.APIPath, cfg.Label))
		cmd.AddCommand(CrudDeleteCmd(cfg.APIPath, cfg.Label, false))
		cmd.AddCommand(CrudUpdateStatusCmd(cfg.APIPath, cfg.Label))
		parent.AddCommand(cmd)
	}
}

func salePurchaseListCmd(cfg SalePurchaseConfig) *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询" + cfg.Label,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			// 基本筛选字段
			CollectStringFlags(cmd, params,
				"warehouse-id", "product-id", "product-name",
				"no", "status", "type")
			// 时间范围（数组格式 inTime[0]/inTime[1] 或 outTime[0]/outTime[1]）
			collectSalePurchaseTime(cmd, params, cfg.TimeKey)
			// 模块特有字段
			for _, f := range cfg.Filters {
				CollectStringFlag(cmd, params, f.Name)
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
	addSalePurchaseBaseFlags(c)
	addSalePurchaseTimeFlags(c)
	for _, f := range cfg.Filters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}

func salePurchasePageCountCmd(cfg SalePurchaseConfig) *cobra.Command {
	c := &cobra.Command{
		Use:   "page-count",
		Short: "统计" + cfg.Label + "数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			CollectStringFlags(cmd, params,
				"warehouse-id", "product-id", "product-name",
				"no", "status", "type")
			collectSalePurchaseTime(cmd, params, cfg.TimeKey)
			for _, f := range cfg.Filters {
				CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := GetClient().Get(context.Background(), cfg.APIPath+"/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计%s失败: %s", cfg.Label, err), "")
			}
			return OutputJSON(json.RawMessage(resp.Data))
		},
	}
	addSalePurchaseBaseFlags(c)
	addSalePurchaseTimeFlags(c)
	for _, f := range cfg.Filters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
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
