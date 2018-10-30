package sandbox

import (
	"fmt"
	"github.com/Azure/azure-automation-go-worker/internal/configuration"
	"github.com/Azure/azure-automation-go-worker/pkg/executil"
	"os"
	"path/filepath"
)

type Sandbox struct {
	Id                         string
	workingDirectory           string
	workingDirectoryPermission os.FileMode

	command *executil.Command
	isAlive bool

	commandHandler executil.CommandHandler
}

var NewSandbox = func(sandboxId string) Sandbox {
	const permission = 644
	return Sandbox{Id: sandboxId,
		command:                    nil,
		workingDirectory:           filepath.Join(configuration.GetWorkingDirectory(), sandboxId),
		workingDirectoryPermission: permission,
		commandHandler:             executil.GetCommandHandler(),
		isAlive:                    false,
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
	s.isAlive = true
	s.commandHandler.Execute(s.command)
	s.isAlive = false
}

func (s *Sandbox) GetOutput() (string, error) {
	if s.command == nil {
		return "", fmt.Errorf("sandbox process not started")
	}

	return s.command.Stdout.String(), nil // TODO: this will need to be refactored so we can get output async
}

func (s *Sandbox) GetErrorOutput() (string, error) {
	if s.command == nil {
		return "", fmt.Errorf("sandbox process not started")
	}

	return s.command.Stderr.String(), nil // TODO: this will need to be refactored so we can get output async
}

func (s *Sandbox) IsAlive() bool {
	return s.isAlive
}

var getSandboxCommand = func(sandboxId string, workingDirectory string) *executil.Command {
	cmd := executil.NewCommand(configuration.GetSandboxExecutablePath(), sandboxId)
	return &cmd
}
