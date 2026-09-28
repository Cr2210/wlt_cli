package cmdutil

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestTimeRangeFlagsRoundTrip(t *testing.T) {
	c := &cobra.Command{Use: "t"}
	AddTimeRangeFlags(c, "load-date", "realLoadDate", "装货日期")
	for _, f := range []string{"load-date-start", "load-date-end"} {
		if c.Flags().Lookup(f) == nil {
			t.Fatalf("missing flag %q", f)
		}
	}
	_ = c.Flags().Set("load-date-start", "2026-07-01 00:00:00")
	params := map[string]any{}
	CollectTimeRangeFlags(c, params, "load-date", "realLoadDate")
	if params["realLoadDate[0]"] != "2026-07-01 00:00:00" {
		t.Errorf("realLoadDate[0] = %v, want 2026-07-01 00:00:00", params["realLoadDate[0]"])
	}
	if _, ok := params["realLoadDate[1]"]; ok {
		t.Errorf("realLoadDate[1] should be absent when end flag empty")
	}
}
