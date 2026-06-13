package sale

import "github.com/spf13/cobra"

var saleCmd = &cobra.Command{
	Use:   "sale",
	Short: "销售管理",
	Long:  "销售模块操作：销售出库、销售退货。",
}

// Register adds the sale command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(saleCmd)
}
