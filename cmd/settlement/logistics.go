package settlement

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

// logisticsSettlementAPIPath 是物流结算的后端端点根路径。
const logisticsSettlementAPIPath = "/erp/logistics-settlement"

// logisticsSettlementCmd 是物流结算单子命令组。
var logisticsSettlementCmd = &cobra.Command{
	Use:   "logistics",
	Short: "物流结算单管理",
}

func init() {
	settlementCmd.AddCommand(logisticsSettlementCmd)
	logisticsSettlementCmd.AddCommand(
		newLogisticsSettlementListCmd(),
		newLogisticsSettlementPageCountCmd(),
		cmdutil.CrudGetCmd(logisticsSettlementAPIPath, "物流结算单"),
		cmdutil.CrudCreateCmd(logisticsSettlementAPIPath, "物流结算单"),
		cmdutil.CrudUpdateCmd(logisticsSettlementAPIPath, "物流结算单"),
		cmdutil.CrudDeleteCmd(logisticsSettlementAPIPath, "物流结算单", false),
		newLogisticsSettlementUpdateStatusCmd(),
		newLogisticsSettlementMismatchListCmd(),
		newLogisticsSettlementAvailableWaybillCmd(),
		newLogisticsSettlementAvailableWaybillCountCmd(),
		newLogisticsSettlementRecalculateAllCmd(),
		newLogisticsSettlementIgnoreRemainderCmd(),
		newLogisticsSettlementExportExcelCmd(),
	)
}

// ---- list / page-count / export 共用筛选 ----
//
// 后端 /erp/logistics-settlement/page 筛选字段：
//   no / name / carrierName / carrierEnterpriseId / status
//   settlementDate[0]/[1]  结算日期范围（--start-date/--end-date 折叠，与 main 一致）

func registerLogisticsListFlags(c *cobra.Command) {
	c.Flags().String("no", "", "结算单号")
	c.Flags().String("name", "", "结算单名称")
	c.Flags().String("carrier-name", "", "承运商名称")
	c.Flags().String("carrier-enterprise-id", "", "承运商企业 ID")
	c.Flags().String("status", "", "状态")
	cmdutil.AddSettlementDateFlags(c)
}

func buildLogisticsListParams(cmd *cobra.Command, pageNo, pageSize int, withPage bool) map[string]any {
	params := map[string]any{}
	if withPage {
		params["pageNo"] = pageNo
		params["pageSize"] = pageSize
	}
	cmdutil.CollectStringFlags(cmd, params, "no", "name", "carrier-name", "carrier-enterprise-id", "status")
	cmdutil.CollectSettlementDate(cmd, params)
	return params
}

// ---- 分页查询 ----

func newLogisticsSettlementListCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询物流结算单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := buildLogisticsListParams(cmd, pageNo, pageSize, true)
			resp, err := cmdutil.GetClient().Get(context.Background(), logisticsSettlementAPIPath+"/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询物流结算单失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	registerLogisticsListFlags(c)
	return c
}

// ---- 分页计数（Σ含税金额/Σ结算金额/Σ开票金额/Σ忽略额） ----

func newLogisticsSettlementPageCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "page-count",
		Short: "按筛选统计物流结算单金额合计",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := buildLogisticsListParams(cmd, 0, 0, false)
			resp, err := cmdutil.GetClient().Get(context.Background(), logisticsSettlementAPIPath+"/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计物流结算单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	registerLogisticsListFlags(c)
	return c
}

// ---- 更新状态（后端为 @RequestParam，用 PutQuery） ----

