package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	// 手动注册 completion 命令，覆盖 Cobra 自动生成的英文版本
	rootCmd.AddCommand(&cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "生成 Shell 自动补全脚本",
		Long: `为指定的 Shell 生成 wlt 自动补全脚本。

加载方式：

  PowerShell:
    wlt completion powershell | Out-String | Invoke-Expression

  Bash:
    source <(wlt completion bash)

  Zsh:
    wlt completion zsh > "${fpath[1]}/_wlt"

  Fish:
    wlt completion fish | source`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return rootCmd.GenBashCompletion(os.Stdout)
			case "zsh":
				return rootCmd.GenZshCompletion(os.Stdout)
			case "fish":
				return rootCmd.GenFishCompletion(os.Stdout, true)
			case "powershell":
				return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
			default:
				return fmt.Errorf("不支持的 Shell: %s", args[0])
			}
		},
	})
}
