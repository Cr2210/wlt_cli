package stock

import "github.com/spf13/cobra"

var stockCmd = &cobra.Command{
	Use:   "stock",
	Short: "库存管理",
	Long:  "库存模块操作：仓库、查询、入库、出库、调拨、盘点、明细。",
}

func init() {
	// subcommands registered by stock_warehouse.go, stock_query.go, stock_in.go, etc.
}

// Register adds the stock command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(stockCmd)
}
