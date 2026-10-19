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
		newFinanceInvoicePageCmd(),
		newFinanceInvoiceGetCmd(),
		newFinanceInvoiceSummaryCmd(),
		newFinanceInvoiceExportExcelCmd(),
		newFinanceInvoiceCreateCmd(),
		newFinanceInvoiceUpdateCmd(),
		newFinanceInvoiceUpdateStatusCmd(),
		newFinanceInvoiceDeleteCmd(),
	)
}

// financeInvoiceFilters 包含所有可筛选字段（含 headers，用于 page / export）。
var financeInvoiceFilters = []cmdutil.FlagSpec{
	{Name: "no", Usage: "收开票号"},
	{Name: "type", Usage: "业务类别(INVOICE_PAYMENT=收票/INVOICE_RECEIPT=开票/...)"},
	{Name: "invoice-type", Usage: "发票类型(SPECIAL_INVOICE/GENERAL_INVOICE/RECEIPT)"},
	{Name: "invoice-no", Usage: "发票号"},
	{Name: "partner-id", Usage: "合作伙伴 ID"},
	{Name: "partner-name", Usage: "客户/供应商名称"},
	{Name: "service-user-id", Usage: "财务负责人 ID"},
	{Name: "service-user-name", Usage: "财务负责人名称"},
	{Name: "status", Usage: "核销状态(NOT_WRITE_OFF/PARTIAL_WRITE_OFF/FULLY_WRITE_OFF)"},
	{Name: "approve-status", Usage: "审核状态"},
	{Name: "invoice-date", Usage: "收开票日期"},
	{Name: "remark", Usage: "备注"},
	{Name: "creator-name", Usage: "创建人"},
	{Name: "updater-name", Usage: "更新人"},
	{Name: "create-time", Usage: "创建时间"},
	{Name: "update-time", Usage: "更新时间"},
	{Name: "keyword", Usage: "关键字搜索(单号/名称/发票号)"},
	{Name: "custom-order", Usage: "前端自定义排序规则"},
	{Name: "headers", Usage: "自定义导出表头"},
}

// financeInvoiceQueryFilters 去掉了 headers，用于 summary。
var financeInvoiceQueryFilters = []cmdutil.FlagSpec{
	{Name: "no", Usage: "收开票号"},
	{Name: "type", Usage: "业务类别(INVOICE_PAYMENT=收票/INVOICE_RECEIPT=开票/...)"},
	{Name: "invoice-type", Usage: "发票类型(SPECIAL_INVOICE/GENERAL_INVOICE/RECEIPT)"},
	{Name: "invoice-no", Usage: "发票号"},
	{Name: "partner-id", Usage: "合作伙伴 ID"},
	{Name: "partner-name", Usage: "客户/供应商名称"},
	{Name: "service-user-id", Usage: "财务负责人 ID"},
	{Name: "service-user-name", Usage: "财务负责人名称"},
	{Name: "status", Usage: "核销状态(NOT_WRITE_OFF/PARTIAL_WRITE_OFF/FULLY_WRITE_OFF)"},
	{Name: "approve-status", Usage: "审核状态"},
	{Name: "invoice-date", Usage: "收开票日期"},
	{Name: "remark", Usage: "备注"},
	{Name: "creator-name", Usage: "创建人"},
	{Name: "updater-name", Usage: "更新人"},
	{Name: "create-time", Usage: "创建时间"},
	{Name: "update-time", Usage: "更新时间"},
	{Name: "keyword", Usage: "关键字搜索(单号/名称/发票号)"},
	{Name: "custom-order", Usage: "前端自定义排序规则"},
}

// ---- 分页查询 ----

func newFinanceInvoicePageCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "page",
		Short: "分页查询收开票(type=INVOICE_PAYMENT 收票 / INVOICE_RECEIPT 开票)",
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

// ---- 获取详情 ----

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

// ---- 汇总数据 ----

func newFinanceInvoiceSummaryCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "summary",
		Short: "获取收开票合计",
		Long:  "返回 totalCount / unSettledCount / partSettledCount / settledCount / totalAmount / checkedAmount / unCheckedAmount",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			collectFinanceFilters(cmd, params, financeInvoiceQueryFilters)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-invoice/summary", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取收开票合计失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	addFinanceFilterFlags(c, financeInvoiceQueryFilters)
	return c
}

// ---- 导出 Excel ----

func newFinanceInvoiceExportExcelCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "export",
		Short: "导出收开票 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]any{}
			collectFinanceFilters(cmd, params, financeInvoiceFilters)
			return cmdutil.RunExport(cmd, "/erp/finance-invoice/export-excel", params)
		},
	}
	addFinanceFilterFlags(c, financeInvoiceFilters)
	return c
}

// ---- 创建 ----

func newFinanceInvoiceCreateCmd() *cobra.Command {
	var data string
	c := &cobra.Command{
		Use:   "create",
		Short: "创建收开票",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/finance-invoice/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建收开票失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 更新 ----

func newFinanceInvoiceUpdateCmd() *cobra.Command {
	var data string
	c := &cobra.Command{
		Use:   "update",
		Short: "更新收开票",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/finance-invoice/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新收开票失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 更新状态（query 参数：id + status） ----

func newFinanceInvoiceUpdateStatusCmd() *cobra.Command {
	var id int64
	var status int
	c := &cobra.Command{
		Use:   "update-status",
		Short: "更新收开票审核状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().PutQuery(context.Background(), "/erp/finance-invoice/update-status",
				map[string]any{"id": id, "status": status})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新收开票状态失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "收开票 ID")
	c.Flags().IntVar(&status, "status", 0, "状态值")
	_ = c.MarkFlagRequired("id")
	_ = c.MarkFlagRequired("status")
	return c
}

// ---- 删除（query 参数：ids） ----

func newFinanceInvoiceDeleteCmd() *cobra.Command {
	var ids string
	c := &cobra.Command{
		Use:   "delete",
		Short: "删除收开票",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/finance-invoice/delete", map[string]any{"ids": ids})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("删除收开票失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&ids, "ids", "", "ID 列表（逗号分隔）")
	_ = c.MarkFlagRequired("ids")
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
