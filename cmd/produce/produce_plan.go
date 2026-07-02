package produce

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var producePlanCmd = &cobra.Command{
	Use:   "plan",
	Short: "生产计划管理",
}

func init() {
	produceCmd.AddCommand(producePlanCmd)
	producePlanCmd.AddCommand(
		newProducePlanPageCmd(),
		newProducePlanPageCountCmd(),
		newProducePlanSimpleListCmd(),
		newProducePlanGetCmd(),
		newProducePlanCreateCmd(),
		newProducePlanUpdateCmd(),
		newProducePlanDeleteCmd(),
		newProducePlanUpdateStatusCmd(),
		newProducePlanExportExcelCmd(),
	)
}

// ---- 分页查询 ----

func newProducePlanPageCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "page",
		Short: "分页查询生产方案",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlags(cmd, params,
				"plan-no",
				"plan-name",
				"status",
				"product-id",
				"product-name",
				"metrics-name",
				"plan-time",
				"customer-id",
				"customer-name",
				"remark",
				"creator",
				"creator-name",
				"create-time",
				"update-time",
				"updater-name",
				"custom-order",
				"keyword",
				"headers",
			)

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/produce-plan/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询生产方案失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().String("plan-no", "", "方案编号")
	c.Flags().String("plan-name", "", "方案名称")
	c.Flags().String("status", "", "方案状态")
	c.Flags().String("product-id", "", "产品编号")
	c.Flags().String("product-name", "", "产品名称")
	c.Flags().String("metrics-name", "", "产品指标")
	c.Flags().String("plan-time", "", "方案日期")
	c.Flags().String("customer-id", "", "客户编号")
	c.Flags().String("customer-name", "", "客户名称")
	c.Flags().String("remark", "", "备注")
	c.Flags().String("creator", "", "创建者")
	c.Flags().String("creator-name", "", "创建人名称")
	c.Flags().String("create-time", "", "创建时间")
	c.Flags().String("update-time", "", "更新时间")
	c.Flags().String("updater-name", "", "更新人名称")
	c.Flags().String("custom-order", "", "前端自定义排序规则")
	c.Flags().String("keyword", "", "关键字")
	c.Flags().String("headers", "", "自定义导出表头")
	return c
}

// ---- 分页计数 ----

func newProducePlanPageCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "page-count",
		Short: "按筛选统计生产方案数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params,
				"plan-no",
				"plan-name",
				"status",
				"product-id",
				"product-name",
				"metrics-name",
				"plan-time",
				"customer-id",
				"customer-name",
				"remark",
				"creator",
				"creator-name",
				"create-time",
				"update-time",
				"updater-name",
				"custom-order",
				"keyword",
			)

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/produce-plan/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计生产方案失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().String("plan-no", "", "方案编号")
	c.Flags().String("plan-name", "", "方案名称")
	c.Flags().String("status", "", "方案状态")
	c.Flags().String("product-id", "", "产品编号")
	c.Flags().String("product-name", "", "产品名称")
	c.Flags().String("metrics-name", "", "产品指标")
	c.Flags().String("plan-time", "", "方案日期")
	c.Flags().String("customer-id", "", "客户编号")
	c.Flags().String("customer-name", "", "客户名称")
	c.Flags().String("remark", "", "备注")
	c.Flags().String("creator", "", "创建者")
	c.Flags().String("creator-name", "", "创建人名称")
	c.Flags().String("create-time", "", "创建时间")
	c.Flags().String("update-time", "", "更新时间")
	c.Flags().String("updater-name", "", "更新人名称")
	c.Flags().String("custom-order", "", "前端自定义排序规则")
	c.Flags().String("keyword", "", "关键字")
	return c
}

// ---- 简单列表 ----

func newProducePlanSimpleListCmd() *cobra.Command {
	var warehouseId string

	c := &cobra.Command{
		Use:   "simple-list",
		Short: "获取生产计划精简列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			if warehouseId != "" {
				params["warehouseId"] = warehouseId
			}

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/produce-plan/simple-list", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取生产计划列表失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&warehouseId, "warehouse-id", "", "仓库 ID")
	return c
}

// ---- 获取详情 ----

func newProducePlanGetCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "get",
		Short: "获取生产计划详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/produce-plan/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取生产计划失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "计划 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 创建 ----

func newProducePlanCreateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "create",
		Short: "创建生产计划",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/produce-plan/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建生产计划失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 更新 ----

func newProducePlanUpdateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update",
		Short: "更新生产计划",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/produce-plan/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新生产计划失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 删除 ----

func newProducePlanDeleteCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "delete",
		Short: "删除生产计划",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/produce-plan/delete", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("删除生产计划失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "计划 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 更新状态 ----

func newProducePlanUpdateStatusCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update-status",
		Short: "更新生产计划状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/produce-plan/update-status", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新生产计划状态失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据（含 id 和 status）")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 导出 Excel ----

func newProducePlanExportExcelCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "export",
		Short: "导出生产方案 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params,
				"plan-no",
				"plan-name",
				"status",
				"product-id",
				"product-name",
				"metrics-name",
				"plan-time",
				"customer-id",
				"customer-name",
				"remark",
				"creator",
				"creator-name",
				"create-time",
				"update-time",
				"updater-name",
				"custom-order",
				"keyword",
				"headers",
			)

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/produce-plan/export-excel", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出生产方案失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	c.Flags().String("plan-no", "", "方案编号")
	c.Flags().String("plan-name", "", "方案名称")
	c.Flags().String("status", "", "方案状态")
	c.Flags().String("product-id", "", "产品编号")
	c.Flags().String("product-name", "", "产品名称")
	c.Flags().String("metrics-name", "", "产品指标")
	c.Flags().String("plan-time", "", "方案日期")
	c.Flags().String("customer-id", "", "客户编号")
	c.Flags().String("customer-name", "", "客户名称")
	c.Flags().String("remark", "", "备注")
	c.Flags().String("creator", "", "创建者")
	c.Flags().String("creator-name", "", "创建人名称")
	c.Flags().String("create-time", "", "创建时间")
	c.Flags().String("update-time", "", "更新时间")
	c.Flags().String("updater-name", "", "更新人名称")
	c.Flags().String("custom-order", "", "前端自定义排序规则")
	c.Flags().String("keyword", "", "关键字")
	c.Flags().String("headers", "", "自定义导出表头")
	return c
}
