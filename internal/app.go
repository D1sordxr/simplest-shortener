package internal

import (
	"context"
	"log/slog"
	"simplest-shortener/internal/infra"
	"simplest-shortener/internal/present"
	"simplest-shortener/internal/svc"
	"simplest-shortener/pkg"
)

type App interface {
	Run()
}

type app struct {
	log           pkg.Log
	storage       *infra.Storage
	shortenSvc    *svc.ShortenerSvc
	mid           *present.Middleware
	handler       *present.Handler
	router        *present.Router
	dynamicRouter *present.DynamicRouter
}

func NewApp() App {
	log := slog.Default()
	log.Info("Starting application...")

	storage := infra.NewStorage()
	log.Info("Storage initialized")

	dRouter := present.NewDynamicRouter()

	shortenSvc := svc.NewShortenerSvc(log, storage, dRouter)
	log.Info("Shortener service initialized")

	mid := present.NewMiddleware(log)
	handler := present.NewHandler(shortenSvc)
	router := present.NewRouter(mid, handler)

	return &app{
		log:        log,
		storage:    storage,
		shortenSvc: shortenSvc,
		mid:        mid,
		handler:    handler,
		router:     router,
	}
}

func (a *app) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a.log.Info("Starting server...")
	err := a.router.StartServer(":8080")
	if err != nil {
		a.log.Error("Error starting server", "error", err)
	}

	a.dynamicRouter.StartSettingUpRoutes(ctx, a.router.Mux)

	a.log.Info("Server started successfully")
}
