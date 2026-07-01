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

// ---------------------------------------------------------------------------
// 后端 /erp/quality-inspection-weight/{summary-order,page-waybill,...} 完整筛选字段 (共 13 个):
//
//   type                      业务类型 (SALE/PURCHASE)
//   orderType                 订单类型 (SALE_OUT/PURCHASE_IN 等)
//   orderNo                   订单号 (如 CGDD20260402000002)
//   enterpriseId              企业 ID
//   productId                 产品 ID
//   orderTime[0]/[1]          订单时间范围 ( --start-time/--end-time 自动折叠 )
//   productName               产品名称 (模糊)
//   metricsName               规格指标
//   sendInspectionNo          发货质检单号
//   sendMetricsName           发货质检指标名
//   receiveInspectionNo       收货质检单号
//   receiveMetricsName        收货质检指标名
//   sendAddress               发货地
//
// 运单端点(page-waybill/summary-waybill/export-excel-waybill)完全复用同一套筛选字段。
// ---------------------------------------------------------------------------

// registerWeightFilterFlags 把质检称重 6 个查询命令共用的筛选 flag 注册到命令上。
// 仅 order-page / order-summary 需要 --type (必填);waybill-* 端点无需 type,调用方不会
// 传该 flag,后端走 waybill 路由。
func registerWeightFilterFlags(c *cobra.Command, withType bool) {
	if withType {
		c.Flags().String("type", "", "业务类型 (SALE/PURCHASE)")
	}
	c.Flags().String("order-type", "", "订单类型 (SALE_OUT/PURCHASE_IN 等)")
	c.Flags().String("order-no", "", "订单号")
	c.Flags().String("enterprise-id", "", "企业 ID")
	c.Flags().String("product-id", "", "产品 ID")
	c.Flags().String("product-name", "", "产品名称 (模糊)")
	c.Flags().String("metrics-name", "", "规格指标")
	c.Flags().String("send-inspection-no", "", "发货质检单号")
	c.Flags().String("send-metrics-name", "", "发货质检指标名")
	c.Flags().String("receive-inspection-no", "", "收货质检单号")
	c.Flags().String("receive-metrics-name", "", "收货质检指标名")
	c.Flags().String("send-address", "", "发货地")
	cmdutil.AddOrderTimeFlags(c)
}

// buildWeightParams 收集 weight 模块 6 个查询命令共用的 query params。
// withPage=true 时追加 pageNo/pageSize (用于 page 类命令)。
func buildWeightParams(cmd *cobra.Command, withPage bool, pageNo, pageSize int, forceType string) map[string]any {
	params := map[string]any{}
	if withPage {
		params["pageNo"] = pageNo
		params["pageSize"] = pageSize
	}
	if forceType != "" {
		params["type"] = forceType
	} else if v, _ := cmd.Flags().GetString("type"); v != "" {
		params["type"] = v
	}
	cmdutil.CollectStringFlags(cmd, params,
		"order-type", "order-no", "enterprise-id", "product-id",
		"product-name", "metrics-name",
		"send-inspection-no", "send-metrics-name",
		"receive-inspection-no", "receive-metrics-name",
		"send-address")
	cmdutil.CollectOrderTime(cmd, params)
	return params
}

// ---- 订单称重分页 ----

func newQualityWeightOrderPageCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "order-page",
		Short: "分页查询订单称重",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			typ, _ := cmd.Flags().GetString("type")
			params := buildWeightParams(cmd, true, pageNo, pageSize, typ)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/quality-inspection-weight/page-order", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询订单称重失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	registerWeightFilterFlags(c, true)
	return c
}

// ---- 订单称重汇总 ----

func newQualityWeightOrderSummaryCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "order-summary",
		Short: "获取订单称重汇总",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			typ, _ := cmd.Flags().GetString("type")
			params := buildWeightParams(cmd, false, 0, 0, typ)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/quality-inspection-weight/summary-order", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取订单称重汇总失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	registerWeightFilterFlags(c, true)
	return c
}

// ---- 导出订单称重 Excel ----

func newQualityWeightOrderExportExcelCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "order-export",
		Short: "导出订单称重 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := buildWeightParams(cmd, false, 0, 0, "")
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/quality-inspection-weight/export-excel-order", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出订单称重失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	registerWeightFilterFlags(c, true)
	return c
}

// ---- 运单称重分页 ----

func newQualityWeightWaybillPageCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "waybill-page",
		Short: "分页查询运单称重",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := buildWeightParams(cmd, true, pageNo, pageSize, "")
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/quality-inspection-weight/page-waybill", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询运单称重失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	registerWeightFilterFlags(c, false)
	return c
}

// ---- 运单称重汇总 ----

func newQualityWeightWaybillSummaryCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "waybill-summary",
		Short: "获取运单称重汇总",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := buildWeightParams(cmd, false, 0, 0, "")
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/quality-inspection-weight/summary-waybill", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取运单称重汇总失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	registerWeightFilterFlags(c, false)
	return c
}

// ---- 导出运单称重 Excel ----

func newQualityWeightWaybillExportExcelCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "waybill-export",
		Short: "导出运单称重 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := buildWeightParams(cmd, false, 0, 0, "")
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/quality-inspection-weight/export-excel-waybill", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出运单称重失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	registerWeightFilterFlags(c, false)
	return c
}
