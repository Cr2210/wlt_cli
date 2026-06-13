package stock

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	stockCmd.AddCommand(cmdutil.NewCRUDSubCmd("in", "/erp/stock-in", "入库单"))
}
