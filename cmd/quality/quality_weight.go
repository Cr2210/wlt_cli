package quality

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var qualityWeightCmd = &cobra.Command{
	Use:   "weight",
	Short: "质检称重管理",
}

func init() {
	qualityCmd.AddCommand(qualityWeightCmd)
	qualityWeightCmd.AddCommand(
		newQualityWeightOrderPageCmd(),
		newQualityWeightOrderSummaryCmd(),
		newQualityWeightOrderExportExcelCmd(),
		newQualityWeightWaybillPageCmd(),
		newQualityWeightWaybillSummaryCmd(),
		newQualityWeightWaybillExportExcelCmd(),
	)
}

// ---- 订单称重分页 ----

func newQualityWeightOrderPageCmd() *cobra.Command {
	var pageNo, pageSize int
	var orderId, productId, orderType string

	c := &cobra.Command{
		Use:   "order-page",
		Short: "分页查询订单称重",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
				"type":     orderType,
			}
			if orderId != "" {
				params["orderId"] = orderId
			}
			if productId != "" {
				params["productId"] = productId
			}

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/quality-inspection-weight/page-order", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询订单称重失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().StringVar(&orderType, "type", "", "订单类型（SALE/PURCHASE）")
	c.Flags().StringVar(&orderId, "order-id", "", "订单 ID")
	c.Flags().StringVar(&productId, "product-id", "", "产品 ID")
	_ = c.MarkFlagRequired("type")
	return c
}

// ---- 订单称重汇总 ----

func newQualityWeightOrderSummaryCmd() *cobra.Command {
	var orderId, productId, orderType string

	c := &cobra.Command{
		Use:   "order-summary",
		Short: "获取订单称重汇总",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{"type": orderType}
			if orderId != "" {
				params["orderId"] = orderId
			}
			if productId != "" {
				params["productId"] = productId
			}

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/quality-inspection-weight/summary-order", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取订单称重汇总失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&orderType, "type", "", "订单类型（SALE/PURCHASE）")
	c.Flags().StringVar(&orderId, "order-id", "", "订单 ID")
	c.Flags().StringVar(&productId, "product-id", "", "产品 ID")
	_ = c.MarkFlagRequired("type")
	return c
}

// ---- 导出订单称重 Excel ----

func newQualityWeightOrderExportExcelCmd() *cobra.Command {
	var orderId, productId string

	c := &cobra.Command{
		Use:   "order-export",
		Short: "导出订单称重 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			if orderId != "" {
				params["orderId"] = orderId
			}
			if productId != "" {
				params["productId"] = productId
			}

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/quality-inspection-weight/export-excel-order", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出订单称重失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	c.Flags().StringVar(&orderId, "order-id", "", "订单 ID")
	c.Flags().StringVar(&productId, "product-id", "", "产品 ID")
	return c
}

// ---- 运单称重分页 ----

func newQualityWeightWaybillPageCmd() *cobra.Command {
	var pageNo, pageSize int
	var waybillId, productId string

	c := &cobra.Command{
		Use:   "waybill-page",
		Short: "分页查询运单称重",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			if waybillId != "" {
				params["waybillId"] = waybillId
			}
			if productId != "" {
				params["productId"] = productId
			}

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/quality-inspection-weight/page-waybill", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询运单称重失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().StringVar(&waybillId, "waybill-id", "", "运单 ID")
	c.Flags().StringVar(&productId, "product-id", "", "产品 ID")
	return c
}

// ---- 运单称重汇总 ----

func newQualityWeightWaybillSummaryCmd() *cobra.Command {
	var waybillId, productId string

	c := &cobra.Command{
		Use:   "waybill-summary",
		Short: "获取运单称重汇总",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			if waybillId != "" {
				params["waybillId"] = waybillId
			}
			if productId != "" {
				params["productId"] = productId
			}

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/quality-inspection-weight/summary-waybill", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取运单称重汇总失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&waybillId, "waybill-id", "", "运单 ID")
	c.Flags().StringVar(&productId, "product-id", "", "产品 ID")
	return c
}

// ---- 导出运单称重 Excel ----

func newQualityWeightWaybillExportExcelCmd() *cobra.Command {
	var waybillId, productId string

	c := &cobra.Command{
		Use:   "waybill-export",
		Short: "导出运单称重 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			if waybillId != "" {
				params["waybillId"] = waybillId
			}
			if productId != "" {
				params["productId"] = productId
			}

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/quality-inspection-weight/export-excel-waybill", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出运单称重失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	c.Flags().StringVar(&waybillId, "waybill-id", "", "运单 ID")
	c.Flags().StringVar(&productId, "product-id", "", "产品 ID")
	return c
}
