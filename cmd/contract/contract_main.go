package contract

import (
	"github.com/weiliantong/cli/internal/cmdutil"
)

// dateRangeFilters 是四类合同共用的时间范围筛选字段（转 orderDate[0]/[1]、endTime[0]/[1]）。
var dateRangeFilters = []cmdutil.FlagSpec{
	{Name: "keyword", Usage: "关键字搜索（合同编号等）"},
	{Name: "enterprise-id", Usage: "企业 ID"},
}

func init() {
	cmdutil.ContractTypeCmds(contractCmd,
		// 采购长协：/erp/contract?type=PURCHASE_LONG_COOPERATE
		cmdutil.ContractTypeConfig{
			Type:         "PURCHASE_LONG_COOPERATE",
			Label:        "采购长协",
			Filters:      dateRangeFilters,
			HasDateRange: true,
		},
		// 销售合同：/erp/service-contract?type=SALE_CONTRACT
		cmdutil.ContractTypeConfig{
			Type:         "SALE_CONTRACT",
			APIPath:      "/erp/service-contract",
			Label:        "销售合同",
			Filters:      dateRangeFilters,
			HasDateRange: true,
		},
		// 销售长协：/erp/contract?type=SALE_LONG_COOPERATE
		cmdutil.ContractTypeConfig{
			Type:         "SALE_LONG_COOPERATE",
			Label:        "销售长协",
			Filters:      dateRangeFilters,
			HasDateRange: true,
		},
		// 运输合同：/erp/transport-contract?type=TRANSPORT
		cmdutil.ContractTypeConfig{
			Type:         "TRANSPORT",
			APIPath:      "/erp/transport-contract",
			Label:        "运输合同",
			Filters:      dateRangeFilters,
			HasDateRange: true,
		},
		// 运输长协：/erp/contract?type=TRANSPORT_LONG
		cmdutil.ContractTypeConfig{
			Type:         "TRANSPORT_LONG",
			Label:        "运输长协",
			Filters:      dateRangeFilters,
			HasDateRange: true,
		},
		// 服务合同：/erp/provision-contract?type=SERVICE
		// 注意：Use 显式指定避免自动 slug "service" 与旧「业务合同」冲突；
		//       APIPath 走 /erp/provision-contract 而非 /erp/contract。
		cmdutil.ContractTypeConfig{
			Type:         "SERVICE",
			Use:          "service-contract",
			APIPath:      "/erp/provision-contract",
			Label:        "服务合同",
			Filters:      dateRangeFilters,
			HasDateRange: true,
		},
		// 服务长协：/erp/provision-contract?type=SERVICE_LONG
		cmdutil.ContractTypeConfig{
			Type:         "SERVICE_LONG",
			Use:          "service-long",
			APIPath:      "/erp/provision-contract",
			Label:        "服务长协",
			Filters:      dateRangeFilters,
			HasDateRange: true,
		},
	)
}
