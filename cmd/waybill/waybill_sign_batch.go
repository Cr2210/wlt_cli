package waybill

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

// ---------------------------------------------------------------------------
// POST /erp/waybill/sign-batch — 批量运单签收(写入)接口。
// 业务规则:仅当运单处于"可签收"状态时允许签收(状态枚举见 waybill_load.go)。
// 本命令不对每条做 GET 自检(批量场景下 N 次请求代价大),由后端校验;
// 可通过 --data 传入完整 JSON 数组做任意批量。
//
// 请求体:JSON 数组
//   [ { "waybillId": <int64>, "signTime": "YYYY-MM-DD HH:mm:ss" }, ... ]
//
// 响应 CommonResultBoolean { "data": true }
// 顶层命令,与 waybillGetCmd / waybillLoadCmd / waybillUnloadCmd / waybillPageCmd / waybillPageCountCmd 同级。
// ---------------------------------------------------------------------------

var waybillSignBatchCmd = &cobra.Command{
	Use:   "sign-batch",
	Short: "批量运单签收（新版 POST /erp/waybill/sign-batch）",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmdutil.EnsureClient(); err != nil {
			return err
		}
		body, err := buildWaybillSignBatchBody(cmd)
		if err != nil {
			return err
		}
		resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/waybill/sign-batch", body)
		if err != nil {
			return output.NewExitError(5, fmt.Sprintf("批量运单签收失败: %s", err), "")
		}
		return cmdutil.OutputJSON(json.RawMessage(resp.Data))
	},
}

// buildWaybillSignBatchBody 收集批量签收请求体(JSON 数组)。
// 支持两种互斥输入方式:
//   1) --data '[{...},{...}]'    直接透传 JSON 数组(推荐批量场景)
//   2) --waybill-id 1,2,3 --sign-time "..."   逗号分隔多个 ID,自动包成数组(便捷)
func buildWaybillSignBatchBody(cmd *cobra.Command) (any, error) {
	// 方式 1: --data 透传 JSON 数组
	if v, _ := cmd.Flags().GetString("data"); v != "" {
		var arr []map[string]any
		if err := json.Unmarshal([]byte(v), &arr); err != nil {
			return nil, output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err),
				"data 应为 JSON 数组,如 [{\"waybillId\":1001,\"signTime\":\"2026-07-02 20:00:00\"}]")
		}
		if len(arr) == 0 {
			return nil, output.NewExitError(4, "data 数组为空,至少需要一条签收记录", "")
		}
		return arr, nil
	}

	// 方式 2: --waybill-id + --sign-time 自动包成数组
	rawIDs, _ := cmd.Flags().GetString("waybill-id")
	if rawIDs == "" {
		return nil, output.NewExitError(4, "缺少必填参数 waybill-id 或 --data",
			"方式 1: --waybill-id 1001 --sign-time \"2026-07-02 20:00:00\"\n方式 2(批量): --data '[{\"waybillId\":1001,\"signTime\":\"...\"},...]'")
	}
	signTime, _ := cmd.Flags().GetString("sign-time")
	if signTime == "" {
		return nil, output.NewExitError(4, "缺少必填参数 sign-time",
			"--sign-time \"YYYY-MM-DD HH:mm:ss\"")
	}

	// 支持单个或逗号分隔的多个 ID
	parts := strings.Split(rawIDs, ",")
	arr := make([]map[string]any, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			return nil, output.NewExitError(4, fmt.Sprintf("非法的 waybill-id: %q", p), "请使用整数或逗号分隔的整数")
		}
		arr = append(arr, map[string]any{
			"waybillId": id,
			"signTime": signTime,
		})
	}
	if len(arr) == 0 {
		return nil, output.NewExitError(4, "未解析到任何有效的 waybill-id", "")
	}
	return arr, nil
}

func init() {
	waybillCmd.AddCommand(waybillSignBatchCmd)
	waybillSignBatchCmd.Flags().String("waybill-id", "", "运单 ID,多个用逗号分隔 (与 --sign-time 组合;与 --data 互斥)")
	waybillSignBatchCmd.Flags().String("sign-time", "", "签收时间 (必填, YYYY-MM-DD HH:mm:ss)")
	waybillSignBatchCmd.Flags().String("data", "", "JSON 数组透传,如 [{\"waybillId\":1001,\"signTime\":\"...\"}] (与上方 flag 互斥)")
}
