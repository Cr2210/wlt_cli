package settlement

import "testing"

func TestSettlementMainIgnoreRemainderCmdFlags(t *testing.T) {
	c := newSettlementMainIgnoreRemainderCmd()
	for _, f := range []string{"id", "type", "reason"} {
		if c.Flags().Lookup(f) == nil {
			t.Errorf("ignore-remainder cmd missing flag %q", f)
		}
	}
}
