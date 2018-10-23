package executil

import "bytes"

type Command struct {
	Name      string
	Arguments []string

	Stdout   bytes.Buffer
	Stderr   bytes.Buffer
	ExitCode int

	IsSuccessful bool
	CommandError error
}

func NewCommand(name string, arguments ...string) Command {
	return Command{Name: name, Arguments: arguments, IsSuccessful: false}
}
