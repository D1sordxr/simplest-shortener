package present

import (
	"context"
	"simplest-shortener/pkg"
)

type Server struct {
	log           pkg.Log
	router        *Router
	dynamicRouter *DynamicRouter
}

func NewServer(log pkg.Log, shortenSvc svc) *Server {
	mid := NewMiddleware(log)
	handler := NewHandler(shortenSvc)
	router := NewRouter(mid, handler)
	dRouter := NewDynamicRouter()

	return &Server{
		log:           log,
		router:        router,
		dynamicRouter: dRouter,
	}
}

func (s *Server) Run(ctx context.Context) error {
	errChan := make(chan error, 1)

	go func() {
		errChan <- s.router.StartServer(":8080")
	}()

	go s.dynamicRouter.StartSettingUpRoutes(ctx, s.router.Mux)

	select {
	case <-ctx.Done():
		s.log.Info("Shutting down server gracefully...")
		return nil
	case err := <-errChan:
		s.log.Error("Server error", "err", err)
		return err
	}
}
