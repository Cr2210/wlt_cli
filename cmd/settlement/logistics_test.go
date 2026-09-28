package settlement

import "testing"

func TestLogisticsSettlementCmdSubcommands(t *testing.T) {
	want := map[string]bool{
		"list": false, "page-count": false, "get": false, "create": false,
		"update": false, "delete": false, "update-status": false,
		"mismatch-list": false, "available-waybill": false,
		"available-waybill-count": false, "recalculate-all": false,
		"ignore-remainder": false, "export": false,
	}
	for _, sub := range logisticsSettlementCmd.Commands() {
		if _, ok := want[sub.Name()]; ok {
			want[sub.Name()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("settlement logistics missing subcommand %q", name)
		}
	}
}

func TestLogisticsSettlementListCmdFlags(t *testing.T) {
	c := newLogisticsSettlementListCmd()
	for _, f := range []string{"no", "name", "carrier-name", "carrier-enterprise-id", "status", "start-date", "end-date", "page-no", "page-size"} {
		if c.Flags().Lookup(f) == nil {
			t.Errorf("logistics list cmd missing flag %q", f)
		}
	}
}

func TestLogisticsSettlementIgnoreRemainderCmdFlags(t *testing.T) {
	c := newLogisticsSettlementIgnoreRemainderCmd()
	for _, f := range []string{"id", "type", "reason"} {
		if c.Flags().Lookup(f) == nil {
			t.Errorf("logistics ignore-remainder cmd missing flag %q", f)
		}
	}
}
