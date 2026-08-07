package cmdutil

import (
	"testing"

	"github.com/spf13/cobra"
)

// 子命令名断言辅助
func mustHaveSub(t *testing.T, parent *cobra.Command, names ...string) {
	t.Helper()
	got := map[string]bool{}
	for _, c := range parent.Commands() {
		got[c.Name()] = true
	}
	for _, n := range names {
		if !got[n] {
			t.Errorf("missing subcommand %q; got %v", n, got)
		}
	}
}

func TestCollectDocumentTimeRange(t *testing.T) {
	cmd := &cobra.Command{}
	addDocumentTimeRangeFlags(cmd)
	_ = cmd.Flags().Set("start-time", "2026-07-01 00:00:00")
	_ = cmd.Flags().Set("end-time", "2026-07-31 23:59:59")

	params := map[string]any{}
	collectDocumentTimeRange(cmd, params, "inTime")
	if params["inTime[0]"] != "2026-07-01 00:00:00" {
		t.Errorf("inTime[0] = %v, want start-time", params["inTime[0]"])
	}
	if params["inTime[1]"] != "2026-07-31 23:59:59" {
		t.Errorf("inTime[1] = %v, want end-time", params["inTime[1]"])
	}
}

func TestCollectDocumentTimeRangeEmpty(t *testing.T) {
	cmd := &cobra.Command{}
	addDocumentTimeRangeFlags(cmd) // 不 set 任何值
	params := map[string]any{}
	collectDocumentTimeRange(cmd, params, "outTime")
	if _, ok := params["outTime[0]"]; ok {
		t.Errorf("outTime[0] should be absent when flags empty; got %v", params["outTime[0]"])
	}
}

func TestDocumentCmdsStructure(t *testing.T) {
	parent := &cobra.Command{Use: "test"}
	DocumentCmds(parent, DocumentConfig{
		Name:    "in",
		APIPath: "/erp/stock-in",
		Label:   "入库单",
		TimeKey: "inTime",
		Headers: true,
		PageAlias: true,
		Filters: []FlagSpec{
			{Name: "no", Usage: "单号"},
			{Name: "supplier-id", Usage: "供应商"},
		},
	})

	// 子命令齐全
	in := subNamed(parent, "in")
	if in == nil {
		t.Fatal("missing subcommand 'in'")
	}
	mustHaveSub(t, in, "list", "page-count", "get", "create", "update", "delete", "update-status")

	// list 有 page 别名
	list := subNamed(in, "list")
	for _, a := range list.Aliases {
		if a == "page" {
			goto ok
		}
	}
	t.Errorf("list Aliases = %v, want contain 'page'", list.Aliases)
ok:

	// list 注册了分页/时间/headers/filter flag
	for _, f := range []string{"page-no", "page-size", "start-time", "end-time", "headers", "no", "supplier-id"} {
		if list.Flags().Lookup(f) == nil {
			t.Errorf("list missing flag %q", f)
		}
	}
}

func TestDocumentCmdsNoExtraTimeFlags(t *testing.T) {
	// TimeKey 为空时不注册时间 flag
	parent := &cobra.Command{Use: "t"}
	DocumentCmds(parent, DocumentConfig{Name: "x", APIPath: "/erp/x", Label: "X", Filters: nil})
	x := subNamed(parent, "x")
	list := subNamed(x, "list")
	if list.Flags().Lookup("start-time") != nil {
		t.Errorf("start-time should not be registered when TimeKey empty")
	}
}

// subNamed 在 parent 的直接子命令里按名查找。
func subNamed(parent *cobra.Command, name string) *cobra.Command {
	for _, c := range parent.Commands() {
		if c.Name() == name {
			return c
		}
	}
	return nil
}
