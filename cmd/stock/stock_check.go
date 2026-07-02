package stock

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

func init() {
	stockCmd.AddCommand(newStockCheckCmd())
}

func newStockCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "盘点单管理",
		Long:  "盘点单操作：分页、详情、创建、更新、删除、更新状态。",
	}
	cmd.AddCommand(
		newStockCheckPageCmd(),
		newStockCheckPageCountCmd(),
		newStockCheckGetCmd(),
		newStockCheckCreateCmd(),
		newStockCheckUpdateCmd(),
		newStockCheckDeleteCmd(),
		newStockCheckUpdateStatusCmd(),
	)
	return cmd
}

func newStockCheckPageCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "page",
		Short: "分页查询盘点单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlags(cmd, params,
				"no",
				"warehouse-id",
				"check-time",
				"status",
				"remark",
				"creator",
				"creator-name",
				"create-time",
				"update-time",
				"updater",
				"updater-name",
				"product-id",
				"product-name",
				"metrics-name",
				"batch-no",
				"user-id",
				"custom-order",
				"keyword",
				"headers",
			)

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock-check/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询盘点单失败: %s", err), "")
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
	c.Flags().String("no", "", "盘点单号")
	c.Flags().String("warehouse-id", "", "仓库编号")
	c.Flags().String("check-time", "", "盘点时间")
	c.Flags().String("status", "", "状态")
	c.Flags().String("remark", "", "备注")
	c.Flags().String("creator", "", "创建者")
	c.Flags().String("creator-name", "", "创建者姓名")
	c.Flags().String("create-time", "", "创建时间")
	c.Flags().String("update-time", "", "更新时间")
	c.Flags().String("updater", "", "更新者编号")
	c.Flags().String("updater-name", "", "更新者姓名")
	c.Flags().String("product-id", "", "产品编号")
	c.Flags().String("product-name", "", "产品名称")
	c.Flags().String("metrics-name", "", "产品指标")
	c.Flags().String("batch-no", "", "批次号")
	c.Flags().String("user-id", "", "业务负责人ID")
	c.Flags().String("custom-order", "", "前端自定义排序规则")
	c.Flags().String("keyword", "", "关键字")
	c.Flags().String("headers", "", "自定义导出表头")
	return c
}

func newStockCheckPageCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "page-count",
		Short: "按筛选统计盘点单数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params,
				"no",
				"warehouse-id",
				"check-time",
				"status",
				"remark",
				"creator",
				"creator-name",
				"create-time",
				"update-time",
				"updater",
				"updater-name",
				"product-id",
				"product-name",
				"metrics-name",
				"batch-no",
				"user-id",
				"custom-order",
				"keyword",
			)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock-check/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计盘点单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().String("no", "", "盘点单号")
	c.Flags().String("warehouse-id", "", "仓库编号")
	c.Flags().String("check-time", "", "盘点时间")
	c.Flags().String("status", "", "状态")
	c.Flags().String("remark", "", "备注")
	c.Flags().String("creator", "", "创建者")
	c.Flags().String("creator-name", "", "创建者姓名")
	c.Flags().String("create-time", "", "创建时间")
	c.Flags().String("update-time", "", "更新时间")
	c.Flags().String("updater", "", "更新者编号")
	c.Flags().String("updater-name", "", "更新者姓名")
	c.Flags().String("product-id", "", "产品编号")
	c.Flags().String("product-name", "", "产品名称")
	c.Flags().String("metrics-name", "", "产品指标")
	c.Flags().String("batch-no", "", "批次号")
	c.Flags().String("user-id", "", "业务负责人ID")
	c.Flags().String("custom-order", "", "前端自定义排序规则")
	c.Flags().String("keyword", "", "关键字")
	return c
}

func newStockCheckGetCmd() *cobra.Command {
	var id int64
	c := &cobra.Command{
		Use:   "get",
		Short: "获取盘点单详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock-check/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取盘点单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "盘点单 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

func newStockCheckCreateCmd() *cobra.Command {
	var data string
	c := &cobra.Command{
		Use:   "create",
		Short: "创建盘点单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/stock-check/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建盘点单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

func newStockCheckUpdateCmd() *cobra.Command {
	var data string
	c := &cobra.Command{
		Use:   "update",
		Short: "更新盘点单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/stock-check/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新盘点单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

func newStockCheckDeleteCmd() *cobra.Command {
	var ids string
	c := &cobra.Command{
		Use:   "delete",
		Short: "删除盘点单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/stock-check/delete", map[string]any{"ids": ids})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("删除盘点单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&ids, "ids", "", "ID 列表（逗号分隔）")
	_ = c.MarkFlagRequired("ids")
	return c
}

func newStockCheckUpdateStatusCmd() *cobra.Command {
	var data string
	c := &cobra.Command{
		Use:   "update-status",
		Short: "更新盘点单状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/stock-check/update-status", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新盘点单状态失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据（含 id 和 status）")
	_ = c.MarkFlagRequired("data")
	return c
}
