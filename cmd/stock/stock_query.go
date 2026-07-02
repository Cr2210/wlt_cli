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
	queryCmd.AddCommand(newStockQueryPageCmd())
	queryCmd.AddCommand(newStockQueryPageCountCmd())
	queryCmd.AddCommand(newStockQueryCountCmd())
	queryCmd.AddCommand(newStockQueryLedgerCmd())
	queryCmd.AddCommand(newStockQueryLedgerCountCmd())
	queryCmd.AddCommand(newStockQueryRecordCountCmd())
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

func newStockQueryPageCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "page",
		Short: "分页查询产品库存",
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
				"warehouse-id",
				"metrics-name",
				"batch-no",
				"send-address",
				"plan-no",
				"supplier-id",
				"supplier-name",
				"is-detail",
				"count-positive",
				"keyword",
				"headers",
			)

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询库存列表失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().String("product-id", "", "产品编号")
	c.Flags().String("warehouse-id", "", "仓库编号")
	c.Flags().String("metrics-name", "", "指标名称")
	c.Flags().String("batch-no", "", "批次号")
	c.Flags().String("send-address", "", "发货地")
	c.Flags().String("plan-no", "", "方案号")
	c.Flags().String("supplier-id", "", "供应商编号")
	c.Flags().String("supplier-name", "", "供应商名称")
	c.Flags().String("is-detail", "", "是否明细")
	c.Flags().String("count-positive", "", "库存为正数")
	c.Flags().String("keyword", "", "关键字")
	c.Flags().String("headers", "", "自定义导出表头")
	return c
}

func newStockQueryCountCmd() *cobra.Command {
	var (
		productId   int64
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
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock/get-count", map[string]any{"productId": productId})
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
	c.Flags().Int64Var(&productId, "product-id", 0, "产品 ID")
	c.Flags().Int64Var(&warehouseId, "warehouse-id", 0, "仓库 ID（可选，获取指定仓库统计）")
	c.Flags().StringVar(&metricName, "metric-name", "", "指标名称（可选）")
	_ = c.MarkFlagRequired("product-id")
	return c
}

func newStockQueryLedgerCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "ledger",
		Short: "查询库存台账（库存行的历史进出明细）",
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
				"warehouse-id",
				"metrics-name",
				"batch-no",
				"send-address",
				"plan-no",
				"supplier-id",
				"supplier-name",
				"is-detail",
				"count-positive",
				"keyword",
			)

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock/batch-detail", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询库存台账失败: %s", err), "")
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
	c.Flags().String("product-id", "", "产品编号")
	c.Flags().String("warehouse-id", "", "仓库编号")
	c.Flags().String("metrics-name", "", "指标名称")
	c.Flags().String("batch-no", "", "批次号")
	c.Flags().String("send-address", "", "发货地")
	c.Flags().String("plan-no", "", "方案号")
	c.Flags().String("supplier-id", "", "供应商编号")
	c.Flags().String("supplier-name", "", "供应商名称")
	c.Flags().String("is-detail", "", "是否明细")
	c.Flags().String("count-positive", "", "库存为正数")
	c.Flags().String("keyword", "", "关键字")
	return c
}

func newStockQueryPageCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "page-count",
		Short: "按筛选统计产品库存数量与总成本",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params,
				"product-id",
				"warehouse-id",
				"metrics-name",
				"batch-no",
				"send-address",
				"plan-no",
				"supplier-id",
				"supplier-name",
				"is-detail",
				"count-positive",
				"keyword",
			)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计库存失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().String("product-id", "", "产品编号")
	c.Flags().String("warehouse-id", "", "仓库编号")
	c.Flags().String("metrics-name", "", "指标名称")
	c.Flags().String("batch-no", "", "批次号")
	c.Flags().String("send-address", "", "发货地")
	c.Flags().String("plan-no", "", "方案号")
	c.Flags().String("supplier-id", "", "供应商编号")
	c.Flags().String("supplier-name", "", "供应商名称")
	c.Flags().String("is-detail", "", "是否明细")
	c.Flags().String("count-positive", "", "库存为正数")
	c.Flags().String("keyword", "", "关键字")
	return c
}

func newStockQueryLedgerCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "ledger-count",
		Short: "按筛选统计库存台账数量与总成本",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params,
				"product-id",
				"warehouse-id",
				"metrics-name",
				"batch-no",
				"send-address",
				"plan-no",
				"supplier-id",
				"supplier-name",
				"is-detail",
				"count-positive",
				"keyword",
			)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock/detail-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计库存台账失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().String("product-id", "", "产品编号")
	c.Flags().String("warehouse-id", "", "仓库编号")
	c.Flags().String("metrics-name", "", "指标名称")
	c.Flags().String("batch-no", "", "批次号")
	c.Flags().String("send-address", "", "发货地")
	c.Flags().String("plan-no", "", "方案号")
	c.Flags().String("supplier-id", "", "供应商编号")
	c.Flags().String("supplier-name", "", "供应商名称")
	c.Flags().String("is-detail", "", "是否明细")
	c.Flags().String("count-positive", "", "库存为正数")
	c.Flags().String("keyword", "", "关键字")
	return c
}

func newStockQueryRecordCountCmd() *cobra.Command {
	var productId, warehouseId, metricsName, batchNo, enterpriseId, sourceEnterpriseId string
	c := &cobra.Command{
		Use:   "stock-record-count",
		Short: "按维度键统计库存明细数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"productId":   productId,
				"warehouseId": warehouseId,
				"metricsName": metricsName,
				"batchNo":     batchNo,
			}
			if enterpriseId != "" {
				params["enterpriseId"] = enterpriseId
			}
			if sourceEnterpriseId != "" {
				params["sourceEnterpriseId"] = sourceEnterpriseId
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/stock/get-stock-record-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计库存明细失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&productId, "product-id", "", "产品 ID")
	c.Flags().StringVar(&warehouseId, "warehouse-id", "", "仓库 ID")
	c.Flags().StringVar(&metricsName, "metrics-name", "", "指标名称")
	c.Flags().StringVar(&batchNo, "batch-no", "", "批次号")
	c.Flags().StringVar(&enterpriseId, "enterprise-id", "", "企业 ID")
	c.Flags().StringVar(&sourceEnterpriseId, "source-enterprise-id", "", "来源企业 ID")
	_ = c.MarkFlagRequired("product-id")
	_ = c.MarkFlagRequired("warehouse-id")
	_ = c.MarkFlagRequired("metrics-name")
	_ = c.MarkFlagRequired("batch-no")
	return c
}
