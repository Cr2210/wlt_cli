package waybill

import "github.com/spf13/cobra"

var waybillCmd = &cobra.Command{
	Use:   "waybill",
	Short: "运单管理",
	Long:  "运单模块操作：运单分页查询（新版 /erp/waybill/*）、推送配置。",
}

// Register adds the waybill command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(waybillCmd)
}
