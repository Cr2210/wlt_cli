package stock

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	stockCmd.AddCommand(cmdutil.NewCRUDSubCmd("out", "/erp/stock-out", "出库单"))
}
