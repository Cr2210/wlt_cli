package sale

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

// saleAccountingAPIPath 是销售核算的后端端点根路径。
const saleAccountingAPIPath = "/erp/sale-accounting"

// saleAccountingCmd 是销售核算单子命令组。
var saleAccountingCmd = &cobra.Command{
	Use:   "accounting",
	Short: "销售核算单管理",
}

func init() {
	saleCmd.AddCommand(saleAccountingCmd)
	saleAccountingCmd.AddCommand(
		newSaleAccountingListCmd(),
		newSaleAccountingPageCountCmd(),
		cmdutil.CrudGetCmd(saleAccountingAPIPath, "销售核算单"),
		cmdutil.CrudCreateCmd(saleAccountingAPIPath, "销售核算单"),
		cmdutil.CrudUpdateCmd(saleAccountingAPIPath, "销售核算单"),
		cmdutil.CrudDeleteCmd(saleAccountingAPIPath, "销售核算单", false),
		newSaleAccountingWaybillListCmd(),
		newSaleAccountingWaybillCountCmd(),
		newSaleAccountingCalculateCmd(),
		newSaleAccountingDefaultConfigCmd(),
		newSaleAccountingExportExcelCmd(),
	)
}

// ---- list / page-count / export 共用筛选 ----
//
// 后端 /erp/sale-accounting/page 筛选字段：
//   no / customerName / waybillNo / periodRange[0]/[1]（核算周期，yyyy-MM-dd）

func registerSaleAccountingListFlags(c *cobra.Command) {
	c.Flags().String("no", "", "核算单号")
	c.Flags().String("customer-name", "", "客户名称")
	c.Flags().String("waybill-no", "", "运单号")
	c.Flags().String("period-start", "", "核算周期起（yyyy-MM-dd，如 2026-07-01）")
	c.Flags().String("period-end", "", "核算周期止（yyyy-MM-dd，如 2026-07-31）")
}

func collectSaleAccountingListParams(cmd *cobra.Command, pageNo, pageSize int, withPage bool) map[string]any {
	params := map[string]any{}
	if withPage {
		params["pageNo"] = pageNo
		params["pageSize"] = pageSize
	}
	cmdutil.CollectStringFlags(cmd, params, "no", "customer-name", "waybill-no")
	if v, _ := cmd.Flags().GetString("period-start"); v != "" {
		params["periodRange[0]"] = v
	}
	if v, _ := cmd.Flags().GetString("period-end"); v != "" {
		params["periodRange[1]"] = v
	}
	return params
}

// ---- 分页查询 ----

func newSaleAccountingListCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询销售核算单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := collectSaleAccountingListParams(cmd, pageNo, pageSize, true)
			resp, err := cmdutil.GetClient().Get(context.Background(), saleAccountingAPIPath+"/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询销售核算单失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	registerSaleAccountingListFlags(c)
	return c
}

// ---- 分页统计（Σ总金额、Σ运单数量） ----

func newSaleAccountingPageCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "page-count",
		Short: "按筛选统计销售核算单（Σ总金额、Σ运单数量）",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := collectSaleAccountingListParams(cmd, 0, 0, false)
			resp, err := cmdutil.GetClient().Get(context.Background(), saleAccountingAPIPath+"/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计销售核算单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	registerSaleAccountingListFlags(c)
	return c
}

// ---- 核算候选运单 ----
//
// 后端 /erp/sale-accounting/waybill-page 筛选字段：
//   customerId / excludeIds / waybillNo / carNumber
//   realLoadDate[0]/[1]、realUnloadDate[0]/[1]（实际装卸时间范围）

func registerSaleAccountingWaybillFlags(c *cobra.Command) {
	c.Flags().Int64("customer-id", 0, "客户 ID")
	c.Flags().String("exclude-ids", "", "排除的运单 ID（逗号分隔）")
	c.Flags().String("waybill-no", "", "运单号")
	c.Flags().String("car-number", "", "车牌号")
	cmdutil.AddTimeRangeFlags(c, "load-date", "realLoadDate", "装货日期")
	cmdutil.AddTimeRangeFlags(c, "unload-date", "realUnloadDate", "卸货日期")
}

func collectSaleAccountingWaybillParams(cmd *cobra.Command, pageNo, pageSize int, withPage bool) map[string]any {
	params := map[string]any{}
	if withPage {
		params["pageNo"] = pageNo
		params["pageSize"] = pageSize
	}
	cmdutil.CollectIntFlags(cmd, params, "customer-id")
	cmdutil.CollectStringFlags(cmd, params, "exclude-ids", "waybill-no", "car-number")
	cmdutil.CollectTimeRangeFlags(cmd, params, "load-date", "realLoadDate")
	cmdutil.CollectTimeRangeFlags(cmd, params, "unload-date", "realUnloadDate")
	return params
}

func newSaleAccountingWaybillListCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "waybill-list",
		Short: "分页查询核算候选运单（游离+已绑定，settlementId 非空标记 locked）",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := collectSaleAccountingWaybillParams(cmd, pageNo, pageSize, true)
			resp, err := cmdutil.GetClient().Get(context.Background(), saleAccountingAPIPath+"/waybill-page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询核算候选运单失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	registerSaleAccountingWaybillFlags(c)
	return c
}

func newSaleAccountingWaybillCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "waybill-count",
		Short: "统计核算候选运单总数",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := collectSaleAccountingWaybillParams(cmd, 0, 0, false)
			resp, err := cmdutil.GetClient().Get(context.Background(), saleAccountingAPIPath+"/waybill-page/count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计核算候选运单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	registerSaleAccountingWaybillFlags(c)
	return c
}

// ---- 试算核算价（纯内存计算，不写表） ----

func newSaleAccountingCalculateCmd() *cobra.Command {
	var data string
	c := &cobra.Command{
		Use:   "calculate",
		Short: "试算核算价",
		Long:  "按规则质检类型取运单最新已审核质检单指标试算，纯内存计算不写任何表。data 结构：{rule: <计价规则>, items: [<运单明细>], settleQtyType, priceType, taxRate}。",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), saleAccountingAPIPath+"/calculate", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("试算核算价失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据（rule/items/settleQtyType/priceType/taxRate）")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 默认结算税率 ----

func newSaleAccountingDefaultConfigCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "default-config",
		Short: "获得默认结算税率（infra 参数 erp.sale-accounting.tax-rate，缺省 13.000）",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), saleAccountingAPIPath+"/default-config", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询默认税率失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

// ---- 导出 Excel ----

func newSaleAccountingExportExcelCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "export",
		Short: "导出销售核算单 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := collectSaleAccountingListParams(cmd, 0, 0, false)
			cmdutil.CollectStringFlags(cmd, params, "headers")
			resp, err := cmdutil.GetClient().Get(context.Background(), saleAccountingAPIPath+"/export-excel", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出销售核算单失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	registerSaleAccountingListFlags(c)
	c.Flags().String("headers", "", "自定义导出表头")
	return c
}
