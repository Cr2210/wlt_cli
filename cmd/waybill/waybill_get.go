package waybill

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var waybillGetCmd = &cobra.Command{
	Use:   "get",
	Short: "获取运单详情（新版 /erp/waybill/get）",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmdutil.EnsureClient(); err != nil {
			return err
		}
		id, _ := cmd.Flags().GetInt64("id")
		resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/waybill/get", map[string]any{"id": id})
		if err != nil {
			return output.NewExitError(5, fmt.Sprintf("获取运单详情失败: %s", err), "")
		}
		return cmdutil.OutputJSON(json.RawMessage(resp.Data))
	},
}

func init() {
	// 与 waybillPageCmd / waybillPushConfigCmd 同级，挂在 waybillCmd 下。
	waybillCmd.AddCommand(waybillGetCmd)
	waybillGetCmd.Flags().Int64("id", 0, "运单 ID")
	_ = waybillGetCmd.MarkFlagRequired("id")
}
