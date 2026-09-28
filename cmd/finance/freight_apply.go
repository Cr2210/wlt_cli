package finance

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

// financeFreightApplyAPIPath 是运费申请（物流付款申请）的后端端点根路径。
const financeFreightApplyAPIPath = "/erp/finance-freight-apply"

// financeFreightApplyCmd 是运费申请子命令组。
var financeFreightApplyCmd = &cobra.Command{
	Use:   "freight-apply",
	Short: "运费申请（物流付款申请）管理",
}

func init() {
	financeCmd.AddCommand(financeFreightApplyCmd)
	financeFreightApplyCmd.AddCommand(
		newFinanceFreightApplyListCmd(),
		cmdutil.CrudGetCmd(financeFreightApplyAPIPath, "运费申请"),
		cmdutil.CrudCreateCmd(financeFreightApplyAPIPath, "运费申请"),
		cmdutil.CrudUpdateCmd(financeFreightApplyAPIPath, "运费申请"),
		cmdutil.CrudDeleteCmd(financeFreightApplyAPIPath, "运费申请", false),
		newFinanceFreightApplyUpdateStatusCmd(),
		newFinanceFreightApplySummaryCmd(),
		newFinanceFreightApplyAvailableWaybillCmd(),
		newFinanceFreightApplyAvailableWaybillCountCmd(),
		newFinanceFreightApplyAvailableSettlementCmd(),
		newFinanceFreightApplyExportExcelCmd(),
	)
}

// ---- list / summary / export 共用筛选 ----
//
// 后端 /erp/finance-freight-apply/page 筛选字段：
//   no / carrierEnterpriseId / carrierName / paymentAccountId / receiptAccountId
//   serviceUserId / approveStatus / payDate[0]/[1]、createTime[0]/[1]

func registerFreightApplyListFlags(c *cobra.Command) {
	c.Flags().String("no", "", "申请单号")
	c.Flags().String("carrier-enterprise-id", "", "承运商企业 ID")
	c.Flags().String("carrier-name", "", "承运商名称")
	c.Flags().String("payment-account-id", "", "付款账户 ID")
	c.Flags().String("receipt-account-id", "", "收款账户 ID")
	c.Flags().String("service-user-id", "", "业务员 ID")
	c.Flags().String("approve-status", "", "审批状态")
	cmdutil.AddTimeRangeFlags(c, "pay-date", "payDate", "付款日期")
	cmdutil.AddTimeRangeFlags(c, "create-time", "createTime", "创建时间")
}

func collectFreightApplyListParams(cmd *cobra.Command, pageNo, pageSize int, withPage bool) map[string]any {
	params := map[string]any{}
	if withPage {
		params["pageNo"] = pageNo
		params["pageSize"] = pageSize
	}
	cmdutil.CollectStringFlags(cmd, params,
		"no", "carrier-enterprise-id", "carrier-name",
		"payment-account-id", "receipt-account-id", "service-user-id", "approve-status")
	cmdutil.CollectTimeRangeFlags(cmd, params, "pay-date", "payDate")
	cmdutil.CollectTimeRangeFlags(cmd, params, "create-time", "createTime")
	return params
}

// ---- 分页查询 ----

func newFinanceFreightApplyListCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询运费申请",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := collectFreightApplyListParams(cmd, pageNo, pageSize, true)
			resp, err := cmdutil.GetClient().Get(context.Background(), financeFreightApplyAPIPath+"/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询运费申请失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	registerFreightApplyListFlags(c)
	return c
}

// ---- 处理运费申请（后端为 @RequestParam，用 PutQuery） ----

