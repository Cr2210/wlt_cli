package waybill

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

// ---------------------------------------------------------------------------
// 新版运单分页端点 (/erp/waybill/*),共 18 个筛选字段:
//
//   waybillNo                 运单号
//   carNumber                 车牌号
//   orderType                 订单类型 (SALE_OUT/PURCHASE_IN 等)
//   addressName               起始地/目的地
//   status                    状态
//   mediumName                货物名称
//   metricsName               规格指标
//   realLoadDate[0]/[1]       实际装车时间范围
//   realUnloadDate[0]/[1]     实际卸车时间范围
//   capacityName              承运商
//   userName                  业务负责人
//   projectName               项目名称
//   createTime[0]/[1]         创建时间范围
//   inputType                 录入方式 (ACQUIRE/MANUAL 等)
//   dataSource                数据来源 (WLCLOUDS/...)
//   outWaybillNo              外部数据编号
//
// list / page-count 与 waybillGetCmd / waybillPushConfigCmd 同级,挂在 waybillCmd 下。
// ---------------------------------------------------------------------------

// collectDateRange 把 --{prefix}-start/--{prefix}-end flag 折叠成 {key[0]}/{key[1]} 数组参数。
// 例如 prefix="real-load-date", key="realLoadDate" 生成 realLoadDate[0]=...&realLoadDate[1]=...
func collectDateRange(cmd *cobra.Command, params map[string]any, flagPrefix, key string) {
	if v, _ := cmd.Flags().GetString(flagPrefix + "-start"); v != "" {
		params[key+"[0]"] = v
	}
	if v, _ := cmd.Flags().GetString(flagPrefix + "-end"); v != "" {
		params[key+"[1]"] = v
	}
}

// registerWaybillFilterFlags 注册新版运单 page 端点 18 个筛选字段。withPage=true 时追加 pageNo/pageSize。
func registerWaybillFilterFlags(c *cobra.Command, withPage bool) {
	c.Flags().String("waybill-no", "", "运单号")
	c.Flags().String("car-number", "", "车牌号")
	c.Flags().String("order-type", "", "订单类型 (SALE_OUT/PURCHASE_IN 等)")
	c.Flags().String("address-name", "", "起始地/目的地")
	c.Flags().String("status", "", "状态")
	c.Flags().String("medium-name", "", "货物名称")
	c.Flags().String("metrics-name", "", "规格指标")
	c.Flags().String("capacity-name", "", "承运商")
	c.Flags().String("user-name", "", "业务负责人")
	c.Flags().String("project-name", "", "项目名称")
	c.Flags().String("input-type", "", "录入方式 (ACQUIRE/MANUAL 等)")
	c.Flags().String("data-source", "", "数据来源 (WLCLOUDS/...)")
	c.Flags().String("out-waybill-no", "", "外部数据编号")
	c.Flags().String("real-load-date-start", "", "实际装车起始时间 (YYYY-MM-DD HH:mm:ss)")
	c.Flags().String("real-load-date-end", "", "实际装车截止时间 (YYYY-MM-DD HH:mm:ss)")
	c.Flags().String("real-unload-date-start", "", "实际卸车起始时间 (YYYY-MM-DD HH:mm:ss)")
	c.Flags().String("real-unload-date-end", "", "实际卸车截止时间 (YYYY-MM-DD HH:mm:ss)")
	c.Flags().String("create-time-start", "", "创建起始时间 (YYYY-MM-DD HH:mm:ss)")
	c.Flags().String("create-time-end", "", "创建截止时间 (YYYY-MM-DD HH:mm:ss)")
	if withPage {
		c.Flags().Int("page-no", 1, "页码")
		c.Flags().Int("page-size", 20, "每页数量")
	}
}

// buildWaybillPageParams 收集新版运单 page 端点筛选参数。
func buildWaybillPageParams(cmd *cobra.Command) map[string]any {
	pageNo, _ := cmd.Flags().GetInt("page-no")
	pageSize, _ := cmd.Flags().GetInt("page-size")

	params := map[string]any{}
	cmdutil.CollectStringFlags(cmd, params,
		"waybill-no", "car-number", "order-type", "address-name",
		"status", "medium-name", "metrics-name", "capacity-name",
		"user-name", "project-name", "input-type", "data-source",
		"out-waybill-no")
	collectDateRange(cmd, params, "real-load-date", "realLoadDate")
	collectDateRange(cmd, params, "real-unload-date", "realUnloadDate")
	collectDateRange(cmd, params, "create-time", "createTime")
	if pageNo > 0 {
		params["pageNo"] = pageNo
	}
	if pageSize > 0 {
		params["pageSize"] = pageSize
	}
	return params
}

// ---- 分页查询 ----

var waybillPageCmd = &cobra.Command{
	Use:   "page",
	Short: "分页查询运单（新版 /erp/waybill/page，含 18 个筛选字段）",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmdutil.EnsureClient(); err != nil {
			return err
		}
		params := buildWaybillPageParams(cmd)
		resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/waybill/page", params)
		if err != nil {
			return output.NewExitError(5, fmt.Sprintf("查询运单失败: %s", err), "")
		}
		pageNo, _ := cmd.Flags().GetInt("page-no")
		pageSize, _ := cmd.Flags().GetInt("page-size")
		return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
	},
}

// ---- 分页计数 ----

var waybillPageCountCmd = &cobra.Command{
	Use:   "page-count",
	Short: "统计运单数量（新版 /erp/waybill/page-count）",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmdutil.EnsureClient(); err != nil {
			return err
		}
		params := buildWaybillPageParams(cmd)
		resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/waybill/page-count", params)
		if err != nil {
			return output.NewExitError(5, fmt.Sprintf("统计运单数量失败: %s", err), "")
		}
		return cmdutil.OutputJSON(json.RawMessage(resp.Data))
	},
}

func init() {
	// page / page-count 与 waybillGetCmd / waybillPushConfigCmd 同级,挂在 waybillCmd 下。
	waybillCmd.AddCommand(waybillPageCmd, waybillPageCountCmd)

	registerWaybillFilterFlags(waybillPageCmd, true)
	registerWaybillFilterFlags(waybillPageCountCmd, false)
}
