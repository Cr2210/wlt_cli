package finance

import "testing"

func TestFinanceReceiptPaymentPageCmdFlags(t *testing.T) {
	c := newFinanceReceiptPaymentPageCmd()
	for _, f := range []string{"is-prepaid", "page-no", "page-size"} {
		if c.Flags().Lookup(f) == nil {
			t.Errorf("receipt-payment page cmd missing flag %q", f)
		}
	}
}
