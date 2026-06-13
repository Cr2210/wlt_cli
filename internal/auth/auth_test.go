package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/weiliantong/cli/internal/config"
)

func newTestManager(t *testing.T, handler http.HandlerFunc) (*Manager, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	// Use temp path to avoid writing to real ~/.wlt/config.yaml
	dir := t.TempDir()
	cfg := config.NewManagerWithPath(filepath.Join(dir, "config.yaml"))
	_ = cfg.Load()
	// Override profile to point to test server
	p, _ := cfg.ActiveProfile()
	p.BaseURL = ts.URL
	p.APIPrefix = ""
	p.TenantID = "1"
	return &Manager{cfg: cfg}, ts
}

func TestLogin_Success(t *testing.T) {
	mgr, ts := newTestManager(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/system/auth/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		resp, _ := json.Marshal(map[string]any{
			"code": 0,
			"msg":  "",
			"data": map[string]any{
				"userId":       "1",
				"accessToken":  "at-123",
				"refreshToken": "rt-456",
				"expiresTime":  9999999999999,
			},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	})

	result, err := mgr.Login(context.Background(), ts.URL, "1", "admin", "password")
	if err != nil {
		t.Fatal(err)
	}
	if result.AccessToken != "at-123" {
		t.Errorf("expected at-123, got %s", result.AccessToken)
	}
}

func TestLogin_Failure(t *testing.T) {
	mgr, ts := newTestManager(t, func(w http.ResponseWriter, r *http.Request) {
		resp, _ := json.Marshal(map[string]any{
			"code": 401,
			"msg":  "用户名或密码错误",
			"data": nil,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	})

	_, err := mgr.Login(context.Background(), ts.URL, "1", "admin", "wrong")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestStatus(t *testing.T) {
	mgr, _ := newTestManager(t, nil)
	loggedIn, _, err := mgr.Status()
	if err != nil {
		t.Fatal(err)
	}
	if loggedIn {
		t.Error("should not be logged in initially")
	}
}
