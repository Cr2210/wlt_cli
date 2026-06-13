package contract

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var provisionFilters = []cmdutil.FlagSpec{
	{Name: "name", Usage: "合同名称"},
	{Name: "status", Usage: "状态"},
	{Name: "no", Usage: "合同编号"},
	{Name: "partner-id", Usage: "业务伙伴 ID"},
}

func init() {
	provisionCmd := &cobra.Command{
		Use:   "provision",
		Short: "供货合同管理",
	}
	provisionCmd.AddCommand(
		cmdutil.CrudListCmd("/erp/provision-contract", "供货合同", provisionFilters),
		cmdutil.CrudCreateCmd("/erp/provision-contract", "供货合同"),
		cmdutil.CrudUpdateCmd("/erp/provision-contract", "供货合同"),
		cmdutil.CrudDeleteCmd("/erp/provision-contract", "供货合同", true),
		newProvisionGetCmd(),
		newProvisionUpdateStatusCmd(),
		cmdutil.CrudPageCountCmd("/erp/provision-contract", "供货合同", provisionFilters),
		newProvisionDeleteBatchCmd(),
		newProvisionFromLongCmd(),
	)
	contractCmd.AddCommand(provisionCmd)
}

func newProvisionGetCmd() *cobra.Command {
	var id int64
	var no string
	c := &cobra.Command{
		Use:   "get",
		Short: "获取供货合同详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{"id": id}
			if no != "" {
				params["no"] = no
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/provision-contract/get", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取供货合同失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "合同 ID")
	c.Flags().StringVar(&no, "no", "", "合同编号")
	_ = c.MarkFlagRequired("id")
	return c
}

func newProvisionUpdateStatusCmd() *cobra.Command {
	var id int64
	var status int
	c := &cobra.Command{
		Use:   "update-status",
		Short: "更新供货合同状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().PutQuery(context.Background(), "/erp/provision-contract/update-status",
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

func newProvisionDeleteBatchCmd() *cobra.Command {
	var ids string
	c := &cobra.Command{
		Use:   "delete-batch",
		Short: "批量删除供货合同",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/provision-contract/delete-batch",
				map[string]any{"ids": ids})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("批量删除失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&ids, "ids", "", "ID 列表（逗号分隔）")
	_ = c.MarkFlagRequired("ids")
	return c
}

func newProvisionFromLongCmd() *cobra.Command {
	var longContractId int64
	c := &cobra.Command{
		Use:   "from-long",
		Short: "从长协合同生成供货合同",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/provision-contract/from-long",
				map[string]any{"longContractId": longContractId})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&longContractId, "long-contract-id", 0, "长协合同 ID")
	_ = c.MarkFlagRequired("long-contract-id")
	return c
}
