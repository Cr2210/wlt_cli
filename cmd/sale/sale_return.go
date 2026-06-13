package sale

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	saleCmd.AddCommand(cmdutil.NewCRUDSubCmd("return", "/erp/sale-return", "销售退货"))
}
