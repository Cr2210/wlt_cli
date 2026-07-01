package purchase

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	cmdutil.SalePurchaseCmds(purchaseCmd,
		cmdutil.SalePurchaseConfig{
			Name:    "in",
			APIPath: "/erp/purchase-in",
			Label:   "采购入库",
			TimeKey: "inTime",
			Filters: []cmdutil.FlagSpec{
				{Name: "supplier-id", Usage: "供应商 ID"},
				{Name: "metrics-name", Usage: "检测指标名（如 含水）"},
			},
		},
	)
}
