package response

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/hanzhuoxian/mall/pkg/errors"
	"github.com/hanzhuoxian/mall/pkg/logger"
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
			logger.Errorf("gRPC error: code=%s message=%s", st.Code(), st.Message())
			c.JSON(grpcToHTTP(st.Code()), Response{
				Code:    int(st.Code()),
				Message: st.Message(),
			})
			return
		}
		logger.Errorf("%v, %+v", err, err)
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
		Data:    marshalData(data),
	})
}

// marshalData 将 proto 消息用 protojson 序列化为 camelCase 字段、
// 时间戳为 RFC3339 字符串，供前端直接消费；非 proto 数据原样返回。
func marshalData(data any) any {
	msg, ok := data.(proto.Message)
	if !ok {
		return data
	}
	b, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(msg)
	if err != nil {
		return data
	}
	return json.RawMessage(b)
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
