package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/build"
	"github.com/weiliantong/cli/internal/cmdutil"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	// Bypass PersistentPreRunE so version works without config file
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("wlt version %s\n", build.Version)
		return cmdutil.OutputJSON(map[string]any{
			"version": build.Version,
		})
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
