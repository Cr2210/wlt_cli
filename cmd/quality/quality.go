package quality

import "github.com/spf13/cobra"

var qualityCmd = &cobra.Command{
	Use:   "quality",
	Short: "质检管理",
	Long:  "质检模块操作：质检单CRUD、质检称重。",
}

// Register adds the quality command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(qualityCmd)
}
