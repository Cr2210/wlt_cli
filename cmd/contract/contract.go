package contract

import "github.com/spf13/cobra"

var contractCmd = &cobra.Command{
	Use:   "contract",
	Short: "合同管理",
	Long:  "合同模块操作：合同、合同条款、服务协议、运输协议。",
}

// Register adds the contract command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(contractCmd)
}
