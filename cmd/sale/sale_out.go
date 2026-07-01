package sale

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	cmdutil.SalePurchaseCmds(saleCmd,
		cmdutil.SalePurchaseConfig{
			Name:    "out",
			APIPath: "/erp/sale-out",
			Label:   "销售出库",
			TimeKey: "outTime",
			Filters: []cmdutil.FlagSpec{
				{Name: "customer-id", Usage: "客户 ID"},
				{Name: "batch-no", Usage: "批次号"},
			},
		},
	)
}
