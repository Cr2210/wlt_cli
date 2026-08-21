package order

import "github.com/spf13/cobra"

var orderCmd = &cobra.Command{
	Use:   "order",
	Short: "订单管理",
	Long:  "订单模块：订单CRUD（main）、订单计划（plan purchase/sale + CRUD）、订单预付关联（prepayment-relation）。",
}

// Register adds the order command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(orderCmd)
}
