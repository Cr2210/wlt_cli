package stats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "数据总览",
	Long:  "数据总览模块：经营总览、采购分析、销售分析、财务分析、生产分析、库存分析。",
}

func init() {
	// 经营总览
	statsCmd.AddCommand(newStatsOverviewCmd())
	// 库存分析
	statsCmd.AddCommand(cmdutil.NewStatsGetCmd("stock", "/erp/homepage", "库存分析", dashboard6Flags))
}

var dashboard6Flags = []cmdutil.FlagSpec{
	{Name: "product-id", Usage: "产品 ID"},
	{Name: "warehouse-id", Usage: "仓库 ID"},
}

func newStatsOverviewCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "overview",
		Short: "经营总览",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectTimeRangeFlags(cmd, params)
			cmdutil.CollectStringFlag(cmd, params, "product-id")
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/homepage/dashboard2", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取经营总览失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	cmdutil.AddStatsFlags(c)
	c.Flags().String("product-id", "", "产品 ID")
	return c
}

// Register adds the stats command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(statsCmd)
}
