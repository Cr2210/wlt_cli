package cmdutil

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/weiliantong/cli/internal/output"
)

// OutputJSON writes data as JSON to stdout.
func OutputJSON(data any) error {
	return output.WriteJSON(os.Stdout, data)
}

// OutputPagedJSON writes paginated data as JSON to stdout.
func OutputPagedJSON(list any, total int64, pageNo, pageSize int) error {
	return output.WritePagedJSON(os.Stdout, list, total, pageNo, pageSize)
}

// OutputRaw writes raw bytes to stdout.
func OutputRaw(data []byte) {
	output.WriteRaw(os.Stdout, data)
}

// ParsePagedJSON extracts paginated data from API response and outputs it.
func ParsePagedJSON(data json.RawMessage, pageNo, pageSize int) error {
	var paged struct {
		List  json.RawMessage `json:"list"`
		Total int64           `json:"total"`
	}
	if err := json.Unmarshal(data, &paged); err != nil {
		return output.NewExitError(5, fmt.Sprintf("解析响应失败: %s", err), "")
	}
	var list any
	if err := json.Unmarshal(paged.List, &list); err != nil {
		list = []any{}
	}
	return OutputPagedJSON(list, paged.Total, pageNo, pageSize)
}
