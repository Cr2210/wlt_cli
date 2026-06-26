package finance

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var financePrepaymentApplyCmd = &cobra.Command{
	Use:   "prepayment-apply",
	Short: "预付申请管理",
}

func init() {
	financeCmd.AddCommand(financePrepaymentApplyCmd)
	financePrepaymentApplyCmd.AddCommand(
		newFinancePrepaymentApplyListCmd(),
		newFinancePrepaymentApplyGetCmd(),
		newFinancePrepaymentApplySummaryCmd(),
	)
}

var financePrepaymentApplyFilters = []cmdutil.FlagSpec{
	{Name: "no", Usage: "申请单号"},
	{Name: "partner-id", Usage: "合作伙伴 ID"},
	{Name: "partner-name", Usage: "合作伙伴名称"},
	{Name: "receipt-account-id", Usage: "收款账户 ID"},
	{Name: "service-user-id", Usage: "财务负责人 ID"},
	{Name: "approve-status", Usage: "审核状态"},
	{Name: "pay-date", Usage: "付款日期"},
	{Name: "keyword", Usage: "关键字"},
}

func newFinancePrepaymentApplyListCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询预付申请",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			collectFinanceFilters(cmd, params, financePrepaymentApplyFilters)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-prepayment-apply/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询预付申请失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	addFinanceFilterFlags(c, financePrepaymentApplyFilters)
	return c
}

func newFinancePrepaymentApplyGetCmd() *cobra.Command {
	var id int64
	c := &cobra.Command{
		Use:   "get",
		Short: "获取预付申请详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-prepayment-apply/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取预付申请失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "预付申请 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

func newFinancePrepaymentApplySummaryCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "summary",
		Short: "获取预付申请合计",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			collectFinanceFilters(cmd, params, financePrepaymentApplyFilters)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-prepayment-apply/summary", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取预付申请合计失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	addFinanceFilterFlags(c, financePrepaymentApplyFilters)
	return c
}
