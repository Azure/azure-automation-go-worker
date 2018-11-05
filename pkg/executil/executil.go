// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package executil

import (
	"os/exec"
	"syscall"
)

type Handler interface {
	Execute(command *Command)
}

type CommandHandler struct {
}

func (h CommandHandler) Execute(cmd *Command) {
	executeCommand(cmd)
}

func GetCommandHandler() CommandHandler {
	return CommandHandler{}
}

func executeCommand(command *Command) {
	cmd := exec.Command(command.Name, command.Arguments...)
	cmd.Stdout = &command.Stdout
	cmd.Stderr = &command.Stderr

	err := cmd.Run()
	exitError, _ := err.(*exec.ExitError)
	if err != nil && exitError == nil {
		command.CommandError = err
		return
	}
	if exitError != nil {
		waitStatus := exitError.Sys().(syscall.WaitStatus)
		command.ExitCode = waitStatus.ExitStatus()
	} else {
		command.ExitCode = 0
	}

	command.IsSuccessful = true
}
