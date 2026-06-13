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
		newProduceMainListCmd(),
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

func newProduceMainListCmd() *cobra.Command {
	var pageNo, pageSize int
	var produceNo, status, warehouseId, productId string

	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询生产单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlags(cmd, params, "produce-no", "status", "warehouse-id", "product-id")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/produce/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询生产单失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().StringVar(&produceNo, "produce-no", "", "生产单号")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&warehouseId, "warehouse-id", "", "仓库 ID")
	c.Flags().StringVar(&productId, "product-id", "", "产品 ID")
	return c
}

// ---- 分页计数 ----

func newProduceMainPageCountCmd() *cobra.Command {
	var produceNo, status, warehouseId, productId string

	c := &cobra.Command{
		Use:   "page-count",
		Short: "统计生产单数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "produce-no", "status", "warehouse-id", "product-id")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/produce/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计生产单数量失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&produceNo, "produce-no", "", "生产单号")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&warehouseId, "warehouse-id", "", "仓库 ID")
	c.Flags().StringVar(&productId, "product-id", "", "产品 ID")
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
	var produceNo, status, warehouseId, productId string

	c := &cobra.Command{
		Use:   "export",
		Short: "导出生产单 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "produce-no", "status", "warehouse-id", "product-id")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/produce/export-excel", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出生产单失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	c.Flags().StringVar(&produceNo, "produce-no", "", "生产单号")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&warehouseId, "warehouse-id", "", "仓库 ID")
	c.Flags().StringVar(&productId, "product-id", "", "产品 ID")
	return c
}

// ---- 质检分页 ----

func newProduceMainQualityPageCmd() *cobra.Command {
	var pageNo, pageSize int
	var produceId string

	c := &cobra.Command{
		Use:   "quality-page",
		Short: "分页查询生产单质检",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			if produceId != "" {
				params["produceId"] = produceId
			}

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/produce/quality-page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询生产单质检失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().StringVar(&produceId, "produce-id", "", "生产单 ID")
	return c
}

// ---- 质检计数 ----

func newProduceMainQualityCountCmd() *cobra.Command {
	var produceId string

	c := &cobra.Command{
		Use:   "quality-count",
		Short: "统计生产单质检数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			if produceId != "" {
				params["produceId"] = produceId
			}

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/produce/quality-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计质检数量失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&produceId, "produce-id", "", "生产单 ID")
	return c
}

// ---- 导出质检 Excel ----

func newProduceMainQualityExportExcelCmd() *cobra.Command {
	var produceId string

	c := &cobra.Command{
		Use:   "quality-export",
		Short: "导出生产单质检 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			if produceId != "" {
				params["produceId"] = produceId
			}

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/produce/quality-export-excel", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出质检失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	c.Flags().StringVar(&produceId, "produce-id", "", "生产单 ID")
	return c
}
