package system

import "github.com/spf13/cobra"

var SystemCmd = &cobra.Command{
	Use:   "system",
	Short: "系统管理",
	Long:  "系统管理模块操作：用户、部门、角色、权限、菜单、字典、租户、OAuth2、通知、公告、邮件、社交、日志、应用。",
}

func Register(parent *cobra.Command) {
	parent.AddCommand(SystemCmd)
}
