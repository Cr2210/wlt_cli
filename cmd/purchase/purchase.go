package purchase

import "github.com/spf13/cobra"

var purchaseCmd = &cobra.Command{
	Use:   "purchase",
	Short: "采购管理",
	Long:  "采购模块操作：采购入库、采购退货。",
}

// Register adds the purchase command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(purchaseCmd)
}
