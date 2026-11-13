package finance

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

// financeSettlementAPIPath 是结算单的后端端点根路径。
const financeSettlementAPIPath = "/erp/finance-settlement"

// financeSettlementListFilters 是结算单列表的筛选字段。
var financeSettlementListFilters = []cmdutil.FlagSpec{
	{Name: "settlement-no", Usage: "结算单号"},
	{Name: "customer-id", Usage: "客户 ID"},
	{Name: "supplier-id", Usage: "供应商 ID"},
	{Name: "status", Usage: "状态"},
	{Name: "account-id", Usage: "账户 ID"},
	{Name: "reviewer-id", Usage: "审核人 ID"},
}

func init() {
	financeCmd.AddCommand(financeSettlementCmd)
	financeSettlementCmd.AddCommand(
		cmdutil.CrudListCmd(financeSettlementAPIPath, "结算单", financeSettlementListFilters),
		cmdutil.CrudGetCmd(financeSettlementAPIPath, "结算单"),
		cmdutil.CrudCreateCmd(financeSettlementAPIPath, "结算单"),
		cmdutil.CrudUpdateCmd(financeSettlementAPIPath, "结算单"),
		cmdutil.CrudDeleteCmd(financeSettlementAPIPath, "结算单", true),
		newFinanceSettlementUpdateStatusCmd(),
		newFinanceSettlementSummaryCmd(),
		newFinanceSettlementExportExcelCmd(),
		newFinanceSettlementFinanceRecordCmd(),
	)
}

var financeSettlementCmd = &cobra.Command{
	Use:   "settlement",
	Short: "结算单管理",
}

// ---- 分页查询 ----

// ---- 更新状态 ----

func newFinanceSettlementUpdateStatusCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update-status",
		Short: "更新结算单状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/finance-settlement/update-status", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新结算单状态失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据（含 id 和 status）")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 汇总数据 ----

func newFinanceSettlementSummaryCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "summary",
		Short: "获取结算单汇总数据",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-settlement/summary", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取结算单汇总失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

// ---- 导出 Excel ----

func newFinanceSettlementExportExcelCmd() *cobra.Command {
	var settlementNo, customerId, supplierId, status, accountId, reviewerId string

	c := &cobra.Command{
		Use:   "export",
		Short: "导出结算单 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "settlement-no", "customer-id", "supplier-id", "status", "account-id", "reviewer-id")

			return cmdutil.RunExport(cmd, "/erp/finance-settlement/export-excel", params)
		},
	}
	c.Flags().StringVar(&settlementNo, "settlement-no", "", "结算单号")
	c.Flags().StringVar(&customerId, "customer-id", "", "客户 ID")
	c.Flags().StringVar(&supplierId, "supplier-id", "", "供应商 ID")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&accountId, "account-id", "", "账户 ID")
	c.Flags().StringVar(&reviewerId, "reviewer-id", "", "审核人 ID")
	return c
}

// ---- 财务记录 ----

func newFinanceSettlementFinanceRecordCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "finance-record",
		Short: "获取结算单的财务记录",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/finance-settlement/finance-record", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取财务记录失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "结算单 ID")
	_ = c.MarkFlagRequired("id")
	return c
}
