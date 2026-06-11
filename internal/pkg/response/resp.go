package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

// grpcToHTTP maps gRPC status codes to HTTP status codes.
func grpcToHTTP(code codes.Code) int {
	switch code {
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// Write writes a unified response. On error, it logs and uses the coder's
// HTTP status and code. On success, it returns HTTP 200 with data.
func Write(c *gin.Context, err error, data any) {
	if err != nil {
		if st, ok := status.FromError(err); ok {
			log.Errorf("gRPC error: code=%s message=%s", st.Code(), st.Message())
			c.JSON(grpcToHTTP(st.Code()), Response{
				Code:    int(st.Code()),
				Message: st.Message(),
			})
			return
		}
		log.Errorf("%v, %+v", err, err)
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
