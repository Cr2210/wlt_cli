package invoice

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

// ---- 发票明细（/erp/invoice-item/*） ----

var invoiceItemCmd = &cobra.Command{
	Use:   "item",
	Short: "发票明细",
	Long:  "发票明细查询（/erp/invoice-item/list，返回明细数组）。",
}

var invoiceItemListCmd = &cobra.Command{
	Use:   "list",
	Short: "查询发票明细",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmdutil.EnsureClient(); err != nil {
			return err
		}
		params := map[string]any{}
		cmdutil.CollectStringFlags(cmd, params, "no", "invoice-type", "type", "enterprise-name",
			"user-id", "write-off-status", "status", "order-no")
		cmdutil.CollectTimeRangeFlags(cmd, params, "invoice-date", "invoiceDate")
		cmdutil.CollectTimeRangeFlags(cmd, params, "create-time", "createTime")
		if v, _ := cmd.Flags().GetInt("page-no"); v > 0 {
			params["pageNo"] = v
		}
		if v, _ := cmd.Flags().GetInt("page-size"); v > 0 {
			params["pageSize"] = v
		}
		// 响应为裸数组（List<InvoiceItemOrderDetailVO>），不走 ParsePagedJSON。
		resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/invoice-item/list", params)
		if err != nil {
			return output.NewExitError(5, fmt.Sprintf("查询发票明细失败: %s", err), "")
		}
		return cmdutil.OutputJSON(json.RawMessage(resp.Data))
	},
}

func init() {
	invoiceItemCmd.AddCommand(invoiceItemListCmd)
	invoiceCmd.AddCommand(invoiceItemCmd)

	f := invoiceItemListCmd.Flags()
	f.String("no", "", "发票号")
	f.String("invoice-type", "", "发票种类")
	f.String("type", "", "业务类型")
	f.String("enterprise-name", "", "企业名称")
	f.String("user-id", "", "经办人 ID")
	f.String("write-off-status", "", "核销状态")
	f.String("status", "", "状态")
	f.String("order-no", "", "订单号")
	cmdutil.AddTimeRangeFlags(invoiceItemListCmd, "invoice-date", "invoiceDate", "开票日期")
	cmdutil.AddTimeRangeFlags(invoiceItemListCmd, "create-time", "createTime", "创建时间")
	f.Int("page-no", 1, "页码")
	f.Int("page-size", 20, "每页数量")
}
