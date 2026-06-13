package finance

import "github.com/spf13/cobra"

var financeCmd = &cobra.Command{
	Use:   "finance",
	Short: "财务管理",
	Long:  "财务模块操作：账户、付款单、收款单、退款、结算单、核销、开票申请、收付款。",
}

// Register adds the finance command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(financeCmd)
}
