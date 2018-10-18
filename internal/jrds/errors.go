package jrds

type RequestError struct {
	message string
}

type RequestInvalidStatusError struct {
	message string
}

type RequestAuthorizationError struct {
	message string
}

func NewRequestError(message string) *RequestError {
	return &RequestError{
		message: message,
	}
}

func NewRequestInvalidStatusError(message string) *RequestInvalidStatusError {
	return &RequestInvalidStatusError{
		message: message,
	}
}

func NewRequestAuthorizationError(message string) *RequestAuthorizationError {
	return &RequestAuthorizationError{
		message: message,
	}
}

func (e *RequestError) Error() string {
	return e.message
}

func (e *RequestInvalidStatusError) Error() string {
	return e.message
}

func (e *RequestAuthorizationError) Error() string {
	return e.message
}
