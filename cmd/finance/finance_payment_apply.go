package finance

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

// financePaymentApplyAPIPath 是付款申请的后端端点根路径。
const financePaymentApplyAPIPath = "/erp/finance-payment-apply"

var financePaymentApplyCmd = &cobra.Command{
	Use:   "payment-apply",
	Short: "付款申请管理",
}

func init() {
	financeCmd.AddCommand(financePaymentApplyCmd)
	financePaymentApplyCmd.AddCommand(
		newFinancePaymentApplyPageCmd(),
		cmdutil.CrudGetCmd(financePaymentApplyAPIPath, "付款申请"),
		newFinancePaymentApplySummaryCmd(),
		newFinancePaymentApplyExportExcelCmd(),
		cmdutil.CrudCreateCmd(financePaymentApplyAPIPath, "付款申请"),
		cmdutil.CrudUpdateCmd(financePaymentApplyAPIPath, "付款申请"),
		newFinancePaymentApplyUpdateStatusCmd(),
		cmdutil.CrudDeleteCmd(financePaymentApplyAPIPath, "付款申请", false),
	)
}

var financePaymentApplyFilters = []cmdutil.FlagSpec{
	{Name: "no", Usage: "申请单号"},
	{Name: "partner-id", Usage: "合作伙伴 ID"},
	{Name: "partner-name", Usage: "合作伙伴名称"},
	{Name: "receipt-account-id", Usage: "收款账户 ID"},
	{Name: "receipt-account-name", Usage: "收款账户名称"},
	{Name: "receipt-account-no", Usage: "收款账户账号"},
	{Name: "service-user-id", Usage: "财务负责人 ID"},
	{Name: "service-user-name", Usage: "财务负责人名称"},
	{Name: "approve-status", Usage: "审核状态"},
	{Name: "pay-date", Usage: "付款日期"},
	{Name: "remark", Usage: "备注"},
	{Name: "creator-name", Usage: "创建人"},
	{Name: "updater-name", Usage: "更新人"},
	{Name: "create-time", Usage: "创建时间"},
	{Name: "update-time", Usage: "更新时间"},
	{Name: "custom-order", Usage: "前端自定义排序规则"},
	{Name: "keyword", Usage: "关键字"},
	{Name: "headers", Usage: "自定义导出表头"},
}

func newFinancePaymentApplyPageCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "page",
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

// ---- 导出 Excel ----

func newFinancePaymentApplyExportExcelCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "export",
		Short: "导出付款申请 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]any{}
			collectFinanceFilters(cmd, params, financePaymentApplyFilters)
			return cmdutil.RunExport(cmd, "/erp/finance-payment-apply/export-excel", params)
		},
	}
	addFinanceFilterFlags(c, financePaymentApplyFilters)
	return c
}

// ---- 处理申请（query 参数：id + approveStatus） ----

func newFinancePaymentApplyUpdateStatusCmd() *cobra.Command {
	var id int64
	var approveStatus int
	c := &cobra.Command{
		Use:   "update-status",
		Short: "处理付款申请",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().PutQuery(context.Background(), "/erp/finance-payment-apply/update-status",
				map[string]any{"id": id, "approveStatus": approveStatus})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("处理付款申请失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "付款申请 ID")
	c.Flags().IntVar(&approveStatus, "approve-status", 0, "审核状态值")
	_ = c.MarkFlagRequired("id")
	_ = c.MarkFlagRequired("approve-status")
	return c
}
