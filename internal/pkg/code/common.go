package code

import (
	"net/http"

	"github.com/hanzhuoxian/mall/pkg/errors"
)

const (
	ErrSuccess errors.ErrorCode = 0
	ErrBind                     = iota + 1
	ErrValidation
	ErrTokenInvalid
	ErrPageNotFound
)

const (
	ErrDatabase errors.ErrorCode = iota + 100100
)

const (
	ErrEncrypt errors.ErrorCode = iota + 100200
	ErrSignatureInvalid
	ErrExpired
	ErrInvalidAuthHeader
	ErrMissingHeader
	ErrPasswordIncorrect
	ErrPermissionDenied
)

func init() {
	register(ErrSuccess, http.StatusOK, "", "")
}

func register(code errors.ErrorCode, status int, msg string, ref string) {
	errors.MustRegister(errors.NewCoder(code, status, msg, ref))
}
