package present

import "net/http"

type Health struct{}

func (h *Health) Check(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
