package product

import (
	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
)

func init() {
	categoryCmd := &cobra.Command{
		Use:   "category",
		Short: "产品分类管理",
	}
	categoryCmd.AddCommand(
		cmdutil.CrudListAllCmd("/erp/product-category", "产品分类", []cmdutil.FlagSpec{
			{Name: "name", Usage: "分类名称"},
			{Name: "status", Usage: "状态"},
			{Name: "parent-id", Usage: "父分类 ID"},
		}),
		cmdutil.CrudGetCmd("/erp/product-category", "产品分类"),
		cmdutil.CrudCreateCmd("/erp/product-category", "产品分类"),
		cmdutil.CrudUpdateCmd("/erp/product-category", "产品分类"),
		cmdutil.CrudDeleteCmd("/erp/product-category", "产品分类", true),
		cmdutil.CrudUpdateStatusCmd("/erp/product-category", "产品分类"),
		cmdutil.CrudSimpleListCmd("/erp/product-category", "产品分类"),
	)
	productCmd.AddCommand(categoryCmd)
}
