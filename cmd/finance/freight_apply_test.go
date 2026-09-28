package finance

import "testing"

func TestFinanceFreightApplyCmdSubcommands(t *testing.T) {
	want := map[string]bool{
		"list": false, "get": false, "create": false, "update": false,
		"delete": false, "update-status": false, "summary": false,
		"available-waybill": false, "available-waybill-count": false,
		"available-settlement": false, "export": false,
	}
	for _, sub := range financeFreightApplyCmd.Commands() {
		if _, ok := want[sub.Name()]; ok {
			want[sub.Name()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("finance freight-apply missing subcommand %q", name)
		}
	}
}

func TestFinanceFreightApplyListCmdFlags(t *testing.T) {
	c := newFinanceFreightApplyListCmd()
	for _, f := range []string{"no", "carrier-enterprise-id", "carrier-name", "payment-account-id", "receipt-account-id", "service-user-id", "approve-status", "pay-date-start", "pay-date-end", "create-time-start", "create-time-end", "page-no", "page-size"} {
		if c.Flags().Lookup(f) == nil {
			t.Errorf("freight-apply list cmd missing flag %q", f)
		}
	}
}

func TestFinanceFreightApplyAvailableSettlementCmdFlags(t *testing.T) {
	c := newFinanceFreightApplyAvailableSettlementCmd()
	for _, f := range []string{"no", "carrier-enterprise-id", "carrier-name", "apply-id", "settlement-date-start", "settlement-date-end", "page-no", "page-size"} {
		if c.Flags().Lookup(f) == nil {
			t.Errorf("available-settlement cmd missing flag %q", f)
		}
	}
}
