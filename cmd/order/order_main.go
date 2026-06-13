package order

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var orderMainCmd = &cobra.Command{
	Use:   "main",
	Short: "订单管理",
}

func init() {
	orderCmd.AddCommand(orderMainCmd)
	orderMainCmd.AddCommand(
		newOrderMainListCmd(),
		newOrderMainPageCountCmd(),
		newOrderMainGetCmd(),
		newOrderMainGetLinkOrderCmd(),
		newOrderMainCreateCmd(),
		newOrderMainUpdateCmd(),
		newOrderMainDeleteCmd(),
		newOrderMainUpdateStatusCmd(),
		newOrderMainCancelCmd(),
		newOrderMainReopenCmd(),
		newOrderMainCompleteCmd(),
		newOrderMainLinkWaybillCmd(),
		newOrderMainUnlinkWaybillCmd(),
		newOrderMainExportExcelCmd(),
	)
}

// ---- 分页查询 ----

func newOrderMainListCmd() *cobra.Command {
	var pageNo, pageSize int
	var orderNo, customerId, status, warehouseId string

	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询订单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlags(cmd, params, "order-no", "customer-id", "status", "warehouse-id")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/order/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询订单失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().StringVar(&orderNo, "order-no", "", "订单号")
	c.Flags().StringVar(&customerId, "customer-id", "", "客户 ID")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&warehouseId, "warehouse-id", "", "仓库 ID")
	return c
}

// ---- 分页计数 ----

func newOrderMainPageCountCmd() *cobra.Command {
	var orderNo, customerId, status, warehouseId string

	c := &cobra.Command{
		Use:   "page-count",
		Short: "统计订单数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "order-no", "customer-id", "status", "warehouse-id")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/order/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计订单数量失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&orderNo, "order-no", "", "订单号")
	c.Flags().StringVar(&customerId, "customer-id", "", "客户 ID")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&warehouseId, "warehouse-id", "", "仓库 ID")
	return c
}

// ---- 获取详情 ----

func newOrderMainGetCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "get",
		Short: "获取订单详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/order/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取订单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "订单 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 获取关联订单 ----

func newOrderMainGetLinkOrderCmd() *cobra.Command {
	var orderId int64

	c := &cobra.Command{
		Use:   "get-linkorder-by-orderId",
		Short: "获取订单的关联运单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/order/get-linkorder-by-orderId", map[string]any{"orderId": orderId})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取关联运单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&orderId, "order-id", 0, "订单 ID")
	_ = c.MarkFlagRequired("order-id")
	return c
}

// ---- 创建 ----

func newOrderMainCreateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "create",
		Short: "创建订单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/order/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建订单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 更新 ----

func newOrderMainUpdateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update",
		Short: "更新订单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/order/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新订单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 删除 ----

func newOrderMainDeleteCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "delete",
		Short: "删除订单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/order/delete", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("删除订单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "订单 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 更新状态 ----

func newOrderMainUpdateStatusCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update-status",
		Short: "更新订单状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/order/update-status", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新订单状态失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据（含 id 和 status）")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 取消订单 ----

func newOrderMainCancelCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "cancel",
		Short: "取消订单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/order/cancel", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("取消订单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 重新打开 ----

func newOrderMainReopenCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "reopen",
		Short: "重新打开订单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/order/reopen", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("重新打开订单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 完成订单 ----

func newOrderMainCompleteCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "complete",
		Short: "完成订单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/order/complete", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("完成订单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 关联运单 ----

func newOrderMainLinkWaybillCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "link-waybill",
		Short: "关联运单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/order/link-waybill", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("关联运单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 取消关联运单 ----

func newOrderMainUnlinkWaybillCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "unlink-waybill",
		Short: "取消关联运单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/order/unlink-waybill", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("取消关联运单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 导出 Excel ----

func newOrderMainExportExcelCmd() *cobra.Command {
	var orderNo, customerId, status, warehouseId string

	c := &cobra.Command{
		Use:   "export",
		Short: "导出订单 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "order-no", "customer-id", "status", "warehouse-id")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/order/export-excel", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出订单失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	c.Flags().StringVar(&orderNo, "order-no", "", "订单号")
	c.Flags().StringVar(&customerId, "customer-id", "", "客户 ID")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&warehouseId, "warehouse-id", "", "仓库 ID")
	return c
}
