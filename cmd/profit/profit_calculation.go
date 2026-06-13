package profit

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var profitCalculationCmd = &cobra.Command{
	Use:   "profit-calculation",
	Short: "利润计算管理",
}

func init() {
	profitCalculationCmd.AddCommand(newProfitCalculationBatchRecalculateAllCmd())
}

// ---- 批量重新计算所有 ----

func newProfitCalculationBatchRecalculateAllCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "batch-recalculate-all",
		Short: "批量重新计算所有利润",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/profit-calculation/batch-recalculate-all", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("批量重新计算利润失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}
