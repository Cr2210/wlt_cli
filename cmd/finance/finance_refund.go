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
	financeCmd.AddCommand(financeRefundCmd)
	financeRefundCmd.AddCommand(
		newFinanceRefundPageCmd(),
		newFinanceRefundPageCountCmd(),
		newFinanceRefundGetCmd(),
		newFinanceRefundCreateCmd(),
		newFinanceRefundUpdateCmd(),
		newFinanceRefundDeleteCmd(),
		newFinanceRefundUpdateStatusCmd(),
		newFinanceRefundSummaryCmd(),
		newFinanceRefundExportExcelCmd(),
	)
}

var financeRefundCmd = &cobra.Command{
	Use:   "refund",
	Short: "退款单管理",
}

// financeRefundFilters 包含所有可筛选字段（含 headers，用于 page / export）。
// 供应商退款: type=SUPPLIER_REFUND
// 客户退款: type=CUSTOMER_REFUND
var financeRefundFilters = []cmdutil.FlagSpec{
	{Name: "no", Usage: "退款单编号"},
	{Name: "type", Usage: "业务类别(SUPPLIER_REFUND=供应商退款/CUSTOMER_REFUND=客户退款/...)"},
	{Name: "pay-date", Usage: "付款日期"},
	{Name: "account-id", Usage: "账户 ID"},
	{Name: "account-name", Usage: "账户名称"},
	{Name: "account-no", Usage: "账户账号"},
	{Name: "partner-id", Usage: "供应商/客户 ID"},
	{Name: "partner-name", Usage: "供应商/客户名称"},
	{Name: "service-user-id", Usage: "财务负责人 ID"},
	{Name: "service-user-name", Usage: "财务负责人名称"},
	{Name: "status", Usage: "核销状态(NOT_WRITE_OFF/PARTIAL_WRITE_OFF/FULLY_WRITE_OFF)"},
	{Name: "approve-status", Usage: "审核状态"},
	{Name: "remark", Usage: "备注"},
	{Name: "creator-name", Usage: "创建人"},
	{Name: "updater-name", Usage: "更新人"},
	{Name: "create-time", Usage: "创建时间"},
	{Name: "update-time", Usage: "更新时间"},
	{Name: "keyword", Usage: "关键字搜索(退款单编号/账户/供应商/客户)"},
	{Name: "custom-order", Usage: "前端自定义排序规则"},
	{Name: "headers", Usage: "自定义导出表头"},
}

// financeRefundPageCountFilters 去掉了 headers。
var financeRefundPageCountFilters = []cmdutil.FlagSpec{
	{Name: "no", Usage: "退款单编号"},
	{Name: "type", Usage: "业务类别(SUPPLIER_REFUND=供应商退款/CUSTOMER_REFUND=客户退款/...)"},
	{Name: "pay-date", Usage: "付款日期"},
	{Name: "account-id", Usage: "账户 ID"},
	{Name: "account-name", Usage: "账户名称"},
	{Name: "account-no", Usage: "账户账号"},
	{Name: "partner-id", Usage: "供应商/客户 ID"},
	{Name: "partner-name", Usage: "供应商/客户名称"},
	{Name: "service-user-id", Usage: "财务负责人 ID"},
	{Name: "service-user-name", Usage: "财务负责人名称"},
	{Name: "status", Usage: "核销状态(NOT_WRITE_OFF/PARTIAL_WRITE_OFF/FULLY_WRITE_OFF)"},
	{Name: "approve-status", Usage: "审核状态"},
	{Name: "remark", Usage: "备注"},
	{Name: "creator-name", Usage: "创建人"},
	{Name: "updater-name", Usage: "更新人"},
	{Name: "create-time", Usage: "创建时间"},
	{Name: "update-time", Usage: "更新时间"},
	{Name: "keyword", Usage: "关键字搜索(退款单编号/账户/供应商/客户)"},
	{Name: "custom-order", Usage: "前端自定义排序规则"},
}

// ---- 分页查询 ----

func newFinanceRefundPageCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "page",
		Short: "分页查询退款单(type=SUPPLIER_REFUND 供应商退款 / CUSTOMER_REFUND 客户退款)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			collectFinanceFilters(cmd, params, financeRefundFilters)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-refund/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询退款单失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	addFinanceFilterFlags(c, financeRefundFilters)
	return c
}

// ---- 分页计数 ----

func newFinanceRefundPageCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "page-count",
		Short: "按筛选统计退款单数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			collectFinanceFilters(cmd, params, financeRefundPageCountFilters)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-refund/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计退款单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	addFinanceFilterFlags(c, financeRefundPageCountFilters)
	return c
}

// ---- 获取详情 ----

func newFinanceRefundGetCmd() *cobra.Command {
	var id int64
	var no string

	c := &cobra.Command{
		Use:   "get",
		Short: "获取退款单详情（id 或 no 任选其一）",
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
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-refund/get", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取退款单详情失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "退款单 ID")
	c.Flags().StringVar(&no, "no", "", "退款单编号")
	return c
}

// ---- 创建 ----

func newFinanceRefundCreateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "create",
		Short: "创建退款单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/finance-refund/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建退款单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 更新 ----

func newFinanceRefundUpdateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update",
		Short: "更新退款单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/finance-refund/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新退款单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 删除 ----

func newFinanceRefundDeleteCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "delete",
		Short: "删除退款单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/finance-refund/delete", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("删除退款单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "退款单 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 更新状态 ----

func newFinanceRefundUpdateStatusCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update-status",
		Short: "更新退款单状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/finance-refund/update-status", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新退款单状态失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据（含 id 和 status）")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 汇总数据 ----

func newFinanceRefundSummaryCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "summary",
		Short: "获取退款单汇总数据（支持按筛选条件）",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			collectFinanceFilters(cmd, params, financeRefundPageCountFilters)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-refund/summary", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取退款单汇总失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	addFinanceFilterFlags(c, financeRefundPageCountFilters)
	return c
}

// ---- 导出 Excel ----

func newFinanceRefundExportExcelCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "export",
		Short: "导出退款单 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			collectFinanceFilters(cmd, params, financeRefundFilters)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-refund/export-excel", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出退款单失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	addFinanceFilterFlags(c, financeRefundFilters)
	return c
}
