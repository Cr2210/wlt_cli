package report

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var stockReportFilters = []cmdutil.FlagSpec{
	{Name: "product-id", Usage: "产品 ID"},
	{Name: "category-id", Usage: "分类 ID"},
	{Name: "warehouse-id", Usage: "仓库 ID"},
	{Name: "start-time", Usage: "开始时间"},
	{Name: "end-time", Usage: "结束时间"},
}

func init() {
	stockReportCmd := &cobra.Command{
		Use:   "stock",
		Short: "库存报表管理",
	}
	stockReportCmd.AddCommand(
		newStockReportDetailCmd(),
		newStockReportWarehouseCmd(),
		newStockReportBuySendCmd(),
		newStockReportFinanceCmd(),
		newStockReportProduceCmd(),
	)
	reportCmd.AddCommand(stockReportCmd)
}

func newStockReportDetailCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "detail",
		Short: "库存报表明细查询",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			for _, f := range stockReportFilters {
				cmdutil.CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock-report/detail-page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询库存报表明细失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	for _, f := range stockReportFilters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}

func newStockReportWarehouseCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "warehouse",
		Short: "库存报表仓库查询",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			for _, f := range stockReportFilters {
				cmdutil.CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock-report/warehouse", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询库存报表仓库失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	for _, f := range stockReportFilters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}

func newStockReportBuySendCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "buy-send",
		Short: "库存报表购销查询",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			for _, f := range stockReportFilters {
				cmdutil.CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock-report/buy-send", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询库存报表购销失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	for _, f := range stockReportFilters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}

func newStockReportFinanceCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "finance",
		Short: "库存报表财务查询",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			for _, f := range stockReportFilters {
				cmdutil.CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock-report/finance", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询库存报表财务失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	for _, f := range stockReportFilters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}

func newStockReportProduceCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "produce",
		Short: "库存报表生产查询",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			for _, f := range stockReportFilters {
				cmdutil.CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock-report/produce", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询库存报表生产失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	for _, f := range stockReportFilters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}
