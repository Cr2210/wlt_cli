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
	financeCmd.AddCommand(financeAccountCmd)
	financeAccountCmd.AddCommand(
		newFinanceAccountListCmd(),
		newFinanceAccountGetCmd(),
		newFinanceAccountCreateCmd(),
		newFinanceAccountUpdateCmd(),
		newFinanceAccountDeleteCmd(),
		newFinanceAccountUpdateStatusCmd(),
		newFinanceAccountUpdateDefaultStatusCmd(),
		newFinanceAccountSimpleListCmd(),
		newFinanceAccountAdjustCmd(),
		newFinanceAccountTransferCmd(),
		newFinanceAccountSettlementPageCmd(),
		newFinanceAccountExportExcelCmd(),
	)
}

var financeAccountCmd = &cobra.Command{
	Use:   "account",
	Short: "结算账户管理",
	Long:  "结算账户操作：账户CRUD、调整、转账、结算分页。",
}

// ---- 分页查询 ----

func newFinanceAccountListCmd() *cobra.Command {
	var pageNo, pageSize int
	var accountId, accountNo, name, status, remark string

	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询结算账户",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlags(cmd, params, "account-id", "account-no", "name", "status", "remark")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/account/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询结算账户失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().StringVar(&accountId, "account-id", "", "账户 ID")
	c.Flags().StringVar(&accountNo, "account-no", "", "账户编号")
	c.Flags().StringVar(&name, "name", "", "账户名称")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&remark, "remark", "", "备注")
	return c
}

// ---- 获取详情 ----

func newFinanceAccountGetCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "get",
		Short: "获取结算账户详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/account/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取结算账户失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "账户 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 创建 ----

func newFinanceAccountCreateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "create",
		Short: "创建结算账户",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/account/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建结算账户失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 更新 ----

func newFinanceAccountUpdateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update",
		Short: "更新结算账户",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/account/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新结算账户失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 删除 ----

func newFinanceAccountDeleteCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "delete",
		Short: "删除结算账户",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/account/delete", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("删除结算账户失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "账户 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 更新状态 ----

func newFinanceAccountUpdateStatusCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update-status",
		Short: "更新结算账户状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/account/update-status", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新结算账户状态失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据（含 id 和 status）")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 更新默认状态 ----

func newFinanceAccountUpdateDefaultStatusCmd() *cobra.Command {
	var id int64
	var isDefault bool

	c := &cobra.Command{
		Use:   "update-default-status",
		Short: "更新结算账户默认状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body := map[string]any{
				"id":         id,
				"isDefault": isDefault,
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/account/update-default-status", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新默认状态失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "账户 ID")
	c.Flags().BoolVar(&isDefault, "default", false, "是否默认")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 简单列表 ----

func newFinanceAccountSimpleListCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "simple-list",
		Short: "获取结算账户精简列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/account/simple-list", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取结算账户列表失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

// ---- 账户调整 ----

func newFinanceAccountAdjustCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "adjust",
		Short: "结算账户调整",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/account/adjust", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("账户调整失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 账户转账 ----

func newFinanceAccountTransferCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "transfer",
		Short: "结算账户转账",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/account/transfer", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("账户转账失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 结算分页 ----

func newFinanceAccountSettlementPageCmd() *cobra.Command {
	var pageNo, pageSize int
	var accountId string

	c := &cobra.Command{
		Use:   "settlement-page",
		Short: "分页查询账户结算明细",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			if accountId != "" {
				params["accountId"] = accountId
			}

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/account/settlement-page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询账户结算明细失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().StringVar(&accountId, "account-id", "", "账户 ID")
	return c
}

// ---- 导出 Excel ----

func newFinanceAccountExportExcelCmd() *cobra.Command {
	var accountId, accountNo, name, status, remark string

	c := &cobra.Command{
		Use:   "export",
		Short: "导出结算账户 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "account-id", "account-no", "name", "status", "remark")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/account/export-excel", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出结算账户失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	c.Flags().StringVar(&accountId, "account-id", "", "账户 ID")
	c.Flags().StringVar(&accountNo, "account-no", "", "账户编号")
	c.Flags().StringVar(&name, "name", "", "账户名称")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&remark, "remark", "", "备注")
	return c
}
