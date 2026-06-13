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
		newFinanceRefundListCmd(),
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

// ---- 分页查询 ----

func newFinanceRefundListCmd() *cobra.Command {
	var pageNo, pageSize int
	var refundNo, customerId, supplierId, status, accountId, reviewerId string

	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询退款单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlags(cmd, params, "refund-no", "customer-id", "supplier-id", "status", "account-id", "reviewer-id")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-refund/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询退款单失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().StringVar(&refundNo, "refund-no", "", "退款单号")
	c.Flags().StringVar(&customerId, "customer-id", "", "客户 ID")
	c.Flags().StringVar(&supplierId, "supplier-id", "", "供应商 ID")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&accountId, "account-id", "", "账户 ID")
	c.Flags().StringVar(&reviewerId, "reviewer-id", "", "审核人 ID")
	return c
}

// ---- 获取详情 ----

func newFinanceRefundGetCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "get",
		Short: "获取退款单详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-refund/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取退款单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "退款单 ID")
	_ = c.MarkFlagRequired("id")
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
		Short: "获取退款单汇总数据",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-refund/summary", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取退款单汇总失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

// ---- 导出 Excel ----

func newFinanceRefundExportExcelCmd() *cobra.Command {
	var refundNo, customerId, supplierId, status, accountId, reviewerId string

	c := &cobra.Command{
		Use:   "export",
		Short: "导出退款单 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "refund-no", "customer-id", "supplier-id", "status", "account-id", "reviewer-id")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-refund/export-excel", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出退款单失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	c.Flags().StringVar(&refundNo, "refund-no", "", "退款单号")
	c.Flags().StringVar(&customerId, "customer-id", "", "客户 ID")
	c.Flags().StringVar(&supplierId, "supplier-id", "", "供应商 ID")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&accountId, "account-id", "", "账户 ID")
	c.Flags().StringVar(&reviewerId, "reviewer-id", "", "审核人 ID")
	return c
}
