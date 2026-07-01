package homepage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var inventoryFlags = []cmdutil.FlagSpec{
	{Name: "product-id", Usage: "产品 ID"},
	{Name: "warehouse-id", Usage: "仓库 ID"},
}

var homepageCmd = &cobra.Command{
	Use:   "homepage",
	Short: "首页数据总览",
	Long:  "首页仪表盘数据：基础统计、库存积压、产品排行。",
}

func init() {
	// dashboard1 - 无参数
	homepageCmd.AddCommand(newHomepageNoParamCmd("dashboard1", "首页仪表盘（基础统计）"))

	// dashboard2 业务概览已迁至 wlt stats overview(2026-07 复核决定:避免重复注册、stats 业务语义更贴)
	// dashboard6 库存分析已迁至 wlt stats stock(2026-07 复核决定:同上,避免 stats/homepage 重复)
	// inventory-backlog / product-ranking - 带库存筛选
	homepageCmd.AddCommand(cmdutil.NewStatsGetCmd("inventory-backlog", "/erp/homepage", "库存积压", inventoryFlags))
	homepageCmd.AddCommand(cmdutil.NewStatsGetCmd("product-ranking", "/erp/homepage", "产品排行", inventoryFlags))
}

func newHomepageNoParamCmd(name, label string) *cobra.Command {
	c := &cobra.Command{
		Use:   name,
		Short: label,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/homepage/"+name, nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("%s失败: %s", label, err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

// Register adds the homepage command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(homepageCmd)
}
