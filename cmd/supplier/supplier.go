package supplier

import (
	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/cmd/partner"
)

var supplierCmd = &cobra.Command{
	Use:   "supplier",
	Short: "供应商管理",
	Long:  "供应商模块操作：供应商信息、开票信息、结算信息。",
}

func init() {
	partner.AddPartnerCommands(supplierCmd, "2", "供应商")
	supplierCmd.AddCommand(partner.NewInvoiceGroup())
	supplierCmd.AddCommand(partner.NewSettlementGroup())
}

// Register adds the supplier command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(supplierCmd)
}
