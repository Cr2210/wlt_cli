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
	recordCmd.AddCommand(newStockPagedCmd("page", "/erp/stock-record", "page", "分页查询出入库明细", recordPageFilters, true))
	recordCmd.AddCommand(newStockCountCmd("page-count", "/erp/stock-record", "page-count", "按筛选统计出入库明细数量", recordPageFilters))
	recordCmd.AddCommand(newRecordGetCmd())
	recordCmd.AddCommand(newRecordCountCmd())
	recordCmd.AddCommand(newStockPagedCmd("record-page", "/erp/stock-record", "record-page", "按维度分页查询库存明细", recordDimensionFilters, false))
	recordCmd.AddCommand(newStockCountCmd("total-cost", "/erp/stock-record", "total-cost", "按维度统计库存明细总成本", recordDimensionFilters))
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
