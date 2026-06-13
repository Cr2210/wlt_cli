package waybill

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var waybillSourceCmd = &cobra.Command{
	Use:   "source",
	Short: "运单管理",
}

func init() {
	waybillCmd.AddCommand(waybillSourceCmd)
	waybillSourceCmd.AddCommand(
		newWaybillSourceListCmd(),
		newWaybillSourcePageCountCmd(),
		newWaybillSourceGetCmd(),
		newWaybillSourceGetEventsCmd(),
		newWaybillSourceGetOrderInfoCmd(),
		newWaybillSourceCreateCmd(),
		newWaybillSourceUpdateCmd(),
		newWaybillSourceDeleteCmd(),
		newWaybillSourceDeleteListCmd(),
		newWaybillSourceLoadCmd(),
		newWaybillSourceUnloadCmd(),
		newWaybillSourceSignCmd(),
		newWaybillSourceSignBatchCmd(),
		newWaybillSourceBatchSignCmd(),
		newWaybillSourceImportCmd(),
		newWaybillSourceExportExcelCmd(),
	)
}

// ---- 分页查询 ----

func newWaybillSourceListCmd() *cobra.Command {
	var pageNo, pageSize int
	var waybillNo, licensePlate, status, warehouseId, customerId, supplierId string

	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询运单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlags(cmd, params, "waybill-no", "license-plate", "status", "warehouse-id", "customer-id", "supplier-id")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/waybill-source/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询运单失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().StringVar(&waybillNo, "waybill-no", "", "运单号")
	c.Flags().StringVar(&licensePlate, "license-plate", "", "车牌号")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&warehouseId, "warehouse-id", "", "仓库 ID")
	c.Flags().StringVar(&customerId, "customer-id", "", "客户 ID")
	c.Flags().StringVar(&supplierId, "supplier-id", "", "供应商 ID")
	return c
}

// ---- 分页计数 ----

func newWaybillSourcePageCountCmd() *cobra.Command {
	var waybillNo, licensePlate, status, warehouseId, customerId, supplierId string

	c := &cobra.Command{
		Use:   "page-count",
		Short: "统计运单数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "waybill-no", "license-plate", "status", "warehouse-id", "customer-id", "supplier-id")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/waybill-source/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计运单数量失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&waybillNo, "waybill-no", "", "运单号")
	c.Flags().StringVar(&licensePlate, "license-plate", "", "车牌号")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&warehouseId, "warehouse-id", "", "仓库 ID")
	c.Flags().StringVar(&customerId, "customer-id", "", "客户 ID")
	c.Flags().StringVar(&supplierId, "supplier-id", "", "供应商 ID")
	return c
}

// ---- 获取详情 ----

func newWaybillSourceGetCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "get",
		Short: "获取运单详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/waybill-source/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取运单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "运单 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 获取运单事件 ----

func newWaybillSourceGetEventsCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "get-events",
		Short: "获取运单事件",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/waybill-source/get-events", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取运单事件失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "运单 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 获取订单信息 ----

func newWaybillSourceGetOrderInfoCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "get-order-info",
		Short: "获取运单的订单信息",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/waybill-source/get-order-info", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取订单信息失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "运单 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 创建 ----

func newWaybillSourceCreateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "create",
		Short: "创建运单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/waybill-source/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建运单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 更新 ----

func newWaybillSourceUpdateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update",
		Short: "更新运单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/waybill-source/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新运单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 删除单个 ----

func newWaybillSourceDeleteCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "delete",
		Short: "删除运单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/waybill-source/delete", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("删除运单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "运单 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 批量删除 ----

func newWaybillSourceDeleteListCmd() *cobra.Command {
	var ids string

	c := &cobra.Command{
		Use:   "delete-list",
		Short: "批量删除运单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/waybill-source/delete-list", map[string]any{"ids": ids})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("批量删除运单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&ids, "ids", "", "运单 ID 列表（逗号分隔）")
	_ = c.MarkFlagRequired("ids")
	return c
}

// ---- 装车 ----

func newWaybillSourceLoadCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "load",
		Short: "运单装车",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/waybill-source/load", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("运单装车失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 卸车 ----

func newWaybillSourceUnloadCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "unload",
		Short: "运单卸车",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/waybill-source/unload", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("运单卸车失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 签收 ----

func newWaybillSourceSignCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "sign",
		Short: "运单签收",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/waybill-source/sign", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("运单签收失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 批量签收 ----

func newWaybillSourceSignBatchCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "sign-batch",
		Short: "批量签收运单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/waybill-source/sign-batch", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("批量签收失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 批量签收（列表） ----

func newWaybillSourceBatchSignCmd() *cobra.Command {
	var ids string

	c := &cobra.Command{
		Use:   "batch-sign",
		Short: "批量签收运单（按ID列表）",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/waybill-source/batch-sign", map[string]any{"ids": ids})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("批量签收失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&ids, "ids", "", "运单 ID 列表（逗号分隔）")
	_ = c.MarkFlagRequired("ids")
	return c
}

// ---- 导入 ----

func newWaybillSourceImportCmd() *cobra.Command {
	var filePath string

	c := &cobra.Command{
		Use:   "import",
		Short: "导入运单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			// 读取文件内容
			data, err := os.ReadFile(filePath)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("读取文件失败: %s", err), "")
			}

			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/waybill-source/import", data)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导入运单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&filePath, "file", "", "Excel 文件路径")
	_ = c.MarkFlagRequired("file")
	return c
}

// ---- 导出 Excel ----

func newWaybillSourceExportExcelCmd() *cobra.Command {
	var waybillNo, licensePlate, status, warehouseId, customerId, supplierId string

	c := &cobra.Command{
		Use:   "export",
		Short: "导出运单 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "waybill-no", "license-plate", "status", "warehouse-id", "customer-id", "supplier-id")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/waybill-source/export-excel", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出运单失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	c.Flags().StringVar(&waybillNo, "waybill-no", "", "运单号")
	c.Flags().StringVar(&licensePlate, "license-plate", "", "车牌号")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&warehouseId, "warehouse-id", "", "仓库 ID")
	c.Flags().StringVar(&customerId, "customer-id", "", "客户 ID")
	c.Flags().StringVar(&supplierId, "supplier-id", "", "供应商 ID")
	return c
}
