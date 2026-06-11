package errors

import "fmt"

type ErrorCode int

var DefaultCode ErrorCode = 1

type codeError struct {
	err   error
	code  ErrorCode
	cause error
	*stack
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

func New(code ErrorCode, format string, args ...any) error {
	return &codeError{
		code:  code,
		err:   fmt.Errorf(format, args...),
		stack: callers(),
	}
}

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
