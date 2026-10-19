package screen

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

// 大屏补充查询（/erp/screen/*），均带 TimeTypeVO 基础参数（timeType/entId），非分页。
// flag 名与后端 VO 字段一一对应（send-city-name → sendCityName）。

func newScreenGetCmd(use, short, apiPath string, extra func(c *cobra.Command)) *cobra.Command {
	c := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			// Collect* 对未注册的 flag 名安全跳过，故统一收集全部候选。
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "time-type", "metrics-name", "send-city-name", "receive-city-name")
			cmdutil.CollectIntFlags(cmd, params, "ent-id", "product-id")
			resp, err := cmdutil.GetClient().Get(context.Background(), apiPath, params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("%s失败: %s", short, err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	f := c.Flags()
	f.String("time-type", "", "时间维度")
	f.Int64("ent-id", 0, "企业 ID")
	if extra != nil {
		extra(c)
	}
	return c
}

func init() {
	screenCmd.AddCommand(
		newScreenGetCmd("line", "获取大屏线路", "/erp/screen/line", nil),
		newScreenGetCmd("purchase-enterprise", "获取采购企业排名", "/erp/screen/purchase-enterprise", nil),
		newScreenGetCmd("settle-amount", "获取应收应付金额", "/erp/screen/settle-amount", nil),
		newScreenGetCmd("line-price", "获取线路价格", "/erp/screen/line-price", func(c *cobra.Command) {
			c.Flags().String("send-city-name", "", "发货城市")
			c.Flags().String("receive-city-name", "", "收货城市")
		}),
		newScreenGetCmd("purchase-in-price", "获取采购价趋势", "/erp/screen/purchase-in-price", func(c *cobra.Command) {
			c.Flags().Int64("product-id", 0, "产品 ID")
			c.Flags().String("metrics-name", "", "指标名称")
		}),
		newScreenGetCmd("purchase-in-product", "获取采购价趋势产品明细", "/erp/screen/purchase-in-product", func(c *cobra.Command) {
			c.Flags().Int64("product-id", 0, "产品 ID")
			c.Flags().String("metrics-name", "", "指标名称")
		}),
	)
}
