package settlement

import "github.com/spf13/cobra"

var settlementCmd = &cobra.Command{
	Use:   "settlement",
	Short: "结算管理",
	Long:  "结算模块操作：结算单CRUD、取消结算运单。",
}

// Register adds the settlement command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(settlementCmd)
}
