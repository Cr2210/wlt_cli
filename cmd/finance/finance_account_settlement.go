package finance

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

// financeAccountSettlementAPIPath 是账户结算的后端端点根路径。
const financeAccountSettlementAPIPath = "/erp/account-settlement"

// financeAccountSettlementListFilters 是账户结算列表的筛选字段。
var financeAccountSettlementListFilters = []cmdutil.FlagSpec{
	{Name: "account-id", Usage: "账户 ID"},
	{Name: "settlement-no", Usage: "结算单号"},
	{Name: "business-no", Usage: "业务单号"},
	{Name: "business-id", Usage: "业务 ID"},
	{Name: "type", Usage: "类型"},
}

func init() {
	financeCmd.AddCommand(financeAccountSettlementCmd)
	financeAccountSettlementCmd.AddCommand(
		cmdutil.CrudListCmdWithFixed(financeAccountSettlementAPIPath, "账户结算", financeAccountSettlementListFilters, nil),
		newFinanceAccountSettlementGetCmd(),
		cmdutil.CrudCreateCmd(financeAccountSettlementAPIPath, "账户结算"),
		cmdutil.CrudUpdateCmd(financeAccountSettlementAPIPath, "账户结算"),
		newFinanceAccountSettlementDeleteCmd(),
		newFinanceAccountSettlementExportExcelCmd(),
	)
}

var financeAccountSettlementCmd = &cobra.Command{
	Use:   "account-settlement",
	Short: "账户结算管理",
}

// ---- 获取详情 ----

func newFinanceAccountSettlementGetCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "get",
		Short: "获取账户结算详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/account-settlement/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取账户结算失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "结算 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 删除 ----

func newFinanceAccountSettlementDeleteCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "delete",
		Short: "删除账户结算",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/account-settlement/delete", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("删除账户结算失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "结算 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 导出 Excel ----

func newFinanceAccountSettlementExportExcelCmd() *cobra.Command {
	var accountId, settlementNo, businessNo, businessId, type_ string

	c := &cobra.Command{
		Use:   "export",
		Short: "导出账户结算 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "account-id", "settlement-no", "business-no", "business-id", "type")

			return cmdutil.RunExport(cmd, "/erp/account-settlement/export-excel", params)
		},
	}
	c.Flags().StringVar(&accountId, "account-id", "", "账户 ID")
	c.Flags().StringVar(&settlementNo, "settlement-no", "", "结算单号")
	c.Flags().StringVar(&businessNo, "business-no", "", "业务单号")
	c.Flags().StringVar(&businessId, "business-id", "", "业务 ID")
	c.Flags().StringVar(&type_, "type", "", "类型")
	return c
}
