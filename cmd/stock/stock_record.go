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
	recordCmd.AddCommand(newRecordPageCmd())
	recordCmd.AddCommand(newRecordPageCountCmd())
	recordCmd.AddCommand(newRecordGetCmd())
	recordCmd.AddCommand(newRecordCountCmd())
	recordCmd.AddCommand(newRecordRecordPageCmd())
	recordCmd.AddCommand(newRecordTotalCostCmd())
}

func newRecordPageCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "page",
		Short: "分页查询出入库明细",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlags(cmd, params,
				"product-id",
				"category-id",
				"warehouse-id",
				"biz-type",
				"biz-no",
				"create-time",
				"in-time",
				"metrics-name",
				"product-name",
				"batch-no",
				"keyword",
				"headers",
			)

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock-record/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询出入库明细失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().String("product-id", "", "产品编号")
	c.Flags().String("category-id", "", "产品分类编号")
	c.Flags().String("warehouse-id", "", "仓库编号")
	c.Flags().String("biz-type", "", "业务类型")
	c.Flags().String("biz-no", "", "业务单号")
	c.Flags().String("create-time", "", "操作时间")
	c.Flags().String("in-time", "", "出入库时间")
	c.Flags().String("metrics-name", "", "指标名称")
	c.Flags().String("product-name", "", "产品名称")
	c.Flags().String("batch-no", "", "批次号")
	c.Flags().String("keyword", "", "关键字")
	c.Flags().String("headers", "", "自定义导出表头")
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

func newRecordPageCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "page-count",
		Short: "按筛选统计出入库明细数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params,
				"product-id",
				"category-id",
				"warehouse-id",
				"biz-type",
				"biz-no",
				"create-time",
				"in-time",
				"metrics-name",
				"product-name",
				"batch-no",
				"keyword",
			)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock-record/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计出入库明细失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().String("product-id", "", "产品编号")
	c.Flags().String("category-id", "", "产品分类编号")
	c.Flags().String("warehouse-id", "", "仓库编号")
	c.Flags().String("biz-type", "", "业务类型")
	c.Flags().String("biz-no", "", "业务单号")
	c.Flags().String("create-time", "", "操作时间")
	c.Flags().String("in-time", "", "出入库时间")
	c.Flags().String("metrics-name", "", "指标名称")
	c.Flags().String("product-name", "", "产品名称")
	c.Flags().String("batch-no", "", "批次号")
	c.Flags().String("keyword", "", "关键字")
	return c
}

func newRecordRecordPageCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "record-page",
		Short: "按维度分页查询库存明细",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlags(cmd, params, "product-id", "warehouse-id", "metrics-name", "batch-no", "supplier-id", "source-supplier-id", "out-count")
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock-record/record-page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询明细失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().String("product-id", "", "产品 ID")
	c.Flags().String("warehouse-id", "", "仓库 ID")
	c.Flags().String("metrics-name", "", "指标名称")
	c.Flags().String("batch-no", "", "批次号")
	c.Flags().String("supplier-id", "", "供应商 ID")
	c.Flags().String("source-supplier-id", "", "关联供应商 ID")
	c.Flags().String("out-count", "", "出库数量")
	return c
}

func newRecordTotalCostCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "total-cost",
		Short: "按维度统计库存明细总成本",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "product-id", "warehouse-id", "metrics-name", "batch-no", "supplier-id", "source-supplier-id", "out-count")
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock-record/total-cost", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计明细成本失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().String("product-id", "", "产品 ID")
	c.Flags().String("warehouse-id", "", "仓库 ID")
	c.Flags().String("metrics-name", "", "指标名称")
	c.Flags().String("batch-no", "", "批次号")
	c.Flags().String("supplier-id", "", "供应商 ID")
	c.Flags().String("source-supplier-id", "", "关联供应商 ID")
	c.Flags().String("out-count", "", "出库数量")
	return c
}
