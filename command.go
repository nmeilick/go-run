package run

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// Cmd represents an external command being prepared or run.
//
// A Cmd cannot be reused after calling its Run, Output or CombinedOutput
// methods.
type Cmd struct {
	exec.Cmd

	Timeout time.Duration

	outbuf bytes.Buffer
	errbuf bytes.Buffer
}

func Command(command string, args ...string) *Cmd {
	cmd := &Cmd{
		Cmd: *exec.Command(command, args...),
	}
	cmd.Stdout = &cmd.outbuf
	cmd.Stderr = &cmd.errbuf
	return cmd
}

func (cmd *Cmd) WithDir(dir string) *Cmd {
	cmd.Dir = dir
	return cmd
}

func (cmd *Cmd) WithEnv(env []string) *Cmd {
	cmd.Env = env
	return cmd
}

func (cmd *Cmd) EmptyEnv() *Cmd {
	cmd.Env = []string{}
	return cmd
}

func (cmd *Cmd) WithStdin(stdin io.Reader) *Cmd {
	cmd.Stdin = stdin
	return cmd
}

func (cmd *Cmd) WithStdout(stdout io.Writer) *Cmd {
	cmd.Stdout = stdout
	return cmd
}

func (cmd *Cmd) WithStderr(stderr io.Writer) *Cmd {
	cmd.Stderr = stderr
	return cmd
}

func (cmd *Cmd) WithTimeout(d time.Duration) *Cmd {
	cmd.Timeout = d
	return cmd
}

func (cmd *Cmd) filterEnv(patterns []string, removeMatched bool) *Cmd {
	newEnv := []string{}
	for _, line := range cmd.Env {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) < 1 {
			continue
		}

		matched := matchEnv(parts[0], patterns)
		if removeMatched {
			matched = !matched
		}

		if matched {
			newEnv = append(newEnv, line)
		}
	}
	cmd.Env = newEnv
	return cmd
}

// RemoveEnv removes env variables whose name matches any of the given wildcard patterns.
func (cmd *Cmd) RemoveEnv(patterns ...string) *Cmd {
	return cmd.filterEnv(patterns, true)
}

// LimitEnv removes all env variables whose name does not match any of the given wildcard patterns.
func (cmd *Cmd) LimitEnv(patterns ...string) *Cmd {
	return cmd.filterEnv(patterns, false)
}

// Run runs the command and returns a result structure, or an error if the command could
// not be run or execution timed out.
func (cmd *Cmd) Run() (*Result, error) {
	if err := cmd.Cmd.Start(); err != nil {
		return nil, err
	}

	errC := make(chan error, 1)
	go func() { errC <- cmd.Wait() }()
	tmoutC := make(chan struct{}, 1)
	if cmd.Timeout > 0 {
		time.AfterFunc(cmd.Timeout, func() { close(tmoutC) })
	}

	var err error

	select {
	case err = <-errC:
	case <-tmoutC:
		return nil, errors.New("timeout")
	}

	r := &Result{}
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			r.ExitCode = ws.ExitStatus()
		} else {
			return nil, err
		}
	} else {
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		r.ExitCode = ws.ExitStatus()
	}

	r.Stdout = string(cmd.outbuf.Bytes())
	r.Stderr = string(cmd.errbuf.Bytes())
	return r, nil
}
