// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package executil

import (
	"github.com/Azure/azure-extension-foundation/errorhelper"
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

	workingDirectory string
	environment      []string
	stdout_f         func(str string)
	stderr_f         func(str string)
	stdoutPipe       io.Reader
	stderrPipe       io.Reader
}

func NewAsyncCommand(stdout func(str string), stderr func(str string), workingDirectory string, environment []string, name string, arguments ...string) AsyncCommand {
	command := AsyncCommand{Name: name,
		Arguments:        arguments,
		stdout_f:         stdout,
		stderr_f:         stderr,
		workingDirectory: workingDirectory,
		environment:      environment,
		IsSuccessful:     false,
		IsRunning:        false}

	return command
}

func (cmd *AsyncCommand) Kill() error {
	if cmd.cmd == nil {
		return errorhelper.NewErrorWithStack("nil cmd")
	}
	return errorhelper.AddStackToError(cmd.cmd.Process.Kill())
}
