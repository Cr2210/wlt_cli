package produce

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var produceMainCmd = &cobra.Command{
	Use:   "main",
	Short: "生产单管理",
}

func init() {
	produceCmd.AddCommand(produceMainCmd)
	produceMainCmd.AddCommand(
		newProduceMainPageCmd(),
		newProduceMainPageCountCmd(),
		newProduceMainGetCmd(),
		newProduceMainCreateCmd(),
		newProduceMainUpdateCmd(),
		newProduceMainDeleteCmd(),
		newProduceMainUpdateStatusCmd(),
		newProduceMainExportExcelCmd(),
		newProduceMainQualityPageCmd(),
		newProduceMainQualityCountCmd(),
		newProduceMainQualityExportExcelCmd(),
	)
}

// ---- 分页查询 ----

func newProduceMainPageCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "page",
		Short: "分页查询生产任务",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlags(cmd, params,
				"no",
				"project-name",
				"status",
				"produce-time",
				"batch-no",
				"metrics-name",
				"warehouse-id",
				"warehouse-name",
				"product-id",
				"product-name",
				"plan-no",
				"plan-name",
				"user-id",
				"user-name",
				"order-id",
				"order-no",
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

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/produce/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询生产任务失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().String("no", "", "生产编号")
	c.Flags().String("project-name", "", "项目名称")
	c.Flags().String("status", "", "出库状态")
	c.Flags().String("produce-time", "", "生产日期")
	c.Flags().String("batch-no", "", "批次号")
	c.Flags().String("metrics-name", "", "指标名称")
	c.Flags().String("warehouse-id", "", "仓库编号")
	c.Flags().String("warehouse-name", "", "仓库名称")
	c.Flags().String("product-id", "", "产品编号")
	c.Flags().String("product-name", "", "产品名称")
	c.Flags().String("plan-no", "", "方案号")
	c.Flags().String("plan-name", "", "方案名称")
	c.Flags().String("user-id", "", "操作人员")
	c.Flags().String("user-name", "", "操作人员名称")
	c.Flags().String("order-id", "", "订单编号")
	c.Flags().String("order-no", "", "订单号")
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

func newProduceMainPageCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "page-count",
		Short: "按筛选统计生产任务数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params,
				"no",
				"project-name",
				"status",
				"produce-time",
				"batch-no",
				"metrics-name",
				"warehouse-id",
				"warehouse-name",
				"product-id",
				"product-name",
				"plan-no",
				"plan-name",
				"user-id",
				"user-name",
				"order-id",
				"order-no",
				"remark",
				"creator",
				"creator-name",
				"create-time",
				"update-time",
				"updater-name",
				"custom-order",
				"keyword",
			)

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/produce/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计生产任务失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().String("no", "", "生产编号")
	c.Flags().String("project-name", "", "项目名称")
	c.Flags().String("status", "", "出库状态")
	c.Flags().String("produce-time", "", "生产日期")
	c.Flags().String("batch-no", "", "批次号")
	c.Flags().String("metrics-name", "", "指标名称")
	c.Flags().String("warehouse-id", "", "仓库编号")
	c.Flags().String("warehouse-name", "", "仓库名称")
	c.Flags().String("product-id", "", "产品编号")
	c.Flags().String("product-name", "", "产品名称")
	c.Flags().String("plan-no", "", "方案号")
	c.Flags().String("plan-name", "", "方案名称")
	c.Flags().String("user-id", "", "操作人员")
	c.Flags().String("user-name", "", "操作人员名称")
	c.Flags().String("order-id", "", "订单编号")
	c.Flags().String("order-no", "", "订单号")
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

// ---- 获取详情 ----

func newProduceMainGetCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "get",
		Short: "获取生产单详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/produce/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取生产单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "生产单 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 创建 ----

func newProduceMainCreateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "create",
		Short: "创建生产单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/produce/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建生产单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 更新 ----

func newProduceMainUpdateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update",
		Short: "更新生产单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/produce/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新生产单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 删除 ----

func newProduceMainDeleteCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "delete",
		Short: "删除生产单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/produce/delete", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("删除生产单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "生产单 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 更新状态 ----

func newProduceMainUpdateStatusCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update-status",
		Short: "更新生产单状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/produce/update-status", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新生产单状态失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据（含 id 和 status）")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 导出 Excel ----

func newProduceMainExportExcelCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "export",
		Short: "导出生产任务 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params,
				"no",
				"project-name",
				"status",
				"produce-time",
				"batch-no",
				"metrics-name",
				"warehouse-id",
				"warehouse-name",
				"product-id",
				"product-name",
				"plan-no",
				"plan-name",
				"user-id",
				"user-name",
				"order-id",
				"order-no",
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

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/produce/export-excel", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出生产任务失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	c.Flags().String("no", "", "生产编号")
	c.Flags().String("project-name", "", "项目名称")
	c.Flags().String("status", "", "出库状态")
	c.Flags().String("produce-time", "", "生产日期")
	c.Flags().String("batch-no", "", "批次号")
	c.Flags().String("metrics-name", "", "指标名称")
	c.Flags().String("warehouse-id", "", "仓库编号")
	c.Flags().String("warehouse-name", "", "仓库名称")
	c.Flags().String("product-id", "", "产品编号")
	c.Flags().String("product-name", "", "产品名称")
	c.Flags().String("plan-no", "", "方案号")
	c.Flags().String("plan-name", "", "方案名称")
	c.Flags().String("user-id", "", "操作人员")
	c.Flags().String("user-name", "", "操作人员名称")
	c.Flags().String("order-id", "", "订单编号")
	c.Flags().String("order-no", "", "订单号")
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

// ---- 质检分页 ----

func newProduceMainQualityPageCmd() *cobra.Command {
	var pageNo, pageSize int
	c := &cobra.Command{
		Use:   "quality-page",
		Short: "分页查询生产质检",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlags(cmd, params,
				"no",
				"produce-time",
				"product-name",
				"metrics-name",
				"inspection-nos-str",
				"inspection-result",
				"warehouse-name",
				"plan-name",
				"user-name",
				"project-name",
				"plan-no",
				"order-no",
				"batch-no",
				"remark",
				"status",
				"creator",
				"create-time",
				"updater",
				"update-time",
				"keyword",
				"headers",
			)

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/produce/quality-page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询生产质检失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().String("no", "", "生产任务编号")
	c.Flags().String("produce-time", "", "生产日期")
	c.Flags().String("product-name", "", "货物名称")
	c.Flags().String("metrics-name", "", "规格指标")
	c.Flags().String("inspection-nos-str", "", "关联质检单编号")
	c.Flags().String("inspection-result", "", "质检结果")
	c.Flags().String("warehouse-name", "", "仓库名称")
	c.Flags().String("plan-name", "", "方案名称")
	c.Flags().String("user-name", "", "业务负责人")
	c.Flags().String("project-name", "", "项目名称")
	c.Flags().String("plan-no", "", "关联方案编号")
	c.Flags().String("order-no", "", "关联销售订单")
	c.Flags().String("batch-no", "", "批次编号")
	c.Flags().String("remark", "", "备注")
	c.Flags().String("status", "", "审核状态")
	c.Flags().String("creator", "", "创建人")
	c.Flags().String("create-time", "", "创建时间")
	c.Flags().String("updater", "", "更新人")
	c.Flags().String("update-time", "", "更新时间")
	c.Flags().String("keyword", "", "关键字")
	c.Flags().String("headers", "", "自定义导出表头")
	return c
}

// ---- 质检计数 ----

func newProduceMainQualityCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "quality-count",
		Short: "按筛选统计生产质检数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params,
				"no",
				"produce-time",
				"product-name",
				"metrics-name",
				"inspection-nos-str",
				"inspection-result",
				"warehouse-name",
				"plan-name",
				"user-name",
				"project-name",
				"plan-no",
				"order-no",
				"batch-no",
				"remark",
				"status",
				"creator",
				"create-time",
				"updater",
				"update-time",
				"keyword",
			)

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/produce/quality-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计生产质检失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().String("no", "", "生产任务编号")
	c.Flags().String("produce-time", "", "生产日期")
	c.Flags().String("product-name", "", "货物名称")
	c.Flags().String("metrics-name", "", "规格指标")
	c.Flags().String("inspection-nos-str", "", "关联质检单编号")
	c.Flags().String("inspection-result", "", "质检结果")
	c.Flags().String("warehouse-name", "", "仓库名称")
	c.Flags().String("plan-name", "", "方案名称")
	c.Flags().String("user-name", "", "业务负责人")
	c.Flags().String("project-name", "", "项目名称")
	c.Flags().String("plan-no", "", "关联方案编号")
	c.Flags().String("order-no", "", "关联销售订单")
	c.Flags().String("batch-no", "", "批次编号")
	c.Flags().String("remark", "", "备注")
	c.Flags().String("status", "", "审核状态")
	c.Flags().String("creator", "", "创建人")
	c.Flags().String("create-time", "", "创建时间")
	c.Flags().String("updater", "", "更新人")
	c.Flags().String("update-time", "", "更新时间")
	c.Flags().String("keyword", "", "关键字")
	return c
}

// ---- 导出质检 Excel ----

func newProduceMainQualityExportExcelCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "quality-export",
		Short: "导出生产质检 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params,
				"no",
				"produce-time",
				"product-name",
				"metrics-name",
				"inspection-nos-str",
				"inspection-result",
				"warehouse-name",
				"plan-name",
				"user-name",
				"project-name",
				"plan-no",
				"order-no",
				"batch-no",
				"remark",
				"status",
				"creator",
				"create-time",
				"updater",
				"update-time",
				"keyword",
				"headers",
			)

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/produce/quality-export-excel", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出生产质检失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	c.Flags().String("no", "", "生产任务编号")
	c.Flags().String("produce-time", "", "生产日期")
	c.Flags().String("product-name", "", "货物名称")
	c.Flags().String("metrics-name", "", "规格指标")
	c.Flags().String("inspection-nos-str", "", "关联质检单编号")
	c.Flags().String("inspection-result", "", "质检结果")
	c.Flags().String("warehouse-name", "", "仓库名称")
	c.Flags().String("plan-name", "", "方案名称")
	c.Flags().String("user-name", "", "业务负责人")
	c.Flags().String("project-name", "", "项目名称")
	c.Flags().String("plan-no", "", "关联方案编号")
	c.Flags().String("order-no", "", "关联销售订单")
	c.Flags().String("batch-no", "", "批次编号")
	c.Flags().String("remark", "", "备注")
	c.Flags().String("status", "", "审核状态")
	c.Flags().String("creator", "", "创建人")
	c.Flags().String("create-time", "", "创建时间")
	c.Flags().String("updater", "", "更新人")
	c.Flags().String("update-time", "", "更新时间")
	c.Flags().String("keyword", "", "关键字")
	c.Flags().String("headers", "", "自定义导出表头")
	return c
}
