package stock

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var warehouseCmd = &cobra.Command{
	Use:   "warehouse",
	Short: "仓库管理",
	Long:  "管理 ERP 仓库：查询、创建、更新、删除、状态管理。",
}

func init() {
	stockCmd.AddCommand(warehouseCmd)

	warehouseCmd.AddCommand(newWarehouseListCmd())
	warehouseCmd.AddCommand(newWarehouseGetCmd())
	warehouseCmd.AddCommand(newWarehouseSimpleListCmd())
	warehouseCmd.AddCommand(newWarehouseCreateCmd())
	warehouseCmd.AddCommand(newWarehouseUpdateCmd())
	warehouseCmd.AddCommand(newWarehouseDeleteCmd())
	warehouseCmd.AddCommand(newWarehouseUpdateStatusCmd())
}

func newWarehouseListCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询仓库",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlags(cmd, params, "name", "status", "type")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/warehouse/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询仓库失败: %s", err), "")
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
			return cmdutil.OutputPagedJSON(list, paged.Total, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().String("name", "", "仓库名称（模糊）")
	c.Flags().String("status", "", "状态")
	c.Flags().String("type", "", "类型")
	return c
}

func newWarehouseGetCmd() *cobra.Command {
	var id int64
	c := &cobra.Command{
		Use:   "get",
		Short: "获取仓库详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/warehouse/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取仓库失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "仓库 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

func newWarehouseSimpleListCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "simple-list",
		Short: "获取仓库精简列表（选择器用）",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/warehouse/simple-list", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取仓库列表失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

func newWarehouseCreateCmd() *cobra.Command {
	var data string
	c := &cobra.Command{
		Use:   "create",
		Short: "创建仓库",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/warehouse/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建仓库失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

func newWarehouseUpdateCmd() *cobra.Command {
	var data string
	c := &cobra.Command{
		Use:   "update",
		Short: "更新仓库",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/warehouse/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新仓库失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

func newWarehouseDeleteCmd() *cobra.Command {
	var ids string
	c := &cobra.Command{
		Use:   "delete",
		Short: "删除仓库",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/warehouse/delete", map[string]any{"ids": ids})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("删除仓库失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&ids, "ids", "", "ID 列表（逗号分隔）")
	_ = c.MarkFlagRequired("ids")
	return c
}

func newWarehouseUpdateStatusCmd() *cobra.Command {
	var data string
	c := &cobra.Command{
		Use:   "update-status",
		Short: "更新仓库状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/warehouse/update-status", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新状态失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据（含 id 和 status）")
	_ = c.MarkFlagRequired("data")
	return c
}
