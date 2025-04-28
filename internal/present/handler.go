package present

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

const (
	ctxTimeoutEnv = time.Second * 5
)

type svc interface {
	Create(ctx context.Context, url string) string
	Get(ctx context.Context, code string) (string, bool)
}

type Handler struct {
	svc
}

func NewHandler(svc svc) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeoutEnv)
	defer cancel()

	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}
	code := h.svc.Create(ctx, url)
	response := map[string]string{"shortened_url": code}

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error getting url", http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(data)
}

func (h *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeoutEnv)
	defer cancel()

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code is required", http.StatusBadRequest)
		return
	}

	url, ok := h.svc.Get(ctx, code)
	if !ok {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	response := map[string]string{"url": url}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
