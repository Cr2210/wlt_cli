package stock

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "库存明细",
	Long:  "查询库存明细记录：分页、详情、统计。",
}

func init() {
	stockCmd.AddCommand(recordCmd)
	recordCmd.AddCommand(newRecordListCmd())
	recordCmd.AddCommand(newRecordGetCmd())
	recordCmd.AddCommand(newRecordCountCmd())
}

func newRecordListCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询库存明细",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectIntFlags(cmd, params, "product-id", "warehouse-id")
			cmdutil.CollectStringFlags(cmd, params, "product-name", "type", "start-time", "end-time")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock-record/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询明细失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().Int64("product-id", 0, "产品 ID")
	c.Flags().Int64("warehouse-id", 0, "仓库 ID")
	c.Flags().String("product-name", "", "产品名称（模糊）")
	c.Flags().String("type", "", "类型")
	c.Flags().String("start-time", "", "开始时间")
	c.Flags().String("end-time", "", "结束时间")
	return c
}

func newRecordGetCmd() *cobra.Command {
	var id int64
	c := &cobra.Command{
		Use:   "get",
		Short: "获取库存明细详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock-record/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取明细失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "明细 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

func newRecordCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "count",
		Short: "获取库存明细总数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock-record/total-count", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计明细失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}
