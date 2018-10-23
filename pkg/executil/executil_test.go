package executil

import (
	"strings"
	"testing"
)

const (
	cmd_echo       = "echo"
	cmd_ls         = "ls"
	cmd_unknowncmd = "unknowncmd"

	test_string          = "this is a test"
	invalid_ls_parameter = "-deeee"
)

func TestRunCommandReturnsSuccessfulExitCode(t *testing.T) {
	cmd := Command{Name: cmd_echo, Arguments: []string{test_string}}
	handler := GetCommandHandler()
	handler.Execute(&cmd)

	if !cmd.IsSuccessful {
		t.Errorf("unexpected command unsuccessful")
	}
	if cmd.ExitCode != 0 {
		t.Errorf("unexpected non zero exit code [%v]", cmd.CommandError)
	}
}

func TestRunCommandReturnsUnSuccessfulExitCode(t *testing.T) {
	cmd := NewCommand(cmd_ls, invalid_ls_parameter)
	handler := GetCommandHandler()
	handler.Execute(&cmd)

	if !cmd.IsSuccessful {
		t.Errorf("unexpected command unsuccessful")
	}
	if cmd.ExitCode == 0 {
		t.Error("unexpected zero exit code ")
	}
}

func TestExecuteUnknownCommandReturnsUnSuccessfulExitCode(t *testing.T) {
	cmd := NewCommand(cmd_unknowncmd)
	handler := GetCommandHandler()
	handler.Execute(&cmd)

	if cmd.IsSuccessful {
		t.Errorf("unexpected command successful")
	}
	if cmd.CommandError == nil {
		t.Errorf("unexpected empty command error")
	}
}

func TestStdout(t *testing.T) {
	cmd := NewCommand(cmd_echo, test_string)
	handler := GetCommandHandler()
	handler.Execute(&cmd)

	if !cmd.IsSuccessful {
		t.Errorf("unexpected command unsuccessful")
	}
	if cmd.ExitCode != 0 {
		t.Errorf("unexpected non zero exit code [%v]", cmd.CommandError)
	}

	if strings.TrimSpace(cmd.Stdout.String()) != test_string {
		t.Error("unexpected invalid stdout")
	}
}

func TestStderr(t *testing.T) {
	cmd := NewCommand(cmd_ls, invalid_ls_parameter)
	handler := GetCommandHandler()
	handler.Execute(&cmd)

	if !cmd.IsSuccessful {
		t.Errorf("unexpected command unsuccessful")
	}
	if cmd.ExitCode == 0 {
		t.Errorf("unexpected non zero exit code [%v]", cmd.CommandError)
	}

	if !strings.Contains(cmd.Stderr.String(), "unknown option") && !strings.Contains(cmd.Stderr.String(), "invalid option") {
		t.Error("unexpected invalid stderr")
	}

}
