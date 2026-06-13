package order

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var orderPlanCmd = &cobra.Command{
	Use:   "plan",
	Short: "订单计划管理",
}

func init() {
	orderCmd.AddCommand(orderPlanCmd)
	orderPlanCmd.AddCommand(
		newOrderPlanListCmd(),
		newOrderPlanPageCountCmd(),
		newOrderPlanGetCmd(),
		newOrderPlanCreateCmd(),
		newOrderPlanUpdateCmd(),
		newOrderPlanDeleteCmd(),
		newOrderPlanUpdateStatusCmd(),
		newOrderPlanCancelCmd(),
		newOrderPlanReopenCmd(),
		newOrderPlanCompleteCmd(),
		newOrderPlanExportExcelCmd(),
	)
}

// ---- 分页查询 ----

func newOrderPlanListCmd() *cobra.Command {
	var pageNo, pageSize int
	var planNo, customerId, status, warehouseId string

	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询订单计划",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlags(cmd, params, "plan-no", "customer-id", "status", "warehouse-id")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/order-plan/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询订单计划失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().StringVar(&planNo, "plan-no", "", "计划单号")
	c.Flags().StringVar(&customerId, "customer-id", "", "客户 ID")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&warehouseId, "warehouse-id", "", "仓库 ID")
	return c
}

// ---- 分页计数 ----

func newOrderPlanPageCountCmd() *cobra.Command {
	var planNo, customerId, status, warehouseId string

	c := &cobra.Command{
		Use:   "page-count",
		Short: "统计订单计划数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "plan-no", "customer-id", "status", "warehouse-id")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/order-plan/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计订单计划数量失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&planNo, "plan-no", "", "计划单号")
	c.Flags().StringVar(&customerId, "customer-id", "", "客户 ID")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&warehouseId, "warehouse-id", "", "仓库 ID")
	return c
}

// ---- 获取详情 ----

func newOrderPlanGetCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "get",
		Short: "获取订单计划详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/order-plan/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取订单计划失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "计划 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 创建 ----

func newOrderPlanCreateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "create",
		Short: "创建订单计划",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/order-plan/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建订单计划失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 更新 ----

func newOrderPlanUpdateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update",
		Short: "更新订单计划",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/order-plan/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新订单计划失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 删除 ----

func newOrderPlanDeleteCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "delete",
		Short: "删除订单计划",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/order-plan/delete", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("删除订单计划失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "计划 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 更新状态 ----

func newOrderPlanUpdateStatusCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update-status",
		Short: "更新订单计划状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/order-plan/update-status", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新订单计划状态失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据（含 id 和 status）")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 取消订单计划 ----

func newOrderPlanCancelCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "cancel",
		Short: "取消订单计划",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/order-plan/cancel", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("取消订单计划失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 重新打开 ----

func newOrderPlanReopenCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "reopen",
		Short: "重新打开订单计划",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/order-plan/reopen", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("重新打开订单计划失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 完成订单计划 ----

func newOrderPlanCompleteCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "complete",
		Short: "完成订单计划",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/order-plan/complete", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("完成订单计划失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 导出 Excel ----

func newOrderPlanExportExcelCmd() *cobra.Command {
	var planNo, customerId, status, warehouseId string

	c := &cobra.Command{
		Use:   "export",
		Short: "导出订单计划 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "plan-no", "customer-id", "status", "warehouse-id")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/order-plan/export-excel", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出订单计划失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	c.Flags().StringVar(&planNo, "plan-no", "", "计划单号")
	c.Flags().StringVar(&customerId, "customer-id", "", "客户 ID")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&warehouseId, "warehouse-id", "", "仓库 ID")
	return c
}
