package product

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

func init() {
	metricsCmd := cmdutil.NewCRUDGroup(cmdutil.CRUDConfig{
		Name:         "metrics",
		APIPath:      "/erp/product-metrics",
		Label:        "产品指标",
		SingleDelete: true,
		ListFilters: []cmdutil.FlagSpec{
			{Name: "name", Usage: "指标名称"},
			{Name: "status", Usage: "状态"},
		},
	})
	metricsCmd.AddCommand(cmdutil.CrudSimpleListCmdWithFlags("/erp/product-metrics", "产品指标",
		[]cmdutil.FlagSpec{{Name: "name", Usage: "指标名称"}}))
	metricsCmd.AddCommand(newMetricsOrderItemListCmd())
	metricsCmd.AddCommand(newMetricsItemListCmd())
	metricsCmd.AddCommand(newMetricsAddItemsCmd())
	productCmd.AddCommand(metricsCmd)
}

func newMetricsOrderItemListCmd() *cobra.Command {
	var id int64
	c := &cobra.Command{
		Use:   "order-item-list",
		Short: "获取订单项指标列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/product-metrics/order-item-list", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取订单项指标失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "订单项 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

func newMetricsItemListCmd() *cobra.Command {
	var metricId int64
	c := &cobra.Command{
		Use:   "item-list",
		Short: "获取指标常用值列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/product-metrics/item/list-by-metricId",
				map[string]any{"metricId": metricId})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取指标值列表失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&metricId, "metric-id", 0, "指标 ID")
	_ = c.MarkFlagRequired("metric-id")
	return c
}

func newMetricsAddItemsCmd() *cobra.Command {
	var data string
	c := &cobra.Command{
		Use:   "add-metrics-items",
		Short: "保存订单项指标",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/product-metrics/addMetricsItems", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("保存指标失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}
