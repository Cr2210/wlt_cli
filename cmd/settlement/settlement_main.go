package settlement

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

func init() {
	settlementCmd.AddCommand(settlementMainCmd)
}

var settlementMainCmd = &cobra.Command{
	Use:   "main",
	Short: "结算单管理",
}

func init() {
	settlementMainCmd.AddCommand(
		newSettlementMainListCmd(),
		newSettlementMainPageCountCmd(),
		newSettlementMainGetCmd(),
		newSettlementMainCreateCmd(),
		newSettlementMainUpdateCmd(),
		newSettlementMainDeleteCmd(),
		newSettlementMainUpdateStatusCmd(),
		newSettlementMainUnsettleWaybillCmd(),
		newSettlementMainUnsettleWaybillCountCmd(),
		newSettlementMainExportExcelCmd(),
	)
}

// ---- 结算 list / page-count 共用筛选字段 ----
//
// 后端 /erp/settlement/page 完整筛选字段：
//   no                        结算单号
//   name                      结算单名称
//   customerId                客户 ID            (销售结算)
//   supplierId                供应商 ID          (采购结算)
//   userId                    业务员 ID          (采购结算)
//   enterpriseId              企业 ID
//   projectId                 项目 ID
//   projectName               项目名称
//   settleStatus              结算状态 (如 PART_SETTLED；多选用逗号分隔)
//   invoiceStatus             发票状态 (如 PART_INVOICED；多选用逗号分隔)
//   type                      结算类型 (SALE_SETTLEMENT / PURCHASE_SETTLEMENT)
//   settleType                结算方式 (SALE / PURCHASE)
//   metricsName               检测指标名 / 规格指标
//   settlementDate[0]/[1]     结算日期范围 (数组参数，由 --start-date/--end-date 折叠)
//
// 注: 原 CLI 的 --settlement-no 发出的后端字段 "settlementNo" 与后端实际字段 "no" 不符，
// 现已修正为 --no (mapping 自 "no")，其余字段为本次补齐。

func buildSettlementListParams(cmd *cobra.Command, pageNo, pageSize int, withPage bool) map[string]any {
	params := map[string]any{}
	if withPage {
		params["pageNo"] = pageNo
		params["pageSize"] = pageSize
	}
	cmdutil.CollectStringFlags(cmd, params,
		"no", "name", "customer-id", "supplier-id", "user-id",
		"enterprise-id", "project-id", "project-name",
		"settle-status", "invoice-status", "type", "settle-type",
		"metrics-name")
	cmdutil.CollectSettlementDate(cmd, params)
	return params
}

func registerSettlementListFlags(c *cobra.Command) {
	c.Flags().String("no", "", "结算单号")
	c.Flags().String("name", "", "结算单名称")
	c.Flags().String("customer-id", "", "客户 ID")
	c.Flags().String("supplier-id", "", "供应商 ID")
	c.Flags().String("user-id", "", "业务员 ID")
	c.Flags().String("enterprise-id", "", "企业 ID")
	c.Flags().String("project-id", "", "项目 ID")
	c.Flags().String("project-name", "", "项目名称")
	c.Flags().String("settle-status", "", "结算状态 (如 PART_SETTLED；多选用逗号分隔)")
	c.Flags().String("invoice-status", "", "发票状态 (如 PART_INVOICED；多选用逗号分隔)")
	c.Flags().String("type", "", "结算类型 (SALE_SETTLEMENT / PURCHASE_SETTLEMENT)")
	c.Flags().String("settle-type", "", "结算方式 (SALE / PURCHASE)")
	c.Flags().String("metrics-name", "", "检测指标名 / 规格指标")
	cmdutil.AddSettlementDateFlags(c)
}

// ---- 分页查询 ----

func newSettlementMainListCmd() *cobra.Command {
	var pageNo, pageSize int

	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询结算单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := buildSettlementListParams(cmd, pageNo, pageSize, true)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/settlement/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询结算单失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	registerSettlementListFlags(c)
	return c
}

// ---- 分页计数 ----

func newSettlementMainPageCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "page-count",
		Short: "统计结算单数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := buildSettlementListParams(cmd, 0, 0, false)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/settlement/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计结算单数量失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	registerSettlementListFlags(c)
	return c
}

// ---- 获取详情 ----

func newSettlementMainGetCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "get",
		Short: "获取结算单详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/settlement/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取结算单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "结算单 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 创建 ----

func newSettlementMainCreateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "create",
		Short: "创建结算单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/settlement/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建结算单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 更新 ----

func newSettlementMainUpdateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update",
		Short: "更新结算单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/settlement/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新结算单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 删除 ----

func newSettlementMainDeleteCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "delete",
		Short: "删除结算单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/settlement/delete", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("删除结算单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "结算单 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 更新状态 ----

func newSettlementMainUpdateStatusCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update-status",
		Short: "更新结算单状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/settlement/update-status", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新结算单状态失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据（含 id 和 status）")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 未结算运单 ----

func newSettlementMainUnsettleWaybillCmd() *cobra.Command {
	var customerId, supplierId, warehouseId string

	c := &cobra.Command{
		Use:   "unsettle-waybill",
		Short: "获取未结算运单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "customer-id", "supplier-id", "warehouse-id")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/settlement/unsettle/waybill", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取未结算运单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&customerId, "customer-id", "", "客户 ID")
	c.Flags().StringVar(&supplierId, "supplier-id", "", "供应商 ID")
	c.Flags().StringVar(&warehouseId, "warehouse-id", "", "仓库 ID")
	return c
}

// ---- 未结算运单计数 ----

func newSettlementMainUnsettleWaybillCountCmd() *cobra.Command {
	var customerId, supplierId, warehouseId string

	c := &cobra.Command{
		Use:   "unsettle-waybill-count",
		Short: "统计未结算运单数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "customer-id", "supplier-id", "warehouse-id")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/settlement/unsettle/waybill/count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计未结算运单数量失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&customerId, "customer-id", "", "客户 ID")
	c.Flags().StringVar(&supplierId, "supplier-id", "", "供应商 ID")
	c.Flags().StringVar(&warehouseId, "warehouse-id", "", "仓库 ID")
	return c
}

// ---- 导出 Excel ----

func newSettlementMainExportExcelCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "export",
		Short: "导出结算单 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := buildSettlementListParams(cmd, 0, 0, false)
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/settlement/export-excel", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出结算单失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	registerSettlementListFlags(c)
	return c
}
