package configcmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "设置配置项",
	Long:  "设置指定环境的配置项。key 格式：<profile>.<field>，例如 sit.base_url。",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		parts := strings.SplitN(key, ".", 2)
		if len(parts) != 2 {
			return output.NewExitError(4, "key 格式错误，应为 <profile>.<field>", "例如: sit.base_url https://example.com")
		}
		profileName, field := parts[0], parts[1]

		// Validate field name before attempting update
		switch field {
		case "base_url", "api_prefix", "enterprise_type":
			// valid
		default:
			return output.NewExitError(4, fmt.Sprintf("未知配置项: %s", field), "支持: base_url, api_prefix, enterprise_type")
		}

		if err := cmdutil.CfgMgr.UpdateProfileField(profileName, field, value); err != nil {
			return output.NewExitError(2, fmt.Sprintf("设置失败: %s", err), "运行 wlt config init 创建")
		}

		if !getQuiet(cmd) {
			fmt.Fprintf(cmd.ErrOrStderr(), "已设置 %s = %s\n", key, value)
		}
		return cmdutil.OutputJSON(map[string]any{"key": key, "value": value})
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
}
