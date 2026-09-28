package sale

import "testing"

func TestSaleAccountingCmdSubcommands(t *testing.T) {
	want := map[string]bool{
		"list": false, "page-count": false, "get": false, "create": false,
		"update": false, "delete": false, "waybill-list": false,
		"waybill-count": false, "calculate": false,
		"default-config": false, "export": false,
	}
	for _, sub := range saleAccountingCmd.Commands() {
		if _, ok := want[sub.Name()]; ok {
			want[sub.Name()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("sale accounting missing subcommand %q", name)
		}
	}
}

func TestSaleAccountingListCmdFlags(t *testing.T) {
	c := newSaleAccountingListCmd()
	for _, f := range []string{"no", "customer-name", "waybill-no", "period-start", "period-end", "page-no", "page-size"} {
		if c.Flags().Lookup(f) == nil {
			t.Errorf("accounting list cmd missing flag %q", f)
		}
	}
}

func TestSaleAccountingWaybillListCmdFlags(t *testing.T) {
	c := newSaleAccountingWaybillListCmd()
	for _, f := range []string{"customer-id", "exclude-ids", "waybill-no", "car-number", "load-date-start", "load-date-end", "unload-date-start", "unload-date-end", "page-no", "page-size"} {
		if c.Flags().Lookup(f) == nil {
			t.Errorf("waybill-list cmd missing flag %q", f)
		}
	}
}
