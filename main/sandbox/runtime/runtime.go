// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package runtime

import (
	"github.com/Azure/azure-automation-go-worker/internal/jrds"
	"github.com/Azure/azure-automation-go-worker/pkg/executil"
	"github.com/Azure/azure-extension-foundation/errorhelper"
	"os"
	"path/filepath"
)

var runbookError string // TODO : temporary

type Runtime struct {
	runbook          Runbook
	language         Language
	jobData          jrds.JobData
	workingDirectory string

	runbookCmd       *executil.AsyncCommand
	isRunbookRunning *bool
}

func NewRuntime(language Language, runbook Runbook, jobData jrds.JobData, workingDirectory string) Runtime {
	false := false
	return Runtime{
		runbook:          runbook,
		language:         language,
		jobData:          jobData,
		workingDirectory: workingDirectory,
		isRunbookRunning: &false}
}

func (runtime *Runtime) Initialize() error {
	runbookPath := getRunbookPathOnDisk(runtime.workingDirectory, runtime.runbook)
	err := writeRunbookToDisk(runbookPath, runtime.runbook)
	return errorhelper.AddStackToError(err)
}

func (runtime *Runtime) IsSupported() bool {
	return runtime.language.interpreter.isSupported()
}

func (runtime *Runtime) StartRunbookAsync(streamHandler func(string)) {
	arguments := append(runtime.language.interpreter.arguments, getRunbookPathOnDisk(runtime.workingDirectory, runtime.runbook))
	handler := executil.GetAsyncCommandHandler()
	cmd := executil.NewAsyncCommand(
		streamHandler,
		rbStderr,
		runtime.workingDirectory,
		nil,
		runtime.language.interpreter.commandName,
		arguments...)
	handler.ExecuteAsync(&cmd)

	runtime.runbookCmd = &cmd
	runtime.isRunbookRunning = &cmd.IsRunning
}

func (runtime *Runtime) IsRunbookRunning() bool {
	return *runtime.isRunbookRunning
}

func (runtime *Runtime) StopRunbook() error {
	if runtime.runbookCmd == nil {
		return nil
	}

	return runtime.runbookCmd.Kill()
}

func (runtime *Runtime) ExitCode() int {
	return runtime.runbookCmd.ExitCode
}

func (runtime *Runtime) GetRunbookError() string {
	return runbookError
}

func (runtime *Runtime) IsRunbookExecutionSuccessful() bool {
	return runtime.runbookCmd.ExitCode == 0
}

var getRunbookPathOnDisk = func(workingDirectory string, runbook Runbook) string {
	return filepath.Join(workingDirectory, runbook.FileName)
}

var writeRunbookToDisk = func(path string, runbook Runbook) error {
	const permission = 0640
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, permission)
	if err != nil {
		return errorhelper.AddStackToError(err)
	}

	_, err = file.Write([]byte(runbook.Definition))
	if err != nil {
		return errorhelper.AddStackToError(err)
	}

	file.Close()
	return nil
}

var rbStderr = func(message string) {
	runbookError += message
}
