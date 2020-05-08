// +build !windows

package run

import "path/filepath"

// matchEnv checks the given environment name against a list of wildcard patterns
// and returns true if there's any match. The matching is case-sensitive.
func matchEnv(name string, patterns []string) bool {
	for _, p := range patterns {
		if ok, err := filepath.Match(name, p); err == nil && ok {
			return true
		}
	}
	return false
}
