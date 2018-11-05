// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package gpg

import (
	"fmt"
	"github.com/Azure/azure-automation-go-worker/pkg/executil"
)

const (
	GPG                           = "gpg"
	GPG_BATCH_OPTION              = "--batch"
	GPG_DECRYPT_OPTION            = "--decrypt"
	GPG_KEYRING_OPTION            = "--keyring"
	GPG_NO_DEFAULT_KEYRING_OPTION = "--no-default-keyring"
	GPG_OUTPUT_OPTION             = "--output"
	GPG_YES_OPTION                = "--yes"
	GPG_DEFAULT_KEYRING_PATH      = "" // TODO: Get this value from configuration
)

var cmdHandler executil.Handler = executil.GetCommandHandler()

// Verifies a files's signature
func IsSignatureValid(signedFilePath string, outputFilePath string, keyrings []string) (bool, error) {
	if (len(keyrings) == 0) || (len(keyrings) == 1 && keyrings[0] == GPG_DEFAULT_KEYRING_PATH) {
		return false, NewKeyringNotConfiguredError("GPG kerying path was empty")
	}

	for _, keyringPath := range keyrings {
		if keyringPath == "" || keyringPath == GPG_DEFAULT_KEYRING_PATH {
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
		err := cmd.CommandError

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
	//No GPG keyring was able to verify the signed file
	return false, nil
}
