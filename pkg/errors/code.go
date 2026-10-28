package errors

import (
	"fmt"
	"net/http"
	"sync"
)

var unknownCoder = defaultCoder{
	Ext:  "An internal server error occurred",
	C:    DefaultCode,
	HTTP: http.StatusInternalServerError,
}

// Coder maps an application error code to its HTTP response details.
type Coder interface {
	// HTTP status that should be used for the associated error code.
	HTTPStatus() int

	// External (user) facing error text.
	String() string

	// Reference returns the detail documents for user.
	Reference() string

	// Code returns the code of the coder
	Code() ErrorCode
}

type defaultCoder struct {
	Ext  string
	Ref  string
	C    ErrorCode
	HTTP int
}

// NewCoder returns a Coder with the supplied code, HTTP status, message, and reference.
func NewCoder(code ErrorCode, httpStatus int, msg, ref string) Coder {
	return defaultCoder{Ext: msg, Ref: ref, C: code, HTTP: httpStatus}
}

func (coder defaultCoder) Code() ErrorCode {
	return coder.C
}

func (coder defaultCoder) String() string {
	return coder.Ext
}

func (coder defaultCoder) Reference() string {
	return coder.Ref
}

func (coder defaultCoder) HTTPStatus() int {
	if coder.HTTP == 0 {
		return http.StatusInternalServerError
	}
	return coder.HTTP
}

// IsCode reports whether err contains the specified application error code.
func IsCode(err error, code ErrorCode) bool {
	if v, ok := err.(*codeError); ok {
		if v.code == code {
			return true
		}
		if v.cause != nil {
			return IsCode(err, code)
		}
	}

	return false
}

// ParseCoder returns the registered coder for err, or the default coder.
func ParseCoder(err error) Coder {
	if err == nil {
		return nil
	}
	if v, ok := err.(*codeError); ok {
		if coder, ok := codes[v.code]; ok {
			return coder
		}
	}
	return unknownCoder
}

var codes = map[ErrorCode]Coder{}
var codeMux sync.Mutex

// Register associates coder with its error code.
func Register(coder Coder) {
	if coder.Code() == 0 {
		panic("code `0` is reserved ")
	}
	codeMux.Lock()
	defer codeMux.Unlock()
	codes[coder.Code()] = coder
}

// MustRegister associates coder with its error code and panics on duplicates.
func MustRegister(coder Coder) {
	if coder.Code() == 0 {
		panic("code `0` is reserved ")
	}
	codeMux.Lock()
	defer codeMux.Unlock()
	if _, ok := codes[coder.Code()]; ok {
		panic(fmt.Sprintf("code: %d already exist", coder.Code()))
	}
	codes[coder.Code()] = coder
}

func init() {
	Register(unknownCoder)
}
