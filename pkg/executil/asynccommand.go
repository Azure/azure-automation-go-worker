// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package executil

import (
	"io"
	"os/exec"
)

type AsyncCommand struct {
	Name      string
	Arguments []string

	ExitCode int

	IsRunning    bool
	IsSuccessful bool
	CommandError error

	cmd *exec.Cmd

	stdout_f   func(str string)
	stderr_f   func(str string)
	stdoutPipe io.Reader
	stderrPipe io.Reader
}

func NewAsyncCommand(stdout func(str string), stderr func(str string), name string, arguments ...string) AsyncCommand {
	return AsyncCommand{Name: name,
		Arguments:    arguments,
		stdout_f:     stdout,
		stderr_f:     stderr,
		IsSuccessful: false,
		IsRunning:    false}
}
