package purchase

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	purchaseCmd.AddCommand(cmdutil.NewCRUDSubCmd("in", "/erp/purchase-in", "采购入库"))
}
