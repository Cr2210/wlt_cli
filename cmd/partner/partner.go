package partner

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

// AddPartnerCommands adds standard partner CRUD commands to a parent command.
func AddPartnerCommands(parent *cobra.Command, partnerType, label string) {
	parent.AddCommand(
		NewPartnerListCmd(partnerType, label),
		cmdutil.CrudGetCmd("/erp/business-partner", label),
		cmdutil.CrudCreateCmd("/erp/business-partner", label),
		cmdutil.CrudUpdateCmd("/erp/business-partner", label),
		cmdutil.CrudDeleteCmd("/erp/business-partner", label, true),
		cmdutil.CrudUpdateStatusCmd("/erp/business-partner", label),
		NewPartnerSimpleListCmd(partnerType, label),
		NewPartnerPageCountCmd(partnerType, label),
		NewPartnerUpdateAuditStatusCmd(label),
		NewPartnerDeleteListCmd(label),
	)
}

func NewPartnerListCmd(partnerType, label string) *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "list",
		Short: fmt.Sprintf("分页查询%s", label),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
				"type":     partnerType,
			}
			cmdutil.CollectStringFlag(cmd, params, "name")
			cmdutil.CollectStringFlag(cmd, params, "status")
			cmdutil.CollectStringFlag(cmd, params, "category-id")
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/business-partner/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询%s失败: %s", label, err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().String("name", "", "名称")
	c.Flags().String("status", "", "状态")
	c.Flags().String("category-id", "", "分类 ID")
	return c
}

func NewPartnerSimpleListCmd(partnerType, label string) *cobra.Command {
	c := &cobra.Command{
		Use:   "simple-list",
		Short: fmt.Sprintf("获取%s精简列表", label),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/business-partner/simple-list",
				map[string]any{"type": partnerType})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取%s列表失败: %s", label, err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

func NewPartnerPageCountCmd(partnerType, label string) *cobra.Command {
	c := &cobra.Command{
		Use:   "page-count",
		Short: fmt.Sprintf("统计%s数量", label),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{"type": partnerType}
			cmdutil.CollectStringFlag(cmd, params, "name")
			cmdutil.CollectStringFlag(cmd, params, "status")
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/business-partner/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计%s失败: %s", label, err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().String("name", "", "名称")
	c.Flags().String("status", "", "状态")
	return c
}

func NewPartnerUpdateAuditStatusCmd(label string) *cobra.Command {
	var id int64
	var status int
	c := &cobra.Command{
		Use:   "update-audit-status",
		Short: fmt.Sprintf("更新%s审核状态", label),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().PutQuery(context.Background(), "/erp/business-partner/update-audit-status",
				map[string]any{"id": id, "status": status})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新审核状态失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, fmt.Sprintf("%s ID", label))
	c.Flags().IntVar(&status, "status", 0, "审核状态")
	_ = c.MarkFlagRequired("id")
	_ = c.MarkFlagRequired("status")
	return c
}

func NewPartnerDeleteListCmd(label string) *cobra.Command {
	var ids string
	c := &cobra.Command{
		Use:   "delete-list",
		Short: fmt.Sprintf("批量删除%s", label),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/business-partner/delete-list",
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
