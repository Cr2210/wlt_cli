package partner

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var settlementFilters = []cmdutil.FlagSpec{
	{Name: "partner-id", Usage: "业务伙伴 ID"},
	{Name: "name", Usage: "名称"},
	{Name: "status", Usage: "状态"},
}

// NewSettlementGroup creates the settlement sub-command group for a partner.
func NewSettlementGroup() *cobra.Command {
	settlementCmd := cmdutil.NewCRUDGroup(cmdutil.CRUDConfig{
		Name:         "settlement",
		APIPath:      "/erp/business-partner-settlement",
		Label:        "结算信息",
		SingleDelete: true,
		SkipStatus:   true,
		ListFilters:  settlementFilters,
	})
	settlementCmd.AddCommand(cmdutil.CrudPageCountCmd("/erp/business-partner-settlement", "结算信息", settlementFilters))
	settlementCmd.AddCommand(newSettlementListByPartnerCmd())
	settlementCmd.AddCommand(newSettlementDeleteListCmd())
	return settlementCmd
}

func newSettlementListByPartnerCmd() *cobra.Command {
	var partnerId int64
	c := &cobra.Command{
		Use:   "list-by-partner",
		Short: "按业务伙伴查询结算信息",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/business-partner-settlement/list-by-partner",
				map[string]any{"businessPartnerId": partnerId})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询结算信息失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&partnerId, "partner-id", 0, "业务伙伴 ID")
	_ = c.MarkFlagRequired("partner-id")
	return c
}

func newSettlementDeleteListCmd() *cobra.Command {
	var ids string
	c := &cobra.Command{
		Use:   "delete-list",
		Short: "批量删除结算信息",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			idList := parseIDList(ids)
			resp, err := cmdutil.GetClient().DeleteWithBody(context.Background(), "/erp/business-partner-settlement/delete-list", idList)
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

func parseIDList(s string) []int64 {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]int64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if id, err := strconv.ParseInt(p, 10, 64); err == nil && id > 0 {
			result = append(result, id)
		}
	}
	return result
}
