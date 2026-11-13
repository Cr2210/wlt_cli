package finance

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

// financeReceiptPaymentAPIPath 是收付款的后端端点根路径。
const financeReceiptPaymentAPIPath = "/erp/finance-receipt-payment"

func init() {
	financeCmd.AddCommand(financeReceiptPaymentCmd)
	financeReceiptPaymentCmd.AddCommand(
		newFinanceReceiptPaymentPageCmd(),
		newFinanceReceiptPaymentPageCountCmd(),
		newFinanceReceiptPaymentGetCmd(),
		cmdutil.CrudCreateCmd(financeReceiptPaymentAPIPath, "收付款"),
		cmdutil.CrudUpdateCmd(financeReceiptPaymentAPIPath, "收付款"),
		cmdutil.CrudDeleteCmd(financeReceiptPaymentAPIPath, "收付款", true),
		cmdutil.CrudUpdateStatusCmd(financeReceiptPaymentAPIPath, "收付款"),
		newFinanceReceiptPaymentSummaryCmd(),
		newFinanceReceiptPaymentExportExcelCmd(),
	)
}

var financeReceiptPaymentFilters = []cmdutil.FlagSpec{
	{Name: "no", Usage: "单号"},
	{Name: "type", Usage: "业务类别（RECEIPT 收款 / PAYMENT 付款）"},
	{Name: "pay-date", Usage: "收付款日期"},
	{Name: "account-id", Usage: "收付款账户"},
	{Name: "account-name", Usage: "账户名称"},
	{Name: "account-no", Usage: "账户账号"},
	{Name: "partner-id", Usage: "客户/供应商"},
	{Name: "partner-name", Usage: "客户/供应商名称"},
	{Name: "service-user-id", Usage: "财务负责人"},
	{Name: "service-user-name", Usage: "财务负责人名称"},
	{Name: "status", Usage: "核销状态"},
	{Name: "approve-status", Usage: "审批状态"},
	{Name: "remark", Usage: "备注"},
	{Name: "creator-name", Usage: "创建人"},
	{Name: "updater-name", Usage: "更新人"},
	{Name: "create-time", Usage: "创建时间"},
	{Name: "update-time", Usage: "更新时间"},
	{Name: "keyword", Usage: "关键字"},
	{Name: "custom-order", Usage: "前端自定义排序规则"},
	{Name: "is-prepaid", Usage: "预付标识（true 仅预付 / false 非预付）"},
	{Name: "headers", Usage: "自定义导出表头"},
}

var financeReceiptPaymentPageCountFilters = []cmdutil.FlagSpec{
	{Name: "no", Usage: "单号"},
	{Name: "type", Usage: "业务类别（RECEIPT 收款 / PAYMENT 付款）"},
	{Name: "pay-date", Usage: "收付款日期"},
	{Name: "account-id", Usage: "收付款账户"},
	{Name: "account-name", Usage: "账户名称"},
	{Name: "account-no", Usage: "账户账号"},
	{Name: "partner-id", Usage: "客户/供应商"},
	{Name: "partner-name", Usage: "客户/供应商名称"},
	{Name: "service-user-id", Usage: "财务负责人"},
	{Name: "service-user-name", Usage: "财务负责人名称"},
	{Name: "status", Usage: "核销状态"},
	{Name: "approve-status", Usage: "审批状态"},
	{Name: "remark", Usage: "备注"},
	{Name: "creator-name", Usage: "创建人"},
	{Name: "updater-name", Usage: "更新人"},
	{Name: "create-time", Usage: "创建时间"},
	{Name: "update-time", Usage: "更新时间"},
	{Name: "keyword", Usage: "关键字"},
	{Name: "custom-order", Usage: "前端自定义排序规则"},
	{Name: "is-prepaid", Usage: "预付标识（true 仅预付 / false 非预付）"},
}

var financeReceiptPaymentCmd = &cobra.Command{
	Use:   "receipt-payment",
	Short: "收付款管理",
}

// ---- 分页查询 ----

func newFinanceReceiptPaymentPageCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "page",
		Short: "分页查询收付款记录（type=RECEIPT 收款 / PAYMENT 付款）",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			collectFinanceFilters(cmd, params, financeReceiptPaymentFilters)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-receipt-payment/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询收付款失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	addFinanceFilterFlags(c, financeReceiptPaymentFilters)
	return c
}

// ---- 分页计数 ----

func newFinanceReceiptPaymentPageCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "page-count",
		Short: "按筛选统计收付款记录数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			collectFinanceFilters(cmd, params, financeReceiptPaymentPageCountFilters)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-receipt-payment/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计收付款失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	addFinanceFilterFlags(c, financeReceiptPaymentPageCountFilters)
	return c
}

// ---- 获取详情 ----

func newFinanceReceiptPaymentGetCmd() *cobra.Command {
	var id int64
	var no string

	c := &cobra.Command{
		Use:   "get",
		Short: "获取收付款详情（id 或 no 任选其一）",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			if id > 0 {
				params["id"] = id
			}
			if no != "" {
				params["no"] = no
			}
			if len(params) == 0 {
				return output.NewExitError(4, "请指定 --id 或 --no", "")
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-receipt-payment/get", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取收付款详情失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "收付款 ID")
	c.Flags().StringVar(&no, "no", "", "收付款单号")
	return c
}

// ---- 汇总数据 ----

func newFinanceReceiptPaymentSummaryCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "summary",
		Short: "获取收付款汇总数据",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			collectFinanceFilters(cmd, params, financeReceiptPaymentPageCountFilters)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-receipt-payment/summary", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取收付款汇总失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	addFinanceFilterFlags(c, financeReceiptPaymentPageCountFilters)
	return c
}

// ---- 导出 Excel ----

func newFinanceReceiptPaymentExportExcelCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "export",
		Short: "导出收付款记录 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			collectFinanceFilters(cmd, params, financeReceiptPaymentFilters)
			return cmdutil.RunExport(cmd, "/erp/finance-receipt-payment/export-excel", params)
		},
	}
	addFinanceFilterFlags(c, financeReceiptPaymentFilters)
	return c
}
