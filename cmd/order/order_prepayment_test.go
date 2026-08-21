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
