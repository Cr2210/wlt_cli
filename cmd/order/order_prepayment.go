package order

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

// orderPrepaymentRelationAPIPath 是订单预付关联的后端端点根路径。
const orderPrepaymentRelationAPIPath = "/erp/order-prepayment-relation"

var orderPrepaymentRelationCmd = &cobra.Command{
	Use:   "prepayment-relation",
	Short: "订单预付关联管理",
}

func init() {
	orderCmd.AddCommand(orderPrepaymentRelationCmd)
	orderPrepaymentRelationCmd.AddCommand(
		newOrderPrepayRelationListCmd(),
		newOrderPrepayRelationCreateCmd(),
		newOrderPrepayRelationUpdateAmountCmd(),
	)
}

// ---- 分页查询关联记录 ----

func newOrderPrepayRelationListCmd() *cobra.Command {
	var orderId int64
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询预付关联",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"orderId":  orderId,
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlag(cmd, params, "relation-type")
			cmdutil.CollectStringFlag(cmd, params, "headers")
			resp, err := cmdutil.GetClient().Get(context.Background(), orderPrepaymentRelationAPIPath+"/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询预付关联失败: %s", err), "")
			}
			// 该接口响应 data 为裸数组（无 list/total 包裹），直接输出。
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&orderId, "order-id", 0, "订单 ID")
	_ = c.MarkFlagRequired("order-id")
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().String("relation-type", "", "关联类型：PAYMENT-付款单，SUPPLIER-供应商期初")
	c.Flags().String("headers", "", "自定义导出表头")
	return c
}

// ---- 创建关联 ----

func newOrderPrepayRelationCreateCmd() *cobra.Command {
	var orderId, relationId int64
	var relationType string
	var relationAmount float64
	c := &cobra.Command{
		Use:   "create",
		Short: "创建预付关联",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body := map[string]any{
				"orderId":        orderId,
				"relationType":   relationType,
				"relationId":     relationId,
				"relationAmount": relationAmount,
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), orderPrepaymentRelationAPIPath+"/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建预付关联失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&orderId, "order-id", 0, "订单 ID")
	c.Flags().StringVar(&relationType, "relation-type", "", "关联类型：PAYMENT-付款单，SUPPLIER-供应商期初")
	c.Flags().Int64Var(&relationId, "relation-id", 0, "关联 ID（付款单 ID 或供应商期初 ID）")
	c.Flags().Float64Var(&relationAmount, "relation-amount", 0, "关联预付金额（单位：元）")
	_ = c.MarkFlagRequired("order-id")
	_ = c.MarkFlagRequired("relation-type")
	_ = c.MarkFlagRequired("relation-id")
	_ = c.MarkFlagRequired("relation-amount")
	return c
}

// ---- 更新关联金额 ----

func newOrderPrepayRelationUpdateAmountCmd() *cobra.Command {
	var id int64
	var relationAmount float64
	c := &cobra.Command{
		Use:   "update-amount",
		Short: "更新关联金额",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body := map[string]any{
				"id":             id,
				"relationAmount": relationAmount,
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), orderPrepaymentRelationAPIPath+"/update-amount", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新关联金额失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "关联记录 ID")
	c.Flags().Float64Var(&relationAmount, "relation-amount", 0, "关联预付金额（单位：元）")
	_ = c.MarkFlagRequired("id")
	_ = c.MarkFlagRequired("relation-amount")
	return c
}
