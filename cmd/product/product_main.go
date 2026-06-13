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
	// Standard CRUD directly on product command
	cmdutil.AddCRUDToParent(productCmd, cmdutil.CRUDConfig{
		APIPath:      "/erp/product",
		Label:        "产品",
		SingleDelete: true,
		ListFilters: []cmdutil.FlagSpec{
			{Name: "name", Usage: "产品名称"},
			{Name: "category-id", Usage: "分类 ID"},
			{Name: "status", Usage: "状态"},
			{Name: "bar-code", Usage: "条码"},
			{Name: "spec", Usage: "规格"},
		},
	})

	// Extra commands
	productCmd.AddCommand(cmdutil.CrudSimpleListCmd("/erp/product", "产品"))
	productCmd.AddCommand(newProductGetMetricsCmd())
	productCmd.AddCommand(newProductGetHistoryMetricsCmd())
}

func newProductGetMetricsCmd() *cobra.Command {
	var id int64
	c := &cobra.Command{
		Use:   "get-metrics",
		Short: "获取产品指标列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/product/get-metrics", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取产品指标失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "产品 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

func newProductGetHistoryMetricsCmd() *cobra.Command {
	var id int64
	var name string
	c := &cobra.Command{
		Use:   "get-history-metrics",
		Short: "获取产品历史指标",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{"id": id}
			if name != "" {
				params["name"] = name
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/product/get-history-metrics", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取历史指标失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "产品 ID")
	c.Flags().StringVar(&name, "name", "", "指标名称")
	_ = c.MarkFlagRequired("id")
	return c
}
