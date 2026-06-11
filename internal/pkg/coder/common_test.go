package coder

import "testing"

func TestErr(t *testing.T) {
	t.Fatal(ErrSignatureInvalid)
}
