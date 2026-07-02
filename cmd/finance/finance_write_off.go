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
	financeCmd.AddCommand(financeWriteOffCmd)
	financeWriteOffCmd.AddCommand(
		newFinanceWriteOffPageCmd(),
		newFinanceWriteOffPageCountCmd(),
		newFinanceWriteOffGetCmd(),
		newFinanceWriteOffCreateCmd(),
		newFinanceWriteOffUpdateCmd(),
		newFinanceWriteOffDeleteCmd(),
		newFinanceWriteOffUpdateStatusCmd(),
		newFinanceWriteOffSummaryCmd(),
		newFinanceWriteOffGetItemRelateCmd(),
		newFinanceWriteOffExportExcelCmd(),
	)
}

var financeWriteOffCmd = &cobra.Command{
	Use:   "write-off",
	Short: "核销单管理",
}

// financeWriteOffFilters 包含所有可筛选字段（含 headers，用于 page / export）。
// 资金核销: type=WRITE_OFF_PURCHASE/WRITE_OFF_SALE, writeOffType=AMOUNT
// 发票核销: type=WRITE_OFF_IN/WRITE_OFF_OUT, writeOffType=INVOICE
var financeWriteOffFilters = []cmdutil.FlagSpec{
	{Name: "no", Usage: "核销单号"},
	{Name: "type", Usage: "业务类别(WRITE_OFF_PURCHASE/WRITE_OFF_SALE/WRITE_OFF_IN/WRITE_OFF_OUT/...)"},
	{Name: "write-off-type", Usage: "核销类型(AMOUNT=资金核销/INVOICE=发票核销)"},
	{Name: "write-off-date", Usage: "核销日期"},
	{Name: "approve-status", Usage: "审核状态"},
	{Name: "partner-id", Usage: "合作伙伴 ID"},
	{Name: "partner-name", Usage: "客户/供应商名称"},
	{Name: "service-user-id", Usage: "财务负责人 ID"},
	{Name: "service-user-name", Usage: "财务负责人名称"},
	{Name: "source-no", Usage: "源单号"},
	{Name: "invoice-no", Usage: "发票号"},
	{Name: "target-no", Usage: "目标单号"},
	{Name: "remark", Usage: "备注"},
	{Name: "creator-name", Usage: "创建人"},
	{Name: "updater-name", Usage: "更新人"},
	{Name: "create-time", Usage: "创建时间"},
	{Name: "update-time", Usage: "更新时间"},
	{Name: "keyword", Usage: "关键字搜索(单号/名称/发票号)"},
	{Name: "custom-order", Usage: "前端自定义排序规则"},
	{Name: "headers", Usage: "自定义导出表头"},
}

// financeWriteOffPageCountFilters 去掉了 headers。
var financeWriteOffPageCountFilters = []cmdutil.FlagSpec{
	{Name: "no", Usage: "核销单号"},
	{Name: "type", Usage: "业务类别(WRITE_OFF_PURCHASE/WRITE_OFF_SALE/WRITE_OFF_IN/WRITE_OFF_OUT/...)"},
	{Name: "write-off-type", Usage: "核销类型(AMOUNT=资金核销/INVOICE=发票核销)"},
	{Name: "write-off-date", Usage: "核销日期"},
	{Name: "approve-status", Usage: "审核状态"},
	{Name: "partner-id", Usage: "合作伙伴 ID"},
	{Name: "partner-name", Usage: "客户/供应商名称"},
	{Name: "service-user-id", Usage: "财务负责人 ID"},
	{Name: "service-user-name", Usage: "财务负责人名称"},
	{Name: "source-no", Usage: "源单号"},
	{Name: "invoice-no", Usage: "发票号"},
	{Name: "target-no", Usage: "目标单号"},
	{Name: "remark", Usage: "备注"},
	{Name: "creator-name", Usage: "创建人"},
	{Name: "updater-name", Usage: "更新人"},
	{Name: "create-time", Usage: "创建时间"},
	{Name: "update-time", Usage: "更新时间"},
	{Name: "keyword", Usage: "关键字搜索(单号/名称/发票号)"},
	{Name: "custom-order", Usage: "前端自定义排序规则"},
}

// ---- 分页查询 ----

func newFinanceWriteOffPageCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "page",
		Short: "分页查询核销单(资金核销 WRITE_OFF_PURCHASE/SALE + 发票核销 WRITE_OFF_IN/OUT)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			collectFinanceFilters(cmd, params, financeWriteOffFilters)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-write-off/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询核销单失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	addFinanceFilterFlags(c, financeWriteOffFilters)
	return c
}

// ---- 分页计数 ----

func newFinanceWriteOffPageCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "page-count",
		Short: "按筛选统计核销单数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			collectFinanceFilters(cmd, params, financeWriteOffPageCountFilters)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-write-off/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计核销单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	addFinanceFilterFlags(c, financeWriteOffPageCountFilters)
	return c
}

// ---- 获取详情 ----

func newFinanceWriteOffGetCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "get",
		Short: "获取核销单详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-write-off/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取核销单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "核销单 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 创建 ----

func newFinanceWriteOffCreateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "create",
		Short: "创建核销单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/finance-write-off/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建核销单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 更新 ----

func newFinanceWriteOffUpdateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update",
		Short: "更新核销单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/finance-write-off/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新核销单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 删除 ----

func newFinanceWriteOffDeleteCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "delete",
		Short: "删除核销单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/finance-write-off/delete", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("删除核销单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "核销单 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 更新状态 ----

func newFinanceWriteOffUpdateStatusCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update-status",
		Short: "更新核销单状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/finance-write-off/update-status", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新核销单状态失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据（含 id 和 status）")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 汇总数据 ----

func newFinanceWriteOffSummaryCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "summary",
		Short: "获取核销单汇总数据（支持按筛选条件）",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			collectFinanceFilters(cmd, params, financeWriteOffPageCountFilters)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-write-off/summary", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取核销单汇总失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	addFinanceFilterFlags(c, financeWriteOffPageCountFilters)
	return c
}

// ---- 获取项目关联 ----

func newFinanceWriteOffGetItemRelateCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "get-item-relate",
		Short: "获取核销单的项目关联",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-write-off/get-item-relate", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取项目关联失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "核销单 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 导出 Excel ----

func newFinanceWriteOffExportExcelCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "export",
		Short: "导出核销单 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			collectFinanceFilters(cmd, params, financeWriteOffFilters)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-write-off/export-excel", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出核销单失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	addFinanceFilterFlags(c, financeWriteOffFilters)
	return c
}
