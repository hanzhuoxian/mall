package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/pkg/errors"
	"github.com/hanzhuoxian/mall/pkg/log"
)

// Response is the unified envelope for all API responses.
type Response struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Reference string `json:"reference,omitempty"`
	Data      any    `json:"data,omitempty"`
}

// Write writes a unified response. On error, it logs and uses the coder's
// HTTP status and code. On success, it returns HTTP 200 with data.
func Write(c *gin.Context, err error, data any) {
	if err != nil {
		log.Errorf("%#+v", err)
		coder := errors.ParseCoder(err)
		c.JSON(coder.HTTPStatus(), Response{
			Code:      int(coder.Code()),
			Message:   coder.String(),
			Reference: coder.Reference(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Success writes a HTTP 200 response with data.
func Success(c *gin.Context, data any) {
	Write(c, nil, data)
}

// Fail writes an error response derived from err.
func Fail(c *gin.Context, err error) {
	Write(c, err, nil)
}

// FailWith writes an error response with optional data payload.
func FailWith(c *gin.Context, err error, data any) {
	Write(c, err, data)
}
