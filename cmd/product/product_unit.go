package product

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	unitCmd := cmdutil.NewCRUDGroup(cmdutil.CRUDConfig{
		Name:    "unit",
		APIPath: "/erp/product-unit",
		Label:   "产品单位",
		ListFilters: []cmdutil.FlagSpec{
			{Name: "name", Usage: "单位名称"},
			{Name: "status", Usage: "状态"},
		},
		SingleDelete: true,
	})
	unitCmd.AddCommand(cmdutil.CrudSimpleListCmd("/erp/product-unit", "产品单位"))
	productCmd.AddCommand(unitCmd)
}
