package waybill

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

// ---- 运单节点事件 ----

var waybillEventsCmd = &cobra.Command{
	Use:   "events",
	Short: "获取运单节点事件列表（/erp/waybill/get-events）",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmdutil.EnsureClient(); err != nil {
			return err
		}
		id, _ := cmd.Flags().GetInt64("id")
		resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/waybill/get-events", map[string]any{"id": id})
		if err != nil {
			return output.NewExitError(5, fmt.Sprintf("获取运单节点事件失败: %s", err), "")
		}
		return cmdutil.OutputJSON(json.RawMessage(resp.Data))
	},
}

// ---- 运单关联订单信息 ----

var waybillOrderInfoCmd = &cobra.Command{
	Use:   "order-info",
	Short: "获取运单关联的订单信息（/erp/waybill/get-order-info）",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmdutil.EnsureClient(); err != nil {
			return err
		}
		id, _ := cmd.Flags().GetInt64("id")
		resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/waybill/get-order-info", map[string]any{"id": id})
		if err != nil {
			return output.NewExitError(5, fmt.Sprintf("获取运单关联订单失败: %s", err), "")
		}
		return cmdutil.OutputJSON(json.RawMessage(resp.Data))
	},
}

func init() {
	// 与 waybillGetCmd 同级的查询命令，挂在 waybillCmd 下。
	waybillCmd.AddCommand(waybillEventsCmd, waybillOrderInfoCmd)
	for _, c := range []*cobra.Command{waybillEventsCmd, waybillOrderInfoCmd} {
		c.Flags().Int64("id", 0, "运单 ID")
		_ = c.MarkFlagRequired("id")
	}
}
