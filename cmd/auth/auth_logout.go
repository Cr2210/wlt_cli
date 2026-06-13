package auth

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "清除登录凭证",
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr := newAuthManager()
		if err := mgr.Logout(); err != nil {
			return output.NewExitError(2, fmt.Sprintf("登出失败: %s", err), "")
		}
		if !getQuiet(cmd) {
			fmt.Fprintln(cmd.ErrOrStderr(), "已登出")
		}
		return cmdutil.OutputJSON(map[string]any{"status": "logged_out"})
	},
}

func init() {
	authCmd.AddCommand(logoutCmd)
}
