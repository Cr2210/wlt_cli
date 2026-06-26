package configcmd

import (
	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
)

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "显示当前配置",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := cmdutil.CfgMgr.GetConfig()
		// Mask sensitive fields
		result := map[string]any{
			"active":   cfg.Active,
			"profiles": map[string]any{},
		}
		profiles := result["profiles"].(map[string]any)
		for name, p := range cfg.Profiles {
			profiles[name] = map[string]any{
				"base_url":        p.BaseURL,
				"api_prefix":      p.APIPrefix,
				"enterprise_type": p.EnterpriseType,
			}
		}
		return cmdutil.OutputJSON(result)
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
}
