package run

import (
	"errors"
	"strconv"
	"strings"
)

// Result contains the result of executing an external command.
type Result struct {
	// ExitCode contains the exit status.
	ExitCode int

	// Stdout contains the captured stdout unless the Stdout
	// Writer was redirected.
	Stdout string

	// Stdout contains the captured stderr unless the Stderr
	// Writer was redirected.
	Stderr string
}

// Failed returns true for a non-zero exit code.
func (r *Result) Failed() bool {
	return r == nil || r.ExitCode != 0
}

// Error returns an error if the exit code is non-zero.
// As error text, the first non-empty line of stderr is used,
// otherwise an unknown error is returned.
func (r *Result) Error() error {
	switch {
	case r == nil:
		return errors.New("result is nil")
	case r.ExitCode == 0:
		return nil
	}

	var lines []string
	for _, l := range strings.Split(r.Stderr, "\n") {
		l = strings.TrimSpace(l)
		if l != "" {
			lines = append(lines, l)
			if !strings.HasSuffix(l, ":") {
				break
			}
		}
	}
	if len(lines) == 0 {
		lines = append(lines, "Unknown error: "+strconv.Itoa(r.ExitCode))
	}
	return errors.New(strings.Join(lines, " "))
}
