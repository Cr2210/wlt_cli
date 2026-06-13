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
	Long:  "首页仪表盘数据：业务概览、库存分析、库存积压、产品排行。",
}

func init() {
	// dashboard1 - 无参数
	homepageCmd.AddCommand(newHomepageNoParamCmd("dashboard1", "首页仪表盘（基础统计）"))

	// dashboard2 - 带时间范围
	homepageCmd.AddCommand(cmdutil.NewStatsGetCmd("dashboard2", "/erp/homepage", "业务概览", nil))

	// dashboard6 / inventory-backlog / product-ranking - 带库存筛选
	homepageCmd.AddCommand(cmdutil.NewStatsGetCmd("dashboard6", "/erp/homepage", "库存分析", inventoryFlags))
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
