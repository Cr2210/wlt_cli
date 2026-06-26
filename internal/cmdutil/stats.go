package cmdutil

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/output"
)

// NewStatsGetCmd creates a stats/dashboard subcommand with time range and optional extra flags.
func NewStatsGetCmd(name, apiPath, label string, extraFlags []FlagSpec) *cobra.Command {
	c := &cobra.Command{
		Use:   name,
		Short: label,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			CollectTimeRangeFlags(cmd, params)
			for _, f := range extraFlags {
				CollectStringFlag(cmd, params, f.Name)
			}
			resp, err := GetClient().Get(context.Background(), apiPath+"/"+name, params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("%s失败: %s", label, err), "")
			}
			return OutputJSON(json.RawMessage(resp.Data))
		},
	}
	AddStatsFlags(c)
	for _, f := range extraFlags {
		c.Flags().String(f.Name, "", f.Usage)
	}
	return c
}

// CollectTimeRangeFlags collects time range and sort flags from the command into params.
func CollectTimeRangeFlags(cmd *cobra.Command, params map[string]any) {
	t, _ := cmd.Flags().GetString("type")
	if t == "" {
		t = "month"
	}
	params["type"] = t

	startTime, _ := cmd.Flags().GetString("start-time")
	endTime, _ := cmd.Flags().GetString("end-time")
	if startTime == "" || endTime == "" {
		defStart, defEnd := DefaultMonthRange()
		if startTime == "" {
			startTime = defStart
		}
		if endTime == "" {
			endTime = defEnd
		}
	}
	params["startTime"] = startTime
	params["endTime"] = endTime

	sortBy, _ := cmd.Flags().GetString("sort-by")
	if sortBy == "" {
		sortBy = "amount"
	}
	params["sortBy"] = sortBy
}

// DefaultMonthRange returns the start of the current month and now.
func DefaultMonthRange() (string, string) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return start.Format("2006-01-02 15:04:05"), now.Format("2006-01-02 15:04:05")
}

// AddStatsFlags adds standard time range and sort flags to a stats command.
func AddStatsFlags(c *cobra.Command) {
	c.Flags().String("type", "month", "时间类型（day/month/year）")
	c.Flags().String("start-time", "", "开始时间，默认当月1号")
	c.Flags().String("end-time", "", "结束时间，默认当前时间")
	c.Flags().String("sort-by", "amount", "排序字段，默认 amount")
}
