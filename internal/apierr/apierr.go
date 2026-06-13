// Package apierr defines structured error types for API responses.
package apierr

import (
	"encoding/json"
	"fmt"
)

// APIError represents an error returned by the backend API.
// The backend returns CommonResult<T> with code != 0 on failure.
type APIError struct {
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	HTTPStatus int    `json:"-"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (http=%d, code=%d): %s", e.HTTPStatus, e.Code, e.Msg)
}

// New creates a new APIError.
func New(httpStatus, code int, msg string) *APIError {
	return &APIError{Code: code, Msg: msg, HTTPStatus: httpStatus}
}

// FromCommonResult parses a CommonResult JSON body and returns an APIError
// when code != 0. Returns nil if code == 0.
func FromCommonResult(httpStatus int, body []byte) *APIError {
	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return New(httpStatus, -1, fmt.Sprintf("failed to parse response: %s", err))
	}
	if result.Code == 0 {
		return nil
	}
	return New(httpStatus, result.Code, result.Msg)
}

// Exit codes for CLI.
const (
	ExitGeneral        = 1
	ExitConfig         = 2
	ExitAuthentication = 3
	ExitValidation     = 4
	ExitAPI            = 5
	ExitNetwork        = 6
)

// ExitType returns the error type string for a given exit code.
func ExitType(code int) string {
	switch code {
	case ExitConfig:
		return "config"
	case ExitAuthentication:
		return "authentication"
	case ExitValidation:
		return "validation"
	case ExitAPI:
		return "api_error"
	case ExitNetwork:
		return "network"
	default:
		return "general"
	}
}
