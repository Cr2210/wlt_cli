package stock

import (
	"testing"

	"github.com/weiliantong/cli/internal/cmdutil"
)

func TestNewStockPagedCmdFlags(t *testing.T) {
	c := newStockPagedCmd("page", "/erp/stock", "page", "分页查询产品库存",
		[]cmdutil.FlagSpec{{Name: "product-id"}, {Name: "keyword"}}, true)
	for _, f := range []string{"page-no", "page-size", "product-id", "keyword", "headers"} {
		if c.Flags().Lookup(f) == nil {
			t.Errorf("paged cmd missing flag %q", f)
		}
	}
	if c.Name() != "page" {
		t.Errorf("name = %q, want page", c.Name())
	}
}

func TestNewStockCountCmdFlags(t *testing.T) {
	c := newStockCountCmd("page-count", "/erp/stock", "page-count", "统计",
		[]cmdutil.FlagSpec{{Name: "warehouse-id"}})
	if c.Flags().Lookup("warehouse-id") == nil {
		t.Errorf("count cmd missing flag warehouse-id")
	}
	for _, f := range []string{"page-no", "page-size", "headers"} {
		if c.Flags().Lookup(f) != nil {
			t.Errorf("count cmd should not have flag %q", f)
		}
	}
}
