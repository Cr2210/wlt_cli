package stats

import (
	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
)

var purchaseFlags = []cmdutil.FlagSpec{
	{Name: "product-id", Usage: "产品 ID"},
}

func init() {
	purchaseStats := &cobra.Command{
		Use:   "purchase",
		Short: "采购分析",
	}
	purchaseStats.AddCommand(
		cmdutil.NewStatsGetCmd("data-overview", "/erp/purchase-statistics", "采购数据总览", purchaseFlags),
		cmdutil.NewStatsGetCmd("supplier-rankings", "/erp/purchase-statistics", "供应商排行", purchaseFlags),
		cmdutil.NewStatsGetCmd("product-rankings", "/erp/purchase-statistics", "产品排行", purchaseFlags),
		cmdutil.NewStatsGetCmd("employee-rankings", "/erp/purchase-statistics", "员工排行", purchaseFlags),
		cmdutil.NewStatsGetCmd("region-rankings", "/erp/purchase-statistics", "区域排行", purchaseFlags),
	)
	statsCmd.AddCommand(purchaseStats)
}
