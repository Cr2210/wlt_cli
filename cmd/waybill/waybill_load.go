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
// 运单状态枚举(后端 String):
//
//   FINISHED       已完成
//   UNDEFINED      未定义
//   UN_LOAD        待发货   ← 唯一可进入装货的状态
//   ON_LOAD        运输中
//   DELIVERED      已送达
//   SIGNED_FOR     已签收
//
// POST /erp/waybill/load 装货(写入)接口要求运单必须处于 UN_LOAD(待发货)状态。
// 调用方通过 GET /erp/waybill/get 先自检,非 UN_LOAD 直接阻断。
// ---------------------------------------------------------------------------

var waybillLoadCmd = &cobra.Command{
	Use:   "load",
	Short: "运单装货（新版 POST /erp/waybill/load,要求运单状态为 UN_LOAD 待发货）",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmdutil.EnsureClient(); err != nil {
			return err
		}
		force, _ := cmd.Flags().GetBool("force")

		// 必填字段先校验(クライ端校验,避免无谓的网络调用);
		// 之后才进入状态预检(会发 GET 详情网络请求)。
		if err := validateLoadRequired(cmd); err != nil {
			return err
		}

		// 状态前置校验 — 只有 UN_LOAD(待发货)才可装货。
		if !force {
			warn, err := preCheckLoadStatus(cmd)
			if err != nil {
				return err
			}
			if warn != "" {
				fmt.Fprintln(cmd.OutOrStdout(), warn)
			}
		}

		body, err := buildWaybillLoadBody(cmd)
		if err != nil {
			return err
		}
		resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/waybill/load", body)
		if err != nil {
			return output.NewExitError(5, fmt.Sprintf("运单装货失败: %s", err), "")
		}
		return cmdutil.OutputJSON(json.RawMessage(resp.Data))
	},
}

// validateLoadRequired 校验必填参数(waybill-id / load-time)。
func validateLoadRequired(cmd *cobra.Command) error {
	waybillID, _ := cmd.Flags().GetInt64("waybill-id")
	if waybillID == 0 {
		return output.NewExitError(4, "缺少必填参数 waybill-id", "请通过 --waybill-id 指定运单 ID")
	}
	if v, _ := cmd.Flags().GetString("load-time"); v == "" {
		return output.NewExitError(4, "缺少必填参数 load-time", "请通过 --load-time 指定装货时间 (YYYY-MM-DD HH:mm:ss)")
	}
	return nil
}

// preCheckLoadStatus 调用 GET /erp/waybill/get 自检运单状态。
// 返回值:
//   warn == ""  — 状态 OK(UN_LOAD / 其它可放行状态),提示信息包含放行说明;
//   warn != ""  — 状态不符时的提示(此时 err != nil);
//   err != nil  — 阻断(非 UN_LOAD / 查询失败)。
func preCheckLoadStatus(cmd *cobra.Command) ( string, error) {
	waybillID, _ := cmd.Flags().GetInt64("waybill-id")
	if waybillID == 0 {
		return "", nil // --data 透传模式时不预检
	}

	resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/waybill/get", map[string]any{"id": waybillID})
	if err != nil {
		return "", output.NewExitError(5, fmt.Sprintf("装货前状态查询失败(ID=%d): %s", waybillID, err),
			"可用 --force 跳过状态自检")
	}
	// resp.data 已经被 client 层 unwrap 了"CommonResult",直接读 data.<字段>
	status := gjson.Get(string(resp.Data), "status").String()
	switch status {
	case "UN_LOAD":
		// 合法:放行
		return fmt.Sprintf("状态自检通过(ID=%d, status=UN_LOAD 待发货),继续装货", waybillID), nil
	case "":
		return "", output.NewExitError(5, fmt.Sprintf("无法获取运单状态(ID=%d):status 字段缺失", waybillID),
			"请确认运单 ID 是否存在,或用 --force 跳过自检")
	default:
		return "", output.NewExitError(4, fmt.Sprintf("运单(ID=%d)当前状态为 %q,不允许装货。只有 UN_LOAD(待发货)才可装货", waybillID, status),
			"可用 --force 跳过状态自检(慎用)")
	}
}

// buildWaybillLoadBody 收集装货请求体。waybillId/loadTime 必填,loadWeight/loadedImgUrl 选填。
func buildWaybillLoadBody(cmd *cobra.Command) (map[string]any, error) {
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
	loadTime, _ := cmd.Flags().GetString("load-time")
	if loadTime == "" {
		return nil, output.NewExitError(4, "缺少必填参数 load-time", "请通过 --load-time 指定装货时间 (YYYY-MM-DD HH:mm:ss)")
	}

	body := map[string]any{
		"waybillId": waybillID,
		"loadTime":  loadTime,
	}
	if v, _ := cmd.Flags().GetFloat64("load-weight"); v != 0 {
		body["loadWeight"] = v
	}
	if v, _ := cmd.Flags().GetString("loaded-img-url"); v != "" {
		body["loadedImgUrl"] = v
	}
	return body, nil
}

func init() {
	waybillCmd.AddCommand(waybillLoadCmd)
	waybillLoadCmd.Flags().Int64("waybill-id", 0, "运单主键 ID (必填)")
	waybillLoadCmd.Flags().String("load-time", "", "装货时间 (必填, YYYY-MM-DD HH:mm:ss)")
	waybillLoadCmd.Flags().Float64("load-weight", 0, "装货重量,单位:吨 (选填)")
	waybillLoadCmd.Flags().String("loaded-img-url", "", "装货磅单图片 URL/OSS 路径 (选填)")
	waybillLoadCmd.Flags().String("data", "", "JSON 透传 (与上方 flag 互斥,用于脚本批量)")
	waybillLoadCmd.Flags().Bool("force", false, "跳过 UN_LOAD 状态自检(慎用)")
}
