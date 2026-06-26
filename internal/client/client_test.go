package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	// Stateless client: auth context supplied directly, no config/auth manager.
	cli := NewClient(RequestContext{
		BaseURL:  ts.URL,
		APIPrefix: "",
		TenantID: "1",
		Token:    "test-token",
	})
	return cli, ts
}

func TestGet_Success(t *testing.T) {
	cli, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/erp/warehouse/page" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		resp, _ := json.Marshal(map[string]any{
			"code": 0,
			"msg":  "",
			"data": map[string]any{
				"list":  []map[string]any{{"id": 1, "name": "仓库A"}},
				"total": 1,
			},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	})

	resp, err := cli.Get(context.Background(), "/erp/warehouse/page", map[string]any{"pageNo": 1})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data == nil {
		t.Fatal("expected data")
	}
}

func TestGet_APIError(t *testing.T) {
	cli, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		resp, _ := json.Marshal(map[string]any{
			"code": 401,
			"msg":  "token 已过期",
			"data": nil,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	})

	_, err := cli.Get(context.Background(), "/test", nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPost_Success(t *testing.T) {
	cli, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		resp, _ := json.Marshal(map[string]any{
			"code": 0,
			"msg":  "",
			"data": 42,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	})

	resp, err := cli.Post(context.Background(), "/erp/warehouse/create", map[string]any{"name": "新仓库"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data == nil {
		t.Fatal("expected data")
	}
}

func TestDelete_Success(t *testing.T) {
	cli, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		resp, _ := json.Marshal(map[string]any{
			"code": 0,
			"msg":  "",
			"data": true,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	})

	_, err := cli.Delete(context.Background(), "/erp/warehouse/delete", map[string]any{"ids": "1,2"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDoRaw(t *testing.T) {
	cli, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"raw":true}`))
	})

	raw, err := cli.DoRaw(context.Background(), http.MethodGet, "/test", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != `{"raw":true}` {
		t.Errorf("unexpected raw: %s", raw)
	}
}
