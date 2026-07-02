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
	stockCmd.AddCommand(newStockOutCmd())
}

func newStockOutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "out",
		Short: "出库单管理",
		Long:  "出库单操作：分页、详情、创建、更新、删除、更新状态。",
	}
	cmd.AddCommand(
		newStockOutPageCmd(),
		newStockOutPageCountCmd(),
		newStockOutGetCmd(),
		newStockOutCreateCmd(),
		newStockOutUpdateCmd(),
		newStockOutDeleteCmd(),
		newStockOutUpdateStatusCmd(),
	)
	return cmd
}

func newStockOutPageCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "page",
		Short: "分页查询出库单",
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
				"customer-id",
				"customer-name",
				"out-time",
				"status",
				"remark",
				"creator",
				"product-id",
				"product-name",
				"warehouse-id",
				"warehouse-name",
				"metrics-name",
				"creator-name",
				"user-id",
				"receive-address",
				"send-address",
				"batch-no",
				"create-time",
				"updater-name",
				"update-time",
				"custom-order",
				"keyword",
				"headers",
			)

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock-out/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询出库单失败: %s", err), "")
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
	c.Flags().String("no", "", "出库单号")
	c.Flags().String("customer-id", "", "客户编号")
	c.Flags().String("customer-name", "", "客户")
	c.Flags().String("out-time", "", "出库时间")
	c.Flags().String("status", "", "状态")
	c.Flags().String("remark", "", "备注")
	c.Flags().String("creator", "", "创建者")
	c.Flags().String("product-id", "", "产品编号")
	c.Flags().String("product-name", "", "产品名")
	c.Flags().String("warehouse-id", "", "仓库编号")
	c.Flags().String("warehouse-name", "", "仓库名")
	c.Flags().String("metrics-name", "", "产品指标名称")
	c.Flags().String("creator-name", "", "创建人名称")
	c.Flags().String("user-id", "", "业务人员ID")
	c.Flags().String("receive-address", "", "收货地址")
	c.Flags().String("send-address", "", "发货地址")
	c.Flags().String("batch-no", "", "批次号")
	c.Flags().String("create-time", "", "创建时间")
	c.Flags().String("updater-name", "", "更新人名称")
	c.Flags().String("update-time", "", "更新时间")
	c.Flags().String("custom-order", "", "前端自定义排序规则")
	c.Flags().String("keyword", "", "关键字")
	c.Flags().String("headers", "", "自定义导出表头")
	return c
}

func newStockOutPageCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "page-count",
		Short: "按筛选统计出库单数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params,
				"no",
				"customer-id",
				"customer-name",
				"out-time",
				"status",
				"remark",
				"creator",
				"product-id",
				"product-name",
				"warehouse-id",
				"warehouse-name",
				"metrics-name",
				"creator-name",
				"user-id",
				"receive-address",
				"send-address",
				"batch-no",
				"create-time",
				"updater-name",
				"update-time",
				"custom-order",
				"keyword",
			)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock-out/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计出库单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().String("no", "", "出库单号")
	c.Flags().String("customer-id", "", "客户编号")
	c.Flags().String("customer-name", "", "客户")
	c.Flags().String("out-time", "", "出库时间")
	c.Flags().String("status", "", "状态")
	c.Flags().String("remark", "", "备注")
	c.Flags().String("creator", "", "创建者")
	c.Flags().String("product-id", "", "产品编号")
	c.Flags().String("product-name", "", "产品名")
	c.Flags().String("warehouse-id", "", "仓库编号")
	c.Flags().String("warehouse-name", "", "仓库名")
	c.Flags().String("metrics-name", "", "产品指标名称")
	c.Flags().String("creator-name", "", "创建人名称")
	c.Flags().String("user-id", "", "业务人员ID")
	c.Flags().String("receive-address", "", "收货地址")
	c.Flags().String("send-address", "", "发货地址")
	c.Flags().String("batch-no", "", "批次号")
	c.Flags().String("create-time", "", "创建时间")
	c.Flags().String("updater-name", "", "更新人名称")
	c.Flags().String("update-time", "", "更新时间")
	c.Flags().String("custom-order", "", "前端自定义排序规则")
	c.Flags().String("keyword", "", "关键字")
	return c
}

func newStockOutGetCmd() *cobra.Command {
	var id int64
	c := &cobra.Command{
		Use:   "get",
		Short: "获取出库单详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock-out/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取出库单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "出库单 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

func newStockOutCreateCmd() *cobra.Command {
	var data string
	c := &cobra.Command{
		Use:   "create",
		Short: "创建出库单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/stock-out/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建出库单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

func newStockOutUpdateCmd() *cobra.Command {
	var data string
	c := &cobra.Command{
		Use:   "update",
		Short: "更新出库单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/stock-out/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新出库单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

func newStockOutDeleteCmd() *cobra.Command {
	var ids string
	c := &cobra.Command{
		Use:   "delete",
		Short: "删除出库单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/stock-out/delete", map[string]any{"ids": ids})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("删除出库单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&ids, "ids", "", "ID 列表（逗号分隔）")
	_ = c.MarkFlagRequired("ids")
	return c
}

func newStockOutUpdateStatusCmd() *cobra.Command {
	var data string
	c := &cobra.Command{
		Use:   "update-status",
		Short: "更新出库单状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/stock-out/update-status", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新出库单状态失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据（含 id 和 status）")
	_ = c.MarkFlagRequired("data")
	return c
}
