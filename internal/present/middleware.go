package present

import (
	"net/http"
	"simplest-shortener/pkg"
	"time"
)

type Middleware struct {
	log pkg.Log
	gen pkg.Generator
}

func NewMiddleware(log pkg.Log) *Middleware {
	return &Middleware{
		log: log,
		gen: new(pkg.SharedGenerator),
	}
}

const (
	newLenEnv = 10
)

func (m *Middleware) Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()

		logToken := m.gen.GenerateRandomString(newLenEnv)
		m.log.Info("Starting request...", "token", logToken, "time", now.String())

		handler.ServeHTTP(w, r)

		m.log.Info("Request finished", "token", logToken, "time", time.Since(now).String())
	})
}
