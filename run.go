package run

// Run executed the given external command using default settings and returns
// a Result containing exit status and captured output, or an error if execution
// was unsuccesful or timed out.
func Run(name string, args ...string) *Result {
	r, err := Command(name, args...).Run()
	if err == nil {
		return r
	}

	return &Result{
		ExitCode: 255,
		Stdout:   "",
		Stderr:   err.Error(),
	}
}
