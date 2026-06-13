package operatelog

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var operateLogCmd = &cobra.Command{
	Use:   "operate-log",
	Short: "操作日志查询",
}

func init() {
	operateLogCmd.AddCommand(newOperateLogPageCmd())
}

// ---- 分页查询 ----

func newOperateLogPageCmd() *cobra.Command {
	var pageNo, pageSize int
	var module, type_, userName string

	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询操作日志",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlags(cmd, params, "module", "type", "user-name")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/operate-log/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询操作日志失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().StringVar(&module, "module", "", "模块名")
	c.Flags().StringVar(&type_, "type", "", "操作类型")
	c.Flags().StringVar(&userName, "user-name", "", "用户名")
	return c
}

// Register adds the operate-log command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(operateLogCmd)
}
