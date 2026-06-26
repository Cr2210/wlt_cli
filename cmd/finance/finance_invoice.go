package finance

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var financeInvoiceCmd = &cobra.Command{
	Use:   "invoice",
	Short: "收开票管理",
}

func init() {
	financeCmd.AddCommand(financeInvoiceCmd)
	financeInvoiceCmd.AddCommand(
		newFinanceInvoiceListCmd(),
		newFinanceInvoiceGetCmd(),
		newFinanceInvoiceSummaryCmd(),
	)
}

var financeInvoiceFilters = []cmdutil.FlagSpec{
	{Name: "no", Usage: "收开票号"},
	{Name: "type", Usage: "业务类别"},
	{Name: "invoice-type", Usage: "发票类型(专票/普票/收据)"},
	{Name: "invoice-no", Usage: "发票号"},
	{Name: "partner-id", Usage: "合作伙伴 ID"},
	{Name: "partner-name", Usage: "客户/供应商名称"},
	{Name: "service-user-id", Usage: "财务负责人 ID"},
	{Name: "status", Usage: "核销状态"},
	{Name: "approve-status", Usage: "审核状态"},
	{Name: "keyword", Usage: "关键字"},
}

func newFinanceInvoiceListCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询收开票",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			collectFinanceFilters(cmd, params, financeInvoiceFilters)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-invoice/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询收开票失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	addFinanceFilterFlags(c, financeInvoiceFilters)
	return c
}

func newFinanceInvoiceGetCmd() *cobra.Command {
	var id int64
	c := &cobra.Command{
		Use:   "get",
		Short: "获取收开票详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-invoice/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取收开票失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "收开票 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

func newFinanceInvoiceSummaryCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "summary",
		Short: "获取收开票合计",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			collectFinanceFilters(cmd, params, financeInvoiceFilters)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-invoice/summary", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取收开票合计失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	addFinanceFilterFlags(c, financeInvoiceFilters)
	return c
}

// addFinanceFilterFlags registers filter flags on a finance command.
func addFinanceFilterFlags(c *cobra.Command, filters []cmdutil.FlagSpec) {
	for _, f := range filters {
		c.Flags().String(f.Name, "", f.Usage)
	}
}

// collectFinanceFilters collects non-empty filter flags into the params map.
func collectFinanceFilters(cmd *cobra.Command, params map[string]any, filters []cmdutil.FlagSpec) {
	for _, f := range filters {
		cmdutil.CollectStringFlag(cmd, params, f.Name)
	}
}
