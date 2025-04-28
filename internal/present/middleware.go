package present

import (
	"math/rand"
	"net/http"
	"simplest-shortener/pkg"
	"time"
)

type Middleware struct {
	log pkg.Log
}

func NewMiddleware(log pkg.Log) *Middleware {
	return &Middleware{log: log}
}

func (m *Middleware) Log(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()

		logToken := generateRandomToken()
		m.log.Info("Starting request...", "token", logToken, "time", now.String())

		handler.ServeHTTP(w, r)

		m.log.Info("Request finished", "token", logToken, "time", time.Since(now).String())
	}
}

const (
	newLength = 8
	charset   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func generateRandomToken() string {
	newLen := rand.Intn(newLength) + 1

	code := make([]byte, newLen)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
