package stats

import (
	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
)

var saleFlags = []cmdutil.FlagSpec{
	{Name: "product-id", Usage: "产品 ID"},
}

func init() {
	saleStats := &cobra.Command{
		Use:   "sale",
		Short: "销售分析",
	}
	saleStats.AddCommand(
		cmdutil.NewStatsGetCmd("data-overview", "/erp/sale-statistics", "销售数据总览", saleFlags),
		cmdutil.NewStatsGetCmd("customer-rankings", "/erp/sale-statistics", "客户排行", saleFlags),
		cmdutil.NewStatsGetCmd("product-rankings", "/erp/sale-statistics", "产品排行", saleFlags),
		cmdutil.NewStatsGetCmd("employee-rankings", "/erp/sale-statistics", "员工排行", saleFlags),
		cmdutil.NewStatsGetCmd("region-rankings", "/erp/sale-statistics", "区域排行", saleFlags),
	)
	statsCmd.AddCommand(saleStats)
}
