package product

import "github.com/spf13/cobra"

var productCmd = &cobra.Command{
	Use:   "product",
	Short: "产品管理",
	Long:  "产品模块操作：产品、分类、单位、指标。",
}

// Register adds the product command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(productCmd)
}
