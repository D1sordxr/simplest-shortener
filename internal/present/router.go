package present

import (
	"context"
	"net/http"
	"sync"
)

type middleware interface {
	Log(handler http.HandlerFunc) http.HandlerFunc
}

type shortenerHandler interface {
	ShortenURL(w http.ResponseWriter, r *http.Request)
	GetURL(w http.ResponseWriter, r *http.Request)
}

const (
	shortenPathEnv      = "/shorten"
	getShortenedPathEnv = "/get-shortened"
)

type Router struct {
	Mux  *http.ServeMux
	mid  middleware
	hand shortenerHandler
}

func NewRouter(mid middleware, hand shortenerHandler) *Router {
	return &Router{
		Mux:  http.NewServeMux(),
		mid:  mid,
		hand: hand,
	}
}

func (r *Router) setupRoutes() {
	r.Mux.HandleFunc(shortenPathEnv, r.mid.Log(r.hand.ShortenURL))
	r.Mux.HandleFunc(getShortenedPathEnv, r.mid.Log(r.hand.GetURL))
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
	wg   sync.WaitGroup
	jobs chan job
}

func NewDynamicRouter() *DynamicRouter {
	return &DynamicRouter{
		wg:   sync.WaitGroup{},
		jobs: make(chan job),
	}
}

func (dr *DynamicRouter) AddJob(mainURL, shortenedURL string) {
	dr.jobs <- job{
		mainURL:      mainURL,
		shortenedURL: shortenedURL,
	}
}

func (dr *DynamicRouter) StartSettingUpRoutes(
	ctx context.Context,
	mux *http.ServeMux,
) {
	for i := 0; i < workersCountEnv; i++ {
		dr.wg.Add(1)
		go func(id int) {
			defer dr.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case j, ok := <-dr.jobs:
					if !ok {
						continue
					}
					mux.HandleFunc(j.shortenedURL, func(w http.ResponseWriter, r *http.Request) {
						http.Redirect(w, r, j.mainURL, http.StatusFound)
					})
				}
			}
		}(i)
	}

	close(dr.jobs)
	dr.wg.Wait()
}
