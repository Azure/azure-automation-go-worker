package gpg

import (
	"errors"
	"github.com/Azure/guestaction-extension-linux/pkg/executil"
	"os"
	"testing"
)

type CommandHandlerMock struct {
	commandToExecute func(*executil.Command) ()
}

type TestWithAssert testing.T

func (cmdH CommandHandlerMock) Execute(command *executil.Command) {
	cmdH.commandToExecute(command)
}

func SuccessfulExecution(cmd *executil.Command) () {
	(*cmd).ExitCode = 0
	(*cmd).IsSuccessful = true
	(*cmd).CommandError = nil
}

func FailedExecution(cmd *executil.Command) () {
	(*cmd).ExitCode = 1
	(*cmd).IsSuccessful = false
	(*cmd).CommandError = errors.New("expected test error")
}

func VerifyParametersInRightSpot(cmd *executil.Command, t *TestWithAssert, signedFilePath string, outputFilePath string, keyringpath string) () {
	t.AssertStringsAreEqual(cmd.Name, GPG)
	for index, value := range cmd.Arguments {
		switch value {
		case signedFilePath:
			t.AssertIntAreEqual(index, len(cmd.Arguments)-1)
		case outputFilePath:
			t.AssertStringsAreEqual(cmd.Arguments[index-1], GPG_OUTPUT_OPTION)
		case keyringpath:
			t.AssertStringsAreEqual(cmd.Arguments[index-1], GPG_KEYRING_OPTION)
			t.AssertStringsAreEqual(cmd.Arguments[index-2], GPG_NO_DEFAULT_KEYRING_OPTION)
		}
	}
	(*cmd).ExitCode = 0
	(*cmd).IsSuccessful = true
	(*cmd).CommandError = nil
}

func TestGpgValidationSucceedsMock(t *testing.T) {
	cmdHandler = CommandHandlerMock{commandToExecute: SuccessfulExecution}
	success, err := VerifySignature("mockPath", "mockPath", []string{"keyring1", "keyring2"})
	if err != nil || !success {
		t.Fatal(err.Error())
	}
}

func TestVerifyFailsWithExecutionFails(t *testing.T) {
	cmdHandler = CommandHandlerMock{commandToExecute: FailedExecution}
	success, err := VerifySignature("mockPath", "mockPath", []string{"keyring1", "keyring2"})
	_, typeMatched := err.(*GpgExecuteError)
	if typeMatched && !success {
		return
	}
	t.Fatal("Error was of unexpected type")
}

func TestGpgValidationSucceeds(t *testing.T) {
	// skip this test can't use real gpg keyring
	t.SkipNow()

	signedFilePath := "./testresources/helloworld.py.asc"
	outputFilePath := "./testoutput/helloworld.py"
	keyringPath := "./testresources/testkeyring.gpg"
	FailIfFileNotExist(signedFilePath, t)
	FailIfFileNotExist(keyringPath, t)

	if t.Failed() {
		t.Fatal("Cannot find required files. Test cannot proceed")
	}
	success, err := VerifySignature(signedFilePath, outputFilePath, []string{keyringPath})
	if err != nil || !success {
		t.Fatal(err.Error())
	}
}

func TestGpgKeyringPathEmptyThrowsError(t *testing.T) {
	cmdHandler = CommandHandlerMock{commandToExecute: SuccessfulExecution}
	success, err := VerifySignature("mockPath", "mockPath", nil)
	_, typeMatched := err.(*KeyringNotConfiguredError)
	if typeMatched && !success {
		return
	}
	t.Fatal("Error was of unexpected type")
}

func TestParametesAreProperlyPassed(t *testing.T) {
	signedFilePath := "./testresources/helloworld.py.asc"
	outputFilePath := "./testoutput/helloworld.py"
	keyringPath := "./testresources/testkeyring.gpg"
	tt := TestWithAssert(*t)
	cmdHandler = CommandHandlerMock{commandToExecute: func(command *executil.Command) {
		VerifyParametersInRightSpot(command, &tt, signedFilePath, outputFilePath, keyringPath)
	}}
	success, err := VerifySignature(signedFilePath, outputFilePath, []string{keyringPath})
	if err != nil || !success {
		t.Fatal(err.Error())
	}
}

func FailIfFileNotExist(filepath string, t *testing.T) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		t.Fail()
	}
}

func (t *TestWithAssert) AssertStringsAreEqual(actual string, expected string) {
	if actual != expected {
		t.Errorf("Values are not equal expected: %s, actual: %s", expected, actual)
	}
}

func (t *TestWithAssert) AssertIntAreEqual(actual int, expected int) {
	if actual != expected {
		t.Errorf("Values are not equal expected: %v, actual: %v", expected, actual)
	}
}
