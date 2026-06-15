package cmdutil

import (
	"context"
	"encoding/json"
	"fmt"

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
