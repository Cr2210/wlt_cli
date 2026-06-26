package cmdutil

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/output"
)

// NewCRUDSubCmd creates a parent command for a CRUD subdomain
// (e.g. in/out/move/check) with list/get/create/update/delete/update-status subcommands.
// This is the legacy CRUD factory used by stock, purchase, and sale subdomains.
func NewCRUDSubCmd(name, apiPath, label string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("%s管理", label),
	}
	cmd.AddCommand(
		newLegacyListCmd(name, apiPath, label),
		newLegacyGetCmd(name, apiPath, label),
		newLegacyPageCountCmd(name, apiPath, label),
		newLegacyCreateCmd(name, apiPath, label),
		newLegacyUpdateCmd(name, apiPath, label),
		newLegacyDeleteCmd(name, apiPath, label),
		newLegacyUpdateStatusCmd(name, apiPath, label),
	)
	return cmd
}

// newLegacyPageCountCmd builds a page-count command mirroring the legacy list
// filters (warehouse-id/product-id/no/status/start-time/end-time/type).
func newLegacyPageCountCmd(name, apiPath, label string) *cobra.Command {
	c := &cobra.Command{
		Use:   "page-count",
		Short: fmt.Sprintf("统计%s数量", label),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			CollectStringFlags(cmd, params, "warehouse-id", "product-id", "no", "status",
				"start-time", "end-time", "type")
			resp, err := GetClient().Get(context.Background(), apiPath+"/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计%s失败: %s", label, err), "")
			}
			return OutputJSON(json.RawMessage(resp.Data))
		},
	}
	AddLegacyOptionalFlags(c)
	return c
}

func newLegacyListCmd(name, apiPath, label string) *cobra.Command {
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
			CollectStringFlags(cmd, params, "warehouse-id", "product-id", "no", "status",
				"start-time", "end-time", "type")

			resp, err := GetClient().Get(context.Background(), apiPath+"/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询%s失败: %s", label, err), "使用 --dry-run 调试")
			}
			var paged struct {
				List  json.RawMessage `json:"list"`
				Total int64           `json:"total"`
			}
			if err := json.Unmarshal(resp.Data, &paged); err != nil {
				return output.NewExitError(5, fmt.Sprintf("解析响应失败: %s", err), "")
			}
			var list any
			if err := json.Unmarshal(paged.List, &list); err != nil {
				list = []any{}
			}
			return OutputPagedJSON(list, paged.Total, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	AddLegacyOptionalFlags(c)
	return c
}

func newLegacyGetCmd(name, apiPath, label string) *cobra.Command {
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

func newLegacyCreateCmd(name, apiPath, label string) *cobra.Command {
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

func newLegacyUpdateCmd(name, apiPath, label string) *cobra.Command {
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

func newLegacyDeleteCmd(name, apiPath, label string) *cobra.Command {
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

func newLegacyUpdateStatusCmd(name, apiPath, label string) *cobra.Command {
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
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
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

// AddLegacyOptionalFlags adds common optional filter flags to a list command.
func AddLegacyOptionalFlags(cmd *cobra.Command) {
	cmd.Flags().String("warehouse-id", "", "仓库 ID")
	cmd.Flags().String("product-id", "", "产品 ID")
	cmd.Flags().String("no", "", "单号")
	cmd.Flags().String("status", "", "状态")
	cmd.Flags().String("start-time", "", "开始时间")
	cmd.Flags().String("end-time", "", "结束时间")
	cmd.Flags().String("type", "", "类型")
}
