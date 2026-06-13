package partner

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var creditFilters = []cmdutil.FlagSpec{
	{Name: "customer-id", Usage: "客户 ID"},
	{Name: "status", Usage: "状态"},
}

// NewCreditGroup creates the credit sub-command group (customer only).
func NewCreditGroup() *cobra.Command {
	creditCmd := cmdutil.NewCRUDGroup(cmdutil.CRUDConfig{
		Name:         "credit",
		APIPath:      "/erp/customer-credit",
		Label:        "客户授信",
		SingleDelete: true,
		SkipStatus:   true,
		ListFilters:  creditFilters,
	})
	creditCmd.AddCommand(cmdutil.CrudPageCountCmd("/erp/customer-credit", "客户授信", creditFilters))
	creditCmd.AddCommand(newCreditCancelCmd())
	creditCmd.AddCommand(newCreditValidCmd())
	creditCmd.AddCommand(newCreditDeleteBatchCmd())
	return creditCmd
}

func newCreditCancelCmd() *cobra.Command {
	var id int64
	c := &cobra.Command{
		Use:   "cancel",
		Short: "作废客户授信",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().PutQuery(context.Background(), "/erp/customer-credit/cancel",
				map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("作废授信失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "授信 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

func newCreditValidCmd() *cobra.Command {
	var customerId int64
	c := &cobra.Command{
		Use:   "valid-credit",
		Short: "获取客户有效授信",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/customer-credit/valid-credit",
				map[string]any{"customerId": customerId})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取有效授信失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&customerId, "customer-id", 0, "客户 ID")
	_ = c.MarkFlagRequired("customer-id")
	return c
}

func newCreditDeleteBatchCmd() *cobra.Command {
	var ids string
	c := &cobra.Command{
		Use:   "delete-batch",
		Short: "批量删除授信",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/customer-credit/delete-batch",
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
