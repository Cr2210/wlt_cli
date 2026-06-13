package invoice

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

func init() {
	invoiceCmd.AddCommand(invoiceMainCmd)
}

var invoiceMainCmd = &cobra.Command{
	Use:   "main",
	Short: "发票管理",
}

func init() {
	invoiceMainCmd.AddCommand(
		newInvoiceMainListCmd(),
		newInvoiceMainPageCountCmd(),
		newInvoiceMainGetCmd(),
		newInvoiceMainCreateCmd(),
		newInvoiceMainUpdateCmd(),
		newInvoiceMainDeleteCmd(),
		newInvoiceMainUpdateStatusCmd(),
		newInvoiceMainExportExcelCmd(),
	)
}

// ---- 分页查询 ----

func newInvoiceMainListCmd() *cobra.Command {
	var pageNo, pageSize int
	var invoiceNo, customerId, supplierId, status, type_ string

	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询发票",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlags(cmd, params, "invoice-no", "customer-id", "supplier-id", "status", "type")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/invoice/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询发票失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().StringVar(&invoiceNo, "invoice-no", "", "发票号")
	c.Flags().StringVar(&customerId, "customer-id", "", "客户 ID")
	c.Flags().StringVar(&supplierId, "supplier-id", "", "供应商 ID")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&type_, "type", "", "发票类型")
	return c
}

// ---- 分页计数 ----

func newInvoiceMainPageCountCmd() *cobra.Command {
	var invoiceNo, customerId, supplierId, status, type_ string

	c := &cobra.Command{
		Use:   "page-count",
		Short: "统计发票数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "invoice-no", "customer-id", "supplier-id", "status", "type")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/invoice/page-count", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("统计发票数量失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&invoiceNo, "invoice-no", "", "发票号")
	c.Flags().StringVar(&customerId, "customer-id", "", "客户 ID")
	c.Flags().StringVar(&supplierId, "supplier-id", "", "供应商 ID")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&type_, "type", "", "发票类型")
	return c
}

// ---- 获取详情 ----

func newInvoiceMainGetCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "get",
		Short: "获取发票详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/invoice/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取发票失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "发票 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 创建 ----

func newInvoiceMainCreateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "create",
		Short: "创建发票",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/invoice/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建发票失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 更新 ----

func newInvoiceMainUpdateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update",
		Short: "更新发票",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/invoice/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新发票失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 删除 ----

func newInvoiceMainDeleteCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "delete",
		Short: "删除发票",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/erp/invoice/delete", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("删除发票失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "发票 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 更新状态 ----

func newInvoiceMainUpdateStatusCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update-status",
		Short: "更新发票状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/erp/invoice/update-status", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新发票状态失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据（含 id 和 status）")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 导出 Excel ----

func newInvoiceMainExportExcelCmd() *cobra.Command {
	var invoiceNo, customerId, supplierId, status, type_ string

	c := &cobra.Command{
		Use:   "export",
		Short: "导出发票 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "invoice-no", "customer-id", "supplier-id", "status", "type")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/invoice/export-excel", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出发票失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	c.Flags().StringVar(&invoiceNo, "invoice-no", "", "发票号")
	c.Flags().StringVar(&customerId, "customer-id", "", "客户 ID")
	c.Flags().StringVar(&supplierId, "supplier-id", "", "供应商 ID")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&type_, "type", "", "发票类型")
	return c
}
