package internal

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"simplest-shortener/internal/infra"
	"simplest-shortener/internal/present"
	"simplest-shortener/internal/svc"
	"simplest-shortener/pkg"
	"syscall"
)

type App struct {
	log    pkg.Log
	server *present.Server
}

func NewApp() *App {
	log := slog.Default()
	storage := infra.NewStorage()
	shortenSvc := svc.NewShortenerSvc(log, storage, present.NewDynamicRouter())
	server := present.NewServer(log, shortenSvc)

	return &App{
		log:    log,
		server: server,
	}
}

func (a *App) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := a.server.Run(ctx); err != nil {
		a.log.Error("Critical error during run", "error", err)
	}
}
