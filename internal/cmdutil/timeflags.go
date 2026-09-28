package cmdutil

import "github.com/spf13/cobra"

// AddTimeRangeFlags 注册 --<flagName>-start / --<flagName>-end 时间范围 flag，
// 查询时折叠为 <paramKey>[0] / <paramKey>[1] 数组参数。
// flagName 为 kebab-case（如 "load-date"），paramKey 为 camelCase（如 "realLoadDate"）。
func AddTimeRangeFlags(c *cobra.Command, flagName, paramKey, label string) {
	c.Flags().String(flagName+"-start", "", label+"起（如 2026-07-01 00:00:00）")
	c.Flags().String(flagName+"-end", "", label+"止（如 2026-07-31 23:59:59）")
}

// CollectTimeRangeFlags 把非空的 --<flagName>-start/end 折叠为 <paramKey>[0]/[1]。
func CollectTimeRangeFlags(cmd *cobra.Command, params map[string]any, flagName, paramKey string) {
	if v, _ := cmd.Flags().GetString(flagName + "-start"); v != "" {
		params[paramKey+"[0]"] = v
	}
	if v, _ := cmd.Flags().GetString(flagName + "-end"); v != "" {
		params[paramKey+"[1]"] = v
	}
}
