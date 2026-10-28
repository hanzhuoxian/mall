package coder

import (
	"net/http"
	"testing"

	"github.com/hanzhuoxian/mall/pkg/errors"
)

func TestAuthenticationErrorMapping(t *testing.T) {
	for _, tc := range []struct {
		code   errors.ErrorCode
		status int
	}{{ErrSignatureInvalid, http.StatusUnauthorized}, {ErrTokenInvalid, http.StatusUnauthorized}, {ErrPermissionDenied, http.StatusForbidden}, {ErrCaptchaInvalid, http.StatusBadRequest}} {
		got := errors.ParseCoder(errors.New(tc.code, "test"))
		if got.HTTPStatus() != tc.status {
			t.Errorf("code %d: got HTTP %d, want %d", tc.code, got.HTTPStatus(), tc.status)
		}
	}
}
