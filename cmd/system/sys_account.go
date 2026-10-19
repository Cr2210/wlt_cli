package system

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

// ---- ERP 系统账号绑定（/erp/system/account/*）与称重账号绑定（/erp/system/weight-account/*） ----

var sysAccountCmd = &cobra.Command{
	Use:   "account",
	Short: "ERP 关联账号",
	Long:  "ERP 关联账号（维运通）查询：关联账号列表、企业信息、运单/轨迹/在途位置、司机/车辆/承运商详情。",
}

func init() {
	SystemCmd.AddCommand(sysAccountCmd, sysWeightAccountGetCmd)
	sysAccountCmd.AddCommand(
		newSysAccountListCmd(),
		newAccountGetCmd("enterprise-info", "获取维运通企业信息", "/erp/system/account/getEnterpriseInfo", func(c *cobra.Command) {
			c.Flags().String("code", "", "企业编码")
			c.Flags().String("mobile", "", "手机号")
		}),
		newAccountGetCmd("waybill-source-info", "获取运单信息", "/erp/system/account/getWaybillSourceInfo", func(c *cobra.Command) {
			c.Flags().String("order-no", "", "订单号")
		}),
		newAccountGetCmd("waybill-follow", "获取运单轨迹", "/erp/system/account/getWaybillFollow", waybillCarIDFlag),
		newAccountGetCmd("latest-waybill-follow", "获取运单最新轨迹", "/erp/system/account/getLatestWaybillFollow", waybillCarIDFlag),
		newAccountGetCmd("waybill", "获取在途位置信息", "/erp/system/account/getWaybill", waybillCarIDFlag),
		newAccountGetCmd("driver-detail", "获取司机详情", "/erp/system/account/getDriverDetail", func(c *cobra.Command) {
			c.Flags().Int64("driver-id", 0, "司机 ID")
		}),
		newAccountGetCmd("car-detail", "获取车辆详情", "/erp/system/account/getCarDetail", func(c *cobra.Command) {
			c.Flags().Int64("car-id", 0, "车辆 ID")
			c.Flags().String("car-category", "", "车辆类别")
		}),
		newAccountGetCmd("carrier-detail", "获取承运商详情", "/erp/system/account/getCarrierDetail", func(c *cobra.Command) {
			c.Flags().Int64("carrier-enterprise-id", 0, "承运商企业 ID")
		}),
	)
}

func waybillCarIDFlag(c *cobra.Command) {
	c.Flags().Int64("waybill-car-id", 0, "运单车辆 ID")
}

// newAccountGetCmd 构造单端点详情查询命令；RunE 统一收集全部候选 flag
// （Collect*Flags 对未注册的 flag 名安全跳过）。
func newAccountGetCmd(use, short, apiPath string, addFlags func(*cobra.Command)) *cobra.Command {
	c := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "code", "mobile", "order-no", "car-category")
			cmdutil.CollectIntFlags(cmd, params, "waybill-car-id", "driver-id", "car-id", "carrier-enterprise-id")
			resp, err := cmdutil.GetClient().Get(context.Background(), apiPath, params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("%s失败: %s", short, err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	addFlags(c)
	return c
}

// ---- 关联账号列表（无参，返回裸数组） ----

func newSysAccountListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "关联账号列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/system/account/list", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询关联账号列表失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
}

// ---- 称重账号绑定（/erp/system/weight-account/get，无参） ----

var sysWeightAccountGetCmd = &cobra.Command{
	Use:   "weight-account",
	Short: "称重账号绑定信息",
	Long:  "称重账号绑定信息查询（/erp/system/weight-account/get）。",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmdutil.EnsureClient(); err != nil {
			return err
		}
		resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/system/weight-account/get", nil)
		if err != nil {
			return output.NewExitError(5, fmt.Sprintf("查询称重账号绑定失败: %s", err), "")
		}
		return cmdutil.OutputJSON(json.RawMessage(resp.Data))
	},
}
