package stats

import (
	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
)

var produceFlags = []cmdutil.FlagSpec{
	{Name: "product-id", Usage: "产品 ID"},
}

func init() {
	produceStats := &cobra.Command{
		Use:   "produce",
		Short: "生产分析",
	}
	produceStats.AddCommand(
		cmdutil.NewStatsGetCmd("data-overview", "/erp/produce-statistics", "生产数据总览", produceFlags),
	)
	statsCmd.AddCommand(produceStats)
}
