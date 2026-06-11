package coder

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
	// ErrSuccess(0) is reserved and cannot be registered
	register(ErrBind, http.StatusBadRequest, "Error occurred while binding the request body to the struct", "")
	register(ErrValidation, http.StatusBadRequest, "Validation failed", "")
	register(ErrTokenInvalid, http.StatusUnauthorized, "Token invalid", "")
	register(ErrPageNotFound, http.StatusNotFound, "Page not found", "")

	register(ErrDatabase, http.StatusInternalServerError, "Database error", "")

	register(ErrEncrypt, http.StatusInternalServerError, "Error occurred while encrypting the user password", "")
	register(ErrSignatureInvalid, http.StatusUnauthorized, "Signature is invalid", "")
	register(ErrExpired, http.StatusUnauthorized, "Token expired", "")
	register(ErrInvalidAuthHeader, http.StatusUnauthorized, "Invalid authorization header", "")
	register(ErrMissingHeader, http.StatusUnauthorized, "Authorization header is missing", "")
	register(ErrPasswordIncorrect, http.StatusUnauthorized, "Password is incorrect", "")
	register(ErrPermissionDenied, http.StatusForbidden, "Permission denied", "")
}

func register(code errors.ErrorCode, status int, msg string, ref string) {
	errors.MustRegister(errors.NewCoder(code, status, msg, ref))
}
