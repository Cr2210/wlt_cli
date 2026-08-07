package stock

import "github.com/weiliantong/cli/internal/cmdutil"

func init() {
	cmdutil.DocumentCmds(stockCmd,
		cmdutil.DocumentConfig{
			Name:      "out",
			APIPath:   "/erp/stock-out",
			Label:     "出库单",
			TimeKey:   "outTime",
			Headers:   true,
			PageAlias: true,
			Filters: []cmdutil.FlagSpec{
				{Name: "no", Usage: "出库单号"},
				{Name: "customer-id", Usage: "客户编号"},
				{Name: "customer-name", Usage: "客户"},
				{Name: "status", Usage: "状态"},
				{Name: "remark", Usage: "备注"},
				{Name: "creator", Usage: "创建者"},
				{Name: "product-id", Usage: "产品编号"},
				{Name: "product-name", Usage: "产品名"},
				{Name: "warehouse-id", Usage: "仓库编号"},
				{Name: "warehouse-name", Usage: "仓库名"},
				{Name: "metrics-name", Usage: "产品指标名称"},
				{Name: "creator-name", Usage: "创建人名称"},
				{Name: "user-id", Usage: "业务人员ID"},
				{Name: "receive-address", Usage: "收货地址"},
				{Name: "send-address", Usage: "发货地址"},
				{Name: "batch-no", Usage: "批次号"},
				{Name: "create-time", Usage: "创建时间"},
				{Name: "updater-name", Usage: "更新人名称"},
				{Name: "update-time", Usage: "更新时间"},
				{Name: "custom-order", Usage: "前端自定义排序规则"},
				{Name: "keyword", Usage: "关键字"},
			},
		},
	)
}
