package weight

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var weightWaybillCmd = &cobra.Command{
	Use:   "waybill",
	Short: "运单称重管理",
}

func init() {
	weightCmd.AddCommand(weightWaybillCmd)
	weightWaybillCmd.AddCommand(
		newWeightWaybillPageCmd(),
		newWeightWaybillImageCmd(),
		newWeightWaybillCreateWaybillSourceCmd(),
		newWeightWaybillBatchCreateWaybillSourceCmd(),
		newWeightWaybillMatchLinkWaybillSourceCmd(),
		newWeightWaybillUnlinkWaybillSourceCmd(),
		newWeightWaybillExportExcelCmd(),
	)
}

// ---- 分页查询 ----

func newWeightWaybillPageCmd() *cobra.Command {
	var pageNo, pageSize int
	var waybillNo, licensePlate, weighingNo string

	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询运单称重",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlags(cmd, params, "waybill-no", "license-plate", "weighing-no")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/weight/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询运单称重失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().StringVar(&waybillNo, "waybill-no", "", "运单号")
	c.Flags().StringVar(&licensePlate, "license-plate", "", "车牌号")
	c.Flags().StringVar(&weighingNo, "weighing-no", "", "称重单号")
	return c
}

// ---- 获取图片 ----

func newWeightWaybillImageCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "image",
		Short: "获取称重图片",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/weight/image", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取称重图片失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "称重 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 创建运单称重 ----

func newWeightWaybillCreateWaybillSourceCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "create-waybill-source",
		Short: "创建运单称重",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/weight/createWaybillSource", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建运单称重失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 批量创建运单称重 ----

func newWeightWaybillBatchCreateWaybillSourceCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "batch-create-waybill-source",
		Short: "批量创建运单称重",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/weight/batchCreateWaybillSource", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("批量创建运单称重失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 匹配关联运单 ----

func newWeightWaybillMatchLinkWaybillSourceCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "match-link-waybill-source",
		Short: "匹配关联运单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/weight/matchLinkWaybillSource", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("匹配关联运单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 取消关联运单 ----

func newWeightWaybillUnlinkWaybillSourceCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "unlink-waybill-source",
		Short: "取消关联运单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/erp/weight/unlinkWaybillSource", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("取消关联运单失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 导出 Excel ----

func newWeightWaybillExportExcelCmd() *cobra.Command {
	var waybillNo, licensePlate, weighingNo string

	c := &cobra.Command{
		Use:   "export",
		Short: "导出运单称重 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "waybill-no", "license-plate", "weighing-no")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/weight/export-excel", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出运单称重失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	c.Flags().StringVar(&waybillNo, "waybill-no", "", "运单号")
	c.Flags().StringVar(&licensePlate, "license-plate", "", "车牌号")
	c.Flags().StringVar(&weighingNo, "weighing-no", "", "称重单号")
	return c
}
