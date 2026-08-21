package order

import "testing"

func TestPrepayRelationListCmdFlags(t *testing.T) {
	c := newOrderPrepayRelationListCmd()
	for _, f := range []string{"order-id", "page-no", "page-size", "relation-type", "headers"} {
		if c.Flags().Lookup(f) == nil {
			t.Errorf("list cmd missing flag %q", f)
		}
	}
}

func TestPrepayRelationCreateCmdFlags(t *testing.T) {
	c := newOrderPrepayRelationCreateCmd()
	for _, f := range []string{"order-id", "relation-type", "relation-id", "relation-amount"} {
		if c.Flags().Lookup(f) == nil {
			t.Errorf("create cmd missing flag %q", f)
		}
	}
}

func TestPrepayRelationUpdateAmountCmdFlags(t *testing.T) {
	c := newOrderPrepayRelationUpdateAmountCmd()
	for _, f := range []string{"id", "relation-amount"} {
		if c.Flags().Lookup(f) == nil {
			t.Errorf("update-amount cmd missing flag %q", f)
		}
	}
}

func TestPrepayAvailableCmdFlags(t *testing.T) {
	for _, name := range []string{"available-payments", "available-initials"} {
		c := newPrepayAvailableCmd(name, "查询")
		for _, f := range []string{"supplier-id", "order-id", "biz-date-start", "biz-date-end", "no", "headers", "page-no", "page-size"} {
			if c.Flags().Lookup(f) == nil {
				t.Errorf("%s cmd missing flag %q", name, f)
			}
		}
		if c.Name() != name {
			t.Errorf("name = %q, want %q", c.Name(), name)
		}
	}
}

func TestPrepayRelationCmdSubcommands(t *testing.T) {
	want := map[string]bool{"list": false, "create": false, "update-amount": false, "available-payments": false, "available-initials": false}
	for _, sub := range orderPrepaymentRelationCmd.Commands() {
		if _, ok := want[sub.Name()]; ok {
			want[sub.Name()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("prepayment-relation missing subcommand %q", name)
		}
	}
}
