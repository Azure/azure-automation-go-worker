package gpg

type KeyringNotConfiguredError struct {
	message string
}
type GpgExecuteError struct {
	message string
}
type ValidationFailedForAllKeyringsError struct {
	message string
}

func NewKeyringNotConfiguredError(message string) (*KeyringNotConfiguredError) {
	return &KeyringNotConfiguredError{message: message}
}

func NewGpgExecuteError(message string) (*GpgExecuteError) {
	return &GpgExecuteError{message: message}
}

func NewValidationFailedForAllKeyringsError(message string) (*ValidationFailedForAllKeyringsError) {
	return &ValidationFailedForAllKeyringsError{message: message}
}

func (err *KeyringNotConfiguredError) Error() (string) {
	return err.message
}

func (err *GpgExecuteError) Error() (string) {
	return err.message
}

func (err *ValidationFailedForAllKeyringsError) Error() (string) {
	return err.message
}
