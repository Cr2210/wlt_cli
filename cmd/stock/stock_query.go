package stock

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "库存查询",
	Long:  "查询产品库存：获取、分页、统计、批次明细。",
}

func init() {
	stockCmd.AddCommand(queryCmd)
	queryCmd.AddCommand(newStockQueryGetCmd())
	queryCmd.AddCommand(newStockQueryListCmd())
	queryCmd.AddCommand(newStockQueryCountCmd())
	queryCmd.AddCommand(newStockQueryBatchDetailCmd())
}

func newStockQueryGetCmd() *cobra.Command {
	var (
		id          int64
		productId   int64
		warehouseId int64
	)
	c := &cobra.Command{
		Use:   "get",
		Short: "获取单个产品库存",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			if id > 0 {
				params["id"] = id
			}
			if productId > 0 {
				params["productId"] = productId
			}
			if warehouseId > 0 {
				params["warehouseId"] = warehouseId
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock/get", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询库存失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "库存记录 ID")
	c.Flags().Int64Var(&productId, "product-id", 0, "产品 ID")
	c.Flags().Int64Var(&warehouseId, "warehouse-id", 0, "仓库 ID")
	return c
}

func newStockQueryListCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询产品库存",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectIntFlags(cmd, params, "warehouse-id", "product-id")
			cmdutil.CollectStringFlags(cmd, params, "product-name")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询库存列表失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().Int64("warehouse-id", 0, "仓库 ID")
	c.Flags().Int64("product-id", 0, "产品 ID")
	c.Flags().String("product-name", "", "产品名称（模糊）")
	return c
}

func newStockQueryCountCmd() *cobra.Command {
	var (
		warehouseId int64
		metricName  string
	)
	c := &cobra.Command{
		Use:   "count",
		Short: "统计库存数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock/get-count", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计库存失败: %s", err), "")
			}
			result := map[string]any{
				"stock_count": json.RawMessage(resp.Data),
			}
			if warehouseId > 0 || metricName != "" {
				params := map[string]any{}
				if warehouseId > 0 {
					params["warehouseId"] = warehouseId
				}
				if metricName != "" {
					params["metricName"] = metricName
				}
				resp2, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock/get-warehouse-count", params)
				if err != nil {
					return output.NewExitError(5, fmt.Sprintf("统计仓库库存失败: %s", err), "")
				}
				result["warehouse_count"] = json.RawMessage(resp2.Data)
			}
			return cmdutil.OutputJSON(result)
		},
	}
	c.Flags().Int64Var(&warehouseId, "warehouse-id", 0, "仓库 ID（可选，获取指定仓库统计）")
	c.Flags().StringVar(&metricName, "metric-name", "", "指标名称（可选）")
	return c
}

func newStockQueryBatchDetailCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "batch-detail",
		Short: "获取批次明细",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectIntFlags(cmd, params, "product-id", "warehouse-id")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock/batch-detail", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询批次明细失败: %s", err), "")
			}
			var paged struct {
				List  json.RawMessage `json:"list"`
				Total int64           `json:"total"`
			}
			if err := json.Unmarshal(resp.Data, &paged); err != nil {
				return cmdutil.OutputJSON(json.RawMessage(resp.Data))
			}
			var list any
			if err := json.Unmarshal(paged.List, &list); err != nil {
				list = []any{}
			}
			return cmdutil.OutputPagedJSON(list, paged.Total, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().Int64("product-id", 0, "产品 ID")
	c.Flags().Int64("warehouse-id", 0, "仓库 ID")
	return c
}
