package present

import (
	"context"
	"net/http"
	"sync"
)

type middleware interface {
	Log(handler http.Handler) http.Handler
}

type shortenerHandler interface {
	ShortenURL(w http.ResponseWriter, r *http.Request)
	GetURL(w http.ResponseWriter, r *http.Request)
}

type health interface {
	Check(w http.ResponseWriter, r *http.Request)
}

const (
	shortenPathEnv      = "/shorten"
	getShortenedPathEnv = "/get-shortened"
)

type Router struct {
	Mux  *http.ServeMux
	mid  middleware
	hand shortenerHandler
	health
}

func NewRouter(mid middleware, hand shortenerHandler) *Router {
	return &Router{
		Mux:    http.NewServeMux(),
		mid:    mid,
		hand:   hand,
		health: new(Health),
	}
}

func (r *Router) setupRoutes() {
	r.Mux.Handle("/health", http.HandlerFunc(r.health.Check))

	r.Mux.Handle(shortenPathEnv, r.mid.Log(http.HandlerFunc(r.hand.ShortenURL)))
	r.Mux.Handle(getShortenedPathEnv, r.mid.Log(http.HandlerFunc(r.hand.GetURL)))
}

func (r *Router) StartServer(addr string) error {
	r.setupRoutes()
	return http.ListenAndServe(addr, r.Mux)
}

const (
	workersCountEnv = 5
)

type job struct {
	mainURL      string
	shortenedURL string
}

type DynamicRouter struct {
	jobs chan job
}

func NewDynamicRouter() *DynamicRouter {
	return &DynamicRouter{
		jobs: make(chan job, 100),
	}
}

func (dr *DynamicRouter) AddJob(mainURL, shortenedURL string) {
	dr.jobs <- job{
		mainURL:      mainURL,
		shortenedURL: shortenedURL,
	}
}

func (dr *DynamicRouter) StartSettingUpRoutes(ctx context.Context, mux *http.ServeMux) {
	var wg sync.WaitGroup

	for i := 0; i < workersCountEnv; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case j, ok := <-dr.jobs:
					if !ok {
						return
					}
					mux.HandleFunc(j.shortenedURL, func(w http.ResponseWriter, r *http.Request) {
						http.Redirect(w, r, j.mainURL, http.StatusFound)
					})
				}
			}
		}(i)
	}

	go func() {
		<-ctx.Done()
		close(dr.jobs)
		wg.Wait()
	}()
}
