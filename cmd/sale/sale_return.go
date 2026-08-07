package sale

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	cmdutil.DocumentCmds(saleCmd,
		cmdutil.DocumentConfig{
			Name:    "return",
			APIPath: "/erp/sale-return",
			Label:   "销售退货",
			Filters: []cmdutil.FlagSpec{
				{Name: "warehouse-id", Usage: "仓库 ID"},
				{Name: "product-id", Usage: "产品 ID"},
				{Name: "no", Usage: "单号"},
				{Name: "status", Usage: "状态"},
				{Name: "start-time", Usage: "开始时间"},
				{Name: "end-time", Usage: "结束时间"},
				{Name: "type", Usage: "类型"},
			},
		},
	)
}
