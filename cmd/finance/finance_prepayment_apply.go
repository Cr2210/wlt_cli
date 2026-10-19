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
		newFinancePrepaymentApplyPageCmd(),
		newFinancePrepaymentApplyGetCmd(),
		newFinancePrepaymentApplySummaryCmd(),
		newFinancePrepaymentApplyExportExcelCmd(),
		newFinancePrepaymentApplyCreateCmd(),
		newFinancePrepaymentApplyUpdateCmd(),
		newFinancePrepaymentApplyUpdateStatusCmd(),
		newFinancePrepaymentApplyDeleteCmd(),
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
	{Name: "headers", Usage: "自定义导出表头"},
}

func newFinancePrepaymentApplyPageCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "page",
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

// ---- 导出 Excel ----

func newFinancePrepaymentApplyExportExcelCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "export",
		Short: "导出预付申请 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]any{}
			collectFinanceFilters(cmd, params, financePrepaymentApplyFilters)
			return cmdutil.RunExport(cmd, "/erp/finance-prepayment-apply/export-excel", params)
		},
	}
	addFinanceFilterFlags(c, financePrepaymentApplyFilters)
	return c
}

// ---- 创建 ----

func newFinancePrepaymentApplyCreateCmd() *cobra.Command {
	var data string
	c := &cobra.Command{
		Use:   "create",
		Short: "创建预付申请",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/finance-prepayment-apply/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建预付申请失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 更新 ----

func newFinancePrepaymentApplyUpdateCmd() *cobra.Command {
	var data string
	c := &cobra.Command{
		Use:   "update",
		Short: "更新预付申请",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/finance-prepayment-apply/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新预付申请失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 处理申请（query 参数：id + approveStatus） ----

func newFinancePrepaymentApplyUpdateStatusCmd() *cobra.Command {
	var id int64
	var approveStatus int
	c := &cobra.Command{
		Use:   "update-status",
		Short: "处理预付申请",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().PutQuery(context.Background(), "/erp/finance-prepayment-apply/update-status",
				map[string]any{"id": id, "approveStatus": approveStatus})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("处理预付申请失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "预付申请 ID")
	c.Flags().IntVar(&approveStatus, "approve-status", 0, "审核状态值")
	_ = c.MarkFlagRequired("id")
	_ = c.MarkFlagRequired("approve-status")
	return c
}

// ---- 删除（query 参数：ids） ----

func newFinancePrepaymentApplyDeleteCmd() *cobra.Command {
	var ids string
	c := &cobra.Command{
		Use:   "delete",
		Short: "删除预付申请",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/finance-prepayment-apply/delete", map[string]any{"ids": ids})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("删除预付申请失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&ids, "ids", "", "ID 列表（逗号分隔）")
	_ = c.MarkFlagRequired("ids")
	return c
}
