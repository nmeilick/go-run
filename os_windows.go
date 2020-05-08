// +build windows

package run

import (
	"path/filepath"
	"strings"
)

// matchEnv checks the given environment name against a list of wildcard patterns
// and returns true if there's any match. The matching is case-insensitive.
func matchEnv(name string, patterns []string) bool {
	name = strings.ToLower(name)
	for _, p := range patterns {
		if ok, err := filepath.Match(name, strings.ToLower(p)); err == nil && ok {
			return true
		}
	}
	return false
}
