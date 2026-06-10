package errors

type codeError struct {
	err   error
	code  int
	cause error
	*stack
}
