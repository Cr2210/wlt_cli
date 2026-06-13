package configcmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/config"
	"github.com/weiliantong/cli/internal/output"
)

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "交互式初始化配置",
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			name      string
			baseURL   string
			apiPrefix = "/admin-api"
			tenantID  string
		)

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("配置名称").Value(&name).Placeholder("sit"),
				huh.NewInput().Title("API 地址").Value(&baseURL).Placeholder("https://erpsit.api.w-lian.com"),
				huh.NewInput().Title("API 前缀").Value(&apiPrefix).Placeholder("/admin-api"),
				huh.NewInput().Title("租户 ID").Value(&tenantID).Placeholder("1"),
			),
		)

		if err := form.Run(); err != nil {
			return output.NewExitError(1, fmt.Sprintf("表单输入失败: %s", err), "")
		}

		if err := cmdutil.CfgMgr.SetProfile(name, &config.Profile{
			BaseURL:   baseURL,
			APIPrefix: apiPrefix,
			TenantID:  tenantID,
		}); err != nil {
			return output.NewExitError(2, fmt.Sprintf("保存配置失败: %s", err), "")
		}
		if err := cmdutil.CfgMgr.SetActive(name); err != nil {
			return output.NewExitError(2, fmt.Sprintf("设置活跃环境失败: %s", err), "")
		}
		if err := cmdutil.CfgMgr.Save(); err != nil {
			return output.NewExitError(2, fmt.Sprintf("保存配置失败: %s", err), "")
		}

		if !getQuiet(cmd) {
			fmt.Fprintf(cmd.ErrOrStderr(), "配置已初始化: %s (%s)\n", name, baseURL)
		}
		return cmdutil.OutputJSON(map[string]any{"profile": name, "base_url": baseURL})
	},
}

func init() {
	configCmd.AddCommand(configInitCmd)
}

func getQuiet(cmd *cobra.Command) bool {
	q, _ := cmd.Flags().GetBool("quiet")
	return q
}
