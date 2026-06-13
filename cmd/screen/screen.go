package screen

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

var screenCmd = &cobra.Command{
	Use:   "screen",
	Short: "大屏数据管理",
	Long:  "大屏数据：采购销售、库存数量、金额使用、项目数量。",
}

func init() {
	screenCmd.AddCommand(
		newScreenPurchaseSaleCmd(),
		newScreenStockCountCmd(),
		newScreenAmountUsedCmd(),
		newScreenProjectCountCmd(),
	)
}

// ---- 采购销售数据 ----

func newScreenPurchaseSaleCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "purchase-sale",
		Short: "获取采购销售数据",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/screen/purchase-sale", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取采购销售数据失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

// ---- 库存数量 ----

func newScreenStockCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "stock-count",
		Short: "获取库存数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/screen/stock-count", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取库存数量失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

// ---- 金额使用 ----

func newScreenAmountUsedCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "amount-used",
		Short: "获取金额使用情况",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/screen/amount-used", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取金额使用情况失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

// ---- 项目数量 ----

func newScreenProjectCountCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "project-count",
		Short: "获取项目数量",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/erp/screen/project-count", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取项目数量失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

// Register adds the screen command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(screenCmd)
}
