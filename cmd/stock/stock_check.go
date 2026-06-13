package stock

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	stockCmd.AddCommand(cmdutil.NewCRUDSubCmd("check", "/erp/stock-check", "盘点单"))
}
