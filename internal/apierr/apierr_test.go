package apierr

import (
	"testing"
)

func TestFromCommonResult_Success(t *testing.T) {
	body := []byte(`{"code":0,"data":{"id":1},"msg":""}`)
	err := FromCommonResult(200, body)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestFromCommonResult_Error(t *testing.T) {
	body := []byte(`{"code":401,"data":null,"msg":"token 已过期"}`)
	err := FromCommonResult(200, body)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != 401 {
		t.Errorf("expected code 401, got %d", err.Code)
	}
	if err.Msg != "token 已过期" {
		t.Errorf("unexpected msg: %s", err.Msg)
	}
	if err.HTTPStatus != 200 {
		t.Errorf("expected http 200, got %d", err.HTTPStatus)
	}
}

func TestFromCommonResult_InvalidJSON(t *testing.T) {
	body := []byte(`not json`)
	err := FromCommonResult(500, body)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if err.Code != -1 {
		t.Errorf("expected code -1, got %d", err.Code)
	}
}

func TestExitType(t *testing.T) {
	tests := []struct {
		code int
		want string
	}{
		{1, "general"},
		{2, "config"},
		{3, "authentication"},
		{4, "validation"},
		{5, "api_error"},
		{6, "network"},
		{99, "general"},
	}
	for _, tt := range tests {
		got := ExitType(tt.code)
		if got != tt.want {
			t.Errorf("ExitType(%d) = %q, want %q", tt.code, got, tt.want)
		}
	}
}
