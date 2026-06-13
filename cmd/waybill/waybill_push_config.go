package waybill

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var waybillPushConfigCmd = &cobra.Command{
	Use:   "push-config",
	Short: "运单推送配置管理",
}

func init() {
	waybillCmd.AddCommand(waybillPushConfigCmd)
	waybillPushConfigCmd.AddCommand(
		newWaybillPushConfigGetCmd(),
		newWaybillPushConfigUpdateCmd(),
		newWaybillPushConfigGenerateSecretKeyCmd(),
	)
}

// ---- 获取配置 ----

func newWaybillPushConfigGetCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "get",
		Short: "获取运单推送配置",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/waybill-push-config/get", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取推送配置失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

// ---- 更新配置 ----

func newWaybillPushConfigUpdateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update",
		Short: "更新运单推送配置",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/waybill-push-config/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新推送配置失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 生成密钥 ----

func newWaybillPushConfigGenerateSecretKeyCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "generate-secret-key",
		Short: "生成推送密钥",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/waybill-push-config/generate-secret-key", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("生成密钥失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}
