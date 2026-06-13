package weight

import "github.com/spf13/cobra"

var weightCmd = &cobra.Command{
	Use:   "weight",
	Short: "称重管理",
	Long:  "称重模块操作：运单称重CRUD、匹配关联、取消关联。",
}

// Register adds the weight command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(weightCmd)
}
