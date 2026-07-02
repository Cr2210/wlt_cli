package waybill

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

// ---------------------------------------------------------------------------
// POST /erp/waybill/unload — 运单卸货(写入)接口。
// 业务规则:只有 ON_LOAD(运输中)才可卸货。运单状态枚举见 waybill_load.go 顶部注释。
//
//   waybillId*      运单主键 ID (必填)
//   unloadTime*     卸货时间 (必填, YYYY-MM-DD HH:mm:ss)
//   unloadWeight*   卸货重量/吨 (必填)
//   unloadedImgUrl  卸货磅单图片 URL/OSS 路径 (选填)
//
// 响应 CommonResultBoolean { "data": true }
// 顶层命令,与 waybillGetCmd / waybillLoadCmd / waybillPageCmd / waybillPageCountCmd 同级。
// ---------------------------------------------------------------------------

var waybillUnloadCmd = &cobra.Command{
	Use:   "unload",
	Short: "运单卸货（新版 POST /erp/waybill/unload,要求运单状态为 ON_LOAD 运输中）",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmdutil.EnsureClient(); err != nil {
			return err
		}
		force, _ := cmd.Flags().GetBool("force")

		// 必填字段先校验(客户端校验,避免无谓的网络调用);
		// 之后才进入状态预检(会发 GET 详情网络请求)。
		if err := validateUnloadRequired(cmd); err != nil {
			return err
		}

		// 状态前置校验 — 只有 ON_LOAD(运输中)才可卸货。
		if !force {
			warn, err := preCheckUnloadStatus(cmd)
			if err != nil {
				return err
			}
			if warn != "" {
				fmt.Fprintln(cmd.OutOrStdout(), warn)
			}
		}

		body, err := buildWaybillUnloadBody(cmd)
		if err != nil {
			return err
		}
		resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/waybill/unload", body)
		if err != nil {
			return output.NewExitError(5, fmt.Sprintf("运单卸货失败: %s", err), "")
		}
		return cmdutil.OutputJSON(json.RawMessage(resp.Data))
	},
}

// validateUnloadRequired 校验必填参数(waybill-id / unload-time / unload-weight)。
func validateUnloadRequired(cmd *cobra.Command) error {
	waybillID, _ := cmd.Flags().GetInt64("waybill-id")
	if waybillID == 0 {
		return output.NewExitError(4, "缺少必填参数 waybill-id", "请通过 --waybill-id 指定运单 ID")
	}
	if v, _ := cmd.Flags().GetString("unload-time"); v == "" {
		return output.NewExitError(4, "缺少必填参数 unload-time", "请通过 --unload-time 指定卸货时间 (YYYY-MM-DD HH:mm:ss)")
	}
	return nil
}

// preCheckUnloadStatus 调用 GET /erp/waybill/get 自检运单状态。
// 只有 ON_LOAD(运输中)可放行;其它状态阻断。--force / --data 透传模式跳过。
func preCheckUnloadStatus(cmd *cobra.Command) (string, error) {
	waybillID, _ := cmd.Flags().GetInt64("waybill-id")
	if waybillID == 0 {
		return "", nil // --data 透传模式时不预检
	}

	resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/waybill/get", map[string]any{"id": waybillID})
	if err != nil {
		return "", output.NewExitError(5, fmt.Sprintf("卸货前状态查询失败(ID=%d): %s", waybillID, err),
			"可用 --force 跳过状态自检")
	}
	status := gjson.Get(string(resp.Data), "status").String()
	switch status {
	case "ON_LOAD":
		return fmt.Sprintf("状态自检通过(ID=%d, status=ON_LOAD 运输中),继续卸货", waybillID), nil
	case "":
		return "", output.NewExitError(5, fmt.Sprintf("无法获取运单状态(ID=%d):status 字段缺失", waybillID),
			"请确认运单 ID 是否存在,或用 --force 跳过自检")
	default:
		return "", output.NewExitError(4, fmt.Sprintf("运单(ID=%d)当前状态为 %q,不允许卸货。只有 ON_LOAD(运输中)才可卸货", waybillID, status),
			"可用 --force 跳过状态自检(慎用)")
	}
}

// buildWaybillUnloadBody 收集卸货请求体。waybillId/unloadTime/unloadWeight 必填,unloadedImgUrl 选填。
func buildWaybillUnloadBody(cmd *cobra.Command) (map[string]any, error) {
	// 解析 --data 透传（可选,用于脚本批量作业）
	if v, _ := cmd.Flags().GetString("data"); v != "" {
		m, err := cmdutil.ParseJSONData(v)
		if err != nil {
			return nil, output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
		}
		return m, nil
	}

	waybillID, err := cmd.Flags().GetInt64("waybill-id")
	if err != nil || waybillID == 0 {
		return nil, output.NewExitError(4, "缺少必填参数 waybill-id", "请通过 --waybill-id 指定运单 ID")
	}
	unloadTime, _ := cmd.Flags().GetString("unload-time")
	if unloadTime == "" {
		return nil, output.NewExitError(4, "缺少必填参数 unload-time", "请通过 --unload-time 指定卸货时间 (YYYY-MM-DD HH:mm:ss)")
	}
	unloadWeight, _ := cmd.Flags().GetFloat64("unload-weight")

	body := map[string]any{
		"waybillId":   waybillID,
		"unloadTime":  unloadTime,
		"unloadWeight": unloadWeight,
	}
	if v, _ := cmd.Flags().GetString("unloaded-img-url"); v != "" {
		body["unloadedImgUrl"] = v
	}
	return body, nil
}

func init() {
	waybillCmd.AddCommand(waybillUnloadCmd)
	waybillUnloadCmd.Flags().Int64("waybill-id", 0, "运单主键 ID (必填)")
	waybillUnloadCmd.Flags().String("unload-time", "", "卸货时间 (必填, YYYY-MM-DD HH:mm:ss)")
	waybillUnloadCmd.Flags().Float64("unload-weight", 0, "卸货重量,单位:吨 (必填)")
	waybillUnloadCmd.Flags().String("unloaded-img-url", "", "卸货磅单图片 URL/OSS 路径 (选填)")
	waybillUnloadCmd.Flags().String("data", "", "JSON 透传 (与上方 flag 互斥,用于脚本批量)")
	waybillUnloadCmd.Flags().Bool("force", false, "跳过 ON_LOAD 状态自检(慎用)")
}
