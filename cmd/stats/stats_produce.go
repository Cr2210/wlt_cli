package stats

import (
	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
)

func init() {
	produceStats := &cobra.Command{
		Use:   "produce",
		Short: "生产分析",
	}
	produceStats.AddCommand(
		cmdutil.NewStatsGetCmd("data-overview", "/erp/produce-statistics", "生产数据总览", nil),
	)
	statsCmd.AddCommand(produceStats)
}
