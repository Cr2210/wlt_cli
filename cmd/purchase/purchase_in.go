package purchase

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	cmdutil.DocumentCmds(purchaseCmd,
		cmdutil.DocumentConfig{
			Name:    "in",
			APIPath: "/erp/purchase-in",
			Label:   "采购入库",
			TimeKey: "inTime",
			Headers: true,
			Filters: []cmdutil.FlagSpec{
				{Name: "warehouse-id", Usage: "仓库 ID"},
				{Name: "product-id", Usage: "产品 ID"},
				{Name: "product-name", Usage: "产品名称（模糊搜索）"},
				{Name: "no", Usage: "单号"},
				{Name: "status", Usage: "状态"},
				{Name: "type", Usage: "类型"},
				{Name: "supplier-id", Usage: "供应商 ID"},
				{Name: "metrics-name", Usage: "检测指标名（如 含水）"},
			},
		},
	)
}
