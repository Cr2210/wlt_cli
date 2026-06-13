// Package output handles JSON formatting for CLI stdout/stderr separation.
package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/weiliantong/cli/internal/apierr"
)

// ExitError represents a structured error to be written to stderr as JSON.
type ExitError struct {
	Code    int            `json:"-"`
	Type    string         `json:"type"`
	Message string         `json:"message"`
	Hint    string         `json:"hint,omitempty"`
	Detail  map[string]any `json:"detail,omitempty"`
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// NewExitError creates an ExitError with the exit code and type auto-derived.
func NewExitError(code int, msg string, hint string) *ExitError {
	return &ExitError{
		Code:    code,
		Type:    apierr.ExitType(code),
		Message: msg,
		Hint:    hint,
	}
}

// NewExitErrorWithDetail creates an ExitError with additional detail.
func NewExitErrorWithDetail(code int, msg string, hint string, detail map[string]any) *ExitError {
	return &ExitError{
		Code:    code,
		Type:    apierr.ExitType(code),
		Message: msg,
		Hint:    hint,
		Detail:  detail,
	}
}

// errorEnvelope wraps an ExitError for JSON output.
type errorEnvelope struct {
	Error *errorBody `json:"error"`
}

type errorBody struct {
	Type    string         `json:"type"`
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Hint    string         `json:"hint,omitempty"`
	Detail  map[string]any `json:"detail,omitempty"`
}

// WriteJSON writes a data response to stdout.
func WriteJSON(out io.Writer, data any) error {
	envelope := map[string]any{"data": data}
	enc := json.NewEncoder(out)
	return enc.Encode(envelope)
}

// WritePagedJSON writes a paginated data response to stdout.
func WritePagedJSON(out io.Writer, list any, total int64, pageNo, pageSize int) error {
	envelope := map[string]any{
		"data": list,
		"meta": map[string]any{
			"page_no":   pageNo,
			"page_size": pageSize,
			"total":     total,
		},
	}
	enc := json.NewEncoder(out)
	return enc.Encode(envelope)
}

// WriteError writes a structured error to stderr.
func WriteError(errOut io.Writer, exitErr *ExitError) error {
	env := errorEnvelope{
		Error: &errorBody{
			Type:    exitErr.Type,
			Code:    exitErr.Code,
			Message: exitErr.Message,
			Hint:    exitErr.Hint,
			Detail:  exitErr.Detail,
		},
	}
	enc := json.NewEncoder(errOut)
	return enc.Encode(env)
}

// WriteRaw writes raw bytes to output without any wrapping.
func WriteRaw(out io.Writer, data []byte) {
	out.Write(data)
}
