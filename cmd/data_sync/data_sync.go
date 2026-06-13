package datasync

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var dataSyncCmd = &cobra.Command{
	Use:   "data-sync",
	Short: "数据同步管理",
}

func init() {
	dataSyncCmd.AddCommand(
		newDataSyncPageCmd(),
		newDataSyncGetCmd(),
		newDataSyncResendCmd(),
	)
}

// ---- 分页查询 ----

func newDataSyncPageCmd() *cobra.Command {
	var pageNo, pageSize int
	var status, type_ string

	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询数据同步消息",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlags(cmd, params, "status", "type")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/data-sync-message/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询数据同步消息失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&type_, "type", "", "类型")
	return c
}

// ---- 获取详情 ----

func newDataSyncGetCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "get",
		Short: "获取数据同步消息详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/data-sync-message/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取数据同步消息失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "消息 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 重新发送 ----

func newDataSyncResendCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "resend",
		Short: "重新发送数据同步消息",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/data-sync-message/resend", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("重新发送数据同步消息失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "消息 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// Register adds the data-sync command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(dataSyncCmd)
}
