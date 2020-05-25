package run

// Start executed the given external command using default settings and returns
// the command instance or an error if the start was unsuccesful.
func Start(name string, args ...string) (*Cmd, error) {
	cmd := Command(name, args...)
	return cmd, cmd.Start()
}
