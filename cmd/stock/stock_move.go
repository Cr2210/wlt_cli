package stock

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	stockCmd.AddCommand(cmdutil.NewCRUDSubCmd("move", "/erp/stock-move", "调拨单"))
}
