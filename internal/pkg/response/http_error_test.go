package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFailAdaptsHTTPAndGRPCErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cause := errors.New("private database detail")
	httpErr := &HTTPError{Status: http.StatusConflict, Message: "公开提示", Cause: cause}
	if !errors.Is(httpErr, cause) {
		t.Fatal("HTTPError must preserve its cause")
	}
	for _, tc := range []struct {
		err          error
		message      string
		status, code int
	}{
		{fmt.Errorf("wrapped: %w", httpErr), "公开提示", http.StatusConflict, http.StatusConflict},
		{status.Error(codes.NotFound, "missing user"), "missing user", http.StatusNotFound, int(codes.NotFound)},
	} {
		t.Run(tc.message, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			Fail(c, tc.err)
			var body Response
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if w.Code != tc.status || body.Code != tc.code || body.Message != tc.message {
				t.Fatalf("unexpected response: HTTP %d %s", w.Code, w.Body)
			}
		})
	}
}
