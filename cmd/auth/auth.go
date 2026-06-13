package auth

import "github.com/spf13/cobra"

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "认证管理",
	Long:  "管理登录状态：登录、登出、查看状态。",
}

func Register(parent *cobra.Command) {
	parent.AddCommand(authCmd)
}
