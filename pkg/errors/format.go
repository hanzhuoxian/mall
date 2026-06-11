package errors

import (
	"fmt"
	"io"
	"runtime"
)

// Format implements fmt.Formatter.
// %s / %v  → error message only
// %+v      → message + code + cause chain + stack trace
func (err *codeError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "[code=%d] %s", err.code, err.err.Error())
			if err.cause != nil {
				fmt.Fprintf(s, "\ncaused by: %+v", err.cause)
			}
			if err.stack != nil {
				err.stack.formatTo(s)
			}
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, err.err.Error())
	case 'q':
		fmt.Fprintf(s, "%q", err.err.Error())
	}
}

func (st *stack) formatTo(s fmt.State) {
	frames := runtime.CallersFrames([]uintptr(*st))
	for {
		f, more := frames.Next()
		fmt.Fprintf(s, "\n\t%s\n\t\t%s:%d", f.Function, f.File, f.Line)
		if !more {
			break
		}
	}
}
