// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package runtime

import (
	"github.com/Azure/azure-automation-go-worker/internal/jrds"
	"github.com/Azure/azure-automation-go-worker/internal/tracer"
	"github.com/Azure/azure-automation-go-worker/pkg/executil"
	"os"
	"path/filepath"
	"time"
)

type Runtime struct {
	runbook          Runbook
	language         Language
	jobData          jrds.JobData
	workingDirectory string
}

func NewRuntime(language Language, runbook Runbook, jobData jrds.JobData, workingDirectory string) Runtime {
	return Runtime{
		runbook:          runbook,
		language:         language,
		jobData:          jobData,
		workingDirectory: workingDirectory}
}

func (runtime *Runtime) Initialize() error {
	runbookPath := getRunbookPathOnDisk(runtime.workingDirectory, runtime.runbook)
	err := writeRunbookToDisk(runbookPath, runtime.runbook)
	if err != nil {
		return err
	}

	return nil
}

func (runtime *Runtime) IsSupported() bool {
	return runtime.language.interpreter.isSupported()
}

func (runtime *Runtime) StartRunbook() {
	arguments := append(runtime.language.interpreter.arguments, getRunbookPathOnDisk(runtime.workingDirectory, runtime.runbook))
	handler := executil.GetAsyncCommandHandler()
	cmd := executil.NewAsyncCommand(
		tracer.LogDebugTrace,
		tracer.LogDebugTrace,
		runtime.workingDirectory,
		nil,
		runtime.language.interpreter.commandName,
		arguments...)
	handler.ExecuteAsync(&cmd)

	for cmd.IsRunning {
		time.Sleep(10 * time.Millisecond)
	}

}

var getRunbookPathOnDisk = func(workingDirectory string, runbook Runbook) string {
	return filepath.Join(workingDirectory, runbook.FileName)
}

var writeRunbookToDisk = func(path string, runbook Runbook) error {
	const permission = 0640
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, permission)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte(runbook.Definition))
	if err != nil {
		return err
	}

	file.Close()
	return nil
}
