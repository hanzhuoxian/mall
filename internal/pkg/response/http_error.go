package response

// HTTPError 为已有 HTTP 接口保留状态码和公开提示，Cause 仅用于服务端日志。
type HTTPError struct {
	Cause   error
	Message string
	Status  int
}

func (err *HTTPError) Error() string {
	if err.Cause != nil {
		return err.Cause.Error()
	}
	return err.Message
}

// Unwrap 返回原始错误，便于记录日志和定位问题。
func (err *HTTPError) Unwrap() error { return err.Cause }
