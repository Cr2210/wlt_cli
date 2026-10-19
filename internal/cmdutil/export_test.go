package cmdutil

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestResolveExportPath(t *testing.T) {
	root := &cobra.Command{Use: "root"}
	root.PersistentFlags().StringP("output", "o", "", "导出文件路径")
	var got string
	child := &cobra.Command{Use: "export", RunE: func(cmd *cobra.Command, args []string) error {
		got = resolveExportPath(cmd)
		return nil
	}}
	root.AddCommand(child)
	root.SetArgs([]string{"export", "--output", "/tmp/x.xlsx"})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if got != "/tmp/x.xlsx" {
		t.Errorf("resolveExportPath with --output = %q, want /tmp/x.xlsx", got)
	}
}

func TestResolveExportPathAutoName(t *testing.T) {
	root := &cobra.Command{Use: "root"}
	root.PersistentFlags().StringP("output", "o", "", "导出文件路径")
	child := &cobra.Command{Use: "export"}
	root.AddCommand(child)
	got := resolveExportPath(child)
	want := regexp.MustCompile(`^wlt-export-\d{8}-\d{6}\.xlsx$`)
	if !want.MatchString(got) {
		t.Errorf("resolveExportPath auto name = %q, want wlt-export-YYYYMMDD-HHMMSS.xlsx", got)
	}
}

func TestWriteExportFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.xlsx")
	raw := []byte{0x50, 0x4b, 0x03, 0x04} // xlsx 魔数（zip 头）
	abs, n, err := writeExportFile(path, raw)
	if err != nil {
		t.Fatalf("writeExportFile: %v", err)
	}
	if n != len(raw) {
		t.Errorf("bytes = %d, want %d", n, len(raw))
	}
	got, err := os.ReadFile(path)
	if err != nil || string(got) != string(raw) {
		t.Errorf("file content mismatch: %v", err)
	}
	if !filepath.IsAbs(abs) {
		t.Errorf("returned path %q not absolute", abs)
	}
}

func TestBuildExportEnvelope(t *testing.T) {
	data := buildExportEnvelope("/tmp/a b.xlsx", 42)
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("envelope not JSON: %v", err)
	}
	if m["file"] != "/tmp/a b.xlsx" || m["bytes"] != float64(42) {
		t.Errorf("envelope = %v, want file/bytes", m)
	}
	if !strings.HasPrefix(string(data), "{\"file\":") {
		t.Errorf("envelope should marshal deterministically: %s", data)
	}
}
