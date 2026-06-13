package stats

import (
	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
)

func init() {
	saleStats := &cobra.Command{
		Use:   "sale",
		Short: "销售分析",
	}
	saleStats.AddCommand(
		cmdutil.NewStatsGetCmd("data-overview", "/erp/sale-statistics", "销售数据总览", nil),
		cmdutil.NewStatsGetCmd("customer-rankings", "/erp/sale-statistics", "客户排行", nil),
		cmdutil.NewStatsGetCmd("product-rankings", "/erp/sale-statistics", "产品排行", nil),
		cmdutil.NewStatsGetCmd("employee-rankings", "/erp/sale-statistics", "员工排行", nil),
		cmdutil.NewStatsGetCmd("region-rankings", "/erp/sale-statistics", "区域排行", nil),
	)
	statsCmd.AddCommand(saleStats)
}
