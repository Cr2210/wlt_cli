package stats

import (
	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
)

func init() {
	financeStats := &cobra.Command{
		Use:   "finance",
		Short: "财务分析",
	}
	financeStats.AddCommand(
		cmdutil.NewStatsGetCmd("data-overview", "/erp/finance-statistics", "财务数据总览", nil),
		cmdutil.NewStatsGetCmd("receivable-rankings", "/erp/finance-statistics", "应收款排行", nil),
		cmdutil.NewStatsGetCmd("overdue-receivable-rankings", "/erp/finance-statistics", "逾期应收排行", nil),
		cmdutil.NewStatsGetCmd("payable-rankings", "/erp/finance-statistics", "应付款排行", nil),
		cmdutil.NewStatsGetCmd("overdue-payable-rankings", "/erp/finance-statistics", "逾期应付排行", nil),
	)
	statsCmd.AddCommand(financeStats)
}
