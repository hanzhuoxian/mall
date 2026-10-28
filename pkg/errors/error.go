package errors

import "fmt"

// ErrorCode is an application-specific error identifier.
type ErrorCode int

// DefaultCode is the default application error code.
var DefaultCode ErrorCode = 1

type codeError struct {
	err   error
	cause error
	*stack
	code ErrorCode
}

func (err *codeError) Error() string {
	return err.err.Error()
}

func (err *codeError) Cause() error {
	return err.cause
}

func (err *codeError) Unwrap() error {
	return err.cause
}

// New creates an error with a code, formatted message, and captured stack.
func New(code ErrorCode, format string, args ...any) error {
	return &codeError{
		code:  code,
		err:   fmt.Errorf(format, args...),
		stack: callers(),
	}
}

// Wrap adds a code, formatted message, and stack to err.
func Wrap(err error, code ErrorCode, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return &codeError{
		code:  code,
		err:   fmt.Errorf(format, args...),
		cause: err,
		stack: callers(),
	}
}
