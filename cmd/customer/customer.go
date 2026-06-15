package customer

import (
	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/cmd/partner"
)

var customerCmd = &cobra.Command{
	Use:   "customer",
	Short: "客户管理",
	Long:  "客户模块操作：客户信息、开票信息、结算信息、授信管理。",
}

func init() {
	partner.AddPartnerCommands(customerCmd, "CUSTOMER", "客户")
	customerCmd.AddCommand(partner.NewInvoiceGroup())
	customerCmd.AddCommand(partner.NewSettlementGroup("CUSTOMER"))
	customerCmd.AddCommand(partner.NewCreditGroup())
}

// Register adds the customer command to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(customerCmd)
}
