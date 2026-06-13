package quality

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var qualityInspectionCmd = &cobra.Command{
	Use:   "inspection",
	Short: "质检单管理",
}

func init() {
	qualityCmd.AddCommand(qualityInspectionCmd)
	qualityInspectionCmd.AddCommand(
		newQualityInspectionListCmd(),
		newQualityInspectionGetCmd(),
		newQualityInspectionCreateCmd(),
		newQualityInspectionUpdateCmd(),
		newQualityInspectionDeleteCmd(),
		newQualityInspectionUpdateStatusCmd(),
		newQualityInspectionSummaryCmd(),
		newQualityInspectionRelateListCmd(),
		newQualityInspectionOrderWaybillInspectionCmd(),
		newQualityInspectionExportExcelCmd(),
		newQualityInspectionRefreshAllSummaryCmd(),
		newQualityInspectionImportCmd(),
	)
}

// ---- 分页查询 ----

func newQualityInspectionListCmd() *cobra.Command {
	var pageNo, pageSize int
	var inspectionNo, status, warehouseId, productId, businessId, businessType string

	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询质检单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlags(cmd, params, "inspection-no", "status", "warehouse-id", "product-id", "business-id", "business-type")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/quality-inspection/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询质检单失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().StringVar(&inspectionNo, "inspection-no", "", "质检单号")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&warehouseId, "warehouse-id", "", "仓库 ID")
	c.Flags().StringVar(&productId, "product-id", "", "产品 ID")
	c.Flags().StringVar(&businessId, "business-id", "", "业务 ID")
	c.Flags().StringVar(&businessType, "business-type", "", "业务类型")
	return c
}

// ---- 获取详情 ----

func newQualityInspectionGetCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "get",
		Short: "获取质检单详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/quality-inspection/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取质检单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "质检单 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 创建 ----

func newQualityInspectionCreateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "create",
		Short: "创建质检单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/quality-inspection/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建质检单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 更新 ----

func newQualityInspectionUpdateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update",
		Short: "更新质检单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/quality-inspection/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新质检单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 删除 ----

func newQualityInspectionDeleteCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "delete",
		Short: "删除质检单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/quality-inspection/delete", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("删除质检单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "质检单 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 更新状态 ----

func newQualityInspectionUpdateStatusCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update-status",
		Short: "更新质检单状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/quality-inspection/update-status", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新质检单状态失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据（含 id 和 status）")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 汇总数据 ----

func newQualityInspectionSummaryCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "summary",
		Short: "获取质检单汇总数据",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/quality-inspection/summary", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取质检单汇总失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

// ---- 关联列表 ----

func newQualityInspectionRelateListCmd() *cobra.Command {
	var businessId, businessType string

	c := &cobra.Command{
		Use:   "relate-list",
		Short: "获取质检单关联列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			if businessId != "" {
				params["businessId"] = businessId
			}
			if businessType != "" {
				params["businessType"] = businessType
			}

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/quality-inspection/relate-list", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取关联列表失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&businessId, "business-id", "", "业务 ID")
	c.Flags().StringVar(&businessType, "business-type", "", "业务类型")
	return c
}

// ---- 订单运单质检 ----

func newQualityInspectionOrderWaybillInspectionCmd() *cobra.Command {
	var businessId, businessType string

	c := &cobra.Command{
		Use:   "order-waybill-inspection",
		Short: "获取订单/运单的质检数据",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			if businessId != "" {
				params["businessId"] = businessId
			}
			if businessType != "" {
				params["businessType"] = businessType
			}

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/quality-inspection/order-waybill-inspection", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取质检数据失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&businessId, "business-id", "", "业务 ID")
	c.Flags().StringVar(&businessType, "business-type", "", "业务类型")
	return c
}

// ---- 导出 Excel ----

func newQualityInspectionExportExcelCmd() *cobra.Command {
	var inspectionNo, status, warehouseId, productId, businessId, businessType string

	c := &cobra.Command{
		Use:   "export",
		Short: "导出质检单 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "inspection-no", "status", "warehouse-id", "product-id", "business-id", "business-type")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/quality-inspection/export-excel", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出质检单失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	c.Flags().StringVar(&inspectionNo, "inspection-no", "", "质检单号")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&warehouseId, "warehouse-id", "", "仓库 ID")
	c.Flags().StringVar(&productId, "product-id", "", "产品 ID")
	c.Flags().StringVar(&businessId, "business-id", "", "业务 ID")
	c.Flags().StringVar(&businessType, "business-type", "", "业务类型")
	return c
}

// ---- 刷新所有汇总 ----

func newQualityInspectionRefreshAllSummaryCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "refresh-all-summary",
		Short: "刷新所有质检单汇总数据",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/quality-inspection/refresh-all-summary", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("刷新汇总数据失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

// ---- 导入 ----

func newQualityInspectionImportCmd() *cobra.Command {
	var filePath string

	c := &cobra.Command{
		Use:   "import",
		Short: "导入质检单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			// 读取文件内容
			data, err := os.ReadFile(filePath)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("读取文件失败: %s", err), "")
			}

			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/quality-inspection/import", data)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导入质检单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&filePath, "file", "", "Excel 文件路径")
	_ = c.MarkFlagRequired("file")
	return c
}
