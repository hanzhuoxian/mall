package errors

import (
	"fmt"
	"net/http"
	"sync"
)

var unknownCoder = defaultCoder{DefaultCode, http.StatusInternalServerError, "An internal server error occurred", ""}

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
	// C refers to the integer code of the ErrCode.
	C ErrorCode

	// HTTP status that should be used for the associated error code.
	HTTP int

	// External (user) facing error text.
	Ext string

	// Ref specify the reference document.
	Ref string
}

// pkg/errors/code.go 新增
func NewCoder(code ErrorCode, httpStatus int, msg, ref string) Coder {
	return defaultCoder{code, httpStatus, msg, ref}
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

func Register(coder Coder) {
	if coder.Code() == 0 {
		panic("code `0` is reserved ")
	}
	codeMux.Lock()
	defer codeMux.Unlock()
	codes[coder.Code()] = coder
}

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
