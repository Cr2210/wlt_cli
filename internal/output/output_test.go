package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	err := WriteJSON(&buf, map[string]any{"id": 1, "name": "test"})
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	data, ok := result["data"].(map[string]any)
	if !ok {
		t.Fatal("missing data field")
	}
	if data["name"] != "test" {
		t.Errorf("expected name=test, got %v", data["name"])
	}
}

func TestWritePagedJSON(t *testing.T) {
	var buf bytes.Buffer
	items := []map[string]any{{"id": 1}, {"id": 2}}
	err := WritePagedJSON(&buf, items, 100, 1, 20)
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	meta, ok := result["meta"].(map[string]any)
	if !ok {
		t.Fatal("missing meta field")
	}
	if meta["total"] != float64(100) {
		t.Errorf("expected total=100, got %v", meta["total"])
	}
}

func TestWriteError(t *testing.T) {
	var buf bytes.Buffer
	exitErr := NewExitError(3, "token expired", "请通过 --token 与 --tenant-id 传入有效鉴权信息")
	err := WriteError(&buf, exitErr)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"type":"authentication"`) {
		t.Errorf("expected authentication type, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), `"code":3`) {
		t.Errorf("expected code 3, got: %s", buf.String())
	}
}

func TestWriteRaw(t *testing.T) {
	var buf bytes.Buffer
	WriteRaw(&buf, []byte(`{"raw":true}`))
	if buf.String() != `{"raw":true}` {
		t.Errorf("unexpected output: %s", buf.String())
	}
}
