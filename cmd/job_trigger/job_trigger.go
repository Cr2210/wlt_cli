package jobtrigger

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var jobTriggerCmd = &cobra.Command{
	Use:   "job-trigger",
	Short: "定时任务触发",
	Long:  "定时任务触发：执行产品成本、执行应收余额。",
}

func init() {
	jobTriggerCmd.AddCommand(
		newJobTriggerExecuteProductCostRecordCmd(),
		newJobTriggerExecuteReceivableBalanceCmd(),
	)
}

// ---- 执行产品成本 ----

func newJobTriggerExecuteProductCostRecordCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "execute-product-cost",
		Short: "执行产品成本记录",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/job-trigger/execute-product-cost-record", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("执行产品成本记录失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

// ---- 执行应收余额 ----

func newJobTriggerExecuteReceivableBalanceCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "execute-receivable-balance",
		Short: "执行应收余额计算",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/job-trigger/execute-receivable-balance", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("执行应收余额计算失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

// Register adds the job-trigger command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(jobTriggerCmd)
}
