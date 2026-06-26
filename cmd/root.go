package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	configcmd "github.com/weiliantong/cli/cmd/config"
	"github.com/weiliantong/cli/cmd/contract"
	"github.com/weiliantong/cli/cmd/customer"
	"github.com/weiliantong/cli/cmd/data_sync"
	"github.com/weiliantong/cli/cmd/finance"
	"github.com/weiliantong/cli/cmd/homepage"
	"github.com/weiliantong/cli/cmd/invoice"
	"github.com/weiliantong/cli/cmd/job_trigger"
	"github.com/weiliantong/cli/cmd/operate_log"
	"github.com/weiliantong/cli/cmd/order"
	"github.com/weiliantong/cli/cmd/produce"
	"github.com/weiliantong/cli/cmd/product"
	"github.com/weiliantong/cli/cmd/profit"
	"github.com/weiliantong/cli/cmd/purchase"
	"github.com/weiliantong/cli/cmd/quality"
	"github.com/weiliantong/cli/cmd/report"
	"github.com/weiliantong/cli/cmd/sale"
	"github.com/weiliantong/cli/cmd/screen"
	"github.com/weiliantong/cli/cmd/settlement"
	"github.com/weiliantong/cli/cmd/stats"
	"github.com/weiliantong/cli/cmd/stock"
	"github.com/weiliantong/cli/cmd/supplier"
	"github.com/weiliantong/cli/cmd/system"
	"github.com/weiliantong/cli/cmd/waybill"
	"github.com/weiliantong/cli/cmd/weight"
	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/config"
	"github.com/weiliantong/cli/internal/output"
)

var rootCmd = &cobra.Command{
	Use:   "wlt",
	Short: "维链通 ERP 命令行工具",
	Long:  "weiliantong-cli (wlt) — 维链通 ERP 系统的命令行工具，AI Agent 原生设计。",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		profileName, _ := cmd.Flags().GetString("profile")
		cfgMgr := config.NewManager()
		cfgMgr.SetProfileOverride(profileName)
		if err := cfgMgr.Load(); err != nil {
			return output.NewExitError(2, fmt.Sprintf("加载配置失败: %s", err), "运行 wlt config init 初始化配置")
		}
		cmdutil.InitManagers(cfgMgr)
		// Stateless auth: token + tenant-id (and optional base-url override)
		// are supplied per-call via flags. Validated lazily in EnsureClient so
		// commands that don't hit the API (version/config/...) stay usable.
		token, _ := cmd.Flags().GetString("token")
		tenantID, _ := cmd.Flags().GetString("tenant-id")
		baseURL, _ := cmd.Flags().GetString("base-url")
		cmdutil.SetAuthFlags(token, tenantID, baseURL)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		// Structured error output
		var exitErr *output.ExitError
		if as, ok := err.(*output.ExitError); ok {
			exitErr = as
		} else {
			exitErr = &output.ExitError{
				Code:    1,
				Type:    "general",
				Message: err.Error(),
			}
		}
		_ = output.WriteError(os.Stderr, exitErr)
		os.Exit(exitErr.Code)
	}
	return nil
}

func init() {
	rootCmd.PersistentFlags().String("profile", "sit", "配置环境(sit/prod),提供 base_url/api_prefix")
	rootCmd.PersistentFlags().Bool("quiet", false, "静默模式")
	// 无状态鉴权:每次调用由调用方传入
	rootCmd.PersistentFlags().String("token", "", "访问令牌(对应 Authorization: Bearer 头,必填)")
	rootCmd.PersistentFlags().String("tenant-id", "", "租户 ID(对应 tenant-id 头,必填)")
	rootCmd.PersistentFlags().String("base-url", "", "可选,覆盖 profile 的 base_url")

	// 基础设施
	configcmd.Register(rootCmd)
	system.Register(rootCmd)

	// 业务域
	contract.Register(rootCmd)
	customer.Register(rootCmd)
	datasync.Register(rootCmd)
	finance.Register(rootCmd)
	homepage.Register(rootCmd)
	invoice.Register(rootCmd)
	jobtrigger.Register(rootCmd)
	operatelog.Register(rootCmd)
	order.Register(rootCmd)
	produce.Register(rootCmd)
	product.Register(rootCmd)
	profit.Register(rootCmd)
	purchase.Register(rootCmd)
	quality.Register(rootCmd)
	report.Register(rootCmd)
	sale.Register(rootCmd)
	screen.Register(rootCmd)
	settlement.Register(rootCmd)
	stats.Register(rootCmd)
	stock.Register(rootCmd)
	supplier.Register(rootCmd)
	waybill.Register(rootCmd)
	weight.Register(rootCmd)
}
