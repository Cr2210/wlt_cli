package partner

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

// NewInvoiceGroup creates the invoice sub-command group for a partner.
func NewInvoiceGroup() *cobra.Command {
	invoiceCmd := cmdutil.NewCRUDGroup(cmdutil.CRUDConfig{
		Name:         "invoice",
		APIPath:      "/erp/business-partner-invoice",
		Label:        "开票信息",
		SingleDelete: true,
		SkipStatus:   true,
		ListFilters: []cmdutil.FlagSpec{
			{Name: "partner-id", Usage: "业务伙伴 ID"},
			{Name: "name", Usage: "名称"},
		},
	})
	invoiceCmd.AddCommand(newInvoiceListByPartnerCmd())
	invoiceCmd.AddCommand(newInvoiceDeleteListCmd())
	return invoiceCmd
}

func newInvoiceListByPartnerCmd() *cobra.Command {
	var partnerId int64
	c := &cobra.Command{
		Use:   "list-by-partner",
		Short: "按业务伙伴查询开票信息",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/business-partner-invoice/list",
				map[string]any{"partnerId": partnerId})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询开票信息失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&partnerId, "partner-id", 0, "业务伙伴 ID")
	_ = c.MarkFlagRequired("partner-id")
	return c
}

func newInvoiceDeleteListCmd() *cobra.Command {
	var ids string
	c := &cobra.Command{
		Use:   "delete-list",
		Short: "批量删除开票信息",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/business-partner-invoice/delete-list",
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
