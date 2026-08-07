package sale

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	cmdutil.DocumentCmds(saleCmd,
		cmdutil.DocumentConfig{
			Name:    "out",
			APIPath: "/erp/sale-out",
			Label:   "销售出库",
			TimeKey: "outTime",
			Headers: true,
			Filters: []cmdutil.FlagSpec{
				{Name: "warehouse-id", Usage: "仓库 ID"},
				{Name: "product-id", Usage: "产品 ID"},
				{Name: "product-name", Usage: "产品名称（模糊搜索）"},
				{Name: "no", Usage: "单号"},
				{Name: "status", Usage: "状态"},
				{Name: "type", Usage: "类型"},
				{Name: "customer-id", Usage: "客户 ID"},
				{Name: "batch-no", Usage: "批次号"},
			},
		},
	)
}
