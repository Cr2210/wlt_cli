package produce

import "github.com/spf13/cobra"

var produceCmd = &cobra.Command{
	Use:   "produce",
	Short: "生产管理",
	Long:  "生产模块操作：生产单CRUD、生产计划CRUD、质检分页。",
}

// Register adds the produce command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(produceCmd)
}
