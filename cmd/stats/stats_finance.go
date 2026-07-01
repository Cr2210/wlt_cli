package stats

import (
	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
)

var financeFlags = []cmdutil.FlagSpec{
	{Name: "product-id", Usage: "产品 ID"},
}

func init() {
	financeStats := &cobra.Command{
		Use:   "finance",
		Short: "财务分析",
	}
	financeStats.AddCommand(
		cmdutil.NewStatsGetCmd("data-overview", "/erp/finance-statistics", "财务数据总览", financeFlags),
		cmdutil.NewStatsGetCmd("receivable-rankings", "/erp/finance-statistics", "应收款排行", financeFlags),
		cmdutil.NewStatsGetCmd("overdue-receivable-rankings", "/erp/finance-statistics", "逾期应收排行", financeFlags),
		cmdutil.NewStatsGetCmd("payable-rankings", "/erp/finance-statistics", "应付款排行", financeFlags),
		cmdutil.NewStatsGetCmd("overdue-payable-rankings", "/erp/finance-statistics", "逾期应付排行", financeFlags),
	)
	statsCmd.AddCommand(financeStats)
}
