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
)

// Response wraps the parsed API response data.
type Response struct {
	Data json.RawMessage // the "data" field from CommonResult
}

// RequestContext carries everything needed to build an authenticated request.
// It is resolved once per process from the active profile plus per-call flags
// (see internal/cmdutil). The client holds no other auth state — auth is
// stateless: every request carries the token supplied on the command line.
type RequestContext struct {
	BaseURL        string
	APIPrefix      string
	TenantID       string
	EnterpriseType string
	Token          string // raw access token; sent as "Authorization: Bearer <token>"
}

// Client is the HTTP client for calling backend APIs.
type Client struct {
	ctx        RequestContext
	httpClient *http.Client
}

// NewClient creates a new API client from a resolved RequestContext.
func NewClient(ctx RequestContext) *Client {
	return &Client{
		ctx: ctx,
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
	fullURL := c.ctx.BaseURL + c.ctx.APIPrefix + path
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
	// Mask the token in dry-run output to avoid echoing secrets to stdout.
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
	fullURL := c.ctx.BaseURL + c.ctx.APIPrefix + path
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

	// Stateless auth: use the token supplied on the command line.
	req.Header = c.buildHeaders(c.ctx.Token)
	return req, nil
}

func (c *Client) buildHeaders(token string) http.Header {
	h := http.Header{}
	h.Set("Authorization", "Bearer "+token)
	h.Set("tenant-id", c.ctx.TenantID)
	h.Set("Content-Type", "application/json")
	if c.ctx.EnterpriseType != "" {
		h.Set("enterprise-type", c.ctx.EnterpriseType)
	}
	return h
}
