package position

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var positionCmd = &cobra.Command{
	Use:   "position",
	Short: "位置服务",
	Long:  "位置服务模块：运单车辆轨迹查询（最新/历史）。",
}

// Register adds the position command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(positionCmd)
}

// ---- 最新轨迹 ----

var latestFollowCmd = &cobra.Command{
	Use:   "latest-follow",
	Short: "获取运单车辆最新轨迹",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmdutil.EnsureClient(); err != nil {
			return err
		}
		id, _ := cmd.Flags().GetInt64("id")
		resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/position/get-latest-waybill-follow", map[string]any{"id": id})
		if err != nil {
			return output.NewExitError(5, fmt.Sprintf("获取最新轨迹失败: %s", err), "")
		}
		return cmdutil.OutputJSON(json.RawMessage(resp.Data))
	},
}

// ---- 历史轨迹 ----

var historyFollowCmd = &cobra.Command{
	Use:   "history-follow",
	Short: "获取运单车辆历史轨迹",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmdutil.EnsureClient(); err != nil {
			return err
		}
		id, _ := cmd.Flags().GetInt64("id")
		resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/position/get-history-waybill-follow", map[string]any{"id": id})
		if err != nil {
			return output.NewExitError(5, fmt.Sprintf("获取历史轨迹失败: %s", err), "")
		}
		return cmdutil.OutputJSON(json.RawMessage(resp.Data))
	},
}

func init() {
	positionCmd.AddCommand(latestFollowCmd, historyFollowCmd)
	for _, c := range []*cobra.Command{latestFollowCmd, historyFollowCmd} {
		c.Flags().Int64("id", 0, "运单 ID")
		_ = c.MarkFlagRequired("id")
	}
}
