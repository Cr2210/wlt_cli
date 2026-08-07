package stock

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

// queryPageFilters 是 stock query 的 page/page-count/ledger/ledger-count 共享筛选 flag。
var queryPageFilters = []cmdutil.FlagSpec{
	{Name: "product-id", Usage: "产品编号"},
	{Name: "warehouse-id", Usage: "仓库编号"},
	{Name: "metrics-name", Usage: "指标名称"},
	{Name: "batch-no", Usage: "批次号"},
	{Name: "send-address", Usage: "发货地"},
	{Name: "plan-no", Usage: "方案号"},
	{Name: "supplier-id", Usage: "供应商编号"},
	{Name: "supplier-name", Usage: "供应商名称"},
	{Name: "is-detail", Usage: "是否明细"},
	{Name: "count-positive", Usage: "库存为正数"},
	{Name: "keyword", Usage: "关键字"},
}

// recordPageFilters 是 stock record 的 page/page-count 共享筛选 flag。
var recordPageFilters = []cmdutil.FlagSpec{
	{Name: "product-id", Usage: "产品编号"},
	{Name: "category-id", Usage: "产品分类编号"},
	{Name: "warehouse-id", Usage: "仓库编号"},
	{Name: "biz-type", Usage: "业务类型"},
	{Name: "biz-no", Usage: "业务单号"},
	{Name: "create-time", Usage: "操作时间"},
	{Name: "in-time", Usage: "出入库时间"},
	{Name: "metrics-name", Usage: "指标名称"},
	{Name: "product-name", Usage: "产品名称"},
	{Name: "batch-no", Usage: "批次号"},
	{Name: "keyword", Usage: "关键字"},
}

// recordDimensionFilters 是 stock record 的 record-page/total-cost 共享维度 flag。
var recordDimensionFilters = []cmdutil.FlagSpec{
	{Name: "product-id", Usage: "产品 ID"},
	{Name: "warehouse-id", Usage: "仓库 ID"},
	{Name: "metrics-name", Usage: "指标名称"},
	{Name: "batch-no", Usage: "批次号"},
	{Name: "supplier-id", Usage: "供应商 ID"},
	{Name: "source-supplier-id", Usage: "关联供应商 ID"},
	{Name: "out-count", Usage: "出库数量"},
}

func addFilterFlags(c *cobra.Command, filters []cmdutil.FlagSpec) {
	for _, f := range filters {
		c.Flags().String(f.Name, "", f.Usage)
	}
}

func collectFilterFlags(cmd *cobra.Command, params map[string]any, filters []cmdutil.FlagSpec) {
	for _, f := range filters {
		cmdutil.CollectStringFlag(cmd, params, f.Name)
	}
}

// newStockPagedCmd 构建分页查询命令：GET {apiPath}/{endpoint}，输出 ParsePagedJSON。
// use=命令名；endpoint=page/batch-detail/record-page；withHeaders=是否注册 --headers。
func newStockPagedCmd(use, apiPath, endpoint, short string, filters []cmdutil.FlagSpec, withHeaders bool) *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{"pageNo": pageNo, "pageSize": pageSize}
			collectFilterFlags(cmd, params, filters)
			if withHeaders {
				cmdutil.CollectStringFlag(cmd, params, "headers")
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), apiPath+"/"+endpoint, params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("%s失败: %s", short, err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	addFilterFlags(c, filters)
	if withHeaders {
		c.Flags().String("headers", "", "自定义导出表头")
	}
	return c
}

// newStockCountCmd 构建统计命令：GET {apiPath}/{endpoint}，输出 OutputJSON。
func newStockCountCmd(use, apiPath, endpoint, short string, filters []cmdutil.FlagSpec) *cobra.Command {
	c := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			collectFilterFlags(cmd, params, filters)
			resp, err := cmdutil.GetClient().Get(context.Background(), apiPath+"/"+endpoint, params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("%s失败: %s", short, err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	addFilterFlags(c, filters)
	return c
}
