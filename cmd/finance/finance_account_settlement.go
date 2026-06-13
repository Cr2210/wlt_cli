package finance

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

func init() {
	financeCmd.AddCommand(financeAccountSettlementCmd)
	financeAccountSettlementCmd.AddCommand(
		newFinanceAccountSettlementListCmd(),
		newFinanceAccountSettlementGetCmd(),
		newFinanceAccountSettlementCreateCmd(),
		newFinanceAccountSettlementUpdateCmd(),
		newFinanceAccountSettlementDeleteCmd(),
		newFinanceAccountSettlementExportExcelCmd(),
	)
}

var financeAccountSettlementCmd = &cobra.Command{
	Use:   "account-settlement",
	Short: "账户结算管理",
}

// ---- 分页查询 ----

func newFinanceAccountSettlementListCmd() *cobra.Command {
	var pageNo, pageSize int
	var accountId, settlementNo, businessNo, businessId, type_ string

	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询账户结算",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlags(cmd, params, "account-id", "settlement-no", "business-no", "business-id", "type")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/account-settlement/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询账户结算失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().StringVar(&accountId, "account-id", "", "账户 ID")
	c.Flags().StringVar(&settlementNo, "settlement-no", "", "结算单号")
	c.Flags().StringVar(&businessNo, "business-no", "", "业务单号")
	c.Flags().StringVar(&businessId, "business-id", "", "业务 ID")
	c.Flags().StringVar(&type_, "type", "", "类型")
	return c
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

// ---- 创建 ----

func newFinanceAccountSettlementCreateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "create",
		Short: "创建账户结算",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/account-settlement/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建账户结算失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 更新 ----

func newFinanceAccountSettlementUpdateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update",
		Short: "更新账户结算",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/account-settlement/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新账户结算失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
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

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/account-settlement/export-excel", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出账户结算失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	c.Flags().StringVar(&accountId, "account-id", "", "账户 ID")
	c.Flags().StringVar(&settlementNo, "settlement-no", "", "结算单号")
	c.Flags().StringVar(&businessNo, "business-no", "", "业务单号")
	c.Flags().StringVar(&businessId, "business-id", "", "业务 ID")
	c.Flags().StringVar(&type_, "type", "", "类型")
	return c
}
