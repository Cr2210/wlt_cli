package report

import "github.com/spf13/cobra"

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "报表管理",
	Long:  "报表模块操作：库存报表、采购报表、销售报表。",
}

// Register adds the report command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(reportCmd)
}
