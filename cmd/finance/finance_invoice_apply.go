package finance

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

func init() {
	financeCmd.AddCommand(financeInvoiceApplyCmd)
	financeInvoiceApplyCmd.AddCommand(
		newFinanceInvoiceApplyPageCmd(),
		newFinanceInvoiceApplyGetCmd(),
		newFinanceInvoiceApplyCreateCmd(),
		newFinanceInvoiceApplyUpdateCmd(),
		newFinanceInvoiceApplyDeleteCmd(),
		newFinanceInvoiceApplyUpdateStatusCmd(),
		newFinanceInvoiceApplySummaryCmd(),
		newFinanceInvoiceApplyExportExcelCmd(),
	)
}

var financeInvoiceApplyFilters = []cmdutil.FlagSpec{
	{Name: "no", Usage: "申请单号"},
	{Name: "partner-id", Usage: "合作伙伴 ID"},
	{Name: "partner-name", Usage: "合作伙伴名称"},
	{Name: "invoice-title", Usage: "发票抬头"},
	{Name: "tax-no", Usage: "税号"},
	{Name: "service-user-id", Usage: "财务负责人 ID"},
	{Name: "service-user-name", Usage: "财务负责人名称"},
	{Name: "approve-status", Usage: "审核状态"},
	{Name: "invoice-date", Usage: "开票日期"},
	{Name: "remark", Usage: "备注"},
	{Name: "creator-name", Usage: "创建人"},
	{Name: "updater-name", Usage: "更新人"},
	{Name: "create-time", Usage: "创建时间"},
	{Name: "update-time", Usage: "更新时间"},
	{Name: "custom-order", Usage: "前端自定义排序规则"},
	{Name: "keyword", Usage: "关键字"},
	{Name: "headers", Usage: "自定义导出表头"},
}

var financeInvoiceApplyPageCountFilters = []cmdutil.FlagSpec{
	{Name: "no", Usage: "申请单号"},
	{Name: "partner-id", Usage: "合作伙伴 ID"},
	{Name: "partner-name", Usage: "合作伙伴名称"},
	{Name: "invoice-title", Usage: "发票抬头"},
	{Name: "tax-no", Usage: "税号"},
	{Name: "service-user-id", Usage: "财务负责人 ID"},
	{Name: "service-user-name", Usage: "财务负责人名称"},
	{Name: "approve-status", Usage: "审核状态"},
	{Name: "invoice-date", Usage: "开票日期"},
	{Name: "remark", Usage: "备注"},
	{Name: "creator-name", Usage: "创建人"},
	{Name: "updater-name", Usage: "更新人"},
	{Name: "create-time", Usage: "创建时间"},
	{Name: "update-time", Usage: "更新时间"},
	{Name: "custom-order", Usage: "前端自定义排序规则"},
	{Name: "keyword", Usage: "关键字"},
}

var financeInvoiceApplyCmd = &cobra.Command{
	Use:   "invoice-apply",
	Short: "开票申请管理",
}

// ---- 分页查询 ----

func newFinanceInvoiceApplyPageCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "page",
		Short: "分页查询开票申请",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			collectFinanceFilters(cmd, params, financeInvoiceApplyFilters)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-invoice-apply/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询开票申请失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	addFinanceFilterFlags(c, financeInvoiceApplyFilters)
	return c
}


// ---- 获取详情 ----

func newFinanceInvoiceApplyGetCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "get",
		Short: "获取开票申请详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-invoice-apply/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取开票申请失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "申请单 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 创建 ----

func newFinanceInvoiceApplyCreateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "create",
		Short: "创建开票申请",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/finance-invoice-apply/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建开票申请失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 更新 ----

func newFinanceInvoiceApplyUpdateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update",
		Short: "更新开票申请",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/finance-invoice-apply/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新开票申请失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 删除 ----

func newFinanceInvoiceApplyDeleteCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "delete",
		Short: "删除开票申请",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/finance-invoice-apply/delete", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("删除开票申请失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "申请单 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 更新状态 ----

func newFinanceInvoiceApplyUpdateStatusCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update-status",
		Short: "更新开票申请状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/finance-invoice-apply/update-status", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新开票申请状态失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据（含 id 和 status）")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 汇总数据 ----

func newFinanceInvoiceApplySummaryCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "summary",
		Short: "获取开票申请合计（支持按筛选条件）",
		Long:  "返回 totalCount / checkCount / processCount / approveCount / rejectCount / totalAmount",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			collectFinanceFilters(cmd, params, financeInvoiceApplyPageCountFilters)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-invoice-apply/summary", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取开票申请合计失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	addFinanceFilterFlags(c, financeInvoiceApplyPageCountFilters)
	return c
}

// ---- 导出 Excel ----

func newFinanceInvoiceApplyExportExcelCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "export",
		Short: "导出开票申请 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			collectFinanceFilters(cmd, params, financeInvoiceApplyFilters)

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-invoice-apply/export-excel", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出开票申请失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	addFinanceFilterFlags(c, financeInvoiceApplyFilters)
	return c
}
