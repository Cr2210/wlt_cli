package invoice

import "github.com/spf13/cobra"

var invoiceCmd = &cobra.Command{
	Use:   "invoice",
	Short: "发票管理",
	Long:  "发票模块操作：发票CRUD、更新状态。",
}

// Register adds the invoice command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(invoiceCmd)
}
