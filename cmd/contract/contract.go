package contract

import "github.com/spf13/cobra"

var contractCmd = &cobra.Command{
	Use:   "contract",
	Short: "合同管理",
	Long:  "采购（采购合同、采购长协）、销售（销售合同、销售长协）、运输（运输合同、运输长协）、服务（服务合同、服务长协）。",
}

// Register adds the contract command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(contractCmd)
}
