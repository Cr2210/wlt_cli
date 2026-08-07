package stock

import (
	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
)

var warehouseCmd = &cobra.Command{
	Use:   "warehouse",
	Short: "仓库管理",
	Long:  "管理 ERP 仓库：查询、创建、更新、删除、状态管理。",
}

func init() {
	stockCmd.AddCommand(warehouseCmd)
	cmdutil.AddCRUDToParent(warehouseCmd, cmdutil.CRUDConfig{
		APIPath: "/erp/warehouse",
		Label:   "仓库",
		ListFilters: []cmdutil.FlagSpec{
			{Name: "name", Usage: "仓库名称（模糊）"},
			{Name: "status", Usage: "状态"},
			{Name: "type", Usage: "类型"},
		},
		SingleDelete: false, // DELETE ?ids=
	})
	warehouseCmd.AddCommand(cmdutil.CrudSimpleListCmd("/erp/warehouse", "仓库"))
}
