package server

import "net/http"

// ServeHealthCheck 在指定地址启动一个独立的 HTTP 健康检查服务，
// 响应 path 路由的请求并返回 {"status": "ok"}。
// 注意：此函数会阻塞，通常在独立 goroutine 中调用。
func ServeHealthCheck(path string, address string) {
	http.HandleFunc("/"+path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}
