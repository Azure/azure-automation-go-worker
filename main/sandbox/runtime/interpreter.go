// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package runtime

import (
	"github.com/Azure/azure-automation-go-worker/pkg/executil"
	"runtime"
)

type Interpreter struct {
	language string

	commandName string
	arguments   []string
}

var getPowerShellInterpreter = func() Interpreter {
	commandName := "pwsh"

	if runtime.GOOS == "windows" {
		commandName = "powershell"
	}

	return Interpreter{
		language:    "PowerShell",
		commandName: commandName,
		arguments:   []string{"-File"}}
}

var getPython2Interpreter = func() Interpreter {
	return Interpreter{
		language:    "Python2",
		commandName: "pythonExtension",
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
