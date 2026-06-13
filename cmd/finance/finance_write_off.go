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
		newFinanceWriteOffListCmd(),
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

// ---- 分页查询 ----

func newFinanceWriteOffListCmd() *cobra.Command {
	var pageNo, pageSize int
	var writeOffNo, customerId, supplierId, status, accountId, reviewerId string

	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询核销单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlags(cmd, params, "write-off-no", "customer-id", "supplier-id", "status", "account-id", "reviewer-id")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-write-off/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询核销单失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().StringVar(&writeOffNo, "write-off-no", "", "核销单号")
	c.Flags().StringVar(&customerId, "customer-id", "", "客户 ID")
	c.Flags().StringVar(&supplierId, "supplier-id", "", "供应商 ID")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&accountId, "account-id", "", "账户 ID")
	c.Flags().StringVar(&reviewerId, "reviewer-id", "", "审核人 ID")
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
		Short: "获取核销单汇总数据",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-write-off/summary", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取核销单汇总失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
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
	var writeOffNo, customerId, supplierId, status, accountId, reviewerId string

	c := &cobra.Command{
		Use:   "export",
		Short: "导出核销单 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "write-off-no", "customer-id", "supplier-id", "status", "account-id", "reviewer-id")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-write-off/export-excel", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出核销单失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	c.Flags().StringVar(&writeOffNo, "write-off-no", "", "核销单号")
	c.Flags().StringVar(&customerId, "customer-id", "", "客户 ID")
	c.Flags().StringVar(&supplierId, "supplier-id", "", "供应商 ID")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&accountId, "account-id", "", "账户 ID")
	c.Flags().StringVar(&reviewerId, "reviewer-id", "", "审核人 ID")
	return c
}
