package stats

import (
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "数据总览",
	Long:  "数据总览模块：经营总览、采购分析、销售分析、财务分析、生产分析、库存分析。",
}

// Register adds the stats command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(statsCmd)
}
