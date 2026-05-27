package server

import "net/http"

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
