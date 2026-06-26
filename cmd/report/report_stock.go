package report

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

func init() {
	stockReportCmd := &cobra.Command{
		Use:   "stock",
		Short: "库存报表管理",
	}
	stockReportCmd.AddCommand(
		// 收发明细
		newReportPageCmd("detail", "库存报表明细查询", "/erp/stock-report/detail-page", stockReportDetailFilters),
		newReportCountCmd("detail-count", "库存报表明细统计", "/erp/stock-report/detail-page-count", stockReportDetailFilters),
		// 库存余额
		newReportPageCmd("warehouse", "库存报表余额查询", "/erp/stock-report/stock-warehouse-page", stockReportWarehouseFilters),
		newReportCountCmd("warehouse-count", "库存报表余额统计", "/erp/stock-report/stock-warehouse-count", stockReportWarehouseFilters),
		// 收发汇总
		newReportPageCmd("buy-send", "库存报表收发汇总查询", "/erp/stock-report/buy-send-page", stockReportBuySendFilters),
		newReportCountCmd("buy-send-count", "库存报表收发汇总统计", "/erp/stock-report/buy-send-count", stockReportBuySendFilters),
		// 应收应付明细
		newReportPageCmd("finance", "库存报表应收应付查询", "/erp/stock-report/finance-record-page", stockReportFinanceFilters),
		newReportCountCmd("finance-count", "库存报表应收应付统计", "/erp/stock-report/finance-record-count", stockReportFinanceFilters),
		// 生产利润
		newReportPageCmd("produce", "库存报表生产利润查询", "/erp/stock-report/produce-profit-page", stockReportProduceFilters),
		newReportCountCmd("produce-count", "库存报表生产利润统计", "/erp/stock-report/produce-profit-page-count", stockReportProduceFilters),
		// 原料调价记录
		newStockReportPriceChangeCmd(),
	)
	reportCmd.AddCommand(stockReportCmd)
}

var stockReportDetailFilters = []cmdutil.FlagSpec{
	{Name: "product-id", Usage: "产品 ID"},
	{Name: "metrics-name", Usage: "产品指标"},
	{Name: "warehouse-id", Usage: "仓库 ID"},
	{Name: "batch-no", Usage: "批次号"},
	{Name: "biz-type", Usage: "类型"},
	{Name: "biz-no", Usage: "出入库单号"},
	{Name: "bill-date", Usage: "单据日期"},
}

var stockReportWarehouseFilters = []cmdutil.FlagSpec{
	{Name: "product-id", Usage: "产品 ID"},
	{Name: "metrics-name", Usage: "产品指标"},
	{Name: "warehouse-id", Usage: "仓库 ID"},
	{Name: "end-time", Usage: "截止时间"},
	{Name: "hide-zero-stock", Usage: "隐藏 0 库存数据"},
}

var stockReportBuySendFilters = []cmdutil.FlagSpec{
	{Name: "product-id", Usage: "产品 ID"},
	{Name: "product-name", Usage: "产品名称"},
	{Name: "metrics-name", Usage: "产品指标"},
	{Name: "warehouse-name", Usage: "仓库名称"},
	{Name: "customer-id", Usage: "客户 ID"},
	{Name: "bill-date-start", Usage: "单据时间起"},
	{Name: "bill-date-end", Usage: "单据时间止"},
}

var stockReportFinanceFilters = []cmdutil.FlagSpec{
	{Name: "type", Usage: "应收/应付"},
	{Name: "enterprise-id", Usage: "企业 ID"},
	{Name: "enterprise-name", Usage: "企业名称"},
	{Name: "enterprise-type", Usage: "客户类型(0欠款客户/1所有客户)"},
	{Name: "create-time-start", Usage: "创建时间起"},
	{Name: "create-time-end", Usage: "创建时间止"},
}

var stockReportProduceFilters = []cmdutil.FlagSpec{
	{Name: "customer-name", Usage: "客户名称"},
	{Name: "product-name", Usage: "产品名称"},
	{Name: "product-id", Usage: "产品 ID"},
	{Name: "metrics-name", Usage: "产品指标"},
	{Name: "warehouse-id", Usage: "仓库 ID"},
	{Name: "warehouse-name", Usage: "仓库名称"},
	{Name: "batch-no", Usage: "生产批次号"},
	{Name: "plan-no", Usage: "方案号"},
	{Name: "operate-time", Usage: "生产日期"},
	{Name: "creator-name", Usage: "创建人"},
}

// newReportPageCmd builds a paginated stock-report query command.
func newReportPageCmd(use, short, path string, filters []cmdutil.FlagSpec) *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			for _, f := range filters {
				cmdutil.CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), path, params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("%s失败: %s", short, err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	for _, f := range filters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}

// newReportCountCmd builds a stock-report count/summary command (non-paginated).
func newReportCountCmd(use, short, path string, filters []cmdutil.FlagSpec) *cobra.Command {
	c := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			for _, f := range filters {
				cmdutil.CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), path, params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("%s失败: %s", short, err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	for _, f := range filters {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}

// newStockReportPriceChangeCmd queries 原料调价记录 by produce (生产任务) ID.
func newStockReportPriceChangeCmd() *cobra.Command {
	var produceId int64
	c := &cobra.Command{
		Use:   "price-change-logs",
		Short: "查询生产原料调价记录",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(),
				"/erp/stock-report/produce-price-change-logs", map[string]any{"produceId": produceId})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询原料调价记录失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&produceId, "produce-id", 0, "生产任务 ID")
	_ = c.MarkFlagRequired("produce-id")
	return c
}
