package stats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var dashboard6Flags = []cmdutil.FlagSpec{
	{Name: "product-id", Usage: "产品 ID"},
	{Name: "warehouse-id", Usage: "仓库 ID"},
}

func init() {
	// 库存分析 → /erp/homepage/dashboard6
	statsCmd.AddCommand(newStatsStockCmd())
}

// newStatsStockCmd is the 库存分析 dashboard, exposed as `wlt stats stock` and
// hitting /erp/homepage/dashboard6 (there is no /erp/homepage/stock endpoint).
func newStatsStockCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "stock",
		Short: "库存分析",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectTimeRangeFlags(cmd, params)
			for _, f := range dashboard6Flags {
				cmdutil.CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/homepage/dashboard6", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取库存分析失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	cmdutil.AddStatsFlags(c)
	for _, f := range dashboard6Flags {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}
