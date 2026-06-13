package contract

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

func init() {
	transportCmd := &cobra.Command{
		Use:   "transport",
		Short: "运输合同管理",
	}
	transportCmd.AddCommand(
		cmdutil.CrudListCmd("/erp/transport-contract", "运输合同", contractFilters),
		cmdutil.CrudCreateCmd("/erp/transport-contract", "运输合同"),
		cmdutil.CrudUpdateCmd("/erp/transport-contract", "运输合同"),
		cmdutil.CrudDeleteCmd("/erp/transport-contract", "运输合同", false),
		cmdutil.CrudGetCmd("/erp/transport-contract", "运输合同"),
		newTransportContractUpdateStatusCmd(),
		cmdutil.CrudPageCountCmd("/erp/transport-contract", "运输合同", contractFilters),
	)
	contractCmd.AddCommand(transportCmd)
}

func newTransportContractUpdateStatusCmd() *cobra.Command {
	var id int64
	var status int
	c := &cobra.Command{
		Use:   "update-status",
		Short: "更新运输合同状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().PutQuery(context.Background(), "/erp/transport-contract/update-status",
				map[string]any{"id": id, "status": status})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新状态失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "合同 ID")
	c.Flags().IntVar(&status, "status", 0, "状态")
	_ = c.MarkFlagRequired("id")
	_ = c.MarkFlagRequired("status")
	return c
}
