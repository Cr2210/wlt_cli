package order

import "github.com/spf13/cobra"

var orderCmd = &cobra.Command{
	Use:   "order",
	Short: "订单管理",
	Long:  "订单模块操作：订单CRUD、取消、重新打开、完成、关联运单、取消关联运单。",
}

// Register adds the order command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(orderCmd)
}
