package gpg

import (
	"fmt"
	"github.com/Azure/guestaction-extension-linux/pkg/executil"
)

const (
	GPG                           = "gpg"
	GPG_BATCH_OPTION              = "--batch"
	GPG_DECRYPT_OPTION            = "--decrypt"
	GPG_KEYRING_OPTION            = "--keyring"
	GPG_NO_DEFAULT_KEYRING_OPTION = "--no-default-keyring"
	GPG_OUTPUT_OPTION             = "--output"
	GPG_YES_OPTION                = "--yes"
)

var cmdHandler executil.Handler = executil.GetCommandHandler()

// Overwrite this value to mock the test

func VerifySignature(signedFilePath string, outputFilePath string, keyrings []string) (success bool, err error) {
	/* Verifies a files's signature
	Returns:
		err == nil if there is an error
		err != nil an error was encountered
	*/
	if (len(keyrings) == 0) || (len(keyrings) == 1 && keyrings[0] == "TODO: get default keyring path") {
		return false, NewKeyringNotConfiguredError("GPG kerying path was empty")
	}

	for _, keyringPath := range keyrings {
		if keyringPath == "" || keyringPath == "TODO: get default keyring path" {
			continue
		}
		args := make([]string, 0, 10)

		args = append(args, GPG_BATCH_OPTION, GPG_YES_OPTION, GPG_DECRYPT_OPTION)

		if keyringPath != "" {
			args = append(args, GPG_NO_DEFAULT_KEYRING_OPTION, GPG_KEYRING_OPTION, keyringPath)
		}
		args = append(args, GPG_OUTPUT_OPTION, outputFilePath, signedFilePath)
		cmd := executil.NewCommand(GPG, args...)
		// execute the command
		cmdHandler.Execute(&cmd)

		ret := cmd.ExitCode
		err = cmd.CommandError

		if err != nil {
			// TODO: trace signature validation success
			return false, NewGpgExecuteError(err.Error())
		}
		if ret != 0 {
			return false, NewGpgExecuteError(fmt.Sprintf("Gpg execution returned code: %v", ret))
		}
		return true, nil
		// TODO: trace signature validation failure
	}
	return false, NewValidationFailedForAllKeyringsError("No GPG keyring was able to verify the signed file")
}
