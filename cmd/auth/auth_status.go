package auth

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "查看认证状态",
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr := newAuthManager()
		loggedIn, expiresAt, err := mgr.Status()
		if err != nil {
			return output.NewExitError(2, fmt.Sprintf("获取认证状态失败: %s", err), "运行 wlt config init 初始化")
		}
		status := "not_logged_in"
		if loggedIn {
			status = "logged_in"
		}
		result := map[string]any{
			"status":  status,
			"profile": cmdutil.CfgMgr.GetActive(),
		}
		if loggedIn && expiresAt > 0 {
			result["expiresAt"] = expiresAt
			expiry := time.UnixMilli(expiresAt)
			result["expiresAtFormatted"] = expiry.Format(time.RFC3339)
			if time.Now().After(expiry) {
				result["status"] = "token_expired"
			}
		}
		if !getQuiet(cmd) {
			fmt.Fprintf(cmd.ErrOrStderr(), "状态: %s (环境: %s)\n", status, cmdutil.CfgMgr.GetActive())
		}
		return cmdutil.OutputJSON(result)
	},
}

func init() {
	authCmd.AddCommand(statusCmd)
}
