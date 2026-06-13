package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	// 手动注册 help 命令，覆盖 Cobra 自动生成的英文版本
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:   "help [命令]",
		Short: "显示任意命令的帮助信息",
		Long:  `显示任意命令的帮助信息，包括用法、可用子命令和标志。`,
		DisableFlagsInUseLine: true,
		Args:                  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return rootCmd.Help()
			}
			target, _, err := rootCmd.Find(args)
			if err != nil {
				return fmt.Errorf("未知命令: %s", args[0])
			}
			return target.Help()
		},
	})

	// 全局帮助模板中文化（一次性替换，避免 "Flags:" 误伤 "Global Flags:"）
	tpl := rootCmd.UsageTemplate()
	replacer := strings.NewReplacer(
		"Usage:", "用法：",
		"Available Commands:", "可用命令：",
		"Global Flags:", "全局标志：",
		"Flags:", "标志：",
		`Use "`, `使用 "`,
		"for more information about a command.", "查看命令的详细帮助信息。",
	)
	rootCmd.SetUsageTemplate(replacer.Replace(tpl))
}
