package finance

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var financePaymentApplyCmd = &cobra.Command{
	Use:   "payment-apply",
	Short: "付款申请管理",
}

func init() {
	financeCmd.AddCommand(financePaymentApplyCmd)
	financePaymentApplyCmd.AddCommand(
		newFinancePaymentApplyListCmd(),
		newFinancePaymentApplyGetCmd(),
		newFinancePaymentApplySummaryCmd(),
	)
}

var financePaymentApplyFilters = []cmdutil.FlagSpec{
	{Name: "no", Usage: "申请单号"},
	{Name: "partner-id", Usage: "合作伙伴 ID"},
	{Name: "partner-name", Usage: "合作伙伴名称"},
	{Name: "receipt-account-id", Usage: "收款账户 ID"},
	{Name: "service-user-id", Usage: "财务负责人 ID"},
	{Name: "approve-status", Usage: "审核状态"},
	{Name: "pay-date", Usage: "付款日期"},
	{Name: "keyword", Usage: "关键字"},
}

func newFinancePaymentApplyListCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询付款申请",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			collectFinanceFilters(cmd, params, financePaymentApplyFilters)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-payment-apply/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询付款申请失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	addFinanceFilterFlags(c, financePaymentApplyFilters)
	return c
}

func newFinancePaymentApplyGetCmd() *cobra.Command {
	var id int64
	c := &cobra.Command{
		Use:   "get",
		Short: "获取付款申请详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-payment-apply/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取付款申请失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "付款申请 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

func newFinancePaymentApplySummaryCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "summary",
		Short: "获取付款申请合计",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			collectFinanceFilters(cmd, params, financePaymentApplyFilters)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-payment-apply/summary", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取付款申请合计失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	addFinanceFilterFlags(c, financePaymentApplyFilters)
	return c
}
