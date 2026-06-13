package sale

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	saleCmd.AddCommand(cmdutil.NewCRUDSubCmd("out", "/erp/sale-out", "销售出库"))
}
