package report

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var saleReportFilters = []cmdutil.FlagSpec{
	{Name: "customer-id", Usage: "客户 ID"},
	{Name: "product-id", Usage: "产品 ID"},
	{Name: "category-id", Usage: "分类 ID"},
	{Name: "warehouse-id", Usage: "仓库 ID"},
	{Name: "status", Usage: "状态"},
	{Name: "start-time", Usage: "开始时间"},
	{Name: "end-time", Usage: "结束时间"},
}

func init() {
	saleReportCmd := &cobra.Command{
		Use:   "sale",
		Short: "销售报表管理",
	}
	saleReportCmd.AddCommand(
		newSaleReportDetailCmd(),
		newSaleReportSummerCmd(),
		newSaleReportProfitCmd(),
	)
	reportCmd.AddCommand(saleReportCmd)
}

func newSaleReportDetailCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "detail",
		Short: "销售报表明细查询",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			for _, f := range saleReportFilters {
				cmdutil.CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/sale-report/detail-page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询销售报表明细失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	for _, f := range saleReportFilters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}

func newSaleReportSummerCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "summer",
		Short: "销售报表汇总查询",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			for _, f := range saleReportFilters {
				cmdutil.CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/sale-report/summer-page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询销售报表汇总失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	for _, f := range saleReportFilters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}

func newSaleReportProfitCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "profit",
		Short: "销售报表利润查询",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			for _, f := range saleReportFilters {
				cmdutil.CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/sale-report/profit-page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询销售报表利润失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	for _, f := range saleReportFilters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}
