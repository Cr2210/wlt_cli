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
	serviceCmd := &cobra.Command{
		Use:   "service",
		Short: "业务合同管理",
	}
	serviceCmd.AddCommand(
		cmdutil.CrudListCmd("/erp/service-contract", "业务合同", contractFilters),
		cmdutil.CrudCreateCmd("/erp/service-contract", "业务合同"),
		cmdutil.CrudUpdateCmd("/erp/service-contract", "业务合同"),
		cmdutil.CrudDeleteCmd("/erp/service-contract", "业务合同", false),
		newServiceContractGetCmd(),
		newServiceContractUpdateStatusCmd(),
		cmdutil.CrudPageCountCmd("/erp/service-contract", "业务合同", contractFilters),
	)
	contractCmd.AddCommand(serviceCmd)
}

func newServiceContractGetCmd() *cobra.Command {
	var id int64
	var code string
	c := &cobra.Command{
		Use:   "get",
		Short: "获取业务合同详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{"id": id}
			if code != "" {
				params["code"] = code
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/service-contract/get", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取业务合同失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "合同 ID")
	c.Flags().StringVar(&code, "code", "", "合同编号")
	_ = c.MarkFlagRequired("id")
	return c
}

func newServiceContractUpdateStatusCmd() *cobra.Command {
	var id int64
	var status int
	c := &cobra.Command{
		Use:   "update-status",
		Short: "更新业务合同状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().PutQuery(context.Background(), "/erp/service-contract/update-status",
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
