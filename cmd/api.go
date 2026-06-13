package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var (
	apiData   string
	apiParams string
	apiDryRun bool
)

var apiCmd = &cobra.Command{
	Use:   "api <method> <path>",
	Short: "通用 API 调用",
	Long: `直接调用后端 API 端点，输出原始 JSON 响应。
	支持方法: GET, POST, PUT, DELETE。

	示例:
	  wlt api GET /erp/warehouse/page
	  wlt api POST /erp/warehouse/create --data '{"name":"新仓库"}'
	  wlt api DELETE /erp/warehouse/delete --params '{"ids":"1,2"}'`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		method := args[0]
		path := args[1]

		// Validate method
		switch method {
		case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete:
			// OK
		default:
			return output.NewExitError(4, fmt.Sprintf("不支持的 HTTP 方法: %s", method), "支持: GET, POST, PUT, DELETE")
		}

		if err := cmdutil.EnsureClient(); err != nil {
			return err
		}

		// Parse params
		params, err := parseMapFlags(apiParams)
		if err != nil && apiParams != "" {
			return output.NewExitError(4, fmt.Sprintf("解析 params 失败: %s", err), "params 应为 JSON 对象")
		}

		// Parse body data (accept any valid JSON, not just objects)
		var bodyData any
		if apiData != "" {
			if !json.Valid([]byte(apiData)) {
				return output.NewExitError(4, "data 不是有效的 JSON", "data 应为合法 JSON")
			}
			bodyData = json.RawMessage(apiData)
		}

		// Dry run
		if apiDryRun {
			m, u, h, b, err := cmdutil.GetClient().BuildDryRun(method, path, params, bodyData)
			if err != nil {
				return err
			}
			dryRunOutput := map[string]any{
				"method":  m,
				"url":     u,
				"headers": h,
			}
			if b != nil {
				dryRunOutput["body"] = string(b)
			}
			return cmdutil.OutputJSON(dryRunOutput)
		}

		// Execute
		raw, err := cmdutil.GetClient().DoRaw(context.Background(), method, path, params, bodyData)
		if err != nil {
			return output.NewExitError(5, fmt.Sprintf("API 调用失败: %s", err), "使用 --dry-run 调试请求")
		}
		cmdutil.OutputRaw(raw)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
	apiCmd.Flags().StringVar(&apiData, "data", "", "请求体 JSON（POST/PUT）")
	apiCmd.Flags().StringVar(&apiParams, "params", "", "查询参数 JSON（GET/DELETE）")
	apiCmd.Flags().BoolVar(&apiDryRun, "dry-run", false, "只打印请求信息，不发送")
}

// parseMapFlags parses a JSON string into map[string]any.
func parseMapFlags(s string) (map[string]any, error) {
	if s == "" {
		return nil, nil
	}
	return cmdutil.ParseJSONData(s)
}
