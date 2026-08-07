package stock

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	cmdutil.DocumentCmds(stockCmd,
		cmdutil.DocumentConfig{
			Name:      "move",
			APIPath:   "/erp/stock-move",
			Label:     "调拨单",
			TimeKey:   "moveTime",
			Headers:   true,
			PageAlias: true,
			Filters: []cmdutil.FlagSpec{
				{Name: "no", Usage: "调拨单号"},
				{Name: "create-time", Usage: "创建时间"},
				{Name: "update-time", Usage: "更新时间"},
				{Name: "status", Usage: "状态"},
				{Name: "remark", Usage: "备注"},
				{Name: "creator", Usage: "创建者编号"},
				{Name: "creator-name", Usage: "创建者姓名"},
				{Name: "updater", Usage: "更新者编号"},
				{Name: "updater-name", Usage: "更新者姓名"},
				{Name: "product-id", Usage: "产品编号"},
				{Name: "from-warehouse-id", Usage: "调出仓库编号"},
				{Name: "to-warehouse-id", Usage: "调入仓库编号"},
				{Name: "product-name", Usage: "产品名称"},
				{Name: "metrics-name", Usage: "产品指标"},
				{Name: "batch-no", Usage: "批次号"},
				{Name: "user-id", Usage: "业务负责人ID"},
				{Name: "custom-order", Usage: "前端自定义排序规则"},
				{Name: "keyword", Usage: "关键字"},
			},
		},
	)
}
