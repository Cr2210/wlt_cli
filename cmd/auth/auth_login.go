package auth

import (
	"context"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/auth"
	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "交互式登录",
	Long:  "通过用户名和密码登录维链通 ERP 系统。",
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, err := cmdutil.CfgMgr.ActiveProfile()
		if err != nil {
			return output.NewExitError(2, fmt.Sprintf("当前环境配置不存在: %s", err), "运行 wlt config init 初始化")
		}

		// Interactive form — ask for company name instead of tenant ID
		var (
			baseURL     = profile.BaseURL
			companyName string
			username    string
			password    string
		)

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("API 地址").Value(&baseURL).Placeholder("https://erpsit.api.w-lian.com"),
				huh.NewInput().Title("公司名称").Value(&companyName).Placeholder("请输入公司名称"),
				huh.NewInput().Title("用户名").Value(&username).Placeholder("请输入用户名"),
				huh.NewInput().Title("密码").Value(&password).Placeholder("请输入密码").EchoMode(huh.EchoModePassword),
			),
		)

		if err := form.Run(); err != nil {
			return output.NewExitError(1, fmt.Sprintf("表单输入失败: %s", err), "")
		}

		// Resolve tenant ID from company name
		fullURL := baseURL + profile.APIPrefix
		mgr := newAuthManager()
		tenantID, err := mgr.GetTenantIDByName(context.Background(), fullURL, companyName)
		if err != nil {
			return output.NewExitError(4, fmt.Sprintf("查询租户失败: %s", err), "请检查公司名称是否正确")
		}

		if !getQuiet(cmd) {
			fmt.Fprintf(cmd.ErrOrStderr(), "已匹配公司 %q，租户 ID: %s\n", companyName, tenantID)
		}

		// Update profile with form values
		profile.BaseURL = baseURL
		profile.TenantID = tenantID

		// Login
		result, err := mgr.Login(context.Background(), profile.BaseURL+profile.APIPrefix, profile.TenantID, username, password)
		if err != nil {
			return output.NewExitError(3, fmt.Sprintf("登录失败: %s", err), "请检查用户名和密码")
		}

		if !getQuiet(cmd) {
			fmt.Fprintf(cmd.ErrOrStderr(), "登录成功! userId=%s\n", result.UserID)
		}
		return cmdutil.OutputJSON(map[string]any{
			"userId":      result.UserID,
			"expiresTime": result.ExpiresTime,
		})
	},
}

func init() {
	authCmd.AddCommand(loginCmd)
}

func newAuthManager() *auth.Manager {
	return cmdutil.GetAuthMgr()
}

func getQuiet(cmd *cobra.Command) bool {
	q, _ := cmd.Flags().GetBool("quiet")
	return q
}
