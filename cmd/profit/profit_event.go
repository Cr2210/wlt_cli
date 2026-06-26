package profit

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var profitEventCmd = &cobra.Command{
	Use:   "profit-event",
	Short: "利润事件管理",
}

func init() {
	profitEventCmd.AddCommand(
		newProfitEventPageCmd(),
		newProfitEventStatisticsCmd(),
		newProfitEventTypesCmd(),
		newProfitEventRetryCmd(),
		newProfitEventCleanExpiredCmd(),
		newProfitEventHealthCmd(),
	)
}

// ---- 分页查询 ----

func newProfitEventPageCmd() *cobra.Command {
	var pageNo, pageSize int
	var eventType, status string

	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询利润事件",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlags(cmd, params, "event-type", "status")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/profit/event/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询利润事件失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().StringVar(&eventType, "event-type", "", "事件类型")
	c.Flags().StringVar(&status, "status", "", "状态")
	return c
}

// ---- 统计数据 ----

func newProfitEventStatisticsCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "statistics",
		Short: "获取利润事件统计",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/profit/event/statistics", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取利润事件统计失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

// ---- 事件类型 ----

func newProfitEventTypesCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "types",
		Short: "获取利润事件类型",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/profit/event/types", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取利润事件类型失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

// ---- 重试事件 ----

func newProfitEventRetryCmd() *cobra.Command {
	var eventId int64

	c := &cobra.Command{
		Use:   "retry",
		Short: "重试利润事件",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), fmt.Sprintf("/erp/profit/event/retry/%d", eventId), nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("重试利润事件失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&eventId, "event-id", 0, "事件 ID")
	_ = c.MarkFlagRequired("event-id")
	return c
}

// ---- 清理过期事件 ----

func newProfitEventCleanExpiredCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "clean-expired",
		Short: "清理过期利润事件",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/profit/event/clean-expired", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("清理过期事件失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

// ---- 健康检查 ----

func newProfitEventHealthCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "health",
		Short: "利润事件健康检查",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/profit/event/health", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("健康检查失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}
