package configcmd

import "github.com/spf13/cobra"

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "配置管理",
	Long:  "管理 CLI 配置：初始化、查看、修改。",
}

func Register(parent *cobra.Command) {
	parent.AddCommand(configCmd)
}
