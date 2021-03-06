// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package runtime

import (
	"github.com/Azure/azure-automation-go-worker/pkg/executil"
	"runtime"
)

const windows = "windows"

type Interpreter struct {
	language string

	commandName string
	arguments   []string
}

var getPowerShellInterpreter = func() Interpreter {
	commandName := "pwsh"

	if runtime.GOOS == windows {
		commandName = "powershell"
	}

	return Interpreter{
		language:    "PowerShell",
		commandName: commandName,
		arguments:   []string{"-File"}}
}

var getPython2Interpreter = func() Interpreter {
	python := "python2"
	if runtime.GOOS == windows {
		python = "C:\\python27\\python.exe"
	}

	return Interpreter{
		language:    "Python2",
		commandName: python,
		arguments:   []string{}}
}

var getPython3Interpreter = func() Interpreter {
	return Interpreter{
		language:    "Python3",
		commandName: "python3",
		arguments:   []string{}}
}

var getBashInterpreter = func() Interpreter {
	return Interpreter{
		language:    "Bash",
		commandName: "bash",
		arguments:   []string{}}
}

func (i *Interpreter) isSupported() bool {
	handler := executil.GetCommandHandler()
	command := executil.NewCommand(i.commandName)
	handler.Execute(&command)
	return command.IsSuccessful
}
