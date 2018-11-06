package sandbox

import (
	"fmt"
	"github.com/Azure/azure-automation-go-worker/internal/configuration"
	"github.com/Azure/azure-automation-go-worker/internal/tracer"
	"github.com/Azure/azure-automation-go-worker/pkg/executil"
	"os"
	"path/filepath"
)

type Sandbox struct {
	Id                         string
	workingDirectory           string
	workingDirectoryPermission os.FileMode

	command   *executil.AsyncCommand
	isRunning *bool

	commandHandler executil.AsyncCommandHandler
}

var NewSandbox = func(sandboxId string) Sandbox {
	const permission = 644
	isRunning := false
	return Sandbox{Id: sandboxId,
		command:                    nil,
		workingDirectory:           filepath.Join(configuration.GetWorkingDirectory(), sandboxId),
		workingDirectoryPermission: permission,
		commandHandler:             executil.GetAsyncCommandHandler(),
		isRunning:                  &isRunning,
	}
}

func (s *Sandbox) CreateBaseDirectory() error {
	err := os.MkdirAll(s.workingDirectory, s.workingDirectoryPermission) // TODO: change sb permission
	if err != nil {
		return err
	}

	return nil
}

func (s *Sandbox) Cleanup() error {
	// TODO: do not cleanup if sandbox crashed
	if s.command == nil {
		return fmt.Errorf("sandbox not started")
	}
	err := os.RemoveAll(s.workingDirectory)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sandbox) Start() {
	s.command = getSandboxCommand(s.Id, s.workingDirectory) // TODO: start sandbox command; this is a blocking call will need to become async
	s.isRunning = &s.command.IsRunning
	s.commandHandler.ExecuteAsync(s.command)
}

func (s *Sandbox) IsAlive() bool {
	return *s.isRunning
}

var getSandboxCommand = func(sandboxId string, workingDirectory string) *executil.AsyncCommand {
	cmd := executil.NewAsyncCommand(tracer.LogSandboxStdout, tracer.LogSandboxStderr, configuration.GetSandboxExecutablePath(), sandboxId)
	return &cmd
}
