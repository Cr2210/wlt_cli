package cmdutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/output"
)

// export-excel 端点返回二进制 Excel 流（非 CommonResult 信封），
// 导出命令统一经 RunExport 走 DoRaw + 写文件，stdout 仍输出机器可读信封。

// RunExport 调用 export 端点并把二进制响应写入 --output 指定文件
// （缺省自动命名 wlt-export-<时间戳>.xlsx），成功后输出 {"file":...,"bytes":N} 信封。
func RunExport(cmd *cobra.Command, apiPath string, params map[string]any) error {
	if err := EnsureClient(); err != nil {
		return err
	}
	raw, err := GetClient().DoRaw(cmd.Context(), http.MethodGet, apiPath, params, nil)
	if err != nil {
		return output.NewExitError(5, fmt.Sprintf("导出失败: %s", err), "")
	}
	path := resolveExportPath(cmd)
	abs, n, err := writeExportFile(path, raw)
	if err != nil {
		return output.NewExitError(2, fmt.Sprintf("写入导出文件失败: %s", err), "检查 --output 路径是否可写")
	}
	return OutputJSON(buildExportEnvelope(abs, n))
}

// resolveExportPath 返回 --output 值；未指定时按当前时间生成自动文件名。
func resolveExportPath(cmd *cobra.Command) string {
	if v, _ := cmd.Flags().GetString("output"); v != "" {
		return v
	}
	return fmt.Sprintf("wlt-export-%s.xlsx", time.Now().Format("20060102-150405"))
}

// writeExportFile 落盘并返回绝对路径与字节数。
func writeExportFile(path string, raw []byte) (string, int, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", 0, err
	}
	if err := os.WriteFile(abs, raw, 0o644); err != nil {
		return "", 0, err
	}
	return abs, len(raw), nil
}

// buildExportEnvelope 构造 {"file":...,"bytes":N} 信封（字段序固定，便于断言）。
func buildExportEnvelope(absPath string, n int) json.RawMessage {
	return json.RawMessage(fmt.Sprintf(`{"file":%s,"bytes":%d}`, mustJSONString(absPath), n))
}

func mustJSONString(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}
