package cmdutil

import (
	"encoding/json"
	"fmt"
	"strconv"

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

// CollectStringFlag collects a non-empty string flag into params map.
func CollectStringFlag(cmd *cobra.Command, params map[string]any, name string) {
	v, err := cmd.Flags().GetString(name)
	if err == nil && v != "" {
		params[name] = v
	}
}

// CollectStringFlags collects non-empty string flags into params map.
func CollectStringFlags(cmd *cobra.Command, params map[string]any, names ...string) {
	for _, name := range names {
		v, err := cmd.Flags().GetString(name)
		if err == nil && v != "" {
			params[name] = v
		}
	}
}

// CollectIntFlags collects non-zero int flags into params map.
func CollectIntFlags(cmd *cobra.Command, params map[string]any, names ...string) {
	for _, name := range names {
		v, err := cmd.Flags().GetInt64(name)
		if err == nil && v != 0 {
			params[name] = strconv.FormatInt(v, 10)
		}
	}
}
