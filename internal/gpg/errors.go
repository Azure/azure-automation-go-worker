// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package gpg

type KeyringNotConfiguredError struct {
	message string
}

type GpgExecuteError struct {
	message string
}

func NewKeyringNotConfiguredError(message string) *KeyringNotConfiguredError {
	return &KeyringNotConfiguredError{message: message}
}

func NewGpgExecuteError(message string) *GpgExecuteError {
	return &GpgExecuteError{message: message}
}

func (err *KeyringNotConfiguredError) Error() string {
	return err.message
}

func (err *GpgExecuteError) Error() string {
	return err.message
}