func newFinanceFreightApplyUpdateStatusCmd() *cobra.Command {
	var id int64
	var approveStatus int
	c := &cobra.Command{
		Use:   "update-status",
		Short: "处理运费申请（审批）",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{"id": id, "approveStatus": approveStatus}
			resp, err := cmdutil.GetClient().PutQuery(context.Background(), financeFreightApplyAPIPath+"/update-status", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("处理运费申请失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "运费申请 ID")
	c.Flags().IntVar(&approveStatus, "approve-status", 0, "审批状态")
	_ = c.MarkFlagRequired("id")
	_ = c.MarkFlagRequired("approve-status")
	return c
}

// ---- 合计（按 list 同条件：数量分状态统计 + Σ总金额） ----

func newFinanceFreightApplySummaryCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "summary",
		Short: "按筛选统计运费申请合计",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := collectFreightApplyListParams(cmd, 0, 0, false)
			resp, err := cmdutil.GetClient().Get(context.Background(), financeFreightApplyAPIPath+"/summary", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计运费申请失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	registerFreightApplyListFlags(c)
	return c
}

// ---- 可选运单（选择运单/详情回显/结算单带出明细） ----

func registerFreightApplyWaybillFlags(c *cobra.Command) {
	c.Flags().String("waybill-no", "", "运单号")
	c.Flags().Int64("carrier-id", 0, "承运商 ID")
	c.Flags().Int64("settlement-id", 0, "物流结算单 ID")
	c.Flags().Int64("apply-id", 0, "运费申请 ID")
	cmdutil.AddTimeRangeFlags(c, "load-date", "realLoadDate", "实际装货日期")
	cmdutil.AddTimeRangeFlags(c, "unload-date", "realUnloadDate", "实际卸货日期")
}

func collectFreightApplyWaybillParams(cmd *cobra.Command, pageNo, pageSize int, withPage bool) map[string]any {
	params := map[string]any{}
	if withPage {
		params["pageNo"] = pageNo
		params["pageSize"] = pageSize
	}
	cmdutil.CollectStringFlags(cmd, params, "waybill-no")
	cmdutil.CollectIntFlags(cmd, params, "carrier-id", "settlement-id", "apply-id")
	cmdutil.CollectTimeRangeFlags(cmd, params, "load-date", "realLoadDate")
	cmdutil.CollectTimeRangeFlags(cmd, params, "unload-date", "realUnloadDate")
	return params
}

func newFinanceFreightApplyAvailableWaybillCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "available-waybill",
		Short: "分页查询可选运单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := collectFreightApplyWaybillParams(cmd, pageNo, pageSize, true)
			resp, err := cmdutil.GetClient().Get(context.Background(), financeFreightApplyAPIPath+"/available-waybill", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询可选运单失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	registerFreightApplyWaybillFlags(c)
	return c
}

func newFinanceFreightApplyAvailableWaybillCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "available-waybill-count",
		Short: "统计可选运单（Σ结算重量/Σ含税金额/Σ装货重量/Σ卸货重量）",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := collectFreightApplyWaybillParams(cmd, 0, 0, false)
			resp, err := cmdutil.GetClient().Get(context.Background(), financeFreightApplyAPIPath+"/available-waybill/count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计可选运单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	registerFreightApplyWaybillFlags(c)
	return c
}

// ---- 可选物流结算单（availableCount=0 置灰） ----

func newFinanceFreightApplyAvailableSettlementCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "available-settlement",
		Short: "分页查询可选物流结算单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{"pageNo": pageNo, "pageSize": pageSize}
			cmdutil.CollectStringFlags(cmd, params, "no", "carrier-enterprise-id", "carrier-name")
			cmdutil.CollectIntFlags(cmd, params, "apply-id")
			cmdutil.CollectTimeRangeFlags(cmd, params, "settlement-date", "settlementDate")
			resp, err := cmdutil.GetClient().Get(context.Background(), financeFreightApplyAPIPath+"/available-settlement", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询可选物流结算单失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().String("no", "", "结算单号")
	c.Flags().String("carrier-enterprise-id", "", "承运商企业 ID")
	c.Flags().String("carrier-name", "", "承运商名称")
	c.Flags().Int64("apply-id", 0, "运费申请 ID")
	cmdutil.AddTimeRangeFlags(c, "settlement-date", "settlementDate", "结算日期")
	return c
}

// ---- 导出 Excel ----

func newFinanceFreightApplyExportExcelCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "export",
		Short: "导出运费申请 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := collectFreightApplyListParams(cmd, 0, 0, false)
			cmdutil.CollectStringFlags(cmd, params, "headers")
			resp, err := cmdutil.GetClient().Get(context.Background(), financeFreightApplyAPIPath+"/export-excel", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出运费申请失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	registerFreightApplyListFlags(c)
	c.Flags().String("headers", "", "自定义导出表头")
	return c
}
