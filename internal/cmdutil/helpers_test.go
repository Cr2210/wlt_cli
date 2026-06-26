package cmdutil

import "testing"

func TestFlagToParamKey(t *testing.T) {
	cases := []struct{ in, want string }{
		{"status", "status"},
		{"no", "no"},
		{"supplier-id", "supplierId"},
		{"warehouse-id", "warehouseId"},
		{"approve-status", "approveStatus"},
		{"bill-date-start", "billDateStart"},
		{"product-id", "productId"},
		{"pageNo", "pageNo"}, // already camelCase, unchanged
		{"", ""},
	}
	for _, c := range cases {
		if got := flagToParamKey(c.in); got != c.want {
			t.Errorf("flagToParamKey(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
