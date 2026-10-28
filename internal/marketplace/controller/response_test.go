package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/internal/marketplace/service"
	"github.com/hanzhuoxian/mall/internal/marketplace/store"
)

func TestErrorResponsesPreserveContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tc := range []struct {
		err     error
		message string
		status  int
	}{
		{service.ErrInvalid, service.ErrInvalid.Error(), http.StatusBadRequest},
		{service.ErrStock, service.ErrStock.Error(), http.StatusBadRequest},
		{service.ErrForbidden, service.ErrForbidden.Error(), http.StatusForbidden},
		{store.ErrNotFound, "record not found", http.StatusNotFound},
		{service.ErrPaymentDisabled, service.ErrPaymentDisabled.Error(), http.StatusConflict},
		{errors.New("SQL connection contains secret"), "操作失败，请稍后重试", http.StatusInternalServerError},
	} {
		t.Run(tc.message, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			fail(c, tc.err)
			var body struct {
				Message string
				Code    int
			}
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if w.Code != tc.status || body.Code != tc.status || body.Message != tc.message || !c.IsAborted() {
				t.Fatalf("unexpected response: HTTP %d %s, aborted=%v", w.Code, w.Body, c.IsAborted())
			}
		})
	}
}

func TestEmptySuccessIncludesNullData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	respond(c, nil, nil)
	var body map[string]json.RawMessage
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if w.Code != http.StatusOK || string(body["code"]) != "0" || string(body["data"]) != "null" {
		t.Fatalf("unexpected success response: %s", w.Body)
	}
}
