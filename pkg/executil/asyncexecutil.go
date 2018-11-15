// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package executil

import (
	"bufio"
	"io"
	"os/exec"
	"syscall"
)

type AsyncHandler interface {
	ExecuteAsync(command *AsyncCommand)
}

type AsyncCommandHandler struct {
}

func (h AsyncCommandHandler) ExecuteAsync(cmd *AsyncCommand) error {
	return executeAsyncCommand(cmd)
}

func GetAsyncCommandHandler() AsyncCommandHandler {
	return AsyncCommandHandler{}
}

func executeAsyncCommand(command *AsyncCommand) error {
	cmd := exec.Command(command.Name, command.Arguments...)
	cmd.Env = command.environment
	cmd.Dir = command.workingDirectory

	command.stdoutPipe, _ = cmd.StdoutPipe()
	command.stderrPipe, _ = cmd.StderrPipe()
	command.cmd = cmd

	err := command.cmd.Start()
	if err != nil {
		return err
	}

	command.IsRunning = true
	go startAndMonitorCommand(command)
	return nil
}

func startAndMonitorCommand(command *AsyncCommand) {
	// scan stdout
	if command.stdoutPipe != nil && command.stdout_f != nil {
		newScanner(command.stdout_f, command.stdoutPipe)
	}

	// scan stderr
	if command.stderrPipe != nil && command.stderr_f != nil {
		newScanner(command.stderr_f, command.stderrPipe)
	}

	// wait for command to complete
	err := command.cmd.Wait()
	command.IsRunning = false

	// set command error and exit code
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

func newScanner(print func(str string), reader io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		print(scanner.Text())
	}
	return scanner
}
