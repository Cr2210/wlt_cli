// Package auth handles login, token refresh, and logout.
// It uses the raw net/http client to avoid circular dependency with the client package.
package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/weiliantong/cli/internal/config"
)

// TokenResult holds the parsed token info from login/refresh.
type TokenResult struct {
	UserID       string `json:"userId"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresTime  int64  `json:"expiresTime"` // epoch milliseconds
}

// Manager handles authentication operations.
type Manager struct {
	cfg *config.Manager
}

// NewManager creates a new auth Manager.
func NewManager(cfg *config.Manager) *Manager {
	return &Manager{cfg: cfg}
}

// Login authenticates with username/password and stores tokens.
// NOTE: baseURL should be the fully-qualified prefix (baseURL + apiPrefix),
// e.g. "https://erpsit.api.w-lian.com/admin-api". RefreshToken() handles
// prefix concatenation internally via profile fields.
func (m *Manager) Login(ctx context.Context, baseURL, tenantID, username, password string) (*TokenResult, error) {
	body, _ := json.Marshal(map[string]string{
		"username":            username,
		"password":            password,
		"captchaVerification": "",
	})
	url := baseURL + "/system/auth/login"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("tenant-id", tenantID)

	resp, err := authHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("login request: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Parse CommonResult
	var result struct {
		Code int          `json:"code"`
		Msg  string       `json:"msg"`
		Data *TokenResult `json:"data"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("login failed (code=%d): %s", result.Code, result.Msg)
	}
	if result.Data == nil {
		return nil, fmt.Errorf("login response missing data")
	}

	// Store tokens in profile
	p, err := m.cfg.ActiveProfile()
	if err != nil {
		return nil, fmt.Errorf("get active profile: %w", err)
	}
	p.AccessToken = result.Data.AccessToken
	p.RefreshToken = result.Data.RefreshToken
	p.ExpiresTime = result.Data.ExpiresTime
	if err := m.cfg.Save(); err != nil {
		return nil, fmt.Errorf("save config: %w", err)
	}

	return result.Data, nil
}

// RefreshToken refreshes the access token using the stored refresh token.
func (m *Manager) RefreshToken(ctx context.Context) (*TokenResult, error) {
	p, err := m.cfg.ActiveProfile()
	if err != nil {
		return nil, fmt.Errorf("get active profile: %w", err)
	}
	if p.RefreshToken == "" {
		return nil, fmt.Errorf("no refresh token available, please run 'wlt auth login'")
	}

	url := fmt.Sprintf("%s%s/system/auth/refresh-token?refreshToken=%s",
		p.BaseURL, p.APIPrefix, p.RefreshToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("tenant-id", p.TenantID)

	resp, err := authHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("refresh request: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result struct {
		Code int          `json:"code"`
		Msg  string       `json:"msg"`
		Data *TokenResult `json:"data"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("refresh failed (code=%d): %s", result.Code, result.Msg)
	}
	if result.Data == nil {
		return nil, fmt.Errorf("refresh response missing data")
	}

	// Update stored tokens
	p.AccessToken = result.Data.AccessToken
	p.RefreshToken = result.Data.RefreshToken
	p.ExpiresTime = result.Data.ExpiresTime
	if err := m.cfg.Save(); err != nil {
		return nil, fmt.Errorf("save config: %w", err)
	}

	return result.Data, nil
}

// GetValidToken returns a valid access token, refreshing if expired.
func (m *Manager) GetValidToken(ctx context.Context) (string, error) {
	p, err := m.cfg.ActiveProfile()
	if err != nil {
		return "", fmt.Errorf("get active profile: %w", err)
	}
	if p.AccessToken == "" {
		return "", fmt.Errorf("not logged in, please run 'wlt auth login'")
	}
	// Check if token is expired (with 60s buffer)
	if time.Now().UnixMilli() > p.ExpiresTime-60000 {
		if _, err := m.RefreshToken(ctx); err != nil {
			return "", fmt.Errorf("token expired and refresh failed: %w", err)
		}
	}
	// Re-read profile after potential refresh
	p, err = m.cfg.ActiveProfile()
	if err != nil {
		return "", fmt.Errorf("get active profile: %w", err)
	}
	return p.AccessToken, nil
}

// Logout clears the stored tokens.
func (m *Manager) Logout() error {
	p, err := m.cfg.ActiveProfile()
	if err != nil {
		return fmt.Errorf("get active profile: %w", err)
	}
	p.AccessToken = ""
	p.RefreshToken = ""
	p.ExpiresTime = 0
	return m.cfg.Save()
}

// Status returns the current authentication status.
func (m *Manager) Status() (loggedIn bool, expiresAt int64, err error) {
	p, err := m.cfg.ActiveProfile()
	if err != nil {
		return false, 0, err
	}
	return p.AccessToken != "", p.ExpiresTime, nil
}

// GetTenantIDByName looks up the tenant ID by company name via the backend API.
func (m *Manager) GetTenantIDByName(ctx context.Context, baseURL, companyName string) (string, error) {
	url := baseURL + "/system/tenant/get-id-by-name?name=" + url.QueryEscape(companyName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	resp, err := authHTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request tenant id: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data int64  `json:"data"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}
	if result.Code != 0 {
		return "", fmt.Errorf("查询租户失败 (code=%d): %s", result.Code, result.Msg)
	}
	if result.Data == 0 {
		return "", fmt.Errorf("未找到公司 %q 对应的租户", companyName)
	}
	return fmt.Sprintf("%d", result.Data), nil
}

// authHTTPClient is a dedicated HTTP client for auth requests with a timeout.
var authHTTPClient = &http.Client{Timeout: 30 * time.Second}