func newLogisticsSettlementUpdateStatusCmd() *cobra.Command {
	var id int64
	var status int
	c := &cobra.Command{
		Use:   "update-status",
		Short: "更新物流结算单状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{"id": id, "status": status}
			resp, err := cmdutil.GetClient().PutQuery(context.Background(), logisticsSettlementAPIPath+"/update-status", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新物流结算单状态失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "物流结算单 ID")
	c.Flags().IntVar(&status, "status", 0, "目标状态")
	_ = c.MarkFlagRequired("id")
	_ = c.MarkFlagRequired("status")
	return c
}

// ---- 差异运单（结算金额 vs 运单金额不一致） ----

func newLogisticsSettlementMismatchListCmd() *cobra.Command {
	var pageNo, pageSize int
	var keyword string
	c := &cobra.Command{
		Use:   "mismatch-list",
		Short: "分页查询金额差异运单（运单金额与结算金额不一致）",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{"pageNo": pageNo, "pageSize": pageSize}
			cmdutil.CollectStringFlags(cmd, params, "keyword")
			resp, err := cmdutil.GetClient().Get(context.Background(), logisticsSettlementAPIPath+"/mismatch-page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询差异运单失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().StringVar(&keyword, "keyword", "", "关键字（结算单号/运单号/车牌号）")
	return c
}

// ---- 可选运单（候选/详情回显/结算单带出明细） ----

func registerLogisticsWaybillFlags(c *cobra.Command) {
	c.Flags().Int64("carrier-id", 0, "承运商 ID")
	c.Flags().String("carrier-name", "", "承运商名称")
	c.Flags().String("keyword", "", "关键字")
	c.Flags().Int64("settlement-id", 0, "物流结算单 ID")
}

func collectLogisticsWaybillParams(cmd *cobra.Command, pageNo, pageSize int, withPage bool) map[string]any {
	params := map[string]any{}
	if withPage {
		params["pageNo"] = pageNo
		params["pageSize"] = pageSize
	}
	cmdutil.CollectIntFlags(cmd, params, "carrier-id", "settlement-id")
	cmdutil.CollectStringFlags(cmd, params, "carrier-name", "keyword")
	return params
}

func newLogisticsSettlementAvailableWaybillCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "available-waybill",
		Short: "分页查询可选运单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := collectLogisticsWaybillParams(cmd, pageNo, pageSize, true)
			resp, err := cmdutil.GetClient().Get(context.Background(), logisticsSettlementAPIPath+"/available-waybill", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询可选运单失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	registerLogisticsWaybillFlags(c)
	return c
}

func newLogisticsSettlementAvailableWaybillCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "available-waybill-count",
		Short: "统计可选运单合计（金额/装货重量/卸货重量）",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := collectLogisticsWaybillParams(cmd, 0, 0, false)
			resp, err := cmdutil.GetClient().Get(context.Background(), logisticsSettlementAPIPath+"/available-waybill/count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计可选运单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	registerLogisticsWaybillFlags(c)
	return c
}

// ---- 全量重算 ----

func newLogisticsSettlementRecalculateAllCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "recalculate-all",
		Short: "全量重算所有物流订单的聚合结算明细（数据对齐，低频运维操作）",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), logisticsSettlementAPIPath+"/recalculate-all", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("全量重算失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

// ---- 忽略剩余结算/开票 ----

func newLogisticsSettlementIgnoreRemainderCmd() *cobra.Command {
	var id int64
	var ignoreType, reason string
	c := &cobra.Command{
		Use:   "ignore-remainder",
		Short: "忽略物流结算单剩余结算/开票金额",
		Long:  "忽略剩余结算/开票金额：剩余归零，状态翻转为已结算/已开票，需填原因。仅已审核且 invoiceFlag=无需开票 时可忽略开票剩余。",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body := map[string]any{
				"id":     id,
				"type":   ignoreType,
				"reason": reason,
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), logisticsSettlementAPIPath+"/ignore-remainder", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("忽略剩余结算/开票失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "物流结算单 ID")
	c.Flags().StringVar(&ignoreType, "type", "", "忽略类型：SETTLE-剩余结算，INVOICE-剩余开票")
	c.Flags().StringVar(&reason, "reason", "", "忽略原因")
	_ = c.MarkFlagRequired("id")
	_ = c.MarkFlagRequired("type")
	_ = c.MarkFlagRequired("reason")
	return c
}

// ---- 导出 Excel ----

func newLogisticsSettlementExportExcelCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "export",
		Short: "导出物流结算单 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := buildLogisticsListParams(cmd, 0, 0, false)
			cmdutil.CollectStringFlags(cmd, params, "headers")
			resp, err := cmdutil.GetClient().Get(context.Background(), logisticsSettlementAPIPath+"/export-excel", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出物流结算单失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	registerLogisticsListFlags(c)
	c.Flags().String("headers", "", "自定义导出表头")
	return c
}
