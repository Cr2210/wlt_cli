package cmdutil

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// ParseJSONData parses a JSON string argument.
func ParseJSONData(s string) (map[string]any, error) {
	if s == "" {
		return nil, fmt.Errorf("empty JSON data")
	}
	var result map[string]any
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return result, nil
}

// flagToParamKey converts a kebab-case flag name (e.g. "supplier-id") to the
// camelCase query-param key the backend expects (e.g. "supplierId"). Names
// without a hyphen are returned unchanged, so single-word flags like "status"
// or already-camelCase keys are unaffected.
func flagToParamKey(name string) string {
	if !strings.Contains(name, "-") {
		return name
	}
	var b strings.Builder
	capitalize := false
	for i, r := range name {
		if r == '-' {
			capitalize = true
			continue
		}
		if capitalize && i > 0 {
			b.WriteRune(toUpperASCII(r))
			capitalize = false
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func toUpperASCII(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - 32
	}
	return r
}

// CollectStringFlag collects a non-empty string flag into params map,
// converting the kebab-case flag name to a camelCase query-param key.
func CollectStringFlag(cmd *cobra.Command, params map[string]any, name string) {
	v, err := cmd.Flags().GetString(name)
	if err == nil && v != "" {
		params[flagToParamKey(name)] = v
	}
}

// CollectStringFlags collects non-empty string flags into params map,
// converting each kebab-case flag name to a camelCase query-param key.
func CollectStringFlags(cmd *cobra.Command, params map[string]any, names ...string) {
	for _, name := range names {
		v, err := cmd.Flags().GetString(name)
		if err == nil && v != "" {
			params[flagToParamKey(name)] = v
		}
	}
}

// CollectIntFlags collects non-zero int flags into params map,
// converting each kebab-case flag name to a camelCase query-param key.
func CollectIntFlags(cmd *cobra.Command, params map[string]any, names ...string) {
	for _, name := range names {
		v, err := cmd.Flags().GetInt64(name)
		if err == nil && v != 0 {
			params[flagToParamKey(name)] = strconv.FormatInt(v, 10)
		}
	}
}
