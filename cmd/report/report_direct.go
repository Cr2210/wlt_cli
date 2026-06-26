package report

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var directReportFilters = []cmdutil.FlagSpec{
	{Name: "purchase-order-no", Usage: "采购订单号"},
	{Name: "sale-order-no", Usage: "销售订单号"},
	{Name: "waybill-no", Usage: "运单编号"},
	{Name: "waybill-status", Usage: "运单状态"},
	{Name: "car-number", Usage: "车牌号"},
	{Name: "supplier-name", Usage: "供应商名称"},
	{Name: "customer-name", Usage: "客户名称"},
	{Name: "project-name", Usage: "项目名称"},
	{Name: "product-name", Usage: "产品名称"},
	{Name: "metrics-name", Usage: "指标名称"},
	{Name: "load-date-start", Usage: "装货日期起"},
	{Name: "load-date-end", Usage: "装货日期止"},
}

func init() {
	directReportCmd := &cobra.Command{
		Use:   "direct",
		Short: "直采直销报表管理",
	}
	directReportCmd.AddCommand(
		newDirectReportDetailCmd(),
		newDirectReportDetailCountCmd(),
	)
	reportCmd.AddCommand(directReportCmd)
}

func newDirectReportDetailCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "detail",
		Short: "直采直销报表明细查询",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			for _, f := range directReportFilters {
				cmdutil.CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/direct-purchase-sale-report/detail-page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询直采直销报表失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	for _, f := range directReportFilters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}

func newDirectReportDetailCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "detail-count",
		Short: "直采直销报表明细统计",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			for _, f := range directReportFilters {
				cmdutil.CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/direct-purchase-sale-report/detail-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计直采直销报表失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	for _, f := range directReportFilters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}
