package sandbox

import (
	"fmt"
	"github.com/Azure/azure-automation-go-worker/internal/configuration"
	"github.com/Azure/azure-automation-go-worker/internal/tracer"
	"github.com/Azure/azure-automation-go-worker/pkg/executil"
	"github.com/Azure/azure-extension-foundation/errorhelper"
	"os"
	"path/filepath"
	"strings"
)

const (
	sandboxWorkingDirectoryName = "sandboxes"
)

type Sandbox struct {
	Id               string
	workingDirectory string

	isRunning *bool
	faulted   *bool

	command        executil.AsyncCommand
	commandHandler executil.AsyncCommandHandler
}

var NewSandbox = func(sandboxId string) Sandbox {
	isRunning := false
	faulted := false
	return Sandbox{
		Id:               sandboxId,
		workingDirectory: filepath.Join(configuration.GetWorkingDirectory(), sandboxWorkingDirectoryName, sandboxId),
		isRunning:        &isRunning,
		faulted:          &faulted,
		commandHandler:   executil.GetAsyncCommandHandler(),
	}
}

func (s *Sandbox) CreateBaseDirectory() error {
	const permission = 0750
	err := os.MkdirAll(s.workingDirectory, permission) // TODO: change sb permission
	if err != nil {
		return errorhelper.AddStackToError(err)
	}

	return nil
}

func (sandbox *Sandbox) Start() error {
	// start sandbox
	command, err := getSandboxCommand(tracer.LogSandboxStdout, tracer.LogSandboxStderr, sandbox.Id, sandbox.workingDirectory)
	if err != nil {
		return nil
	}

	sandbox.isRunning = &command.IsRunning
	sandbox.faulted = &command.IsSuccessful
	err = sandbox.commandHandler.ExecuteAsync(command)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sandbox) Cleanup() error {
	if *s.isRunning {
		return fmt.Errorf("sandbox is running")
	}

	tracer.LogWorkerSandboxProcessExited(s.Id, 0, 0)

	// do not clean if sandbox faulted
	if *s.faulted {
		return nil
	}

	err := os.RemoveAll(s.workingDirectory)
	if err != nil {
		return errorhelper.AddStackToError(err)
	}

	return nil
}

func (s *Sandbox) IsAlive() bool {
	return *s.isRunning
}

var getSandboxCommand = func(stdout func(str string), stderr func(str string), sandboxId string, workingDirectory string) (*executil.AsyncCommand, error) {
	environ, err := getSandboxProcessEnvrion(workingDirectory, os.Environ(), configuration.GetConfiguration())
	if err != nil {
		return nil, err
	}
	cmd := executil.NewAsyncCommand(stdout, stderr, workingDirectory, environ, configuration.GetSandboxExecutablePath(), sandboxId)
	return &cmd, nil
}

var getSandboxProcessEnvrion = func(workingDirectory string, environ []string, config configuration.Configuration) ([]string, error) {
	config.WorkerWorkingDirectory = workingDirectory
	config.Component = configuration.Component_sandbox
	serialized, err := configuration.SerializeConfiguration(&config)
	if err != nil {
		return []string{}, err
	}

	for i, v := range environ {
		if strings.Contains(v, configuration.EnvironmentConfigurationKey) {
			environ[i] = fmt.Sprintf("%v=%v", configuration.EnvironmentConfigurationKey, string(serialized))
			break
		}
	}
	return environ, nil
}
