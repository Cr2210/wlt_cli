package report

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var purchaseReportFilters = []cmdutil.FlagSpec{
	{Name: "supplier-id", Usage: "供应商 ID"},
	{Name: "product-id", Usage: "产品 ID"},
	{Name: "category-id", Usage: "分类 ID"},
	{Name: "warehouse-id", Usage: "仓库 ID"},
	{Name: "status", Usage: "状态"},
	{Name: "start-time", Usage: "开始时间"},
	{Name: "end-time", Usage: "结束时间"},
}

func init() {
	purchaseReportCmd := &cobra.Command{
		Use:   "purchase",
		Short: "采购报表管理",
	}
	purchaseReportCmd.AddCommand(
		newPurchaseReportDetailCmd(),
		newPurchaseReportDetailCountCmd(),
		newPurchaseReportSummerCmd(),
		newPurchaseReportSummerCountCmd(),
	)
	reportCmd.AddCommand(purchaseReportCmd)
}

func newPurchaseReportDetailCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "detail",
		Short: "采购报表明细查询",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			for _, f := range purchaseReportFilters {
				cmdutil.CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/purchase-report/detail-page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询采购报表明细失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	for _, f := range purchaseReportFilters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}

func newPurchaseReportSummerCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "summer",
		Short: "采购报表汇总查询",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			for _, f := range purchaseReportFilters {
				cmdutil.CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/purchase-report/summer-page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询采购报表汇总失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	for _, f := range purchaseReportFilters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}

func newPurchaseReportDetailCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "detail-count",
		Short: "采购报表明细统计",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			for _, f := range purchaseReportFilters {
				cmdutil.CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/purchase-report/detail-page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计采购报表明细失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	for _, f := range purchaseReportFilters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}

func newPurchaseReportSummerCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "summer-count",
		Short: "采购报表汇总统计",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			for _, f := range purchaseReportFilters {
				cmdutil.CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/purchase-report/summer-page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计采购报表汇总失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	for _, f := range purchaseReportFilters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}
