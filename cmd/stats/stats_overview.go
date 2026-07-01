package stats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

func init() {
	// 经营总览 → /erp/homepage/dashboard2
	statsCmd.AddCommand(newStatsOverviewCmd())
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
