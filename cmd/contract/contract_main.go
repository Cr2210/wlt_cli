package contract

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var contractFilters = []cmdutil.FlagSpec{
	{Name: "name", Usage: "合同名称"},
	{Name: "status", Usage: "状态"},
	{Name: "code", Usage: "合同编号"},
	{Name: "partner-id", Usage: "业务伙伴 ID"},
}

func init() {
	contractCmd.AddCommand(
		cmdutil.CrudListCmdWithFixed("/erp/contract", "长协合同", contractFilters, map[string]any{"type": "LONG"}),
		cmdutil.CrudCreateCmd("/erp/contract", "长协合同"),
		cmdutil.CrudUpdateCmd("/erp/contract", "长协合同"),
		cmdutil.CrudDeleteCmd("/erp/contract", "长协合同", false),
		newContractGetCmd(),
		newContractUpdateStatusCmd(),
		cmdutil.CrudPageCountCmdWithFixed("/erp/contract", "长协合同", contractFilters, map[string]any{"type": "LONG"}),
	)
}

func newContractGetCmd() *cobra.Command {
	var id int64
	var code string
	c := &cobra.Command{
		Use:   "get",
		Short: "获取长协合同详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{"id": id}
			if code != "" {
				params["code"] = code
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/contract/get", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取长协合同失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "合同 ID")
	c.Flags().StringVar(&code, "code", "", "合同编号")
	_ = c.MarkFlagRequired("id")
	return c
}

func newContractUpdateStatusCmd() *cobra.Command {
	var id int64
	var status int
	c := &cobra.Command{
		Use:   "update-status",
		Short: "更新长协合同状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().PutQuery(context.Background(), "/erp/contract/update-status",
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
