package run

// Run executed the given external command using default settings and returns
// a Result containing exit status and captured output, or an error if execution
// was unsuccesful or timed out.
func Run(name string, args ...string) (*Result, error) {
	return Command(name, args...).Run()
}
