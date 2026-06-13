package purchase

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	purchaseCmd.AddCommand(cmdutil.NewCRUDSubCmd("return", "/erp/purchase-return", "采购退货"))
}
