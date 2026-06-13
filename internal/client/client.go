// Package client provides an HTTP client with auth header injection,
// CommonResult parsing, and retry logic.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/tidwall/gjson"

	"github.com/weiliantong/cli/internal/apierr"
	"github.com/weiliantong/cli/internal/auth"
	"github.com/weiliantong/cli/internal/config"
)

// Response wraps the parsed API response data.
type Response struct {
	Data json.RawMessage // the "data" field from CommonResult
}

// Client is the HTTP client for calling backend APIs.
type Client struct {
	profile    *config.Profile
	authMgr    *auth.Manager
	httpClient *http.Client
}

// NewClient creates a new API client.
func NewClient(profile *config.Profile, authMgr *auth.Manager) *Client {
	return &Client{
		profile: profile,
		authMgr: authMgr,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Get sends a GET request.
func (c *Client) Get(ctx context.Context, path string, params map[string]any) (*Response, error) {
	return c.do(ctx, http.MethodGet, path, params, nil)
}

// Post sends a POST request with a JSON body.
func (c *Client) Post(ctx context.Context, path string, body any) (*Response, error) {
	return c.do(ctx, http.MethodPost, path, nil, body)
}

// Put sends a PUT request with a JSON body.
func (c *Client) Put(ctx context.Context, path string, body any) (*Response, error) {
	return c.do(ctx, http.MethodPut, path, nil, body)
}

// Delete sends a DELETE request with query params.
func (c *Client) Delete(ctx context.Context, path string, params map[string]any) (*Response, error) {
	return c.do(ctx, http.MethodDelete, path, params, nil)
}

// PutQuery sends a PUT request with query params (no body).
// Used for endpoints that take @RequestParam instead of @RequestBody.
func (c *Client) PutQuery(ctx context.Context, path string, params map[string]any) (*Response, error) {
	return c.do(ctx, http.MethodPut, path, params, nil)
}

// DeleteWithBody sends a DELETE request with a JSON body.
// Used for endpoints that take @RequestBody with DELETE.
func (c *Client) DeleteWithBody(ctx context.Context, path string, body any) (*Response, error) {
	return c.do(ctx, http.MethodDelete, path, nil, body)
}

// DoRaw sends a request and returns the raw response body without CommonResult parsing.
// Used by `wlt api` command.
func (c *Client) DoRaw(ctx context.Context, method, path string, params map[string]any, body any) ([]byte, error) {
	req, err := c.buildRequest(ctx, method, path, params, body)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	return raw, nil
}

// BuildDryRun returns the request details without sending.
func (c *Client) BuildDryRun(method, path string, params map[string]any, body any) (string, string, http.Header, []byte, error) {
	fullURL := c.profile.BaseURL + c.profile.APIPrefix + path
	if len(params) > 0 {
		v := url.Values{}
		for k, val := range params {
			v.Set(k, fmt.Sprintf("%v", val))
		}
		fullURL += "?" + v.Encode()
	}
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}
	headers := c.buildHeaders("dry-run-token")
	return method, fullURL, headers, bodyBytes, nil
}

func (c *Client) do(ctx context.Context, method, path string, params map[string]any, body any) (*Response, error) {
	// Retry once on network errors
	var lastErr error
	for attempt := 0; attempt < 2; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(500 * time.Millisecond):
			}
		}
		req, err := c.buildRequest(ctx, method, path, params, body)
		if err != nil {
			return nil, err
		}
		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			continue // retry on network error
		}
		// NOTE: close body inline, NOT defer — defer in a loop leaks until function returns
		raw, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("read response: %w", err)
		}

		// Parse CommonResult
		code := gjson.GetBytes(raw, "code").Int()
		if code != 0 {
			apiErr := apierr.FromCommonResult(resp.StatusCode, raw)
			return nil, apiErr
		}

		dataBytes := gjson.GetBytes(raw, "data").Raw
		return &Response{Data: json.RawMessage(dataBytes)}, nil
	}
	return nil, lastErr
}

func (c *Client) buildRequest(ctx context.Context, method, path string, params map[string]any, body any) (*http.Request, error) {
	fullURL := c.profile.BaseURL + c.profile.APIPrefix + path
	if len(params) > 0 {
		v := url.Values{}
		for k, val := range params {
			v.Set(k, fmt.Sprintf("%v", val))
		}
		fullURL += "?" + v.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Get valid token (refreshes if needed)
	token, err := c.authMgr.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}
	req.Header = c.buildHeaders(token)
	return req, nil
}

func (c *Client) buildHeaders(token string) http.Header {
	h := http.Header{}
	h.Set("Authorization", "Bearer "+token)
	h.Set("tenant-id", c.profile.TenantID)
	h.Set("Content-Type", "application/json")
	if c.profile.EnterpriseType != "" {
		h.Set("enterprise-type", c.profile.EnterpriseType)
	}
	return h
}
